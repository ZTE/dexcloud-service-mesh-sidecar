package controllers

import (
	"apiroute/handlers"

	"github.com/astaxie/beego"
)

type ConfController struct {
	beego.Controller
}

// @Title GetDefaultPorts
// @Description get default ports
// @Success 200 {object} models.HTTPHTTPSDefaultPorts
// @Failure 500 get ports from config file failed
// @router /defaultports [get]
func (c *ConfController) GetDefaultPorts() {
	c.EnableRender = false
	c.Data["json"] = handlers.GetDefaultPorts()
	c.ServeJSON()
}
