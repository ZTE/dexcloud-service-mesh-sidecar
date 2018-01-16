package managers

import (
	"apiroute/logs"
	"apiroute/managers/configmanager"
	"apiroute/models"
	"apiroute/util"
	"encoding/json"
	"errors"
	"strconv"
)

type DataSyncManager struct {
	syncCache *DataSyncCache
	listeners []ServiceChangedListener
}

func NewDataSyncManager() *DataSyncManager {
	dsm := &DataSyncManager{}
	dsm.initDataSyncManager()
	return dsm
}

func (ds *DataSyncManager) initDataSyncManager() {
	ds.syncCache = NewDataSyncCache()
	ds.listeners = make([]ServiceChangedListener, 0)
}

func (ds *DataSyncManager) RegisterListener(scl ServiceChangedListener) {
	ds.listeners = append(ds.listeners, scl)
}

func (ds *DataSyncManager) UpdateOperation(serviceUnits []*util.ServiceUnit) (isSuccess bool, err error) {
	if serviceUnits == nil || len(serviceUnits) == 0 {
		logs.Log.Info("pending update serviceUnits is nil or size = 0")
		isSuccess = true
		return isSuccess, nil
	}

	logs.Log.Info("receive %s:%s update opreation", serviceUnits[0].Namespace, serviceUnits[0].Name)

	//need delete old services filter
	needDeleteServices := ds.filter(serviceUnits)
	if needDeleteServices != nil && len(needDeleteServices) != 0 {
		for _, needDeleteService := range needDeleteServices {
			logs.Log.Info("filter old service to delete:%s-%s-%s", needDeleteService.Namespace,
				needDeleteService.ServiceName, needDeleteService.ServiceVersion)
			isSuccess, err = ds.DeleteSingleOperate(needDeleteService)
			if err != nil {
				logs.Log.Info("delete old service %s failed:%s",
					needDeleteService.Namespace+"-"+needDeleteService.ServiceName, err.Error())
				return isSuccess, err
			}
		}
	}

	versions := make(map[string]string)

	//update and set service
	for _, serviceUnit := range serviceUnits {
		//valid check
		if ok := ds.checkServiceUnitValid(serviceUnit); !ok {
			continue
		}

		//special handle
		//		ds.specialHandleRouterMetrics(serviceUnit)

		//convert to serviceInfo
		var (
			serviceInfo *models.ServiceInfo
			internalErr error
		)

		if serviceInfo, internalErr = AssembleServiceInfo(serviceUnit); internalErr != nil {
			logs.Log.Warn("covert serviceUnit %s to serviceInfo failed:%s",
				serviceUnit.Name, internalErr.Error())
			continue
		}

		//persistence serviceInfo
		sm := GetServiceManager()
		key := CovertServiceKey(serviceInfo.ServicePrefix, serviceInfo.ServiceKey)
		logs.Log.Info("update service:%s", key)
		val, jsonErr := json.Marshal(serviceInfo.ServiceValue)
		if jsonErr != nil {
			logs.Log.Warn("%s convert to json failed:%s", key, jsonErr.Error())
			continue
		}

		if redisErr := sm.SaveService(key, string(val)); redisErr != nil {
			return false, redisErr
		}

		//service changed
		for _, listener := range ds.listeners {
			isSuccess, err = listener.OnSave(serviceInfo)
			if err != nil {
				logs.Log.Warn("service  %s changed handle failed:%s",
					serviceInfo.ServiceKey.Namespace+"-"+serviceInfo.ServiceKey.ServiceName, err.Error())
				return isSuccess, err
			}
		}

		//
		versions[serviceUnit.Version] = serviceUnit.Path
	}

	//update synccache
	ds.syncCache.UpdateCache(serviceUnits[0].Namespace, serviceUnits[0].Name, versions)
	return isSuccess, err
}

func (ds *DataSyncManager) checkServiceUnitValid(serviceUnit *util.ServiceUnit) bool {
	if serviceUnit == nil {
		logs.Log.Info("serviceUnit is nil")
		return false
	}

	if len(serviceUnit.Name) == 0 {
		logs.Log.Info("serviceUnit name is empty")
		return false
	}

	if serviceUnit.Instances == nil || len(serviceUnit.Instances) == 0 {
		logs.Log.Info("serviceUnit:%s instances is empty", serviceUnit.Name)
		return false
	}

	return true
}

func (ds *DataSyncManager) specialHandleRouterMetrics(serviceUnit *util.ServiceUnit) {
	//already checkvaild

	//filter router_metrics service
	if serviceUnit.Name != "router_metrics" {
		return
	}

	var (
		defaultPort = configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport
	)

	//find default port
	newInstances := make([]util.InstanceUnit, 1)
	for _, instance := range serviceUnit.Instances {
		if instance.ServicePort == strconv.Itoa(defaultPort) {
			newInstances[0] = instance
			serviceUnit.Instances = newInstances
			return
		}
	}
}

