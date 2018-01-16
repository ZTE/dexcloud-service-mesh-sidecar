package managers

import (
	"apiroute/cache"
	"apiroute/logs"
	"apiroute/template"
	"apiroute/util"
	"sync"
	"time"
)

const timeToLive = 2

type Dispatcher struct {
	errCh       chan error
	stopCh      chan struct{}
	DigestCache *cache.ServiceDigestCache
	reloadCache *cache.ReloadCache
}

func NewDispatcher(cache *cache.ServiceDigestCache, reload *cache.ReloadCache) *Dispatcher {
	return &Dispatcher{
		errCh:       make(chan error),
		stopCh:      make(chan struct{}),
		DigestCache: cache,
		reloadCache: reload,
	}
}

func (d *Dispatcher) StartDispatch(inCh <-chan []*ServiceDigestData) {
	for {
		select {
		case <-d.stopCh:
			logs.Log.Info("Dispatcher receive message on stopCh, return right away")
			return
		case data := <-inCh:
			logs.Log.Info("Dispatcher receive service digest list, start to process it")
			logs.Log.Info(">>>L7 Port Reference details before dispatch>>>")
			d.reloadCache.PrintL7PortRef()
			//Save data for restore
			var oldCacheData *cache.CacheData
			digestMap := make(map[string]*cache.ServiceDigest, len(data))

			for _, digest := range data {
				var name string
				var ns string
				if digest.Namespace == "default" || digest.Namespace == "" {
					name = digest.Name
					ns = ""
				} else {
					name = digest.Name + "-" + digest.Namespace
					ns = digest.Namespace
				}

				digestMap[name] = &cache.ServiceDigest{
					Namespace:      ns,
					MaxModifyIndex: digest.MaxModifyIndex,
					NumOfInstance:  digest.NumOfInstance,
				}
			}

			var errors error
			digestCacheToRedisSyncLocker.Lock()

			if d.DigestCache.Size() == 0 {
				inRedis, err := GetDataSyncManager().GetNameList()
				if err != nil {
					logs.Log.Warn("GetNameList() returned with error:%v, drop this dispatch", err)
					digestCacheToRedisSyncLocker.Unlock()
					continue
				}
				oldCacheData = d.reloadCache.ResetReloadCache()
				updateList := covertToUpdateUnit(data)

				var emptyDelList, emptyUpdateList []*util.NameAndNamespace
				if len(inRedis) > 0 {
					errors = d.startWorkers(inRedis, emptyUpdateList, false)
				}
				if errors != nil {
					logs.Log.Warn("Delete All from redis returned with error:%v, drop this dispatch", errors)
					digestCacheToRedisSyncLocker.Unlock()
					continue
				}
				errors = d.startWorkers(emptyDelList, updateList, false)
			} else {
				deleteList, updateList := d.DigestCache.Diff(digestMap)
				if len(deleteList) == 0 && len(updateList) == 0 {
					logs.Log.Info("Nothing changed, proceed with next watch")
					digestCacheToRedisSyncLocker.Unlock()
					continue
				}
				oldCacheData = d.reloadCache.GetData()
				errors = d.startWorkers(deleteList, updateList, false)
			}

			if errors != nil {
				d.reloadCache.RestoreAll(oldCacheData)
				logs.Log.Warn("Worker returned with error:%v, restore reload cache and drop this dispatch", errors)
				digestCacheToRedisSyncLocker.Unlock()
				continue
			}

			logs.Log.Info(">>>L7 Port Reference details after dispatch>>>")
			d.reloadCache.PrintL7PortRef()
			// Redis update completes so far, safely override the cache
			d.OverrideDigestCache(digestMap)
			digestCacheToRedisSyncLocker.Unlock()

			// Start to render nginx conf files if needed
			http, https, updateHTTP, updateHTTPS := d.reloadCache.CheckPortUpdate()
			if updateHTTP {
				if err := template.RenderHTTP(http); err != nil {
					d.reloadCache.RestoreAll(oldCacheData)
					logs.Log.Warn("Render http returned with error:%v, restore reload cache", err)
					continue
				}
			}

			if updateHTTPS {
				if err := template.RenderHTTPS(https); err != nil {
					d.reloadCache.RestoreAll(oldCacheData)
					logs.Log.Warn("Render https returned with error:%v, restore reload cache", err)
					continue
				}
			}

			streams, updateStream := d.reloadCache.CheckStreamUpdate()
			if updateStream {
				if err := template.RenderStream(streams); err != nil {
					d.reloadCache.RestoreAll(oldCacheData)
					logs.Log.Warn("Render stream returned with error:%v, restore reload cache", err)
					continue
				}
			}

			// Start to reload nginx if needed
			if updateHTTP || updateHTTPS || updateStream {
				logs.Log.Info("Start to reload nginx")
				if err := cache.Reloading(); err != nil {
					d.reloadCache.RestoreAll(oldCacheData)
					logs.Log.Warn("Reloading returned with error:%v, restore reload cache and proceed with next watch", err)
				}
			}

			// Reset reload flags
			d.reloadCache.ResetFlags(false)
		}
	}
}

