package handlers

import (
	"apiroute/logs"
	"apiroute/managers"
	"apiroute/managers/configmanager"
	"apiroute/models"
	"apiroute/util"
	"encoding/json"
	"path/filepath"
	"strconv"
)

func GetRouteAbstractInfoListByNamespace(namespace, routeWay string) (
	routeAbstractList []*models.RouteAbstractInfo, err error) {

	if configmanager.CfgM.GetConfigInfo().EnableTest {
		return MockGetRouteAbstractInfoListByNamespace(namespace)
	}

	var (
		newNamespace      string
		newRouteWay       string
		relationKeyPatten string
		relationKeys      []string
		relations         []*models.ServiceAndRouteRelationMap
	)

	//calc relation keypattern
	newNamespace, newRouteWay = handleParams(namespace, routeWay)
	if newRouteWay == "" {
		return routeAbstractList, err
	}

	relationKeyPatten = managers.AssembleGetRelationsByNamespacePattern(newNamespace)

	//get abstract info from relation
	relationM := managers.GetReleationManager()
	relationKeys, err = relationM.QueryRelationsByKeyPattern(relationKeyPatten)
	if err != nil {
		logs.Log.Warn("%s get relations failed:%s", relationKeyPatten, err.Error())
		return nil, err
	}

	if relationKeys == nil || len(relationKeys) == 0 {
		logs.Log.Warn("%s don't have routes", namespace)
		return nil, nil
	}

	//change keys to abstracts
	relations = managers.ConvertKeysToRelations(relationKeys)

	//change relations to abstractinfo
	if relations != nil && len(relations) != 0 {
		routeAbstractList = make([]*models.RouteAbstractInfo, 0, len(relations)*2)
		for _, relation := range relations {
			convertRelationToRouteAbstractInfo(relation, &routeAbstractList)
		}
	}

	return routeAbstractList, err
}

func handleParams(namespace, routeWay string) (newNamespace, newRouteWay string) {
	//namespace
	newNamespace = managers.ChangeEmptyToDefault(namespace)

	//routeway
	if routeWay == "" || routeWay == "ip" {
		newRouteWay = managers.RouteWayIP
	}

	if routeWay == "ip|domain" || routeWay == "domain" {
		routeWays := configmanager.CfgM.GetRouteWay()
		support := false
		for _, way := range routeWays {
			if way == managers.RouteWayDomain {
				support = true
				break
			}
		}
		if routeWay == "ip|domain" {
			if support {
				newRouteWay = "*"
			} else {
				newRouteWay = managers.RouteWayIP
			}
		}

		if routeWay == "domain" {
			if support {
				newRouteWay = managers.RouteWayDomain
			} else {
				newRouteWay = ""
			}
		}
	}

	return newNamespace, newRouteWay
}

func convertRelationToRouteAbstractInfo(relation *models.ServiceAndRouteRelationMap,
	pRouteAbstractList *[]*models.RouteAbstractInfo) {

	//check relation
	if relation == nil || pRouteAbstractList == nil {
		return
	}

	//check protocol
	if relation.PublishProtocol == managers.HTTPPublishProtocol ||
		relation.PublishProtocol == managers.HTTPSPublishProtocol {
		//check routeway
		if relation.RouteWay == managers.RouteWayDomain {
			//			routeAbstractInfo := constructDomainRouteAbstractInfo(relation)
			//			if routeAbstractInfo != nil {
			//				*pRouteAbstractList = append(*pRouteAbstractList, routeAbstractInfo)
			//			}
		} else if relation.RouteWay == managers.RouteWayIP {
			//check default or uer-define
			if relation.PublishPort == managers.DefaultPublishPort {
				constructDefaultPortRouteAbstractInfo(relation, pRouteAbstractList)
			} else {
				routeAbstractInfo := constructUserDefinePortRouteAbstractInfo(relation)
				if routeAbstractInfo != nil {
					*pRouteAbstractList = append(*pRouteAbstractList, routeAbstractInfo)
				}
			}
		}
	} else {
		routeAbstractInfo := constructTCPUDPRouteAbstractInfo(relation)
		if routeAbstractInfo != nil {
			*pRouteAbstractList = append(*pRouteAbstractList, routeAbstractInfo)
		}
	}
}

