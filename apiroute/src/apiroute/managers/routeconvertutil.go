package managers

import (
	"apiroute/logs"
	cfg "apiroute/managers/configmanager"
	"apiroute/models"
	"apiroute/util"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// protocol
const (
	ProtocolHTTP   = "HTTP"
	ProtocolREST   = "REST"
	ProtocolUI     = "UI"
	ProtocolUDP    = "UDP"
	ProtocolTCP    = "TCP"
	ProtocolPORTAL = "PORTAL"
)

const (
	CustumProtocol = "portal"
)

const (
	RouteWayIP     = "ip"
	RouteWayDomain = "domain"
)

const (
	RouteIPPrefix     = "msb"
	RouteDomainPrefix = "msb:host"
)

const (
	DefaultPublishPort   = "routing"
	HTTPPublishProtocol  = "http"
	HTTPSPublishProtocol = "https"
)

const (
	RouteKind = "route"
)

const (
	APIPrefix  = "/api/"
	IUIPrefix  = "/iui/"
	APIRegExp1 = "/(api)/([^/]+)/([Vv]\\d+(?:\\.\\d+)*)$"
	APIRegExp2 = "/(api)/([Vv]\\d+(?:\\.\\d+)*)/([^/]+)$"
	APIRegExp3 = "/api/([Vv]\\d+(?:\\.\\d+)*)$"
)

const (
	APIRouteType    = "api"
	IUIRouteType    = "iui"
	CustomRouteType = "custom"
)

const (
	RelationKeyPattern = "namespace:servicename:serviceversion:protocol:port:routeway:routetype:routename:routeversion"
)

////////////////////////////////////////////////////////////////////////////////
func AssembleRouteInfo(serviceInfo *models.ServiceInfo) (routeInfos []*models.RouteInfo,
	relations []*models.ServiceAndRouteRelationMap, err error) {

	//check
	if serviceInfo == nil {
		err = errors.New("assemble route info fialed.serviceInfo is empty")
		logs.Log.Info("assemble route info fialed.serviceInfo is empty")
		return nil, nil, err
	}

	//check protocol
	originalProtocol := serviceInfo.ServiceValue.Spec.Protocol
	if !strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolHTTP) &&
		!strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolREST) &&
		!strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolUI) {
		err = errors.New("protocol is wrong" + serviceInfo.ServiceValue.Spec.Protocol)
		logs.Log.Warn("protocol is wrong" + serviceInfo.ServiceValue.Spec.Protocol)
		return nil, nil, err
	}

	routeInfos = make([]*models.RouteInfo, 0)

	//check routeWay
	routeWays := cfg.CfgM.GetRouteWay()
	for _, routeWay := range routeWays {
		if strings.EqualFold(strings.ToLower(strings.TrimSpace(routeWay)), RouteWayIP) { //ip
			ipRoutes := HandleIPRouteInfo(serviceInfo)
			if ipRoutes == nil || len(ipRoutes) == 0 {
				continue
			}
			for _, v := range ipRoutes {
				routeInfos = append(routeInfos, v)
			}
		} else { //domain
			domainRoute := HanldeDomainRouteInfo(serviceInfo)
			if domainRoute == nil {
				continue
			}
			routeInfos = append(routeInfos, domainRoute)
		}
	}

	relations = make([]*models.ServiceAndRouteRelationMap, 0)
	for _, routeInfo := range routeInfos {
		relation := &models.ServiceAndRouteRelationMap{}
		relation.Namespace = serviceInfo.ServiceKey.Namespace
		relation.ServiceName = serviceInfo.ServiceKey.ServiceName
		relation.ServiceVersion = serviceInfo.ServiceKey.ServiceVersion
		relation.PublishProtocol = routeInfo.RouteValue.Spec.PublishProtocol
		relation.PublishPort = routeInfo.RouteKey.PublishPort
		if routeInfo.RoutePrefix == RouteIPPrefix {
			relation.RouteWay = RouteWayIP
		} else {
			relation.RouteWay = RouteWayDomain
		}
		relation.RouteType = routeInfo.RouteKey.RouteType
		relation.RouteName = routeInfo.RouteKey.RouteName
		relation.RouteVersion = routeInfo.RouteKey.RouteVersion
		relations = append(relations, relation)
	}

	return routeInfos, relations, err
}

