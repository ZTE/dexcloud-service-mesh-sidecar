package managers

import (
	"apiroute/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ReadServiceTestResource() []*models.ServiceInfo {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testResPath := filepath.Join(workPath, "testresource/service.json")

	b, err := ioutil.ReadFile(testResPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	var serviceInfos []*models.ServiceInfo
	err = json.Unmarshal(b, &serviceInfos)
	if err != nil {
		fmt.Println(err.Error())
	}
	return serviceInfos
}

func ReadRouteTestResource() []*models.RouteInfo {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testResPath := filepath.Join(workPath, "testresource/route.json")

	b, err := ioutil.ReadFile(testResPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	var routeInfos []*models.RouteInfo
	err = json.Unmarshal(b, &routeInfos)
	if err != nil {
		fmt.Println(err.Error())
	}
	return routeInfos
}

func TestAssembleRouteInfo(t *testing.T) {
	serviceInfos := ReadServiceTestResource()
	rstRouteInfos := ReadRouteTestResource()
	for _, serviceInfo := range serviceInfos {
		routeInfos, _, _ := AssembleRouteInfo(serviceInfo)
		/*
			for _, routeInfo := range routeInfos {
				if b, err := json.Marshal(routeInfo); err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println(string(b))
				}
			}

			for _, relation := range relations {
				if b, err := json.Marshal(relation); err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println(string(b))
				}
			}
		*/
		for _, routeInfo := range routeInfos {
			find := false
			for _, rstRouteInfo := range rstRouteInfos {
				if routeInfo.RoutePrefix == rstRouteInfo.RoutePrefix &&
					routeInfo.RouteKey == rstRouteInfo.RouteKey &&
					routeInfo.RouteValue.APIVersion == rstRouteInfo.RouteValue.APIVersion &&
					routeInfo.RouteValue.Kind == rstRouteInfo.RouteValue.Kind &&
					routeInfo.RouteValue.Status == rstRouteInfo.RouteValue.Status &&
					routeInfo.RouteValue.MetaData.Name == rstRouteInfo.RouteValue.MetaData.Name &&
					routeInfo.RouteValue.MetaData.Namespace == rstRouteInfo.RouteValue.MetaData.Namespace &&
					len(routeInfo.RouteValue.MetaData.Labels) == len(rstRouteInfo.RouteValue.MetaData.Labels) &&
					routeInfo.RouteValue.Spec.Apijson == rstRouteInfo.RouteValue.Spec.Apijson &&
					routeInfo.RouteValue.Spec.Apijsontype == rstRouteInfo.RouteValue.Spec.Apijsontype &&
					routeInfo.RouteValue.Spec.ConnectTimeOut == rstRouteInfo.RouteValue.Spec.ConnectTimeOut &&
					routeInfo.RouteValue.Spec.ConsulServiceName == rstRouteInfo.RouteValue.Spec.ConsulServiceName &&
					routeInfo.RouteValue.Spec.Control == rstRouteInfo.RouteValue.Spec.Control &&
					routeInfo.RouteValue.Spec.EnableReferMatch == rstRouteInfo.RouteValue.Spec.EnableReferMatch &&
					routeInfo.RouteValue.Spec.Enablessl == rstRouteInfo.RouteValue.Spec.Enablessl &&
					routeInfo.RouteValue.Spec.Host == rstRouteInfo.RouteValue.Spec.Host &&
					routeInfo.RouteValue.Spec.LBPolicy == rstRouteInfo.RouteValue.Spec.LBPolicy &&
					routeInfo.RouteValue.Spec.MetricsURL == rstRouteInfo.RouteValue.Spec.MetricsURL &&
					len(routeInfo.RouteValue.Spec.Nodes) == len(rstRouteInfo.RouteValue.Spec.Nodes) &&
					routeInfo.RouteValue.Spec.PublishPort == rstRouteInfo.RouteValue.Spec.PublishPort &&
					routeInfo.RouteValue.Spec.PublishProtocol == rstRouteInfo.RouteValue.Spec.PublishProtocol &&
					routeInfo.RouteValue.Spec.ReadTimeout == rstRouteInfo.RouteValue.Spec.ReadTimeout &&
					routeInfo.RouteValue.Spec.Scenario == rstRouteInfo.RouteValue.Spec.Scenario &&
					routeInfo.RouteValue.Spec.SendTimeout == rstRouteInfo.RouteValue.Spec.SendTimeout &&
					routeInfo.RouteValue.Spec.URL == rstRouteInfo.RouteValue.Spec.URL &&
					routeInfo.RouteValue.Spec.UseOwnUpstream == rstRouteInfo.RouteValue.Spec.UseOwnUpstream &&
					routeInfo.RouteValue.Spec.VisualRange == rstRouteInfo.RouteValue.Spec.VisualRange {
					find = true
				}
			}

			if find {
				t.Log("covert success")
			} else {
				if b, err := json.Marshal(routeInfo); err != nil {
					t.Error("covert failed:" + err.Error())
				} else {
					t.Error("covert failed:" + string(b))
				}

			}
		}
	}
}

func Test_calcRealPorts(t *testing.T) {
	var publishPort string = ""
	ports := calcRealPorts(publishPort)
	if len(ports) == 0 {
		t.Log("empty sting,len 0")
	}

	publishPort = "28001|"
	ports = calcRealPorts(publishPort)
	if len(ports) == 1 {
		t.Log("28001| sting,len 1")
	}

	publishPort = "28001|28002"
	ports = calcRealPorts(publishPort)
	if len(ports) == 2 {
		t.Log("28001|28002 sting,len 2")
	}

	publishPort = "28001|28002|28003|28004"
	ports = calcRealPorts(publishPort)
	if len(ports) == 2 {
		t.Log("28001|28002|28003|28004 sting,len 2")
	}
}

func TestParsePublishURL(t *testing.T) {
	Convey("ParsePublishURL", t, func() {
		Convey("PublishURL is empty:", func() {
			routeType, routeName, routeVersion := parsePublishURL("")
			So(routeType, ShouldEqual, "")
			So(routeName, ShouldEqual, "")
			So(routeVersion, ShouldEqual, "")
		})

		Convey("PublishURL has prefix api:", func() {
			{
				routeType, routeName, routeVersion := parsePublishURL("/api/huangleibo/v1")
				So(routeType, ShouldEqual, "api")
				So(routeName, ShouldEqual, "huangleibo")
				So(routeVersion, ShouldEqual, "v1")
			}
			{
				routeType, routeName, routeVersion := parsePublishURL("/api/v1/huangleibo")
				So(routeType, ShouldEqual, "api")
				So(routeName, ShouldEqual, "huangleibo")
				So(routeVersion, ShouldEqual, "v1")
			}
			{
				routeType, routeName, routeVersion := parsePublishURL("/api/huangleibo")
				So(routeType, ShouldEqual, "api")
				So(routeName, ShouldEqual, "huangleibo")
				So(routeVersion, ShouldEqual, "")
			}
			{
				routeType, routeName, routeVersion := parsePublishURL("/api/v1.0/v2.0")
				So(routeType, ShouldEqual, "api")
				So(routeName, ShouldEqual, "v1.0")
				So(routeVersion, ShouldEqual, "v2.0")
			}
			{
				routeType, routeName, routeVersion := parsePublishURL("/api/v1.0")
				So(routeType, ShouldEqual, "custom")
				So(routeName, ShouldEqual, "/api/v1.0")
				So(routeVersion, ShouldEqual, "")
			}
			{
				routeType, routeName, routeVersion := parsePublishURL("/api/s1/v1/xxx")
				So(routeType, ShouldEqual, "custom")
				So(routeName, ShouldEqual, "/api/s1/v1/xxx")
				So(routeVersion, ShouldEqual, "")
			}
			{
				routeType, routeName, routeVersion := parsePublishURL("/api/s1/s1")
				So(routeType, ShouldEqual, "custom")
				So(routeName, ShouldEqual, "/api/s1/s1")
				So(routeVersion, ShouldEqual, "")
			}
		})

		Convey("PublishURL is iui:", func() {
			{
				routeType, routeName, routeVersion := parsePublishURL("/iui/huangleibo")
				So(routeType, ShouldEqual, "iui")
				So(routeName, ShouldEqual, "huangleibo")
				So(routeVersion, ShouldEqual, "")
			}
			{
				routeType, routeName, routeVersion := parsePublishURL("/iui/huangleibo/v1")
				So(routeType, ShouldEqual, "custom")
				So(routeName, ShouldEqual, "/iui/huangleibo/v1")
				So(routeVersion, ShouldEqual, "")
			}
		})

		Convey("PublishURL is custom:", func() {
			{
				routeType, routeName, routeVersion := parsePublishURL("/apihuangleibo")
				So(routeType, ShouldEqual, CustomRouteType)
				So(routeName, ShouldEqual, "/apihuangleibo")
				So(routeVersion, ShouldEqual, "")
			}
			{
				routeType, routeName, routeVersion := parsePublishURL("/iuihuangleibo")
				So(routeType, ShouldEqual, CustomRouteType)
				So(routeName, ShouldEqual, "/iuihuangleibo")
				So(routeVersion, ShouldEqual, "")
			}
		})
	})
}
