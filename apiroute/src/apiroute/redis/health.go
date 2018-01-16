package redis

import (
	"apiroute/logs"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

const (
	checkInterval       = 20
	retryIntervalOnFail = 5
	retryCountOnFail    = 12
)

type Result struct {
	message string
	err     error
}

var (
	writeCheckFlag = true
	result         = Result{}
	quit           = make(chan struct{})
	mutex1         = &sync.Mutex{}
	mutex2         = &sync.Mutex{}
)

func (c *Clients) StartHealthCheck() {
	checkTicker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	go func() {
		for {
			select {
			case <-checkTicker.C:
				logs.Log.Debug("Start Health Check at %s", time.Now().Format("2006-01-02T15:04:05.000Z0700"))
				healthCheck(c.RouterClient)
				if ok, err := IsHealthy(); ok {
					logs.Log.Debug("End Health Check with result:OK")
				} else {
					logs.Log.Debug("End Health Check with result:%v", err)
				}
			case <-quit:
				checkTicker.Stop()
				return
			}
		}
	}()
}

func StopHealthCheck() {
	close(quit)
}

func SetWriteCheckFlag(value bool) {
	mutex1.Lock()
	defer mutex1.Unlock()
	writeCheckFlag = value
}

func setResult(msg string, e error) {
	mutex2.Lock()
	defer mutex2.Unlock()
	result = Result{
		message: msg,
		err:     e,
	}
}

func IsHealthy() (bool, error) {
	mutex2.Lock()
	defer mutex2.Unlock()
	if result.err == nil {
		return true, nil
	}
	return false, result.err
}

func healthCheck(client *redis.Client) {
	var (
		failCount   int
		flagChecked bool
	)

	for {
		flagChecked = false
		mutex1.Lock()
		if writeCheckFlag {
			mutex1.Unlock()
			flagChecked = true
			//Write check
			msg, err := client.Set("HealthCheck:checktime", time.Now().Format("2006-01-02T15:04:05.000Z0700"), 0).Result()
			if err == nil {
				SetWriteCheckFlag(false)
			} else {
				failCount++
				if failCount == retryCountOnFail {
					setResult(msg, err)
					return
				}
				time.Sleep(time.Duration(retryIntervalOnFail) * time.Second)
				continue
			}
		}

		if !flagChecked {
			mutex1.Unlock()
		}

		//Read check
		msg, err := client.Get("HealthCheck:checktime").Result()
		if err == nil {
			setResult(msg, err)
			return
		}

		failCount++
		if failCount == retryCountOnFail {
			setResult(msg, err)
			return
		}
		time.Sleep(time.Duration(retryIntervalOnFail) * time.Second)
	}
}
