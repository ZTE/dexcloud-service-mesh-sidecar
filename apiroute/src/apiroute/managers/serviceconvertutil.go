package managers

import (
	"apiroute/logs"
	"apiroute/models"
	"apiroute/util"
	"errors"
	"strings"

	"time"
)

const (
	ServiceKind   = "service"
	ServicePrefix = "discover:microservices"
)

const (
	APIVersion = "v1"
)

const (
	ServiceKeyPattern = "discover:microservices:namespace:servicename:serviceversion"
)

var AssembleServiceInfo = func(serviceUnit *util.ServiceUnit) (serviceInfo *models.ServiceInfo, err error) {
	if serviceUnit == nil {
		err = errors.New("assemble service info failed.serviceUnit is nil")
		logs.Log.Error("assemble service info failed.serviceUnit is nil")
		return nil, err
	}

	serviceInfo = &models.ServiceInfo{}

	//prefix
	serviceInfo.ServicePrefix = ServicePrefix

	//service key
	serviceInfo.ServiceKey.Namespace = ChangeEmptyToDefault(serviceUnit.Namespace)
	serviceInfo.ServiceKey.ServiceName = serviceUnit.Name
	serviceInfo.ServiceKey.ServiceVersion = serviceUnit.Version

	//service detailinfo
	serviceInfo.ServiceValue.Kind = ServiceKind
	serviceInfo.ServiceValue.APIVersion = APIVersion
	serviceInfo.ServiceValue.Status = "1"

	serviceInfo.ServiceValue.MetaData.Name = serviceUnit.Name
	serviceInfo.ServiceValue.MetaData.Version = serviceUnit.Version
	serviceInfo.ServiceValue.MetaData.Namespace = serviceUnit.Namespace
	//UID
	t := time.Now().UTC()
	serviceInfo.ServiceValue.MetaData.UpdateTimestamp = t.Format("2006-01-02T15:04:05Z")
	serviceInfo.ServiceValue.MetaData.Labels = handleLables(serviceUnit.Labels)
	//Annotations

	serviceInfo.ServiceValue.Spec.VisualRange = serviceUnit.VisualRange
	serviceInfo.ServiceValue.Spec.URL = serviceUnit.URL
	serviceInfo.ServiceValue.Spec.Path = serviceUnit.Path
	serviceInfo.ServiceValue.Spec.PublishPort = serviceUnit.PublishPort
	serviceInfo.ServiceValue.Spec.Host = serviceUnit.Host

	//portal
	if strings.EqualFold(strings.TrimSpace(serviceUnit.Protocol), ProtocolPORTAL) {
		serviceInfo.ServiceValue.Spec.Protocol = ProtocolHTTP
		serviceInfo.ServiceValue.Spec.Custom = "portal"
	} else {
		serviceInfo.ServiceValue.Spec.Protocol = serviceUnit.Protocol
	}

	serviceInfo.ServiceValue.Spec.LBPolicy = serviceUnit.LBPolicy
	serviceInfo.ServiceValue.Spec.EnableSSL = serviceUnit.EnableSSL
	serviceInfo.ServiceValue.Spec.EnableTLS = serviceUnit.EnableTLS
	serviceInfo.ServiceValue.Spec.SwaggerURL = serviceUnit.SwaggerURL
	serviceInfo.ServiceValue.Spec.EnableReferMatch = serviceUnit.EnableReferMatch
	serviceInfo.ServiceValue.Spec.ProxyRule = serviceUnit.ProxyRule
	serviceInfo.ServiceValue.Spec.Nodes = calcServiceNodes(serviceUnit.Instances)
	return serviceInfo, nil

}

var calcServiceNodes = func(instances []util.InstanceUnit) (nodes []models.ServiceNodeObject) {
	if instances == nil || len(instances) == 0 {
		nodes = make([]models.ServiceNodeObject, 0, 0)
		return nodes
	}

	nodes = make([]models.ServiceNodeObject, 0, len(instances))

	for _, instance := range instances {
		node := models.ServiceNodeObject{}
		node.IP = instance.ServiceIP
		node.Port = instance.ServicePort
		node.TTL = -1
		node.LBServerParams = instance.LBServerParams
		nodes = append(nodes, node)
	}
	return nodes
}

var ChangeEmptyToDefault = func(namespace string) string {
	if len(namespace) == 0 {
		return "default"
	}
	return namespace
}

//lables [] to map
var handleLables = func(oldLabels []string) (newLables map[string]string) {
	newLables = make(map[string]string)

	if oldLabels == nil || len(oldLabels) == 0 {
		return newLables
	}

	for _, label := range oldLabels {
		kvp := strings.Split(label, ":")
		if len(kvp) < 2 {
			newLables[kvp[0]] = ""
		} else {
			newLables[kvp[0]] = kvp[1]
		}
	}
	return newLables
}

//////////////////////////////////////////////////////////
//discover:microservices:namespace:servicename:serviceversion
var CovertServiceKey = func(servicePrefix string,
	serviceKey models.ServiceKey) (strKey string) {
	strKey = servicePrefix + ":" + serviceKey.Namespace + ":" +
		serviceKey.ServiceName + ":" + serviceKey.ServiceVersion
	return strKey
}

var GetAllServiceKeysKeyPattern = func() string {
	return ServicePrefix + ":*"
}
