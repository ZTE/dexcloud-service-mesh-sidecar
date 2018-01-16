package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["apiroute/controllers:CacheController"] = append(beego.GlobalControllerRouter["apiroute/controllers:CacheController"],
		beego.ControllerComments{
			Method: "GetNameList",
			Router: `/internal/datasynccache/namelist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apiroute/controllers:CacheController"] = append(beego.GlobalControllerRouter["apiroute/controllers:CacheController"],
		beego.ControllerComments{
			Method: "GetVersionsAndPaths",
			Router: `/internal/datasynccache/namespace/:namespace/servicename/:servicename/versions`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apiroute/controllers:ConfController"] = append(beego.GlobalControllerRouter["apiroute/controllers:ConfController"],
		beego.ControllerComments{
			Method: "GetDefaultPorts",
			Router: `/defaultports`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apiroute/controllers:HealthController"] = append(beego.GlobalControllerRouter["apiroute/controllers:HealthController"],
		beego.ControllerComments{
			Method: "GetHealthStatus",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apiroute/controllers:NamespacesController"] = append(beego.GlobalControllerRouter["apiroute/controllers:NamespacesController"],
		beego.ControllerComments{
			Method: "GetAllNamespaces",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apiroute/controllers:RouteListController"] = append(beego.GlobalControllerRouter["apiroute/controllers:RouteListController"],
		beego.ControllerComments{
			Method: "GetRouteAbstractInfoList",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apiroute/controllers:RoutesController"] = append(beego.GlobalControllerRouter["apiroute/controllers:RoutesController"],
		beego.ControllerComments{
			Method: "GetRouteDetailInfoByRouteKey",
			Router: `/internal`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apiroute/controllers:RoutesController"] = append(beego.GlobalControllerRouter["apiroute/controllers:RoutesController"],
		beego.ControllerComments{
			Method: "GetRouteDetailInfoByServiceKey",
			Router: `/servicename/:servicename/version/:serviceversion`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
