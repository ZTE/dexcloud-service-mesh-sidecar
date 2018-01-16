package controllers

import (
	"apiroute/handlers"
	"apiroute/logs"

	"github.com/astaxie/beego"
)

type CacheController struct {
	beego.Controller
}

// @Title GetNameList
// @Description get data cache:name namespace
// @Success 200 {object} util.NameAndNamespace
// @Failure 500 get cache failed
// @router /internal/datasynccache/namelist [get]
func (c *CacheController) GetNameList() {
	c.EnableRender = false
	c.Data["json"] = handlers.GetNameList()
	c.ServeJSON()
}

// @Title GetVersionsAndPaths
// @Description get data cache:version path
// @Param namespace path string false "namespace"
// @Param servicename path string true "servicename"
// @Success 200 {map[string]string} version,path
// @Failure 500 get cache failed
// @router /internal/datasynccache/namespace/:namespace/servicename/:servicename/versions [get]
func (c *CacheController) GetVersionsAndPaths() {
	var (
		namespace   = c.Ctx.Input.Param(":namespace")
		serviceName = c.Ctx.Input.Param(":servicename")
	)

	logs.Log.Info("GetVersionsAndPaths:%s-%s", namespace, serviceName)

	c.EnableRender = false
	c.Data["json"] = handlers.GetVersionsAndPaths(namespace, serviceName)
	c.ServeJSON()
}
