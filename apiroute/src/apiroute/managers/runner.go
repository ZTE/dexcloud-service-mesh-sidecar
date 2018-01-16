package managers

import (
	"apiroute/cache"
	"apiroute/logs"
	"apiroute/managers/configmanager"
	"apiroute/redis"
	"apiroute/util"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

const (
	syncInterval = 300
)

var (
	clientIP                     = configmanager.CfgM.GetConfigInfo().Discover.IP
	clientPort                   = configmanager.CfgM.GetConfigInfo().Discover.Port
	namespace                    = configmanager.CfgM.GetConfigInfo().Apigatewaycfg.Namespace
	labels                       = configmanager.CfgM.GetConfigInfo().Apigatewaycfg.Lables
	customFilterConfig           = configmanager.CfgM.GetConfigInfo().Apigatewaycfg.CustomFilterConfig
	routeWay                     = configmanager.CfgM.GetConfigInfo().Apigatewaycfg.RouteWay
	routeSubdomain               = configmanager.CfgM.GetConfigInfo().Apigatewaycfg.RouteSubdomain
	publishPort                  = configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport
	httpsPort                    = configmanager.CfgM.GetConfigInfo().Listenport.Httpsdefaultport
	serviceIP                    = configmanager.CfgM.GetConfigInfo().Apigatewaycfg.ServiceIP
	metricsIP                    = configmanager.CfgM.GetConfigInfo().Apigatewaycfg.MetricsIP
	registerName                 string
	serviceManager               *ServiceManager
	routeManager                 *RouteManager
	releationManager             *ReleationManager
	namespaceCache               *NamespaceManager
	dataSyncManager              *DataSyncManager
	digestCacheToRedisSyncLocker = &sync.Mutex{}
)

type Runner struct {
	ErrCh        chan error
	watcher      *Watcher
	dispatcher   *Dispatcher
	dataCh       chan []*ServiceDigestData
	syncerStopCh chan struct{}
}

func NewRunner() *Runner {
	return &Runner{
		ErrCh:        make(chan error),
		watcher:      NewWatcher(),
		dataCh:       make(chan []*ServiceDigestData),
		syncerStopCh: make(chan struct{}),
	}
}

func (r *Runner) Init() {
	digestCache := cache.NewServiceDigestCache()
	clients := redis.NewClients()
	reloadCache := cache.NewReloadCache()
	//
	serviceManager = &ServiceManager{clients.TempClient}
	routeManager = &RouteManager{clients.RouterClient}
	releationManager = &ReleationManager{clients.MapperClient}
	//
	InitRoutes()
	//
	namespaceCache = NewNamespaceManager()
	dataSyncManager = NewDataSyncManager()
	dataSyncManager.RegisterListener(RouteListener{})
	dataSyncManager.RegisterListener(NamespaceLister{})
	r.dispatcher = NewDispatcher(digestCache, reloadCache)
	clients.StartHealthCheck()
}

func (r *Runner) Start() {
	logs.Log.Info("Runner starts...")
	logs.Log.Info("Starts Dispatcher...")
	go r.dispatcher.StartDispatch(r.dataCh)
	logs.Log.Info("Starts Watcher...")
	go r.watcher.StartWatch(r.dataCh)
	logs.Log.Info("Starts beego...")
	go beego.Run()
	registerMe()
	r.StartSyncer(r.dispatcher)

	select {
	case err := <-r.watcher.ErrCh:
		r.StopSyncer()
		deregisterMe()
		logs.Log.Error("Watcher returned error:%v, start to terminate Runner", err)
		logs.Log.Info("Start to terminate Dispatcher...")
		r.dispatcher.StopDispatch()
		logs.Log.Info("Stop Health Check...")
		redis.StopHealthCheck()
		namespaceCache.CloseRefreshRoutine()
		logs.Log.Info("Propagate the error to the CLI...")
		r.ErrCh <- err
	}
}

func (r *Runner) Stop() {
	r.StopSyncer()
	deregisterMe()
	logs.Log.Info("Stop Runner...")
	logs.Log.Info("Start to stop Watcher...")
	r.watcher.StopWatch()
	logs.Log.Info("Start to stop Dispatcher...")
	r.dispatcher.StopDispatch()
	logs.Log.Info("Runner stopped")
	logs.Log.Info("Stop Health Check...")
	redis.StopHealthCheck()
	namespaceCache.CloseRefreshRoutine()
}

func (r *Runner) StartSyncer(dsp *Dispatcher) {
	ticker := time.NewTicker(time.Duration(syncInterval) * time.Second)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				logs.Log.Info("Sync Redis with Digest Cache at:%s", t.String())
				digestCacheToRedisSyncLocker.Lock()
				dataInRedis, err := GetDataSyncManager().GetAllServiceDigest()
				if err != nil {
					logs.Log.Warn("Sync failed with error:%v, wait for next Sync", err)
					digestCacheToRedisSyncLocker.Unlock()
					continue
				}
				delList, UpdateList := dsp.DigestCache.Diff2(dataInRedis)
				if len(delList) == 0 && len(UpdateList) == 0 {
					logs.Log.Info("Redis and Digest Cache are the same")
					digestCacheToRedisSyncLocker.Unlock()
					continue
				}

				logs.Log.Info("Redis and Digest Cache are not the same, start to sync them...")
				if err = dsp.startWorkers(delList, UpdateList, true); err != nil {
					logs.Log.Warn("Sync failed with error:%v, wait for next Sync", err)
				} else {
					logs.Log.Info("Sync Redis with Digest Cache Completed")
				}
				digestCacheToRedisSyncLocker.Unlock()

			case <-r.syncerStopCh:
				logs.Log.Info("Stop Digest Cache to Redis Syncer")
				ticker.Stop()
				return
			}
		}
	}()
}