func constructTCPUDPRouteAbstractInfo(relation *models.ServiceAndRouteRelationMap) (
	routeAbstractInfo *models.RouteAbstractInfo) {
	if relation == nil {
		return routeAbstractInfo
	}
	routeAbstractInfo = &models.RouteAbstractInfo{}
	routeAbstractInfo.Namespace = relation.Namespace
	routeAbstractInfo.ServiceName = relation.ServiceName
	routeAbstractInfo.ServiceVersion = relation.ServiceVersion
	routeAbstractInfo.PublishProtocol = relation.PublishProtocol
	port, _ := strconv.Atoi(relation.PublishPort)
	routeAbstractInfo.PublishPort = port

	return routeAbstractInfo
}

func constructDomainRouteAbstractInfo(relation *models.ServiceAndRouteRelationMap) (
	routeAbstractInfo *models.RouteAbstractInfo) {
	if relation == nil {
		return routeAbstractInfo
	}

	routeAbstractInfo = &models.RouteAbstractInfo{}
	routeAbstractInfo.Namespace = relation.Namespace
	routeAbstractInfo.ServiceName = relation.ServiceName
	routeAbstractInfo.ServiceVersion = relation.ServiceVersion
	routeAbstractInfo.PublishProtocol = relation.PublishProtocol
	if routeAbstractInfo.PublishProtocol == managers.HTTPPublishProtocol {
		routeAbstractInfo.PublishPort = configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport
	} else if routeAbstractInfo.PublishProtocol == managers.HTTPSPublishProtocol {
		routeAbstractInfo.PublishPort = configmanager.CfgM.GetConfigInfo().Listenport.Httpsdefaultport
	}

	routeAbstractInfo.RouteType = relation.RouteType
	routeAbstractInfo.RouterName = relation.RouteName
	routeAbstractInfo.RouterVersion = relation.RouteVersion
	routeAbstractInfo.PublishURL = relation.PublishPort

	return routeAbstractInfo
}

func constructUserDefinePortRouteAbstractInfo(relation *models.ServiceAndRouteRelationMap) (
	routeAbstractInfo *models.RouteAbstractInfo) {
	if relation == nil {
		return routeAbstractInfo
	}
	routeAbstractInfo = &models.RouteAbstractInfo{}
	routeAbstractInfo.Namespace = relation.Namespace
	routeAbstractInfo.ServiceName = relation.ServiceName
	routeAbstractInfo.ServiceVersion = relation.ServiceVersion
	routeAbstractInfo.PublishProtocol = relation.PublishProtocol
	port, _ := strconv.Atoi(relation.PublishPort)
	routeAbstractInfo.PublishPort = port
	routeAbstractInfo.RouteType = relation.RouteType
	routeAbstractInfo.RouterName = relation.RouteName
	routeAbstractInfo.RouterVersion = relation.RouteVersion
	routeAbstractInfo.PublishURL = assemblePulishURL(relation.RouteType, relation.RouteName, relation.RouteVersion)
	return routeAbstractInfo
}

func constructDefaultPortRouteAbstractInfo(relation *models.ServiceAndRouteRelationMap,
	pRouteAbstractList *[]*models.RouteAbstractInfo) {
	if relation == nil {
		return
	}
	httpRouteAbstractInfo := &models.RouteAbstractInfo{}
	httpRouteAbstractInfo.Namespace = relation.Namespace
	httpRouteAbstractInfo.ServiceName = relation.ServiceName
	httpRouteAbstractInfo.ServiceVersion = relation.ServiceVersion
	httpRouteAbstractInfo.PublishProtocol = managers.HTTPPublishProtocol
	httpRouteAbstractInfo.PublishPort = configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport
	httpRouteAbstractInfo.RouteType = relation.RouteType
	httpRouteAbstractInfo.RouterName = relation.RouteName
	httpRouteAbstractInfo.RouterVersion = relation.RouteVersion
	httpRouteAbstractInfo.PublishURL = assemblePulishURL(relation.RouteType, relation.RouteName, relation.RouteVersion)
	*pRouteAbstractList = append(*pRouteAbstractList, httpRouteAbstractInfo)

	httpsRouteAbstractInfo := &models.RouteAbstractInfo{}
	httpsRouteAbstractInfo.Namespace = relation.Namespace
	httpsRouteAbstractInfo.ServiceName = relation.ServiceName
	httpsRouteAbstractInfo.ServiceVersion = relation.ServiceVersion
	httpsRouteAbstractInfo.PublishProtocol = managers.HTTPSPublishProtocol
	httpsRouteAbstractInfo.PublishPort = configmanager.CfgM.GetConfigInfo().Listenport.Httpsdefaultport
	httpsRouteAbstractInfo.RouteType = relation.RouteType
	httpsRouteAbstractInfo.RouterName = relation.RouteName
	httpsRouteAbstractInfo.RouterVersion = relation.RouteVersion
	httpsRouteAbstractInfo.PublishURL = assemblePulishURL(relation.RouteType, relation.RouteName, relation.RouteVersion)
	*pRouteAbstractList = append(*pRouteAbstractList, httpsRouteAbstractInfo)
}

