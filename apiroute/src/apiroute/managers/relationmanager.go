package managers

import (
	"apiroute/logs"
	"apiroute/models"
	apiredis "apiroute/redis"

	"github.com/go-redis/redis"
)

type ReleationManager struct {
	releationClient *redis.Client
}

func (rm *ReleationManager) QueryRelationsByServiceKey(serviceKey *models.ServiceKey) (
	relations []*models.ServiceAndRouteRelationMap, err error) {
	if serviceKey == nil {
		return nil, err
	}

	keyPattern := AssembleGetRelationsByServiceKeyPattern(serviceKey)
	keys, redisErr := rm.QueryRelationsByKeyPattern(keyPattern)
	if redisErr != nil {
		return nil, redisErr
	}

	relations = ConvertKeysToRelations(keys)

	return relations, err
}

func (rm *ReleationManager) QueryRelationsByKeyPattern(keyPattern string) (relations []string, err error) {
	keys, redisErr := apiredis.FilterKeys(keyPattern, rm.releationClient)
	if redisErr != nil {
		logs.Log.Warn("QueryRelationsByKeyPattern from redis failed:%s:%s", keyPattern, redisErr.Error())
		return nil, redisErr
	}

	return keys, err
}

func (rm *ReleationManager) SaveRelation(key string, val string) (err error) {
	rst, redisErr := rm.releationClient.Set(key, val, 0).Result()
	if redisErr != nil {
		logs.Log.Warn("SaveReleation %s save to redis failed:%s:%s", key, redisErr.Error(), rst)
		apiredis.SetWriteCheckFlag(true)
		return redisErr
	}

	return err
}

func (rm *ReleationManager) DeleteRelation(key string) (err error) {
	rst, redisErr := rm.releationClient.Del(key).Result()
	if redisErr != nil {
		logs.Log.Warn("DeleteRelation %s failed:%s:%d", key, redisErr.Error(), rst)
		return redisErr
	}

	return err
}
