package managers

import (
	"apiroute/logs"
	"sync"
	"time"
)

const (
	checkInterval = 20
)

type NamespaceManager struct {
	flagMutex     sync.Mutex
	cacheMutex    sync.Mutex
	flag          int64
	namespaceList []string
	quit          chan struct{}
}

func NewNamespaceManager() *NamespaceManager {
	logs.Log.Info("init namespace manager")
	nm := &NamespaceManager{}
	nm.initNamespaceManager()
	return nm
}

func (nm *NamespaceManager) initNamespaceManager() {
	nm.namespaceList = make([]string, 0)
	nm.syncDataToCache()
	nm.quit = make(chan struct{})
	go nm.refreshNamespaceCache()
}

func (nm *NamespaceManager) GetNamespaceList() []string {
	nm.cacheMutex.Lock()
	defer nm.cacheMutex.Unlock()
	return nm.namespaceList
}

func (nm *NamespaceManager) ResetNamespaceList(newList []string) {
	nm.cacheMutex.Lock()
	defer nm.cacheMutex.Unlock()
	if newList != nil {
		nm.namespaceList = newList
	} else {
		nm.namespaceList = make([]string, 0)
	}
}

func (nm *NamespaceManager) NotifyNamespaceListUpdate() {
	nm.flagMutex.Lock()
	defer nm.flagMutex.Unlock()
	nm.flag++
}

func (nm *NamespaceManager) GetFlag() int64 {
	nm.flagMutex.Lock()
	defer nm.flagMutex.Unlock()
	return nm.flag
}

func (nm *NamespaceManager) ResetFlag(oldFlag int64) {
	nm.flagMutex.Lock()
	defer nm.flagMutex.Unlock()
	if oldFlag == nm.flag {
		nm.flag = 0
	}
}

func (nm *NamespaceManager) refreshNamespaceCache() {
	logs.Log.Info("namsespace refresh routine start work")
	checkTicker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	for {
		select {
		case <-checkTicker.C:
			currentflag := nm.GetFlag()
			if currentflag != 0 {
				//refresh namespacelist
				nm.syncDataToCache()
				//reset flag
				nm.ResetFlag(currentflag)
			}
		case <-nm.quit:
			return
		}
	}
}

func (nm *NamespaceManager) syncDataToCache() {
	sm := GetServiceManager()
	pServiceKeys, redisErr := sm.GetServiceKeyList()

	if redisErr != nil {
		logs.Log.Warn("sync namespace failed:%s", redisErr.Error())
		return
	}

	namespaces := make([]string, 0, len(pServiceKeys))

	if pServiceKeys != nil && len(pServiceKeys) != 0 {
		for _, pServiceKey := range pServiceKeys {
			find := false
			for _, namespace := range namespaces {
				if pServiceKey.Namespace == namespace {
					find = true
					break
				}
			}
			if !find {
				namespaces = append(namespaces, pServiceKey.Namespace)
			}
		}
	}

	//
	nm.ResetNamespaceList(namespaces)
}

func (nm *NamespaceManager) CloseRefreshRoutine() {
	logs.Log.Info("close namsespace refresh routine")
	close(nm.quit)
}