func (d *Dispatcher) StopDispatch() {
	close(d.stopCh)
}

func (d *Dispatcher) OverrideDigestCache(data map[string]*cache.ServiceDigest) {

	if d.DigestCache.Size() == 0 {
		d.DigestCache.Merge(data)
	} else {
		d.DigestCache.Replace(data)
	}
}

func (d *Dispatcher) startWorkers(delList, updateList []*util.NameAndNamespace, mode bool) error {
	startTime := time.Now().UnixNano()

	completeCh := make(chan struct{})
	numOfDeleter := len(delList) / 4
	numOfUpdater := len(updateList) / 4

	if numOfDeleter == 0 && len(delList) != 0 {
		numOfDeleter = 1
	}

	if numOfUpdater == 0 && len(updateList) != 0 {
		numOfUpdater = 1
	}

	if numOfDeleter > 100 {
		numOfDeleter = 100
	}

	if numOfUpdater > 100 {
		numOfUpdater = 100
	}

	workerPool := NewPool(numOfDeleter + numOfUpdater)

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(delList) + len(updateList))

		for _, n := range delList {
			deleter := &Deleter{
				errCh:    d.errCh,
				syncMode: mode,
				nameAndNamespace: &util.NameAndNamespace{
					Name:      n.Name,
					Namespace: n.Namespace,
				},
				reloadCache: d.reloadCache,
			}

			go func() {
				workerPool.Run(deleter)
				wg.Done()
			}()
		}

		for _, m := range updateList {

			updater := &Updater{
				errCh:       d.errCh,
				syncMode:    mode,
				name:        m.Name,
				namespace:   m.Namespace,
				reloadCache: d.reloadCache,
			}

			go func() {
				workerPool.Run(updater)
				wg.Done()
			}()
		}

		wg.Wait()
		workerPool.Shutdown()
		endTime := time.Now().UnixNano()
		logs.Log.Info("sync data start-end:%d-%d,duration(ms):%d", startTime, endTime, (endTime-startTime)/1000000)
		completeCh <- struct{}{}
	}()

	select {
	case err := <-d.errCh:
		logs.Log.Error("Dispatcher received error:%v from worker, quit this dispatch right away", err)
		return err
	case <-completeCh:
		logs.Log.Info("Workers completed all the tasks")
		return nil
	}
}

func covertToUpdateUnit(digestList []*ServiceDigestData) (updateList []*util.NameAndNamespace) {
	for _, digest := range digestList {
		ns := digest.Namespace
		if ns == "default" {
			ns = ""
		}
		updateList = append(updateList, &util.NameAndNamespace{
			Name:      digest.Name,
			Namespace: ns,
		})
	}
	return
}

func sendError(ch chan<- error, err error) {
	timeAfter := time.NewTimer(time.Duration(timeToLive) * time.Second)
	defer timeAfter.Stop()
	select {
	case ch <- err:
	case <-timeAfter.C:
	}
}
