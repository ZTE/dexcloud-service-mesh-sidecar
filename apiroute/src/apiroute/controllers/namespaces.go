package controllers

import (
	"apiroute/handlers"

	"github.com/astaxie/beego"
)

type NamespacesController struct {
	beego.Controller
}

// @Title get all namespaces
// @Description get all namespaces
// @Success 200 {[]string} namesapce
// @Failure 500 get namespaces list error
// @router / [get]
func (n *NamespacesController) GetAllNamespaces() {
	n.EnableRender = false
	n.Data["json"] = handlers.GetAllNamespaces()
	n.ServeJSON()
}