//handle domain
func HanldeDomainRouteInfo(serviceInfo *models.ServiceInfo) (routeInfo *models.RouteInfo) {
	var (
		path           = serviceInfo.ServiceValue.Spec.Path
		protocol       = serviceInfo.ServiceValue.Spec.Protocol
		namespace      = serviceInfo.ServiceValue.MetaData.Namespace
		serviceName    = serviceInfo.ServiceValue.MetaData.Name
		serviceVersion = serviceInfo.ServiceValue.MetaData.Version
		host           = serviceInfo.ServiceValue.Spec.Host
		custom         = serviceInfo.ServiceValue.Spec.Custom
		originalURL    = serviceInfo.ServiceValue.Spec.URL
	)

	//route key
	pRouteKey, rewriteURL, err := calcDomainRouteKeyExcludePulishPort(path, protocol,
		serviceName, serviceVersion, originalURL)
	if err != nil {
		logs.Log.Info("calc route key failed:" + serviceName + ":" + serviceVersion)
		return nil
	}

	routeInfo = &models.RouteInfo{}

	routeInfo.RoutePrefix = RouteDomainPrefix
	routeInfo.RouteKey.PublishPort = calcHostAttribute(host, namespace, serviceName)
	routeInfo.RouteKey.RouteType = pRouteKey.RouteType
	routeInfo.RouteKey.RouteName = pRouteKey.RouteName
	routeInfo.RouteKey.RouteVersion = pRouteKey.RouteVersion

	routeInfo.RouteValue = assembleRouteDetailInfo(serviceInfo,
		&routeInfo.RouteKey, RouteWayDomain)

	//special handle PublishPort PublishProtocol
	if strings.EqualFold(custom, CustumProtocol) {
		routeInfo.RouteValue.Spec.PublishProtocol = HTTPSPublishProtocol
	} else {
		routeInfo.RouteValue.Spec.PublishProtocol = HTTPPublishProtocol
	}
	//special handle URL
	routeInfo.RouteValue.Spec.URL = rewriteURL

	if routeInfo.RouteValue.Spec.URL == "/" {
		routeInfo.RouteValue.Spec.URL = ""
	}

	return routeInfo

}

//handle ip
func HandleIPRouteInfo(serviceInfo *models.ServiceInfo) (routeInfos []*models.RouteInfo) {
	routeInfos = make([]*models.RouteInfo, 0)
	ports := calcRealPorts(serviceInfo.ServiceValue.Spec.PublishPort)
	if ports == nil || len(ports) == 0 { //default http|https
		routeInfo := assembleIPRouteInfo(serviceInfo, DefaultPublishPort, HTTPPublishProtocol)
		routeInfos = append(routeInfos, routeInfo)
	} else if len(ports) == 1 { //user-define https port
		routeInfo := assembleIPRouteInfo(serviceInfo, ports[0], HTTPSPublishProtocol)
		routeInfos = append(routeInfos, routeInfo)
	} else { //user-define https|https port
		//first https
		{
			routeInfo := assembleIPRouteInfo(serviceInfo, ports[0], HTTPSPublishProtocol)
			routeInfos = append(routeInfos, routeInfo)
		}
		//second http
		{
			routeInfo := assembleIPRouteInfo(serviceInfo, ports[1], HTTPPublishProtocol)
			routeInfos = append(routeInfos, routeInfo)
		}
	}
	return routeInfos
}

