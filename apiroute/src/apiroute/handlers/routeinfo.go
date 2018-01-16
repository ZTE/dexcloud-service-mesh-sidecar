package handlers

import (
	"apiroute/logs"
	"apiroute/managers"
	"apiroute/managers/configmanager"
	"apiroute/models"
	"apiroute/util"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
)

func CheckServiceKeyValid(serviceKey *models.ServiceKey) error {
	if serviceKey == nil {
		logs.Log.Warn("check servicekey valid:serviceKey is nil")
		err := errors.New("check servicekey valid:serviceKey is nil")
		return err
	}

	if serviceKey.ServiceName == "" {
		logs.Log.Warn("check servicekey valid:the serviceName is empty")
		err := errors.New("check servicekey valid:the serviceName is empty")
		return err
	}

	return nil
}

func CheckRouteKeyValid(routeKey *models.RouteKey) error {
	if routeKey == nil {
		logs.Log.Warn("check routeKey valid:routeKey is nil")
		err := errors.New("check routeKey valid:routeKey is nil")
		return err
	}

	if routeKey.RouteName == "" {
		logs.Log.Warn("check routeKey valid:routeName is empty")
		err := errors.New("check routeKey valid:routeName is empty")
		return err
	}

	if routeKey.PublishPort == "" {
		logs.Log.Warn("check routeKey %s valid:PublishPort is empty", routeKey.RouteName)
		err := errors.New("check routeKey valid:PublishPort is empty")
		return err
	}

	if routeKey.RouteType == "" {
		logs.Log.Warn("check routeKey %s valid:RouteType is empty", routeKey.RouteName)
		err := errors.New("check routeKey valid:RouteType is empty")
		return err
	}

	return nil
}

func GetRouteDetailInfoByRouteKeyInJSON(publishPort, routeWay, routeType, routeName,
	routeVersion string) (routeDetailInfo string, err error) {

	//key
	var (
		key      string
		routeKey models.RouteKey
	)

	routeKey = models.RouteKey{publishPort, routeType, routeName, routeVersion}
	err = CheckRouteKeyValid(&routeKey)
	if err != nil {
		logs.Log.Warn("routekey is not valid:%s", err.Error())
		return "", err
	}

	if routeWay == managers.RouteWayDomain {
		key = managers.ConvertRouteKey(managers.RouteDomainPrefix, routeKey)
	} else if routeWay == managers.RouteWayIP {
		key = managers.ConvertRouteKey(managers.RouteIPPrefix, routeKey)
	}

	//query val
	rm := managers.GetRouteManager()
	routeDetailInfo, err = rm.QueryRouteDetailInfo(key)
	if err != nil {
		logs.Log.Warn("get %s routeinfo falid:%s", key, err.Error())
		return "", err
	}

	return routeDetailInfo, err

}

func GetRouteDetailInfoByRouteKey(publishPort, routeWay, routeType, routeName,
	routeVersion string) (routeDetailInfo *models.RouteDetailInfo, err error) {

	if configmanager.CfgM.GetConfigInfo().EnableTest {
		return MockGetRouteDetailInfoByRouteKey()
	}

	var (
		routeKeyStr string
		routeVal    string
	)

	err = CheckRouteKeyValid(&models.RouteKey{publishPort, routeType, routeName, routeVersion})
	if err != nil {
		logs.Log.Warn("routekey is not valid:%s", err.Error())
		return nil, err
	}

	routeKeyStr = fmt.Sprintf("%s:%s:%s:%s", publishPort,
		routeType, routeName, routeVersion)

	routeVal, err = GetRouteDetailInfoByRouteKeyInJSON(publishPort, routeWay,
		routeType, routeName, routeVersion)
	if err != nil {
		logs.Log.Warn("get %s route failed:%s", routeKeyStr, err.Error())
		return nil, err
	}

	if routeVal == "" {
		logs.Log.Warn("%s route don't exist", routeKeyStr)
		return nil, nil
	}

	err = json.Unmarshal([]byte(routeVal), &routeDetailInfo)
	if err != nil {
		logs.Log.Warn("convert route %s json to struct failed:%s", routeKeyStr, err.Error())
		return nil, err
	}

	return routeDetailInfo, err

}