func (r *Runner) StopSyncer() {
	close(r.syncerStopCh)
}

func registerMe() {
	logs.Log.Info("Start to register apigateway services")
	base := "http://" + clientIP + ":" + strconv.Itoa(int(clientPort)) + "/api/microservices/v1/services"
	var (
		retryCount = 12
		count      int
	)

	metaMap := make(map[string]string)
	if namespace != "" {
		metaMap["namespace"] = namespace
	}

	if routeWay != "" {
		metaMap["routeWay"] = routeWay
	}

	if routeSubdomain != "" {
		metaMap["routeSubdomain"] = routeSubdomain
	}

	metaMap["httpPort"] = strconv.Itoa(int(publishPort))
	metaMap["httpsPort"] = strconv.Itoa(int(httpsPort))

	var (
		vs = "0" //default value
		// 		labelsArr []string
	)

	lbs := strings.Split(labels, ",")
	for _, v := range lbs {
		kv := strings.Split(v, ":")
		if len(kv) == 2 {
			metaMap[kv[0]] = kv[1]
			if kv[0] == "visualRange" {
				vs = kv[1]
			} //else {
			// 				labelsArr = append(labelsArr, v)
			// 			}
		}
	}

	//load registerName
	if vs == "0" {
		registerName = "router"
	} else if vs == "1" {
		registerName = "apigateway"
	} else {
		registerName = "router"
	}

	//Merge with custom filter configs
	customcfgs := strings.Split(customFilterConfig, ",")
	for _, v := range customcfgs {
		kv := strings.Split(v, ":")
		if len(kv) == 2 {
			metaMap[kv[0]] = kv[1]
		}
	}

	metaData := make([]util.MetaUnit, len(metaMap))
	index := 0
	for k, v := range metaMap {
		metaData[index] = util.MetaUnit{
			Key:   k,
			Value: v,
		}
		index++
	}

	if serviceIP == "" {
		serviceIP = util.GetOutboundIP()
	}

	//Register apigateway/router in backend
	instance := make([]util.InstanceUnit, 1)
	instance[0] = util.InstanceUnit{
		ServiceIP:   serviceIP,
		ServicePort: strconv.Itoa(int(publishPort)),
	}

	serviceUnit := util.ServiceUnit{
		Name:        registerName,
		Version:     "v1",
		URL:         "/api/route/v1",
		Protocol:    "REST",
		VisualRange: vs,
		Instances:   instance,
		Metadata:    metaData,
	}

	out, _ := json.Marshal(serviceUnit)
	for {
		if err := util.HTTPPost(base, "", out); err != nil {
			if count < retryCount {
				logs.Log.Error("Register apigateway/router with error:%v, sleep 10 secs and retry the %dth time", err, count+1)
				time.Sleep(10 * time.Second)
				count++
				continue
			}
			logs.Log.Error("Register apigateway/router with error:%v, max retry exhausted", err)
		}
		break
	}
	// 	count = 0
	//
	// 	if registerName == "router" {
	// 		logs.Log.Info("End to register apigateway services")
	// 		return
	// 	}
	// 	//Register Metrics service in backend
	// 	if metricsIP == "" {
	// 		metricsIP = "127.0.0.1"
	// 	}

	//Construct Labels and namespace for metrics, to make it same with the apigateway
	// 	var mlabels []string
	// 	for k, v := range metaMap {
	// 		if k == "namespace" || k == "routeWay" || k == "routeSubdomain" || k == "httpPort" ||
	// 			k == "httpsPort" || k == "visualRange" {
	// 			continue
	// 		}
	// 		mlabels = append(mlabels, k+":"+v)
	// 	}
	//
	// 	var ns string
	// 	if val, ok := metaMap["namespace"]; ok {
	// 		ns = val
	// 	}

	// 	instanceMetrics := make([]util.InstanceUnit, 1)
	// 	instanceMetrics[0] = util.InstanceUnit{
	// 		ServiceIP:   metricsIP,
	// 		ServicePort: strconv.Itoa(int(publishPort)),
	// 	}
	//
	// 	metricsServiceUnit := util.ServiceUnit{
	// 		Name:        registerName + "_metrics",
	// 		Version:     "v1",
	// 		URL:         "/admin/microservices/v1",
	// 		Protocol:    "REST",
	// 		VisualRange: "0|1",
	// 		Namespace:   namespace,
	// 		Instances:   instanceMetrics,
	// 		Labels:      labelsArr,
	// 	}
	//
	// 	out, _ = json.Marshal(metricsServiceUnit)
	// 	for {
	// 		if err := util.HTTPPost(base, "", out); err != nil {
	// 			if count < retryCount {
	// 				logs.Log.Error("Register metrics with error:%v, sleep 10 secs and retry %dth time", err, count+1)
	// 				time.Sleep(10 * time.Second)
	// 				count++
	// 				continue
	// 			}
	// 			logs.Log.Error("Register metrics with error:%v, max retry exhausted", err)
	// 		}
	// 		break
	// 	}

	logs.Log.Info("End to register apigateway services")
}

