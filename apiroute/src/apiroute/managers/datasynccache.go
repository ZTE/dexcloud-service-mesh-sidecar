package managers

import (
	"apiroute/logs"
	"apiroute/models"
	"apiroute/util"
	"sync"
	"time"
)

type DataSyncCache struct {
	sync.Mutex
	cache map[util.NameAndNamespace]map[string]string
}

func NewDataSyncCache() *DataSyncCache {
	dsc := &DataSyncCache{}
	dsc.initDataSyncCache()
	return dsc
}

func (dsc *DataSyncCache) initDataSyncCache() {
	var (
		pServiceKeys []*models.ServiceKey
		redisErr     error
	)

	//read keys from redis
	dsc.cache = make(map[util.NameAndNamespace]map[string]string)
	sm := GetServiceManager()

	for retryCount := 0; retryCount < retryTotalCount; retryCount++ {
		pServiceKeys, redisErr = sm.GetServiceKeyList()
		if redisErr != nil {
			logs.Log.Warn("sync data cache failed:%s", redisErr.Error())
			time.Sleep(retryInterval)
			continue
		}
		break
	}

	if redisErr != nil {
		logs.Log.Warn("scan service keys after retry:%d ,still failed, break", retryTotalCount)
		return
	}

	if pServiceKeys != nil && len(pServiceKeys) != 0 {
		for _, pServiceKey := range pServiceKeys {

			if versions, ok := dsc.cache[util.NameAndNamespace{pServiceKey.ServiceName,
				pServiceKey.Namespace}]; ok {
				versions[pServiceKey.ServiceVersion] = ""
			} else {
				versions = make(map[string]string)
				versions[pServiceKey.ServiceVersion] = ""
				dsc.cache[util.NameAndNamespace{pServiceKey.ServiceName,
					pServiceKey.Namespace}] = versions
			}
		}
	}
}

func (dsc *DataSyncCache) IsExist(namespace string, serviceName string) bool {
	dsc.Lock()
	defer dsc.Unlock()
	changedNamespace := ChangeEmptyToDefault(namespace)
	_, ok := dsc.cache[util.NameAndNamespace{serviceName, changedNamespace}]
	return ok
}

func (dsc *DataSyncCache) GetKeys() (keys []*util.NameAndNamespace) {
	keys = make([]*util.NameAndNamespace, len(dsc.cache))
	var index int = 0
	for key, _ := range dsc.cache {
		keys[index] = &util.NameAndNamespace{key.Name, key.Namespace}
		index++
	}
	return keys
}

func (dsc *DataSyncCache) GetVersionsAndPaths(namespace string, serviceName string) (versions map[string]string) {
	dsc.Lock()
	defer dsc.Unlock()
	changedNamespace := ChangeEmptyToDefault(namespace)
	if val, ok := dsc.cache[util.NameAndNamespace{serviceName, changedNamespace}]; ok {
		versions = val
	}
	return versions
}

func (dsc *DataSyncCache) UpdateCache(namespace string, serviceName string, versions map[string]string) {
	dsc.Lock()
	defer dsc.Unlock()
	if versions != nil {
		changedNamespace := ChangeEmptyToDefault(namespace)
		dsc.cache[util.NameAndNamespace{serviceName,
			changedNamespace}] = versions
	}
}

func (dsc *DataSyncCache) DeleteCache(namespace string, serviceName string) {
	dsc.Lock()
	defer dsc.Unlock()
	changedNamespace := ChangeEmptyToDefault(namespace)
	delete(dsc.cache, util.NameAndNamespace{serviceName,
		changedNamespace})
}
