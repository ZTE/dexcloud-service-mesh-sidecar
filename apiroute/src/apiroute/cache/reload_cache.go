package cache

import (
	"apiroute/logs"
	"os/exec"
	"sync"
)

const (
	OpAdd                  = "add"
	OpDelete               = "delete"
	ProxyTimeout           = "proxy_timeout"
	DefaultUDPProxyTimeout = "20s"
	reloadScript           = "../openresty/reload.sh"
)

type Rule struct {
	Name  string
	Value string
}

type StreamServer struct {
	Addr   string
	Params string
}

type StreamMetaData struct {
	IsUDP       string
	PublishPort string
	LBPolicy    string
	EnableTLS   string
	ProxyRules  []*Rule
	ServerList  []*StreamServer
}

type CacheData struct {
	HTTPPortRef  map[string]int
	HTTPSPortRef map[string]int
	TCPPorts     map[string]string
	UDPPorts     map[string]string
	StreamCache  map[string]*StreamMetaData
}

type ReloadCache struct {
	sync.Mutex
	httpPortRef         map[string]int
	httpsPortRef        map[string]int
	tcpPorts            map[string]string
	udpPorts            map[string]string
	streamCache         map[string]*StreamMetaData
	httpPortNeedUpdate  bool
	httpsPortNeedUpdate bool
	streamNeedUpdate    bool
}

func NewReloadCache() *ReloadCache {
	return &ReloadCache{
		httpPortRef:  make(map[string]int),
		httpsPortRef: make(map[string]int),
		tcpPorts:     make(map[string]string),
		udpPorts:     make(map[string]string),
		streamCache:  make(map[string]*StreamMetaData),
	}
}

func (rc *ReloadCache) copyData() (data *CacheData) {
	data = &CacheData{
		HTTPPortRef:  make(map[string]int),
		HTTPSPortRef: make(map[string]int),
		TCPPorts:     make(map[string]string),
		UDPPorts:     make(map[string]string),
		StreamCache:  make(map[string]*StreamMetaData),
	}

	for k, v := range rc.httpPortRef {
		data.HTTPPortRef[k] = v
	}

	for k, v := range rc.httpsPortRef {
		data.HTTPSPortRef[k] = v
	}

	for k, v := range rc.tcpPorts {
		data.TCPPorts[k] = v
	}

	for k, v := range rc.udpPorts {
		data.UDPPorts[k] = v
	}

	for k, v := range rc.streamCache {
		data.StreamCache[k] = v
	}
	return
}

func (rc *ReloadCache) PrintL7PortRef() {
	logs.Log.Info("HTTP Port Reference details:[port] ==> [reference times]")
	for k, v := range rc.httpPortRef {
		logs.Log.Info("%s ==> %d", k, v)
	}

	logs.Log.Info("HTTPS Port Reference details:[port] ==> [reference times]")
	for k, v := range rc.httpsPortRef {
		logs.Log.Info("%s ==> %d", k, v)
	}
}

func (rc *ReloadCache) GetData() *CacheData {
	return rc.copyData()
}

func (rc *ReloadCache) RestoreAll(cd *CacheData) {
	rc.httpPortRef = cd.HTTPPortRef
	rc.httpsPortRef = cd.HTTPSPortRef
	rc.tcpPorts = cd.TCPPorts
	rc.udpPorts = cd.UDPPorts
	rc.streamCache = cd.StreamCache
	rc.httpPortNeedUpdate = false
	rc.httpsPortNeedUpdate = false
	rc.streamNeedUpdate = false
}

func (rc *ReloadCache) ResetReloadCache() *CacheData {
	if len(rc.httpPortRef) > 0 {
		rc.httpPortRef = make(map[string]int)
	}
	if len(rc.httpsPortRef) > 0 {
		rc.httpsPortRef = make(map[string]int)
	}
	if len(rc.tcpPorts) > 0 {
		rc.tcpPorts = make(map[string]string)
	}
	if len(rc.udpPorts) > 0 {
		rc.udpPorts = make(map[string]string)
	}
	if len(rc.streamCache) > 0 {
		rc.streamCache = make(map[string]*StreamMetaData)
	}
	rc.ResetFlags(true)

	return &CacheData{
		HTTPPortRef:  rc.httpPortRef,
		HTTPSPortRef: rc.httpsPortRef,
		TCPPorts:     rc.tcpPorts,
		UDPPorts:     rc.udpPorts,
		StreamCache:  rc.streamCache,
	}
}

func (rc *ReloadCache) ResetFlags(val bool) {
	rc.httpPortNeedUpdate = val
	rc.httpsPortNeedUpdate = val
	rc.streamNeedUpdate = val
}

