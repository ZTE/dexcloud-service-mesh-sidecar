package managers

import (
	"apiroute/logs"
	"apiroute/models"
	"encoding/json"
	"strings"
)

type ServiceChangedListener interface {
	OnSave(serviceInfo *models.ServiceInfo) (isSuccess bool, err error)
	OnDelete(serviceInfo *models.ServiceInfo) (isSuccess bool, err error)
}

type RouteListener struct {
}

func (rl RouteListener) OnSave(serviceInfo *models.ServiceInfo) (isSuccess bool, err error) {
	var (
		routeInfos  []*models.RouteInfo
		relations   []*models.ServiceAndRouteRelationMap
		internalErr error
	)

	//check protocol:TCP/UDP
	originalProtocol := serviceInfo.ServiceValue.Spec.Protocol
	if !strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolHTTP) &&
		!strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolREST) &&
		!strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolUI) {
		//just save relation
		relation := &models.ServiceAndRouteRelationMap{serviceInfo.ServiceKey.Namespace,
			serviceInfo.ServiceKey.ServiceName, serviceInfo.ServiceKey.ServiceVersion, originalProtocol,
			serviceInfo.ServiceValue.Spec.PublishPort, "", "", "", ""}
		relationM := GetReleationManager()
		key := ConvertReleationKey(relation)
		logs.Log.Info("update relation:%s", key)
		if redisErr := relationM.SaveRelation(key, ""); redisErr != nil {
			return false, redisErr
		}
		return true, err
	}

	//convert routeinfo relation
	if routeInfos, relations, internalErr = AssembleRouteInfo(serviceInfo); internalErr != nil {
		logs.Log.Warn("covert serviceInfo %s to routeInfo failed:%s",
			serviceInfo.ServiceKey.ServiceName, internalErr.Error())
		return false, internalErr
	}

	//persistence routeInfo
	rm := GetRouteManager()
	if routeInfos != nil && len(routeInfos) != 0 {
		for _, routeInfo := range routeInfos {
			key := ConvertRouteKey(routeInfo.RoutePrefix, routeInfo.RouteKey)
			logs.Log.Info("update route:%s", key)
			val, jsonErr := json.Marshal(routeInfo.RouteValue)
			if jsonErr != nil {
				logs.Log.Warn("%s convert to json failed:%s", key, jsonErr.Error())
				continue
			}
			if redisErr := rm.SaveRoute(key, string(val)); redisErr != nil {
				return false, redisErr
			}
		}
	}

	//persistence relation
	relationM := GetReleationManager()
	if relations != nil && len(relations) != 0 {
		for _, relation := range relations {
			key := ConvertReleationKey(relation)
			logs.Log.Info("update relation:%s", key)
			if redisErr := relationM.SaveRelation(key, ""); redisErr != nil {
				return false, redisErr
			}
		}
	}

	return true, err
}

func (rl RouteListener) OnDelete(serviceInfo *models.ServiceInfo) (isSuccess bool, err error) {

	if serviceInfo == nil {
		logs.Log.Warn("delete service:serviceInfo is nil")
		return true, nil
	}

	var (
		redisErr         error
		originalProtocol = serviceInfo.ServiceValue.Spec.Protocol
		namespace        = serviceInfo.ServiceKey.Namespace
		serviceName      = serviceInfo.ServiceKey.ServiceName
		serviceVersion   = serviceInfo.ServiceKey.ServiceVersion
		publishPort      = serviceInfo.ServiceValue.Spec.PublishPort

		routeInfos  []*models.RouteInfo
		relations   []*models.ServiceAndRouteRelationMap
		internalErr error
	)

	//check protocol:TCP/UDP
	if !strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolHTTP) &&
		!strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolREST) &&
		!strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolUI) {
		//just delete relation
		relation := &models.ServiceAndRouteRelationMap{namespace,
			serviceName, serviceVersion, originalProtocol,
			publishPort, "", "", "", ""}
		relationM := GetReleationManager()
		key := ConvertReleationKey(relation)
		logs.Log.Info("delete relation:%s", key)
		if redisErr = relationM.DeleteRelation(key); redisErr != nil {
			return false, redisErr
		}

		return true, nil
	}

	//convert routeinfo relation
	if routeInfos, relations, internalErr = AssembleRouteInfo(serviceInfo); internalErr != nil {
		logs.Log.Warn("prepare delete service:covert serviceInfo %s to routeInfo failed:%s",
			serviceInfo.ServiceKey.ServiceName, internalErr.Error())
		return false, internalErr
	}

	//delete routeInfo
	rm := GetRouteManager()
	if routeInfos != nil && len(routeInfos) != 0 {
		for _, routeInfo := range routeInfos {
			key := ConvertRouteKey(routeInfo.RoutePrefix, routeInfo.RouteKey)
			logs.Log.Info("delete route:%s", key)
			if redisErr = rm.DeleteRoute(key); redisErr != nil {
				return false, redisErr
			}
		}
	}

	//delete relation
	relationM := GetReleationManager()
	if relations != nil && len(relations) != 0 {
		for _, relation := range relations {
			key := ConvertReleationKey(relation)
			logs.Log.Info("delete relation:%s", key)
			if redisErr = relationM.DeleteRelation(key); redisErr != nil {
				return false, redisErr
			}
		}
	}

	return true, err
}

type NamespaceLister struct {
}

func (nl NamespaceLister) OnSave(serviceInfo *models.ServiceInfo) (isSuccess bool, err error) {
	nm := GetNamespaceManager()
	nm.NotifyNamespaceListUpdate()
	return true, nil
}

func (nl NamespaceLister) OnDelete(serviceInfo *models.ServiceInfo) (isSuccess bool, err error) {
	nm := GetNamespaceManager()
	nm.NotifyNamespaceListUpdate()
	return true, nil
}