func assembleIPRouteInfo(serviceInfo *models.ServiceInfo, port, publishProtocol string) *models.RouteInfo {

	var (
		path           = serviceInfo.ServiceValue.Spec.Path
		protocol       = serviceInfo.ServiceValue.Spec.Protocol
		serviceName    = serviceInfo.ServiceValue.MetaData.Name
		serviceVersion = serviceInfo.ServiceValue.MetaData.Version
	)

	//route key
	pRouteKey, err := calcIPRouteKeyExcludePulishPort(path, protocol,
		serviceName, serviceVersion)
	if err != nil {
		logs.Log.Info("calc route key failed:" + serviceName + ":" + serviceVersion)
		return nil
	}

	routeInfo := &models.RouteInfo{}

	routeInfo.RoutePrefix = RouteIPPrefix

	routeInfo.RouteKey.PublishPort = port
	routeInfo.RouteKey.RouteType = pRouteKey.RouteType
	routeInfo.RouteKey.RouteName = pRouteKey.RouteName
	routeInfo.RouteKey.RouteVersion = pRouteKey.RouteVersion

	routeInfo.RouteValue = assembleRouteDetailInfo(serviceInfo,
		&routeInfo.RouteKey, RouteWayIP)

	//special handle publishprot and publishprotocol
	if strings.EqualFold(port, DefaultPublishPort) {
		routeInfo.RouteValue.Spec.PublishPort = ""
		routeInfo.RouteValue.Spec.PublishProtocol = HTTPPublishProtocol
	} else {
		routeInfo.RouteValue.Spec.PublishPort = port
		routeInfo.RouteValue.Spec.PublishProtocol = publishProtocol
	}

	//special handle URL
	if routeInfo.RouteValue.Spec.URL == "/" {
		routeInfo.RouteValue.Spec.URL = ""
	}
	return routeInfo

}

//route detail info
func assembleRouteDetailInfo(serviceInfo *models.ServiceInfo, routekey *models.RouteKey, routeWay string) (routeVal models.RouteDetailInfo) {

	//
	var (
		serviceMetaData = &serviceInfo.ServiceValue.MetaData
		serviceSpec     = &serviceInfo.ServiceValue.Spec
		routeType       = routekey.RouteType
		routeName       = routekey.RouteName
		routeVersion    = routekey.RouteVersion
	)

	routeVal.Kind = RouteKind
	routeVal.APIVersion = APIVersion
	routeVal.Status = serviceInfo.ServiceValue.Status

	routeVal.MetaData.Name = routeName
	routeVal.MetaData.Version = routeVersion
	routeVal.MetaData.Namespace = serviceMetaData.Namespace
	routeVal.MetaData.ServiceName = serviceMetaData.Name
	routeVal.MetaData.ServiceVersion = serviceMetaData.Version
	//UID
	t := time.Now().UTC()
	routeVal.MetaData.UpdateTimestamp = t.Format("2006-01-02T15:04:05Z")
	routeVal.MetaData.Labels = serviceMetaData.Labels
	//Annotations

	routeVal.Spec.VisualRange = calcVisualRange(serviceSpec.VisualRange)

	//may special handle external
	routeVal.Spec.URL = serviceSpec.URL

	//may special handle external
	routeVal.Spec.PublishPort = serviceSpec.PublishPort
	//routeVal.Spec.PublishProtocol

	routeVal.Spec.Host = calcHostAttribute(serviceSpec.Host, serviceMetaData.Namespace,
		serviceMetaData.Name)
	if serviceSpec.SwaggerURL == "" && routeType == APIRouteType {
		routeVal.Spec.Apijson = routeVal.Spec.URL + "/swagger.json"
	} else {
		routeVal.Spec.Apijson = serviceSpec.SwaggerURL
	}

	//"[apiJson Type] 0：local file  1： remote file", allowableValues = "0,1", example = "1"
	routeVal.Spec.Apijsontype = "1"
	routeVal.Spec.MetricsURL = "/admin/metrics"

	routeVal.Spec.LBPolicy = serviceSpec.LBPolicy
	routeVal.Spec.ConsulServiceName = calcConsulServiceNameAttribute(serviceMetaData.Namespace,
		serviceMetaData.Name)
	routeVal.Spec.UseOwnUpstream = calcUseOwnUpstreamAttribute(serviceSpec.LBPolicy)

	routeVal.Spec.Enablessl = calcEnableSSLAttribute(serviceSpec.Custom, routeWay, serviceSpec.EnableSSL)
	// "[control Range] 0：default   1：readonly  2：hidden ", allowableValues = "0,1,2", example = "0"
	routeVal.Spec.Control = "0"
	routeVal.Spec.Scenario = calcScenarioAttribute(serviceMetaData.Labels)
	routeVal.Spec.EnableReferMatch = serviceSpec.EnableReferMatch
	routeVal.Spec.ConnectTimeOut = ""
	routeVal.Spec.ReadTimeout = strings.TrimSuffix(serviceSpec.ProxyRule.HTTPProxy.ReadTimeout, "s")
	routeVal.Spec.SendTimeout = strings.TrimSuffix(serviceSpec.ProxyRule.HTTPProxy.SendTimeout, "s")
	routeVal.Spec.Nodes = calcRouteNodesAttribute(serviceSpec.Nodes,
		serviceSpec.Custom, routeWay,
		serviceMetaData.Namespace,
		serviceMetaData.Name, serviceInfo.ServiceValue.MetaData.Version)
	return routeVal
}