func assemblePulishURL(routeType, routerName, routerVersion string) string {
	if routeType == managers.APIRouteType {
		if routerVersion != "" {
			return "/" + routeType + "/" + routerName + "/" + routerVersion
		}

		return "/" + routeType + "/" + routerName
	}

	if routeType == managers.IUIRouteType {
		return "/" + routeType + "/" + routerName
	}

	if routeType == managers.CustomRouteType {
		return routerName
	}

	return ""
}

func MockGetRouteAbstractInfoListByNamespace(namespace string) (
	routeAbstractList []*models.RouteAbstractInfo, err error) {
	var routeAbstractListTemplet []models.RouteAbstractInfo

	datafile := filepath.Join(util.GetCfgFilePath(), "ext/data/RouteAbstractInfoList.json")
	if err := json.Unmarshal(util.ReadJsonfile(datafile), &routeAbstractListTemplet); err != nil {
		logs.Log.Info("covert json failed:" + err.Error())
	}

	currentPublishPort := 28003
	routeAbstractList = make([]*models.RouteAbstractInfo, 0, 50000)
	//magnification
	for i := range routeAbstractListTemplet {
		if routeAbstractListTemplet[i].PublishProtocol == "Http" ||
			routeAbstractListTemplet[i].PublishProtocol == "Https" {
			MagnifyHTTPHTTPSData(routeAbstractListTemplet[i], &routeAbstractList)

		} else {
			MagnifyTCPUDPData(routeAbstractListTemplet[i], &routeAbstractList, &currentPublishPort)
		}
	}

	return routeAbstractList, err
}

func MagnifyHTTPHTTPSData(templet models.RouteAbstractInfo, routeAbstractList *[]*models.RouteAbstractInfo) {
	for i := 0; i < 10; i++ {
		data := &models.RouteAbstractInfo{}
		data.Namespace = templet.Namespace
		data.ServiceName = templet.ServiceName + strconv.Itoa(i)
		data.ServiceVersion = templet.ServiceVersion
		data.PublishProtocol = templet.PublishProtocol
		data.PublishPort = templet.PublishPort
		data.RouteType = templet.RouteType
		data.RouterName = templet.RouterName
		data.RouterVersion = templet.RouterVersion
		data.PublishURL = templet.PublishURL
		*routeAbstractList = append(*routeAbstractList, data)
	}
}

func MagnifyTCPUDPData(templet models.RouteAbstractInfo,
	routeAbstractList *[]*models.RouteAbstractInfo, currentPublishPort *int) {
	for i := 0; i < 10; i++ {
		data := &models.RouteAbstractInfo{}
		data.Namespace = templet.Namespace
		data.ServiceName = templet.ServiceName + strconv.Itoa(i)
		data.ServiceVersion = templet.ServiceVersion
		data.PublishProtocol = templet.PublishProtocol
		(*currentPublishPort)++
		data.PublishPort = *currentPublishPort
		data.RouteType = templet.RouteType
		data.RouterName = templet.RouterName
		data.RouterVersion = templet.RouterVersion
		data.PublishURL = templet.PublishURL
		*routeAbstractList = append(*routeAbstractList, data)
	}
}
