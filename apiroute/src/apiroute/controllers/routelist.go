package controllers

import (
	"apiroute/handlers"
	"apiroute/logs"
	"apiroute/models"
	"net/http"

	"github.com/astaxie/beego"
)

type RouteListController struct {
	beego.Controller
}

// @Title get route abstract info list
// @Description get {namespace}'s route abstract info list
// @Param namespace query string false "one namespace"
// @Success 200 {object} models.RouteAbstractInfo
// @Failure 500 read redis failed
// @router / [get]
func (s *RouteListController) GetRouteAbstractInfoList() {
	var (
		namespace    string
		routeWay     string
		abstractList []*models.RouteAbstractInfo
		err          error
	)
	//get params
	namespace = s.Ctx.Input.Query("namespace")
	routeWay = s.Ctx.Input.Query("routeway")
	logs.Log.Info("GetRouteAbstractInfoList:namespace:%s,routeway:%s", namespace, routeWay)

	//query abstruct info list
	abstractList, err = handlers.GetRouteAbstractInfoListByNamespace(namespace, routeWay)

	if err != nil {
		logs.Log.Warn("query route abstruct list error:%s", err.Error())
		s.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		s.Ctx.ResponseWriter.Write([]byte(err.Error()))
		s.Abort("")
	}

	if abstractList == nil || len(abstractList) == 0 {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
		s.Ctx.ResponseWriter.Write([]byte("[]"))
		s.Abort("")
	}

	s.Data["json"] = abstractList
	s.EnableRender = false
	s.ServeJSON()

}