//parse pulishPorts
func calcRealPorts(pulishPorts string) (realPorts []string) {
	if len(pulishPorts) == 0 {
		return make([]string, 0)
	}

	ports := strings.Split(pulishPorts, "|")

	if len(ports) == 1 {
		realPorts = ports
	} else {
		//"https|"
		if len(ports[1]) == 0 {
			realPorts = ports[0:1]
		} else {
			realPorts = ports[0:2]
		}
	}
	return realPorts
}

//route key
func calcIPRouteKeyExcludePulishPort(path, protocol, serviceName, serviceVersion string) (
	routeKey *models.RouteKey, err error) {

	publishURL := buildPublishURLByParams(serviceName, serviceVersion, protocol, path)

	if len(publishURL) == 0 {
		err = errors.New("calc route key failed.servicename:" + serviceName + publishURL)
		logs.Log.Info("calc route key failed.servicename:" + serviceName + publishURL)
		return nil, err
	}

	routeKey = &models.RouteKey{}
	//calc route type
	routeKey.RouteType, routeKey.RouteName, routeKey.RouteVersion = parsePublishURL(publishURL)

	return routeKey, err
}

var calcDomainRouteKeyExcludePulishPort = func(path, protocol, serviceName, serviceVersion,
	originalURL string) (routeKey *models.RouteKey, rewriteURL string, err error) {
	publishURL, rewriteURL := buildDomainForwardRuleByParams(serviceName, serviceVersion, protocol, path, originalURL)
	if len(publishURL) == 0 {
		err = errors.New("calc route key failed.servicename:" + serviceName + publishURL)
		logs.Log.Info("calc route key failed.servicename:" + serviceName + publishURL)
		return nil, originalURL, err
	}

	routeKey = &models.RouteKey{}
	routeKey.RouteType, routeKey.RouteName, routeKey.RouteVersion = parsePublishURL(publishURL)
	return routeKey, rewriteURL, err
}

//route url
func buildPublishURLByParams(serviceName, version, protocol, path string) string {
	if protocol == ProtocolTCP || protocol == ProtocolUDP {
		return ""
	}

	//path
	if path != "" && path != "/" {
		return path
	}
	if protocol == ProtocolHTTP {
		if version != "" {
			return "/" + serviceName + "/" + version
		}
		return "/" + serviceName
	}
	if protocol == ProtocolREST {
		if version != "" {
			return "/api/" + serviceName + "/" + version
		}
		return "/api/" + serviceName
	}
	if protocol == ProtocolUI {
		return "/iui/" + serviceName
	}

	return ""
}

//domian route url
var buildDomainForwardRuleByParams = func(serviceName, version, protocol, path,
	originalURL string) (publishURL, rewriteURL string) {

	if protocol == ProtocolTCP || protocol == ProtocolUDP {
		return "", originalURL
	}

	//path
	if path != "" && path != "/" {
		return path, originalURL
	}

	if protocol == ProtocolHTTP {
		if version != "" {
			publishURL = "/" + serviceName + "/" + version
		} else {
			publishURL = "/" + serviceName
		}

	}

	if protocol == ProtocolREST {
		if version != "" {
			publishURL = "/api/" + serviceName + "/" + version
		} else {
			publishURL = "/api/" + serviceName
		}
	}

	if protocol == ProtocolUI {
		publishURL = "/iui/" + serviceName
	}

	rewriteURL = originalURL
	//handle originalURL=="/" & originalURL==publishURL
	if originalURL == "/" || originalURL == publishURL {
		publishURL = "/"
		rewriteURL = "/"
	}

	return publishURL, rewriteURL
}

