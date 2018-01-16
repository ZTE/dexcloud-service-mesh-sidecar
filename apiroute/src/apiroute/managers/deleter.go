package managers

import (
	"apiroute/cache"
	"apiroute/models"
	"apiroute/util"
	"strings"
)

const (
	httpProtocol   = "HTTP"
	restProtocol   = "REST"
	uiProtocol     = "UI"
	portalProtocol = "PORTAL"
	tcpProtocol    = "TCP"
	udpProtocol    = "UDP"
)

type Deleter struct {
	errCh            chan error
	syncMode         bool
	nameAndNamespace *util.NameAndNamespace
	reloadCache      *cache.ReloadCache
}

func (d *Deleter) Task() {
	var (
		err            error
		details        map[models.ServiceKey]*models.PublishInfo
		isHTTP         bool
		isStream       bool
		streamProtocol string
	)

	if exist, err := GetDataSyncManager().ServiceExist(d.nameAndNamespace); err == nil {
		if !exist {
			return
		}
	} else {
		sendError(d.errCh, err)
		return
	}

	if d.syncMode {
		goto SYNC_MODE
	}

	if details, err = GetDataSyncManager().GetPublishInfo(d.nameAndNamespace.Namespace, d.nameAndNamespace.Name, "all"); err != nil {
		sendError(d.errCh, err)
		return
	}

	for _, v := range details {
		protocol := strings.ToUpper(v.Protocol)
		if protocol == httpProtocol ||
			protocol == restProtocol ||
			protocol == uiProtocol ||
			protocol == portalProtocol {
			isHTTP = true
			break
		}

		if protocol == tcpProtocol ||
			protocol == udpProtocol {
			isStream = true
			streamProtocol = protocol
			break
		}
		break
	}

	if isHTTP {
		for _, v := range details {
			ports := strings.Split(v.PublishPort, "|")
			if len(ports) == 2 {
				if d.reloadCache.IsPortAlreadyUsedInHTTP(ports[1]) {
					d.reloadCache.UpdateHTTPPort(ports[1], v.NumOfNode, 0)
				}
				if d.reloadCache.IsPortAlreadyUsedInHTTP(ports[0]) {
					d.reloadCache.UpdateHTTPSPort(ports[0], v.NumOfNode, 0)
				}
			} else {
				if d.reloadCache.IsPortAlreadyUsedInHTTP(ports[0]) {
					d.reloadCache.UpdateHTTPSPort(ports[0], v.NumOfNode, 0)
				}
			}
		}
	}

	if isStream {
		for k, v := range details {
			if d.reloadCache.IsPortAlreadyUsedInHTTP(v.PublishPort) {
				continue
			}
			var key string
			key = k.ServiceName
			if k.Namespace != "" && k.Namespace != "default" {
				key = key + "-" + k.Namespace
			}
			if k.ServiceVersion != "" {
				key = key + "-" + k.ServiceVersion
			}
			d.reloadCache.DeleteStream(key, v.PublishPort, streamProtocol)
		}
	}

SYNC_MODE:
	if err := GetDataSyncManager().Delete(d.nameAndNamespace); err != nil {
		sendError(d.errCh, err)
		return
	}
}