func deregisterMe() {
	logs.Log.Info("Start to deregister apigateway services")

	//Deregister apigateway/router
	base := "http://" + clientIP + ":" + strconv.Itoa(int(clientPort)) + "/api/microservices/v1/services/" +
		registerName + "/version/v1/nodes/" + serviceIP + "/" + strconv.Itoa(int(publishPort))

	if err := util.HTTPDelete(base, ""); err != nil {
		logs.Log.Error("Deregister apigateway/router service with error:%v", err)
	}

	//Deregister metrics service
	// 	base = "http://" + clientIP + ":" + strconv.Itoa(int(clientPort)) + "/api/microservices/v1/services/" +
	// 		registerName + "_metrics" + "/version/v1/nodes/" + metricsIP + "/" + strconv.Itoa(int(publishPort))
	//
	// 	if err := util.HTTPDelete(base, ""); err != nil {
	// 		logs.Log.Error("Deregister metrics service with error:%v", err)
	// 	}

	logs.Log.Info("End to deregister apigateway services")
}

func getMergedLabels() string {
	labelMap := make(map[string]string)
	if namespace != "" {
		labelMap["namespace"] = namespace
	}

	lbs := strings.Split(labels, ",")
	for _, v := range lbs {
		kv := strings.Split(v, ":")
		if len(kv) == 2 {
			labelMap[kv[0]] = kv[1]
		}
	}

	//Merge with custom filter configs
	customcfgs := strings.Split(customFilterConfig, ",")
	for _, v := range customcfgs {
		kv := strings.Split(v, ":")
		if len(kv) == 2 {
			labelMap[kv[0]] = kv[1]
		}
	}

	//Sort keys
	var sortedKeys []string
	for k := range labelMap {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var out string
	for _, k := range sortedKeys {
		out = out + k + ":" + labelMap[k] + ","
	}

	return strings.TrimSuffix(out, ",")
}

func GetServiceManager() *ServiceManager {
	return serviceManager
}

func GetRouteManager() *RouteManager {
	return routeManager
}

func GetReleationManager() *ReleationManager {
	return releationManager
}

func GetNamespaceManager() *NamespaceManager {
	return namespaceCache
}

func GetDataSyncManager() *DataSyncManager {
	return dataSyncManager
}