//routetype,routename,version from publishURL
var parsePublishURL = func(publishURL string) (routeType, routeName, routeVersion string) {
	if publishURL == "" {
		return "", "", ""
	}

	//iui
	if strings.HasPrefix(publishURL, IUIPrefix) {
		strArray := strings.Split(publishURL, "/")
		if len(strArray) == 3 { //"/routertype/routename"
			routeType = IUIRouteType
			routeName = strArray[2]
			routeVersion = ""
			return routeType, routeName, routeVersion
		}
	}

	//api
	if strings.HasPrefix(publishURL, APIPrefix) {

		strArray := strings.Split(publishURL, "/")

		if len(strArray) == 3 { //"/routertype/routename"
			api3 := regexp.MustCompile(APIRegExp3)
			if api3.MatchString(publishURL) {
				goto custom_handle
			}
			routeType = APIRouteType
			routeName = strArray[2]
			routeVersion = ""
			return routeType, routeName, routeVersion
		}

		if len(strArray) == 4 { //"/routertype/routename/routerversion"
			//
			api1 := regexp.MustCompile(APIRegExp1)
			if api1.MatchString(publishURL) {
				routeType = APIRouteType
				routeName = strArray[2]
				routeVersion = strArray[3]
				logs.Log.Info("api1:%s-%s-%s", routeType, routeName, routeVersion)
				return routeType, routeName, routeVersion
			}

			api2 := regexp.MustCompile(APIRegExp2)
			if api2.MatchString(publishURL) {
				routeType = APIRouteType
				routeVersion = strArray[2]
				routeName = strArray[3]
				logs.Log.Info("api2:%s-%s-%s", routeType, routeVersion, routeName)
				return routeType, routeName, routeVersion
			}
		}

	}

custom_handle:
	//custom
	routeType = CustomRouteType
	routeName = publishURL
	routeVersion = ""

	return routeType, routeName, routeVersion
}

//route type
func calcRouteType(publishURL string) string {
	if strings.HasPrefix(publishURL, APIPrefix) {
		return APIRouteType
	} else if strings.HasPrefix(publishURL, IUIPrefix) {
		return IUIRouteType
	} else {
		return CustomRouteType
	}
}

//visualrange
func calcVisualRange(vs string) (realVS string) {
	visualRanges := strings.Split(vs, "|")
	if len(visualRanges) > 1 {
		cfgvs := cfg.CfgM.GetVisualRange()
		if len(cfgvs) > 1 {
			realVS = "0"
		} else {
			realVS = cfgvs[0]
		}
	} else {
		realVS = vs
	}

	return realVS
}

//host
func calcHostAttribute(host string, namespace string, serviceName string) (realHost string) {
	//host empty,servicename-namespace
	if len(host) == 0 {
		if strings.EqualFold(namespace, "default") || len(namespace) == 0 {
			realHost = serviceName
		} else {
			realHost = serviceName + "-" + namespace
		}
	} else {
		realHost = host
	}
	return strings.ToLower(realHost)
}

//ConsulServiceName
func calcConsulServiceNameAttribute(namespace string, serviceName string) (realConsulServiceName string) {
	if strings.EqualFold(namespace, "default") || len(namespace) == 0 {
		realConsulServiceName = serviceName
	} else {
		realConsulServiceName = serviceName + "-" + namespace
	}
	return realConsulServiceName
}

//UseOwnUpstream
func calcUseOwnUpstreamAttribute(LBPolicy string) (realUseOwnUpstream string) {
	if strings.EqualFold(LBPolicy, "ip_hash") {
		realUseOwnUpstream = "1"
	} else {
		realUseOwnUpstream = "0"
	}
	return realUseOwnUpstream
}

//Enablessl
func calcEnableSSLAttribute(portal string, routeWay string, enableSSL bool) bool {
	if strings.EqualFold(portal, CustumProtocol) {
		realEnableSSL := false
		if strings.EqualFold(routeWay, RouteWayDomain) {
			realEnableSSL = true
		} else {
			realEnableSSL = false
		}
		return realEnableSSL
	}

	return enableSSL
}