func (rc *ReloadCache) UpdateHTTPPort(port string, oldVal, newVal int) {
	rc.Lock()
	defer rc.Unlock()
	var (
		delta int
		op    string
	)
	if delta = oldVal - newVal; delta >= 0 {
		op = OpDelete
	} else {
		op = OpAdd
		delta = -delta
	}

	if ref, ok := rc.httpPortRef[port]; ok {
		switch op {
		case OpAdd:
			rc.httpPortRef[port] = ref + delta
		case OpDelete:
			if ref <= delta {
				delete(rc.httpPortRef, port)
				rc.httpPortNeedUpdate = true

			} else {
				rc.httpPortRef[port] = ref - delta
			}
		}
	} else {
		rc.httpPortRef[port] = newVal
		rc.httpPortNeedUpdate = true
	}
}

func (rc *ReloadCache) UpdateHTTPSPort(port string, oldVal, newVal int) {
	rc.Lock()
	defer rc.Unlock()
	var (
		delta int
		op    string
	)
	if delta = oldVal - newVal; delta >= 0 {
		op = OpDelete
	} else {
		op = OpAdd
		delta = -delta
	}

	if ref, ok := rc.httpsPortRef[port]; ok {
		switch op {
		case OpAdd:
			rc.httpsPortRef[port] = ref + delta
		case OpDelete:
			if ref <= delta {
				delete(rc.httpsPortRef, port)
				rc.httpsPortNeedUpdate = true

			} else {
				rc.httpsPortRef[port] = ref - delta
			}
		}
	} else {
		rc.httpsPortRef[port] = newVal
		rc.httpsPortNeedUpdate = true
	}
}

func (rc *ReloadCache) HasStream(key, protocol string) bool {
	rc.Lock()
	defer rc.Unlock()

	switch protocol {
	case "TCP", "UDP":
		if _, ok := rc.streamCache[key]; ok {
			return true
		}
	default:
		logs.Log.Warn("Unknown protocol;%s", protocol)
	}
	return false
}

func (rc *ReloadCache) DeleteStream(name, port, protocol string) {
	rc.Lock()
	defer rc.Unlock()
	if _, ok := rc.streamCache[name]; ok {
		delete(rc.streamCache, name)
		switch protocol {
		case "TCP":
			if _, ok := rc.tcpPorts[port]; ok {
				delete(rc.tcpPorts, port)
			}
		case "UDP":
			if _, ok := rc.udpPorts[port]; ok {
				delete(rc.udpPorts, port)
			}
		}
		rc.streamNeedUpdate = true
	}
}

func (rc *ReloadCache) UpdateStream(name, protocol string, stream *StreamMetaData) {
	rc.Lock()
	defer rc.Unlock()
	rc.streamCache[name] = stream
	rc.streamNeedUpdate = true
	switch protocol {
	case "TCP":
		rc.tcpPorts[stream.PublishPort] = stream.PublishPort
	case "UDP":
		rc.udpPorts[stream.PublishPort] = stream.PublishPort
	}
}

func (rc *ReloadCache) CheckPortUpdate() (httpList, httpsList []string, updateHTTP, updateHTTPS bool) {
	rc.Lock()
	defer rc.Unlock()
	if rc.httpPortNeedUpdate {
		updateHTTP = true
		for name := range rc.httpPortRef {
			httpList = append(httpList, name)
		}
	}

	if rc.httpsPortNeedUpdate {
		updateHTTPS = true
		for name := range rc.httpsPortRef {
			httpsList = append(httpsList, name)
		}
	}
	return
}

func (rc *ReloadCache) CheckStreamUpdate() (streamList map[string]*StreamMetaData, needUpdate bool) {
	rc.Lock()
	defer rc.Unlock()
	if rc.streamNeedUpdate {
		return rc.streamCache, true
	}
	return
}

func (rc *ReloadCache) IsPortAlreadyUsedInStream(port string) bool {
	rc.Lock()
	defer rc.Unlock()
	if _, ok := rc.tcpPorts[port]; ok {
		return true
	}

	if _, ok := rc.udpPorts[port]; ok {
		return true
	}
	return false
}

func (rc *ReloadCache) IsPortAlreadyUsedInHTTP(port string) bool {
	rc.Lock()
	defer rc.Unlock()
	if _, ok := rc.httpPortRef[port]; ok {
		return true
	}
	if _, ok := rc.httpsPortRef[port]; ok {
		return true
	}
	return false
}

func (rc *ReloadCache) IsStreamPortAlreadyUsedInStream(port, protocol string) bool {
	rc.Lock()
	defer rc.Unlock()
	switch protocol {
	case "TCP":
		if _, ok := rc.tcpPorts[port]; ok {
			return true
		}
	case "UDP":
		if _, ok := rc.udpPorts[port]; ok {
			return true
		}
	}
	return false
}

func Reloading() error {
	out, err := exec.Command(reloadScript).CombinedOutput()
	if err != nil {
		logs.Log.Error("Error with reloading:%v", err)
		return err
	}

	logs.Log.Error("Error with reloading:%v msg:%s", err, out)
	return nil
}
