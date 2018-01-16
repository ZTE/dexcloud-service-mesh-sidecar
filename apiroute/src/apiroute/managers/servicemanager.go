package managers

import (
	"apiroute/logs"
	"apiroute/models"
	apiredis "apiroute/redis"
	"encoding/json"
	"strings"

	"github.com/go-redis/redis"
)

type ServiceManager struct {
	serviceClient *redis.Client
}

//*********************query operation*****************//
func (sm *ServiceManager) QueryServiceDetailInfoStruct(key string) (
	servcie *models.ServiceDetailInfo, err error) {
	var (
		val string
	)

	val, err = sm.QueryServiceDetailInfo(key)

	if err != nil {
		logs.Log.Warn("QueryService: %s failed:%s", key, err.Error())
		return nil, err
	}

	if val == "" {
		logs.Log.Warn("QueryService: %s val is empty", key)
		return nil, nil
	}

	err = json.Unmarshal([]byte(val), &servcie)
	if err != nil {
		logs.Log.Warn("QueryService: %s convert val to struct failed:%s", key, err.Error())
		return nil, err
	}

	return servcie, err
}

func (sm *ServiceManager) QueryServiceDetailInfo(key string) (val string, err error) {
	rst, redisErr := sm.serviceClient.Get(key).Result()

	if redisErr == redis.Nil {
		logs.Log.Warn("%s does not exist", key)
		return "", err
	}

	if redisErr != nil {
		logs.Log.Warn("QueryService %s from redis failed:%s:%s", key, redisErr.Error(), rst)
		return "", redisErr
	}
	return rst, err
}

func (sm *ServiceManager) GetServiceKeyList() (serviceKeys []*models.ServiceKey, err error) {
	keys, redisErr := sm.GetAllServiceKeys()
	if redisErr != nil {
		return nil, redisErr
	}

	if keys != nil && len(keys) != 0 {
		serviceKeys = make([]*models.ServiceKey, 0)
		for _, key := range keys {
			strArr := strings.Split(key, ":")
			if len(strArr) == 5 {
				serviceKeys = append(serviceKeys,
					&models.ServiceKey{strArr[2], strArr[3], strArr[4]})
			} else {
				logs.Log.Warn("GetServiceKeyList:the key %s is not valid", key)
			}
		}
	} else {
		logs.Log.Warn("no services in service db")
	}
	return serviceKeys, err
}

func (sm *ServiceManager) GetAllServiceKeys() (keys []string, err error) {
	keyPattern := GetAllServiceKeysKeyPattern()
	keys, redisErr := apiredis.FilterKeys(keyPattern, sm.serviceClient)
	if redisErr != nil {
		logs.Log.Warn("GetAllServiceKeys from redis failed:%s:%s", keyPattern, redisErr.Error())
		return nil, redisErr
	}
	return keys, err
}

func (sm *ServiceManager) SaveService(key string, val string) (err error) {
	rst, redisErr := sm.serviceClient.Set(key, val, 0).Result()
	if redisErr != nil {
		logs.Log.Warn("SaveService %s save to redis failed:%s:%s", key, redisErr.Error(), rst)
		apiredis.SetWriteCheckFlag(true)
		return redisErr
	}
	return err
}

func (sm *ServiceManager) DeleteService(key string) (err error) {
	rst, redisErr := sm.serviceClient.Del(key).Result()
	if redisErr != nil {
		logs.Log.Warn("DeleteService %s from redis failed:%s:%d", key, redisErr.Error(), rst)
		return redisErr
	}
	return err
}
