package managers

import (
	"apiroute/cache"
	"apiroute/logs"
	"apiroute/models"
	"apiroute/util"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

type Updater struct {
	errCh       chan error
	syncMode    bool
	name        string
	namespace   string
	reloadCache *cache.ReloadCache
}

func (u *Updater) Task() {
	url := "http://" + clientIP + ":" + strconv.Itoa(int(clientPort)) + "/api/microservices/v1/services/" + u.name
	queryParams := "?namespace=" + u.namespace

	var details map[models.ServiceKey]*models.PublishInfo

	buf, err := util.HTTPGet(url, queryParams)
	if err != nil {
		logs.Log.Warn("Failed to query service details for %s:%v", u.name, err)
		sendError(u.errCh, err)
		return
	}

	data := []*util.ServiceUnit{}
	err = json.Unmarshal(buf, &data)

	if err != nil {
		logs.Log.Warn("Failed to Unmarshal data to []ServiceUnit:%v", err)
		sendError(u.errCh, err)
		return
	}

	if len(data) == 0 {
		logs.Log.Warn("%s does not exist in backend, ignore the update", u.name)
		return
	}

	protocol := strings.ToUpper(data[0].Protocol)

	if protocol == tcpProtocol || protocol == udpProtocol {
		var proceed bool
		for _, d := range data {
			if len(d.PublishPort) != 0 {
				proceed = true
				break
			}
		}
		if !proceed {
			logs.Log.Info("ignore tcp/udp service:%s with no publish port", data[0].Name)
			return
		}
	}

	if u.syncMode {
		goto SYNC_MODE
	}

	if details, err = GetDataSyncManager().GetPublishInfo(u.namespace, u.name, "all"); err != nil {
		logs.Log.Warn("Failed to GetPublishInfo:%v", err)
		sendError(u.errCh, err)
		return
	}

	if protocol == httpProtocol || protocol == restProtocol ||
		protocol == uiProtocol || protocol == portalProtocol {

		var (
			oldRefMap = make(map[string]int)
			newRefMap = make(map[string]int)
		)

		for _, v := range details {
			if v.PublishPort != "" {
				if num, ok := oldRefMap[v.PublishPort]; ok {
					oldRefMap[v.PublishPort] = num + v.NumOfNode
				} else {
					oldRefMap[v.PublishPort] = v.NumOfNode
				}
			}
		}

		for _, d := range data {
			if d.PublishPort != "" {
				if num, ok := newRefMap[d.PublishPort]; ok {
					newRefMap[d.PublishPort] = num + len(d.Instances)
				} else {
					newRefMap[d.PublishPort] = len(d.Instances)
				}
			}
		}

		for ko, vo := range oldRefMap {
			var (
				vn int
				ok bool
			)
			ports := strings.Split(ko, "|")
			if vn, ok = newRefMap[ko]; !ok {
				vn = 0
			}
			if !u.reloadCache.IsPortAlreadyUsedInStream(ports[0]) {
				u.reloadCache.UpdateHTTPSPort(ports[0], vo, vn)
			} else {
				logs.Log.Error("The port:%s for %s service:%s is already used in stream, ignore it", ports[0], protocol, u.name)
			}

			if len(ports) == 2 {
				if !u.reloadCache.IsPortAlreadyUsedInStream(ports[1]) {
					u.reloadCache.UpdateHTTPPort(ports[1], vo, vn)
				} else {
					logs.Log.Error("The port:%s for %s service:%s is already used in stream, ignore it", ports[1], protocol, u.name)
				}
			}
		}

		for kn, vn := range newRefMap {
			ports := strings.Split(kn, "|")
			if _, ok := oldRefMap[kn]; !ok {
				if !u.reloadCache.IsPortAlreadyUsedInStream(ports[0]) {
					u.reloadCache.UpdateHTTPSPort(ports[0], 0, vn)
				} else {
					logs.Log.Error("The port:%s for %s service:%s is already used in stream, ignore it", ports[0], protocol, u.name)
				}
				if len(ports) == 2 {
					if !u.reloadCache.IsPortAlreadyUsedInStream(ports[1]) {
						u.reloadCache.UpdateHTTPPort(ports[1], 0, vn)
					} else {
						logs.Log.Error("The port:%s for %s service:%s is already used in stream, ignore it", ports[1], protocol, u.name)
					}
				}
			}
		}
	}

SYNC_MODE:
	if err = GetDataSyncManager().Update(data); err != nil {
		logs.Log.Warn("Failed to update redis:%v", err)
		sendError(u.errCh, err)
		return
	}

	if u.syncMode {
		return
	}

	if protocol == tcpProtocol || protocol == udpProtocol {
		if len(details) == 0 {
			for _, d := range data {
				if u.reloadCache.IsPortAlreadyUsedInHTTP(d.PublishPort) {
					logs.Log.Error("The port:%s for %s service:%s:%s is already used in http/https, ignore it", d.PublishPort, protocol, u.name, d.Version)
					continue
				}

				if u.reloadCache.IsStreamPortAlreadyUsedInStream(d.PublishPort, protocol) {
					logs.Log.Error("The port:%s for %s service:%s:%s is already used in stream, ignore it", d.PublishPort, protocol, u.name, d.Version)
					continue
				}
				u.updateStream(d, protocol)
			}
		} else {
			ns := u.namespace
			if ns == "" {
				ns = "default"
			}

			for _, d := range data {
				key := models.ServiceKey{
					Namespace:      ns,
					ServiceName:    u.name,
					ServiceVersion: d.Version,
				}
				if val, ok := details[key]; ok {
					delete(details, key)
					searchName := key.ServiceName
					if key.Namespace != "" && key.Namespace != "default" {
						searchName = searchName + "-" + key.Namespace
					}
					if key.ServiceVersion != "" {
						searchName = searchName + "-" + key.ServiceVersion
					}

					if d.PublishPort == val.PublishPort && u.reloadCache.HasStream(searchName, protocol) {
						u.deleteStream(key, val.PublishPort, val.Protocol)
					} else {
						if protocol != udpProtocol && u.reloadCache.IsPortAlreadyUsedInHTTP(d.PublishPort) {
							logs.Log.Error("The port:%s for %s service:%s:%s is already used in http/https, ignore it", d.PublishPort, protocol, u.name, d.Version)
							continue
						}

						if u.reloadCache.IsStreamPortAlreadyUsedInStream(d.PublishPort, protocol) {
							logs.Log.Error("The port:%s for %s service:%s:%s is already used in stream, ignore it", d.PublishPort, protocol, u.name, d.Version)
							continue
						}
						u.updateStream(d, protocol)
					}
					u.updateStream(d, protocol)
				} else {
					if protocol != udpProtocol && u.reloadCache.IsPortAlreadyUsedInHTTP(d.PublishPort) {
						logs.Log.Error("The port:%s for %s service:%s:%s is already used in http/https, ignore it", d.PublishPort, protocol, u.name, d.Version)
						continue
					}

					if u.reloadCache.IsStreamPortAlreadyUsedInStream(d.PublishPort, protocol) {
						logs.Log.Error("The port:%s for %s service:%s:%s is already used in stream, ignore it", d.PublishPort, protocol, u.name, d.Version)
						continue
					}
					u.updateStream(d, protocol)
				}
			}
			for k, v := range details {
				u.deleteStream(k, v.PublishPort, v.Protocol)
			}
		}
	}
}

func (u *Updater) deleteStream(key models.ServiceKey, port, protocol string) {
	name := key.ServiceName
	if key.Namespace != "" && key.Namespace != "default" {
		name = name + "-" + key.Namespace

	}
	if key.ServiceVersion != "" {
		name = name + "-" + key.ServiceVersion

	}
	u.reloadCache.DeleteStream(name, port, protocol)
}

func (u *Updater) updateStream(su *util.ServiceUnit, protocol string) {
	var (
		udp, tls   string
		proxyRules []*cache.Rule
	)

	if protocol == udpProtocol {
		udp = "true"

	}

	if su.EnableTLS {
		tls = "true"

	}

	svrList := make([]*cache.StreamServer, len(su.Instances))
	name := u.name
	if su.Namespace != "" {
		name = name + "-" + su.Namespace

	}
	if su.Version != "" {
		name = name + "-" + su.Version

	}
	var isProxyTimeoutSet bool
	rules := reflect.ValueOf(su.ProxyRule.StreamProxy)
	for i := 0; i < rules.Type().NumField(); i++ {
		if rules.Field(i).String() != "" {
			rname := strings.Split(rules.Type().Field(i).Tag.Get("json"), ",")[0]
			if rname == cache.ProxyTimeout {
				isProxyTimeoutSet = true

			}
			proxyRules = append(proxyRules, &cache.Rule{
				Name:  rname,
				Value: rules.Field(i).String(),
			})

		}

	}
	if !isProxyTimeoutSet && udp == "true" {
		proxyRules = append(proxyRules, &cache.Rule{
			Name:  cache.ProxyTimeout,
			Value: cache.DefaultUDPProxyTimeout,
		})

	}
	for index := 0; index < len(su.Instances); index++ {
		lbParams := strings.Replace(su.Instances[index].LBServerParams, ",", " ", -1)
		svrList[index] = &cache.StreamServer{
			Addr:   su.Instances[index].ServiceIP + ":" + su.Instances[index].ServicePort,
			Params: lbParams,
		}
	}
	//MSB-1087 && MSB-1095
	lb := su.LBPolicy
	switch lb {
	case "ip_hash":
		lb = "hash $remote_addr"
	case "least_conn":
	default:
		lb = ""
	}

	streamMetaData := &cache.StreamMetaData{
		IsUDP:       udp,
		PublishPort: su.PublishPort,
		LBPolicy:    lb,
		EnableTLS:   tls,
		ProxyRules:  proxyRules,
		ServerList:  svrList,
	}
	u.reloadCache.UpdateStream(name, protocol, streamMetaData)
}