//PublishProtocol
func calcPublishProtocolAttribute(enableSSL bool) string {
	if enableSSL {
		return "https"
	}
	return "http"
}

//Scenario
//1、when： enable_router=true or enable_router不存在 then： scenario | 1
//2、when： enable_cos=true then： scenario | 2
func calcScenarioAttribute(labels map[string]string) int {

	if labels == nil || len(labels) == 0 {
		return 1 //default enable_router
	}

	scenario := 0
	if val, ok := labels["enable_router"]; ok {
		if boolVal, _ := strconv.ParseBool(val); boolVal {
			scenario = scenario | 1
		}
	} else {
		scenario = scenario | 1
	}

	if val, ok := labels["enable_cos"]; ok {
		if boolVal, _ := strconv.ParseBool(val); boolVal {
			scenario = scenario | 2
		}
	}

	return scenario
}

//RouteNodes
func calcRouteNodesAttribute(instances []models.ServiceNodeObject, portal, routeWay,
	namespace, serviceName, serviceVersion string) (nodes []models.RouteNodeObject) {
	if instances == nil || len(instances) == 0 {
		return make([]models.RouteNodeObject, 0, 0)
	}

	//portal domain
	if strings.EqualFold(portal, CustumProtocol) && strings.EqualFold(routeWay, RouteWayDomain) {
		url := "http://" + cfg.CfgM.GetConfigInfo().Discover.IP + ":" +
			strconv.FormatInt(cfg.CfgM.GetConfigInfo().Discover.Port, 10) + "/api/microservices/v1/services/" +
			serviceName + "/version/" + serviceVersion + "/allpublishaddress?namespace=" + namespace +
			"&visualRange=0"

		logs.Log.Info(url)

		//get pulishurl from sdclient
		b, httperr := util.HTTPGet(url, "")
		if httperr != nil {
			logs.Log.Info("get pulishurl from sdclient failed :" + url)
			return make([]models.RouteNodeObject, 0, 0)
		}

		var addresses []models.AllPulishAddressForIP
		if err := json.Unmarshal(b, &addresses); err != nil || len(addresses) == 0 {
			return make([]models.RouteNodeObject, 0, 0)
		}

		nodes = make([]models.RouteNodeObject, 0, len(addresses))
		for _, val := range addresses {
			if len(val.IP) != 0 && strings.EqualFold(val.PublishProtocol, "https") {
				if port, changeerr := strconv.Atoi(val.Port); changeerr == nil {
					nodes = append(nodes, models.RouteNodeObject{val.IP, port, 0, ""})
				}
			}
		}
	}

	nodes = make([]models.RouteNodeObject, 0, len(instances))
	for _, instance := range instances {
		node := models.RouteNodeObject{}
		node.IP = instance.IP
		node.Port, _ = strconv.Atoi(instance.Port)
		node.LBServerParams = instance.LBServerParams
		nodes = append(nodes, node)
	}
	return nodes
}

