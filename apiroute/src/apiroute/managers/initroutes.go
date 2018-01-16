package managers

import (
	"apiroute/logs"
	"apiroute/managers/configmanager"
	"apiroute/models"
	"apiroute/util"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	DiscoverPortStr        = "${discover.port}"
	DiscoverIPStr          = "${discover.ip}"
	NginxPortStr           = "${listenport.httpdefaultport}"
	APPPortStr             = "${app.httpport}"
	ApigatewayMetricsIPStr = "${apigatewaymetrics.ip}"
)

const (
	retryTotalCount = 5
	retryInterval   = 2 * time.Second
)

func InitRoutes() {
	var (
		serviceDetailInfos []*models.ServiceDetailInfo
		err                error
		isSuccess          bool
	)

	serviceDetailInfos, err = GetInitServiceInfo()
	if err != nil || serviceDetailInfos == nil || len(serviceDetailInfos) == 0 {
		logs.Log.Warn("no routes need to initialize")
		return
	}

	routeListener := RouteListener{}
	for _, serviceDetailInfo := range serviceDetailInfos {

		//check
		if serviceDetailInfo == nil {
			continue
		}

		//check node
		if !checkServiceNodesValid(serviceDetailInfo.Spec.Nodes) {
			continue
		}

		//
		serviceInfo := &models.ServiceInfo{}
		serviceInfo.ServicePrefix = ServicePrefix
		serviceInfo.ServiceKey.Namespace = ChangeEmptyToDefault(serviceDetailInfo.MetaData.Namespace)
		serviceInfo.ServiceKey.ServiceName = serviceDetailInfo.MetaData.Name
		serviceInfo.ServiceKey.ServiceVersion = serviceDetailInfo.MetaData.Version
		serviceInfo.ServiceValue = *serviceDetailInfo

		for retryCount := 0; retryCount < retryTotalCount; retryCount++ {
			isSuccess, err = routeListener.OnSave(serviceInfo)
			if err != nil {
				logs.Log.Warn("save routes and relations,have some error:%s", err.Error())
				time.Sleep(retryInterval)
				continue
			}
			break
		}

		if err != nil {
			logs.Log.Warn("[%s-%s-%s] after retry:%d ,still failed break",
				serviceInfo.ServiceKey.Namespace, serviceInfo.ServiceKey.ServiceName, serviceInfo.ServiceKey.ServiceVersion, retryTotalCount)
			break
		}

		if !isSuccess {
			logs.Log.Warn("save routes and relations failed:%s", err.Error())
		}
	}

	return

}

var checkServiceNodesValid = func(nodes []models.ServiceNodeObject) bool {
	if nodes == nil || len(nodes) == 0 {
		return false
	}

	//check ip
	for _, node := range nodes {
		if node.IP == "" {
			return false
		}
	}

	return true
}

func GetInitServiceInfo() (serviceDetailInfos []*models.ServiceDetailInfo, err error) {
	var (
		datafile string
		data     []byte
	)

	//read user-define json file:conf/ext/initRoutes/msb.json
	datafile = filepath.Join(util.GetCfgFilePath(), "ext/initRoutes/msb.json")

	data, err = ioutil.ReadFile(datafile)
	if err != nil {
		logs.Log.Warn("read initroutes:%s file failed:%s", datafile, err.Error())
		return serviceDetailInfos, err
	}

	newData := specialHandleData(string(data))
	err = json.Unmarshal([]byte(newData), &serviceDetailInfos)
	if err != nil {
		logs.Log.Warn("convert initroutes data json to struct failed:%s", err.Error())
		return serviceDetailInfos, err
	}

	return serviceDetailInfos, nil
}

func specialHandleData(data string) string {
	var (
		DiscoverPort        = strconv.FormatInt(configmanager.CfgM.GetConfigInfo().Discover.Port, 10)
		DiscoverIP          = configmanager.CfgM.GetConfigInfo().Discover.IP
		NginxPort           = strconv.Itoa(configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport)
		APPPort             = strconv.Itoa(configmanager.CfgM.GetConfigInfo().BasicConfig.Listen.HTTPPort)
		ApigatewayMetricsIP = configmanager.CfgM.GetConfigInfo().ApigatewayMetrics.IP
	)

	data = strings.Replace(data, DiscoverPortStr, DiscoverPort, -1)
	data = strings.Replace(data, DiscoverIPStr, DiscoverIP, -1)
	data = strings.Replace(data, NginxPortStr, NginxPort, -1)
	data = strings.Replace(data, APPPortStr, APPPort, -1)
	data = strings.Replace(data, ApigatewayMetricsIPStr, ApigatewayMetricsIP, -1)
	return data
}
