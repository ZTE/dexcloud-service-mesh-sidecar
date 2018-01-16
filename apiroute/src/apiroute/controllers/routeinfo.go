package controllers

import (
	"apiroute/handlers"
	"apiroute/logs"
	"apiroute/managers"
	"apiroute/managers/configmanager"
	"apiroute/models"
	_ "apiroute/models" //justifying it
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
)

type RoutesController struct {
	beego.Controller
}

// @Title get route detail info
// @Description get route detail info by router key
// @Param publishport query int false "publish port."
// @Param routetype query string true "api|iui|custom"
// @Param routename query string true "route name"
// @Param routeversion query string false "route version"
// @Success 200 {object} 	models.RouteDetailInfo
// @Failure 403 require fields are nil
// @router /internal [get]
func (s *RoutesController) GetRouteDetailInfoByRouteKey() {
	var (
		publishPort  string
		integerPort  int
		routeType    string
		routeName    string
		routeVersion string
		routeWay     string
		route        *models.RouteDetailInfo
		err          error
	)
	//get params
	publishPort = s.Ctx.Input.Query("publishport")
	routeType = s.Ctx.Input.Query("routetype")
	routeName = s.Ctx.Input.Query("routename")
	routeVersion = s.Ctx.Input.Query("routeversion")
	routeWay = managers.RouteWayIP

	logs.Log.Info("GetRouteDetailInfoByRouteKey %s-%s-%s-%s", publishPort, routeType, routeName, routeVersion)

	//check
	if routeType == "" || routeName == "" {
		logs.Log.Warn("param err:routeType:%s or routeName:%s  is empty", routeType, routeName)
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Ctx.WriteString("param err: there is some empty")
		s.Abort("")
	}

	//handle default value:publishport
	publishPort, integerPort, err = s.handlePublishPort(publishPort)
	if err != nil {
		logs.Log.Info("publishPort is not a int.%s", err.Error())
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Ctx.WriteString("param err: publishPort is not a arabic numerals")
		s.Abort("")
	}

	if routeVersion == "null" {
		routeVersion = ""
	}

	route, err = handlers.GetRouteDetailInfoByRouteKey(publishPort, routeWay, routeType,
		routeName, routeVersion)

	if err != nil {
		logs.Log.Warn("query route info error:%s", err.Error())
		s.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		s.Ctx.ResponseWriter.Write([]byte(err.Error()))
		s.Abort("")
	}

	if route == nil {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		s.Ctx.ResponseWriter.Write([]byte("{}"))
		s.Abort("")
	}

	//specail handle:publishPort==""|"routing"
	s.handleRouteInfo(integerPort, route)

	s.Data["json"] = route
	s.EnableRender = false
	s.ServeJSON()
}

func (s *RoutesController) handlePublishPort(publishPort string) (string, int, error) {
	var (
		httpdefaultport  = configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport
		httpsdefaultport = configmanager.CfgM.GetConfigInfo().Listenport.Httpsdefaultport
	)

	if publishPort == "" { //default
		return managers.DefaultPublishPort, httpdefaultport, nil
	}

	port, err := strconv.Atoi(publishPort)
	if err != nil {
		return publishPort, 0, err
	}

	if port == httpdefaultport ||
		port == httpsdefaultport {
		return managers.DefaultPublishPort, port, nil
	}

	return publishPort, port, nil
}

func (s *RoutesController) handleRouteInfo(port int, route *models.RouteDetailInfo) {
	var (
		httpdefaultport  = configmanager.CfgM.GetConfigInfo().Listenport.Httpdefaultport
		httpsdefaultport = configmanager.CfgM.GetConfigInfo().Listenport.Httpsdefaultport
	)

	if port == httpdefaultport {
		route.Spec.PublishPort = strconv.Itoa(port)
		route.Spec.PublishProtocol = managers.HTTPPublishProtocol
	} else if port == httpsdefaultport {
		route.Spec.PublishPort = strconv.Itoa(port)
		route.Spec.PublishProtocol = managers.HTTPSPublishProtocol
	}
}

// @Title get route detail info
// @Description get route detail info by servicekey
// @Param servicename path string true "servicename"
// @Param serviceversion path string true "serviceversion"
// @Param namespace query string false "namespace"
// @Success 200 {object} models.RouteDetailInfo
// @Failure 403 require fields are nil
// @router /servicename/:servicename/version/:serviceversion [get]
func (s *RoutesController) GetRouteDetailInfoByServiceKey() {

	var (
		serviceName    string
		serviceVersion string
		namespace      string
		routes         string
		err            error
	)

	//get params
	serviceName = s.Ctx.Input.Param(":servicename")
	serviceVersion = s.Ctx.Input.Param(":serviceversion")
	namespace = s.Ctx.Input.Query("namespace")

	logs.Log.Info("GetRouteDetailInfoByServiceKey:%s-%s-%s", namespace, serviceName, serviceVersion)

	if serviceName == "" {
		logs.Log.Warn("param err:serviceName is empty")
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Ctx.WriteString("param err:serviceName is empty")
		s.Abort("")
	}

	if namespace == "" {
		namespace = "default"
	}

	if serviceVersion == "null" {
		serviceVersion = ""
	}

	routes, err = handlers.GetRouteDetailInfoByServiceKey(namespace, serviceName, serviceVersion)
	if err != nil {
		logs.Log.Warn("query routes by service key[%s-%s-%s] error:%s", namespace, serviceName,
			serviceVersion, err.Error())
		s.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		s.Ctx.ResponseWriter.Write([]byte(err.Error()))
		s.Abort("")
	}

	if routes == "" {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
		s.Ctx.ResponseWriter.Write([]byte("[]"))
		s.Abort("")
	}

	s.Ctx.WriteString(routes)
}