//////////////////////////////////////////////////
func AssembleTCPUDPRouteInfo(serviceDetailInfo *models.ServiceDetailInfo) (routeVal *models.TCPUDPRouteDetailInfo, err error) {

	if serviceDetailInfo == nil {
		err = errors.New("assemble tcp/udp route info failed.serviceDetailInfo is empty")
		logs.Log.Info("assemble tcp/udp info failed.serviceDetailInfo is empty")
		return nil, err
	}

	//check protocol
	originalProtocol := serviceDetailInfo.Spec.Protocol
	if !strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolTCP) &&
		!strings.EqualFold(strings.TrimSpace(originalProtocol), ProtocolUDP) {
		err = errors.New("protocol is wrong" + originalProtocol)
		logs.Log.Warn("protocol is wrong" + originalProtocol)
		return nil, err
	}

	routeVal = &models.TCPUDPRouteDetailInfo{}
	routeVal.Kind = RouteKind
	routeVal.APIVersion = APIVersion
	routeVal.Status = serviceDetailInfo.Status

	routeVal.MetaData.Name = serviceDetailInfo.MetaData.Name
	routeVal.MetaData.Version = serviceDetailInfo.MetaData.Version
	routeVal.MetaData.Namespace = serviceDetailInfo.MetaData.Namespace
	routeVal.MetaData.ServiceName = serviceDetailInfo.MetaData.Name
	routeVal.MetaData.ServiceVersion = serviceDetailInfo.MetaData.Version
	routeVal.MetaData.UpdateTimestamp = serviceDetailInfo.MetaData.UpdateTimestamp

	routeVal.Spec.VisualRange = serviceDetailInfo.Spec.VisualRange
	routeVal.Spec.PublishPort = serviceDetailInfo.Spec.PublishPort
	routeVal.Spec.PublishProtocol = serviceDetailInfo.Spec.Protocol
	routeVal.Spec.Host = calcHostAttribute(serviceDetailInfo.Spec.Host,
		serviceDetailInfo.MetaData.Namespace, serviceDetailInfo.MetaData.Name)
	routeVal.Spec.EnableTLS = serviceDetailInfo.Spec.EnableTLS
	routeVal.Spec.LBPolicy = serviceDetailInfo.Spec.LBPolicy
	routeVal.Spec.ConnectTimeOut = ""
	routeVal.Spec.ProxyTimeout = serviceDetailInfo.Spec.ProxyRule.StreamProxy.ProxyTimeout
	routeVal.Spec.ProxyResponse = serviceDetailInfo.Spec.ProxyRule.StreamProxy.ProxyResponse
	routeVal.Spec.Nodes = calcRouteNodesAttribute(serviceDetailInfo.Spec.Nodes, serviceDetailInfo.Spec.Custom, "", "", "", "")
	return routeVal, err
}

//////////////////////////////////////////////////
//prefix:port:routetype:routename:routeversion
func ConvertRouteKey(routePrefix string,
	routeKey models.RouteKey) (strKey string) {
	strKey = routePrefix + ":" + routeKey.PublishPort + ":" + routeKey.RouteType + ":" +
		routeKey.RouteName
	if routeKey.RouteType == APIRouteType {
		strKey = strKey + ":" + routeKey.RouteVersion
	}
	return strKey
}

//namespace:servicename:serviceversion:protocol:port:routeway:routetype:routename:routeversion
func ConvertReleationKey(releation *models.ServiceAndRouteRelationMap) (strKey string) {
	strKey = releation.Namespace + ":" + releation.ServiceName + ":" +
		releation.ServiceVersion + ":" + releation.PublishProtocol + ":" +
		releation.PublishPort + ":" + releation.RouteWay +
		":" + releation.RouteType + ":" + releation.RouteName + ":" + releation.RouteVersion
	return strKey
}

func AssembleGetRelationsByNamespacePattern(namespace string) (strKey string) {
	strKey = namespace + ":*:*:*:*:*:*:*:*"
	return strKey
}

func AssembleGetRelationsByNamespaceRouteWayPattern(namespace, routeWay string) (strKey string) {
	strKey = namespace + ":*:*:*:*:" + routeWay + ":*:*:*"
	return strKey
}

func AssembleGetRelationsByServiceKeyPattern(serviceKey *models.ServiceKey) string {
	strKey := serviceKey.Namespace + ":" + serviceKey.ServiceName + ":" +
		serviceKey.ServiceVersion + ":*:*:*:*:*:*"
	return strKey
}

/////////////////////////////////////////////
func ConvertKeysToRelations(keys []string) (relations []*models.ServiceAndRouteRelationMap) {
	if keys != nil && len(keys) != 0 {
		relations = make([]*models.ServiceAndRouteRelationMap, 0, len(keys))
		for _, key := range keys {
			strArr := strings.Split(key, ":")
			if len(strArr) == 9 {
				relation := &models.ServiceAndRouteRelationMap{strArr[0], strArr[1],
					strArr[2], strArr[3], strArr[4], strArr[5], strArr[6], strArr[7], strArr[8]}
				relations = append(relations, relation)
			} else {
				logs.Log.Warn("CovertKeysToRelations:the key %s is not valid", key)
			}
		}

	}
	return relations
}