func GetTCPUDPRouteDetailInfoByServiceKey(namespace, serviceName,
	serviceVersion string) (tcpudpRouteDetailInfo *models.TCPUDPRouteDetailInfo, err error) {

	var (
		serviceKey        models.ServiceKey
		key               string
		serviceVal        string
		serviceDetailInfo *models.ServiceDetailInfo
	)

	newNamespace := managers.ChangeEmptyToDefault(namespace)
	serviceKey = models.ServiceKey{newNamespace, serviceName, serviceVersion}
	err = CheckServiceKeyValid(&serviceKey)
	if err != nil {
		logs.Log.Warn("servicekey is not valid:%s", err.Error())
		return nil, err
	}

	//query service
	key = managers.CovertServiceKey(managers.ServicePrefix, serviceKey)
	sm := managers.GetServiceManager()
	serviceVal, err = sm.QueryServiceDetailInfo(key)
	if err != nil {
		logs.Log.Warn("get %s serviceInfo falid:%s", key, err.Error())
		return nil, err
	}

	//not exists
	if len(serviceVal) == 0 {
		logs.Log.Info("the servive:%s don't exist", key)
		//		err = errors.New("the servive:" + key + " don't exist")
		return nil, err
	}

	err = json.Unmarshal([]byte(serviceVal), &serviceDetailInfo)
	if err != nil {
		logs.Log.Warn("convert service %s json to struct failed:%s", key, err.Error())
		return nil, err
	}

	tcpudpRouteDetailInfo, err = managers.AssembleTCPUDPRouteInfo(serviceDetailInfo)
	if err != nil {
		logs.Log.Warn("convert serviceDetailInfo %s to TCPUDPRouteDetailInfo failed:%s", key, err.Error())
		return nil, err
	}

	return tcpudpRouteDetailInfo, err
}

func GetRouteDetailInfoByServiceKey(namespace, serviceName, serviceVersion string) (
	routes string, err error) {

	if configmanager.CfgM.GetConfigInfo().EnableTest {
		return MockGetRouteDetailInfoByServiceKey()
	}

	var (
		serviceKey            models.ServiceKey
		serviceKeyStr         string
		relationKeyPattern    string
		relations             []*models.ServiceAndRouteRelationMap
		routesVal             []byte
		tcpudpRouteDetailInfo *models.TCPUDPRouteDetailInfo
		routeKeyStr           string
		routeDetailInfo       *models.RouteDetailInfo
	)

	newNamespace := managers.ChangeEmptyToDefault(namespace)
	serviceKey = models.ServiceKey{newNamespace, serviceName, serviceVersion}
	err = CheckServiceKeyValid(&serviceKey)
	if err != nil {
		logs.Log.Warn("servicekey is not valid:%s", err.Error())
		return "", err
	}
	serviceKeyStr = managers.CovertServiceKey(managers.ServicePrefix, serviceKey)
	relationKeyPattern = managers.AssembleGetRelationsByServiceKeyPattern(&serviceKey)

	//query relation
	relationM := managers.GetReleationManager()
	relations, err = relationM.QueryRelationsByServiceKey(&serviceKey)
	if err != nil {
		logs.Log.Warn("get %s relations falid:%s", relationKeyPattern, err.Error())
		return "", err
	}

	if relations == nil || len(relations) == 0 {
		logs.Log.Info("the servive:%s don't have routes", serviceKeyStr)
		return "", nil
	}

	//judge the protocol
	if relations[0].PublishProtocol == managers.HTTPPublishProtocol ||
		relations[0].PublishProtocol == managers.HTTPSPublishProtocol {

		//http/https,use route key,go query routeAbstractList
		routeArr := make([]*models.RouteDetailInfo, 0, len(relations)*2)
		for _, relation := range relations {

			routeKeyStr = fmt.Sprintf("%s:%s:%s:%s", relation.PublishPort,
				relation.RouteType, relation.RouteName, relation.RouteVersion)

			if relation.RouteWay == managers.RouteWayDomain {
				continue
			}

			routeDetailInfo, err = GetRouteDetailInfoByRouteKey(relation.PublishPort, relation.RouteWay,
				relation.RouteType, relation.RouteName, relation.RouteVersion)

			if err != nil {
				logs.Log.Warn("get %s route failed:%s", routeKeyStr, err.Error())
				continue
			}

			if routeDetailInfo == nil {
				logs.Log.Warn("%s route don't exist", routeKeyStr)
				continue
			}

			if relation.PublishPort == managers.DefaultPublishPort {
				specialHandleDefaultPortRoute(routeDetailInfo, &routeArr)
			} else {
				routeArr = append(routeArr, routeDetailInfo)
			}
		}

		routesVal, err = json.Marshal(routeArr)
		if err != nil {
			logs.Log.Warn("%s convert RouteDetailInfos to json failed:%s", serviceKeyStr, err.Error())
			return "", err
		}
	} else { //tcp/udp ,use service key,go query service,covert to routeAbstractList
		tcpudpRouteDetailInfo, err = GetTCPUDPRouteDetailInfoByServiceKey(namespace, serviceName, serviceVersion)
		if err != nil {
			logs.Log.Warn("get %s TCPUDPRoute failed:%s", serviceKeyStr, err.Error())
			return "", err
		}

		if tcpudpRouteDetailInfo == nil {
			logs.Log.Warn(" %s TCPUDPRoute don't exist", serviceKeyStr)
			return "", nil
		}

		routeArr := make([]*models.TCPUDPRouteDetailInfo, 1)
		routeArr[0] = tcpudpRouteDetailInfo
		routesVal, err = json.Marshal(routeArr)
		if err != nil {
			logs.Log.Warn("%s convert routeArr to json failed:%s", serviceKeyStr, err.Error())
			return "", err
		}
	}

	return string(routesVal), err
}