func (ds *DataSyncManager) filter(serviceUnits []*util.ServiceUnit) (needDeleteServices []*models.ServiceKey) {

	needDeleteServices = make([]*models.ServiceKey, 0)

	if serviceUnits == nil || len(serviceUnits) == 0 {
		return needDeleteServices
	}

	versions := ds.syncCache.GetVersionsAndPaths(serviceUnits[0].Namespace, serviceUnits[0].Name)

	if versions != nil && len(versions) != 0 {
		//find in versions,not in serviceUnit
		for version, path := range versions {
			found := false
			for _, serviceUnit := range serviceUnits {
				if version == serviceUnit.Version {
					//find
					found = true
					//compare path.
					if path != serviceUnit.Path {
						//path different,need delete
						needDeleteService := &models.ServiceKey{serviceUnit.Namespace,
							serviceUnit.Name, serviceUnit.Version}
						needDeleteServices = append(needDeleteServices, needDeleteService)
					}
					break
				}
			}

			if !found {
				needDeleteService := &models.ServiceKey{serviceUnits[0].Namespace,
					serviceUnits[0].Name, version}
				needDeleteServices = append(needDeleteServices, needDeleteService)
			}
		}
	}

	return needDeleteServices
}

func (ds *DataSyncManager) DeleteOperation(nameAndNamespace *util.NameAndNamespace) (isSuccess bool, err error) {
	if nameAndNamespace == nil {
		logs.Log.Warn("delete: nameAndNamespace is nil")
		return true, nil
	}

	logs.Log.Info("receive %s:%s delete opreation", nameAndNamespace.Namespace, nameAndNamespace.Name)

	//get the complete service key
	versions := ds.syncCache.GetVersionsAndPaths(nameAndNamespace.Namespace, nameAndNamespace.Name)

	if versions == nil || len(versions) == 0 {
		logs.Log.Info("%s-%s did not have versions.", nameAndNamespace.Namespace, nameAndNamespace.Name)
		return true, nil
	}

	//loop delete
	for version, _ := range versions {
		isSuccess, err = ds.DeleteSingleOperate(&models.ServiceKey{nameAndNamespace.Namespace,
			nameAndNamespace.Name, version})
		if err != nil {
			return isSuccess, err
		}
	}

	//update synccache
	ds.syncCache.DeleteCache(nameAndNamespace.Namespace, nameAndNamespace.Name)
	return isSuccess, err
}

func (ds *DataSyncManager) DeleteSingleOperate(serviceKey *models.ServiceKey) (isSuccess bool, err error) {
	var (
		redisErr          error
		serviceKeyStr     string
		serviceDetailInfo *models.ServiceDetailInfo
		serviceInfo       *models.ServiceInfo
	)
	//check
	if ok := ds.checkServiceKeyValid(serviceKey); !ok {
		logs.Log.Warn("delete opration.check ServiceKey failed.")
		err = errors.New("delete opration.check ServiceKey failed")
		return false, err
	}

	//change namespace
	newNamespace := ChangeEmptyToDefault(serviceKey.Namespace)
	serviceKey.Namespace = newNamespace

	//query serviceInfo
	sm := GetServiceManager()
	serviceKeyStr = CovertServiceKey(ServicePrefix, *serviceKey)
	serviceDetailInfo, redisErr = sm.QueryServiceDetailInfoStruct(serviceKeyStr)

	if redisErr != nil {
		logs.Log.Warn("prepare delete service %s:get service info from redis failed:%s",
			serviceKeyStr, redisErr.Error())
		return false, redisErr
	}

	if serviceDetailInfo == nil {
		logs.Log.Warn("prepare delete service %s:get service info from redis:the val is empty",
			serviceKeyStr)
		return true, nil
	}

	//delete service
	logs.Log.Info("delete service:%s", serviceKeyStr)
	redisErr = sm.DeleteService(serviceKeyStr)
	if redisErr != nil {
		logs.Log.Warn("delete service %s failed:%s", serviceKeyStr, redisErr.Error())
		return false, redisErr
	}

	//service changed
	serviceInfo = &models.ServiceInfo{ServicePrefix, *serviceKey, *serviceDetailInfo}
	for _, listener := range ds.listeners {
		isSuccess, internalErr := listener.OnDelete(serviceInfo)
		if internalErr != nil {
			return isSuccess, internalErr
		}
	}

	return true, err
}

