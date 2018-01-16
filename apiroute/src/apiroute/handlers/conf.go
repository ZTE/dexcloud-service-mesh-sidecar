package handlers

import (
	"apiroute/managers/configmanager"
	"apiroute/models"
)

func GetDefaultPorts() *models.HTTPHTTPSDefaultPorts {
	ports := &models.HTTPHTTPSDefaultPorts{}
	ports.HTTPPort = configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport
	ports.HTTPSPort = configmanager.CfgM.GetConfigInfo().Listenport.Httpsdefaultport
	return ports
}
