package managers

import (
	"apiroute/logs"
	"apiroute/util"
	"encoding/json"
	"strconv"
	"time"
)

const (
	failRetryCount = 12
)

type ServiceDigestData struct {
	Name           string `json:"serviceName"`
	Namespace      string `json:"namespace"`
	MaxModifyIndex uint64 `json:"maxModifyIndex"`
	NumOfInstance  int    `json:"numOfInstance"`
}

type Watcher struct {
	ErrCh           chan error
	stopCh          chan struct{}
	upstreamTimeOut time.Duration
}

func NewWatcher() *Watcher {
	return &Watcher{
		ErrCh:           make(chan error),
		stopCh:          make(chan struct{}),
		upstreamTimeOut: 10,
	}
}

func (w *Watcher) StartWatch(outCh chan<- []*ServiceDigestData) {
	var (
		failCount    int
		queryIndex   int64
		firstSuccess bool
	)
	filterMeta := getMergedLabels()
	hashcode := util.GetMD5Hash(registerName + filterMeta)
	url := "http://" + clientIP + ":" + strconv.Itoa(int(clientPort)) + "/api/microservices/v1/tenants/" + hashcode + "/services/digest"
	queryStr := "?initial=true&filter-meta=" + filterMeta

	for {
		select {
		case <-w.stopCh:
			logs.Log.Info("Recieve message on the stopCh, return from watch")
			return
		default:
		}

		logs.Log.Info("Start to watch service digest endpoint:%s%s at index:%d", url, queryStr, queryIndex)
		buf, tag, err := util.HTTPGetWithIndex(url, queryStr, strconv.FormatInt(queryIndex, 10))
		logs.Log.Info("Watch service digest return")

		if err != nil {
			if failCount < failRetryCount {
				failCount++
				logs.Log.Warn("Watch service digest returned with error:%v, Sleep 10 secs and retry the %dth time", err, failCount)
				time.Sleep(10 * time.Second)
				continue
			} else {
				if firstSuccess {
					logs.Log.Warn("SDClient is not available, try connect again...")
					failCount = 0
					continue
				}
				logs.Log.Warn("Watch service digest returned with error:%v, failed %d times at total, bail out program", err, failRetryCount)
				w.ErrCh <- err
				return
			}
		}

		failCount = 0
		if !firstSuccess {
			firstSuccess = true
		}

		data := []*ServiceDigestData{}
		err = json.Unmarshal(buf, &data)

		if err != nil {
			logs.Log.Warn("Unmarshal data returned with error:%v", err)
			queryStr = "?filter-meta=" + filterMeta
			continue
		}

		if len(data) == 0 {
			logs.Log.Info("Service Digest list is empty, proceed with full deletion")
		}

		indexReceived, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			logs.Log.Warn("Failed to parse string:%s to int64:%v", tag, err)
		} else {
			queryIndex = indexReceived + 1
		}

		logs.Log.Info("Send Service Digest list to upstream")
		timeAfter := time.NewTimer(time.Duration(w.upstreamTimeOut) * time.Second)

		select {
		case outCh <- data:
			logs.Log.Info("Service Digest list was sent to upstream")
			queryStr = "?filter-meta=" + filterMeta
		case <-timeAfter.C:
			logs.Log.Warn("Upstream has not fetched the Service Digest list within %d seconds, abandon this change notification", w.upstreamTimeOut)
			queryStr = "?initial=true&filter-meta=" + filterMeta
		}
		timeAfter.Stop()
	}
}

func (w *Watcher) StopWatch() {
	close(w.stopCh)
}
