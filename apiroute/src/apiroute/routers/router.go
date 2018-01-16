package routers

import (
	"apiroute/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api/route/v1",
		beego.NSNamespace("/namespaces",
			beego.NSInclude(&controllers.NamespacesController{})),
		beego.NSNamespace("/routelist",
			beego.NSInclude(&controllers.RouteListController{})),
		beego.NSNamespace("/routes",
			beego.NSInclude(&controllers.RoutesController{})),
		beego.NSNamespace("/config",
			beego.NSInclude(&controllers.ConfController{})),
		beego.NSNamespace("/cache",
			beego.NSInclude(&controllers.CacheController{})),
		beego.NSNamespace("/health",
			beego.NSInclude(&controllers.HealthController{})))
	beego.AddNamespace(ns)
}
