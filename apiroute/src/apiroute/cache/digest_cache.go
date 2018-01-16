package cache

import (
	"apiroute/util"
	"strings"
	"sync"
)

type ServiceDigest struct {
	Namespace      string
	MaxModifyIndex uint64
	NumOfInstance  int
}

type ServiceDigestCache struct {
	sync.Mutex
	digestMap map[string]*ServiceDigest
}

func NewServiceDigestCache() *ServiceDigestCache {
	return &ServiceDigestCache{
		digestMap: make(map[string]*ServiceDigest),
	}
}

func (dc *ServiceDigestCache) Size() int {
	dc.Lock()
	defer dc.Unlock()
	return len(dc.digestMap)
}

func (dc *ServiceDigestCache) Replace(data map[string]*ServiceDigest) {
	dc.Lock()
	defer dc.Unlock()
	dc.digestMap = data
}

func (dc *ServiceDigestCache) Merge(data map[string]*ServiceDigest) {
	dc.Lock()
	defer dc.Unlock()
	for k, v := range data {
		dc.digestMap[k] = v
	}
}

func (dc *ServiceDigestCache) Diff(data map[string]*ServiceDigest) ([]*util.NameAndNamespace, []*util.NameAndNamespace) {
	dc.Lock()
	defer dc.Unlock()
	var delList, updateList []*util.NameAndNamespace

	for name, sd := range dc.digestMap {
		if digest, ok := data[name]; ok {
			if digest.MaxModifyIndex > dc.digestMap[name].MaxModifyIndex ||
				digest.NumOfInstance < dc.digestMap[name].NumOfInstance {
				updateList = append(updateList, &util.NameAndNamespace{
					Name:      trimNamespace(name, digest),
					Namespace: digest.Namespace,
				})
			}
			continue
		}

		delList = append(delList, &util.NameAndNamespace{
			Name:      trimNamespace(name, sd),
			Namespace: sd.Namespace,
		})
	}

	for name, digest := range data {
		if _, ok := dc.digestMap[name]; ok {
			continue
		}
		updateList = append(updateList, &util.NameAndNamespace{
			Name:      trimNamespace(name, digest),
			Namespace: digest.Namespace,
		})
	}

	return delList, updateList
}

func (dc *ServiceDigestCache) Diff2(redis map[string]*util.DigestUnit) ([]*util.NameAndNamespace, []*util.NameAndNamespace) {
	dc.Lock()
	defer dc.Unlock()
	var delList, updateList []*util.NameAndNamespace

	for name, sd := range dc.digestMap {
		if _, ok := redis[name]; ok {
			continue
		}
		//In cache not in Redis
		updateList = append(updateList, &util.NameAndNamespace{
			Name:      trimNamespace(name, sd),
			Namespace: sd.Namespace,
		})
	}

	for name, digest := range redis {
		if sd, ok := dc.digestMap[name]; ok {
			if sd.NumOfInstance != digest.NumOfNode {
				updateList = append(updateList, &util.NameAndNamespace{
					Name:      trimNamespace2(name, digest),
					Namespace: digest.Namespace,
				})
			}
			continue
		}
		//In Redis not in Cache
		delList = append(delList, &util.NameAndNamespace{
			Name:      trimNamespace2(name, digest),
			Namespace: digest.Namespace,
		})
	}

	return delList, updateList
}

func trimNamespace(name string, digest *ServiceDigest) string {
	if digest.Namespace == "" {
		return name
	}

	return strings.TrimSuffix(name, "-"+digest.Namespace)
}

func trimNamespace2(name string, digest *util.DigestUnit) string {
	if digest.Namespace == "" {
		return name
	}

	return strings.TrimSuffix(name, "-"+digest.Namespace)
}