func specialHandleDefaultPortRoute(routeDetailInfo *models.RouteDetailInfo,
	pRouteArr *[]*models.RouteDetailInfo) {
	var (
		httpdefaultport  = strconv.Itoa(configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport)
		httpsdefaultport = strconv.Itoa(configmanager.CfgM.GetConfigInfo().Listenport.Httpsdefaultport)
	)

	//http
	routeDetailInfo.Spec.PublishPort = httpdefaultport
	routeDetailInfo.Spec.PublishProtocol = managers.HTTPPublishProtocol
	*pRouteArr = append(*pRouteArr, routeDetailInfo)

	//https
	httpsRouteDetailInfo := &models.RouteDetailInfo{}
	httpsRouteDetailInfo.Kind = routeDetailInfo.Kind
	httpsRouteDetailInfo.APIVersion = routeDetailInfo.APIVersion
	httpsRouteDetailInfo.MetaData = routeDetailInfo.MetaData
	httpsRouteDetailInfo.Spec = routeDetailInfo.Spec
	httpsRouteDetailInfo.Status = routeDetailInfo.Status
	httpsRouteDetailInfo.Spec.PublishPort = httpsdefaultport
	httpsRouteDetailInfo.Spec.PublishProtocol = managers.HTTPSPublishProtocol

	*pRouteArr = append(*pRouteArr, httpsRouteDetailInfo)

}

func MockGetRouteDetailInfoByRouteKey() (routeDetailInfo *models.RouteDetailInfo, err error) {
	datafile := filepath.Join(util.GetCfgFilePath(), "ext/data/RouteDetailInfo.json")
	routeVal := util.ReadJsonfile(datafile)

	err = json.Unmarshal([]byte(routeVal), &routeDetailInfo)
	if err != nil {
		return nil, err
	}

	return routeDetailInfo, nil
}

func MockGetRouteDetailInfoByServiceKey() (string, error) {
	datafile := filepath.Join(util.GetCfgFilePath(), "ext/data/RouteDetailInfos.json")
	return string(util.ReadJsonfile(datafile)), nil
}