func (ds *DataSyncManager) checkServiceKeyValid(serviceKey *models.ServiceKey) bool {
	if serviceKey == nil {
		logs.Log.Info("the serviceKey is nil.")
		return false
	}

	if serviceKey.ServiceName == "" {
		logs.Log.Info("the serviceName is empty.")
		return false
	}

	return true
}

func (ds *DataSyncManager) GetNameList() (nameList []*util.NameAndNamespace, err error) {
	nameList = ds.syncCache.GetKeys()
	return nameList, err
}

func (ds *DataSyncManager) GetVersionsAndPaths(namespace string, serviceName string) (versions map[string]string) {
	return ds.syncCache.GetVersionsAndPaths(namespace, serviceName)
}

func (ds *DataSyncManager) ServiceExist(name *util.NameAndNamespace) (bool, error) {
	isExist := ds.syncCache.IsExist(name.Namespace, name.Name)
	return isExist, nil
}

//get all ,please set serviceVersion "all"
func (ds *DataSyncManager) GetPublishInfo(namespace, serviceName,
	serviceVersion string) (publisInfos map[models.ServiceKey]*models.PublishInfo, err error) {

	var (
		newNamespace string
		sm           *ServiceManager
		serviceKey   models.ServiceKey
		key          string
		servcie      *models.ServiceDetailInfo
		publisInfo   *models.PublishInfo
	)
	sm = GetServiceManager()
	newNamespace = ChangeEmptyToDefault(namespace)

	if serviceVersion == "all" {
		versions := ds.syncCache.GetVersionsAndPaths(newNamespace, serviceName)
		if versions == nil || len(versions) == 0 {
			logs.Log.Info("%s:%s is not in redis.", namespace, serviceName)
			//			err = errors.New(namespace + ":" + serviceName + " didn't have versions")
			return publisInfos, err
		}

		publisInfos = make(map[models.ServiceKey]*models.PublishInfo, len(versions))
		for version, _ := range versions {
			serviceKey = models.ServiceKey{newNamespace, serviceName, version}
			key = CovertServiceKey(ServicePrefix, serviceKey)
			servcie, err = sm.QueryServiceDetailInfoStruct(key)
			if err != nil {
				return publisInfos, err
			}

			if servcie == nil || servcie.Spec.PublishPort == "" {
				continue
			}

			publisInfo = &models.PublishInfo{servcie.Spec.Protocol,
				servcie.Spec.PublishPort, len(servcie.Spec.Nodes)}
			publisInfos[serviceKey] = publisInfo
		}

	} else {
		publisInfos = make(map[models.ServiceKey]*models.PublishInfo, 1)
		serviceKey = models.ServiceKey{newNamespace, serviceName, serviceVersion}
		key = CovertServiceKey(ServicePrefix, serviceKey)
		servcie, err = sm.QueryServiceDetailInfoStruct(key)
		if err != nil {
			return publisInfos, err
		}

		if servcie == nil || servcie.Spec.PublishPort == "" {
			//			err = errors.New(namespace + ":" + serviceName + " didn't have valid val")
			return publisInfos, err
		}

		publisInfo = &models.PublishInfo{servcie.Spec.Protocol,
			servcie.Spec.PublishPort, len(servcie.Spec.Nodes)}
		publisInfos[serviceKey] = publisInfo
	}

	return publisInfos, err
}

func (ds *DataSyncManager) GetAllServiceDigest() (allServiceDigest map[string]*util.DigestUnit, err error) {
	allServiceDigest = make(map[string]*util.DigestUnit)
	sm := GetServiceManager()
	keys, err := sm.GetServiceKeyList()
	if err != nil {
		return
	}

	var service *models.ServiceDetailInfo
	for _, key := range keys {
		service, err = sm.QueryServiceDetailInfoStruct(CovertServiceKey(ServicePrefix, *key))
		if err != nil {
			return
		}

		if service == nil {
			continue
		}

		var (
			nameDashNamespace string
			namespace         string
		)

		if key.Namespace == "" || key.Namespace == "default" {
			nameDashNamespace = key.ServiceName
			namespace = ""
		} else {
			nameDashNamespace = key.ServiceName + "-" + key.Namespace
			namespace = key.Namespace
		}

		if val, exist := allServiceDigest[nameDashNamespace]; exist {
			val.NumOfNode += len(service.Spec.Nodes)
		} else {
			allServiceDigest[nameDashNamespace] = &util.DigestUnit{
				Namespace: namespace,
				NumOfNode: len(service.Spec.Nodes),
			}
		}
	}
	return
}

func (ds *DataSyncManager) Delete(name *util.NameAndNamespace) error {
	_, err := ds.DeleteOperation(name)
	return err
}

func (ds *DataSyncManager) Update(service []*util.ServiceUnit) error {
	_, err := ds.UpdateOperation(service)
	return err
}
