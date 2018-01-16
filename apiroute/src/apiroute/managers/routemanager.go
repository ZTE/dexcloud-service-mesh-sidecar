package managers

import (
	"apiroute/logs"
	apiredis "apiroute/redis"

	"github.com/go-redis/redis"
)

type RouteManager struct {
	routerClient *redis.Client
}

//*********************query operation*****************//
func (rm *RouteManager) QueryRouteDetailInfo(key string) (val string, err error) {
	rst, redisErr := rm.routerClient.Get(key).Result()

	if redisErr == redis.Nil {
		logs.Log.Warn("%s does not exist", key)
		return "", err
	}

	if redisErr != nil {
		logs.Log.Warn("QueryRoute %s from redis failed:%s:%s", key, redisErr.Error(), rst)
		return "", redisErr
	}
	return rst, err
}

func (rm *RouteManager) SaveRoute(key string, val string) (err error) {
	rst, redisErr := rm.routerClient.Set(key, val, 0).Result()
	if redisErr != nil {
		logs.Log.Warn("SaveRoute %s save to redis failed:%s:%s", key, redisErr.Error(), rst)
		apiredis.SetWriteCheckFlag(true)
		return redisErr
	}
	return err
}

func (rm *RouteManager) DeleteRoute(key string) (err error) {
	rst, redisErr := rm.routerClient.Del(key).Result()
	if redisErr != nil {
		logs.Log.Warn("DeleteRoute %s from redis failed:%s:%d", key, redisErr.Error(), rst)
		return redisErr
	}
	return err
}
