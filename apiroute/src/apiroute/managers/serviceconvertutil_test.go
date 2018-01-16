package managers

import (
	"apiroute/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeEmptyToDefault(t *testing.T) {
	Convey("change empty namespace to default", t, func() {
		Convey("empty namespace must return default", func() {
			str := ChangeEmptyToDefault("")
			So(str, ShouldEqual, "default")
		})

		Convey("non-empty namespace must return itself", func() {
			str := ChangeEmptyToDefault("sss")
			So(str, ShouldEqual, "sss")
		})
	})
}

func TestCalcServiceNodes(t *testing.T) {
	Convey("convert serviceUnit instances to serviceInfo nodes", t, func() {
		Convey("instance is nil or empty,must return empty nodes", func() {
			nodes := calcServiceNodes(nil)
			So(nodes, ShouldNotBeNil)
			So(len(nodes), ShouldEqual, 0)

			instances := make([]util.InstanceUnit, 0)
			nodes = calcServiceNodes(instances)
			So(nodes, ShouldNotBeNil)
			So(len(nodes), ShouldEqual, 0)
		})
		Convey("instance is non-nil,check nodes content", func() {
			instances := make([]util.InstanceUnit, 1)
			instance := util.InstanceUnit{}
			instance.ServiceIP = "127.0.0.1"
			instance.ServicePort = "9999"
			instance.LBServerParams = "xxxxxxxx"
			instances[0] = instance
			nodes := calcServiceNodes(instances)
			So(nodes, ShouldNotBeNil)
			So(len(nodes), ShouldEqual, len(instances))
			So(nodes[0].IP, ShouldEqual, instances[0].ServiceIP)
			So(nodes[0].Port, ShouldEqual, instances[0].ServicePort)
			So(nodes[0].LBServerParams, ShouldEqual, instances[0].LBServerParams)
			So(nodes[0].TTL, ShouldEqual, -1)
		})
	})
}

func TestHandleLables(t *testing.T) {
	Convey("convert serviceUnit lables to serviceInfo lables", t, func() {
		Convey("serviceUnit-lables is nil or empty,must return empty serviceInfo lables", func() {
			labels := handleLables(nil)
			So(labels, ShouldNotBeNil)
			So(len(labels), ShouldEqual, 0)

			oldlables := make([]string, 0)
			labels = handleLables(oldlables)
			So(labels, ShouldNotBeNil)
			So(len(labels), ShouldEqual, 0)
		})
		Convey("serviceUnit-lables is non-nil,check serviceInfo-lables content", func() {
			oldlables := []string{"key:val", "", "keyval", "val1:val2:val3"}
			labels := handleLables(oldlables)
			So(labels, ShouldNotBeNil)
			So(len(labels), ShouldEqual, len(oldlables))
			//
			var (
				val  string
				find bool
			)
			val, find = labels["key"]
			So(find, ShouldBeTrue)
			So(val, ShouldEqual, "val")

			val, find = labels[""]
			So(find, ShouldBeTrue)
			So(val, ShouldEqual, "")

			val, find = labels["keyval"]
			So(find, ShouldBeTrue)
			So(val, ShouldEqual, "")

			val, find = labels["val1"]
			So(find, ShouldBeTrue)
			So(val, ShouldEqual, "val2")

			val, find = labels["xxxxx"]
			So(find, ShouldBeFalse)
			So(val, ShouldEqual, "")
		})
	})
}

/*
{
    "serviceName": "",
    "version": "",
    "url": "",
    "protocol": "",
    "visualRange": "",
    "lb_policy": "",
    "publish_port": "",
    "namespace": "",
    "network_plane_type": "",
    "host": "",
    "path": "",
    "nodes": [
        {
            "ip": "",
            "port": "",
            "lb_server_params": "",
            "checkType": "",
            "checkUrl": "",
            "checkInterval": "",
            "ttl": "",
            "checkTimeOut": "",
            "ha_role": "",
            "nodeId": "",
            "status": ""
        }
    ],
    "metadata": [
        {
            "key": "",
            "value": ""
        }
    ],
    "labels": [
        "enable_router:true"
    ],
    "swagger_url": "",
    "is_manual": false,
    "enable_ssl": false,
    "enable_tls": false,
    "enable_refer_match": "",
    "proxy_rule": {
        "http_proxy": {
            "send_timeout": "",
            "read_timeout": ""
        },
        "stream_proxy": {
            "proxy_timeout": "",
            "proxy_responses": ""
        }
    }
}
*/
func ReadServiceUnitTestResource() []*util.ServiceUnit {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testResPath := filepath.Join(workPath, "testresource/serviceunit.json")

	b, err := ioutil.ReadFile(testResPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	var serviceUnits []*util.ServiceUnit
	err = json.Unmarshal(b, &serviceUnits)
	if err != nil {
		fmt.Println(err.Error())
	}
	return serviceUnits
}

func TestAssembleServiceInfo(t *testing.T) {

	Convey("use serviceUnit to assemble ServiceInfo", t, func() {
		Convey("serviceUnit is nil,must return ServiceInfo:nil,err:non-nil", func() {
			serviceInfo, err := AssembleServiceInfo(nil)
			So(serviceInfo, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("serviceUnit is non-nil,check content", func() {

			serviceUnits := ReadServiceUnitTestResource()
			rstServiceInfos := ReadServiceTestResource()

			for _, serviceUnit := range serviceUnits {
				serviceInfo, err := AssembleServiceInfo(serviceUnit)
				//
				So(serviceInfo, ShouldNotBeNil)
				So(err, ShouldBeNil)

				//check content
				find := false
				index := 0
				for index, _ := range rstServiceInfos {
					if rstServiceInfos[index].ServicePrefix == serviceInfo.ServicePrefix &&
						rstServiceInfos[index].ServiceKey == serviceInfo.ServiceKey {
						find = true
					}
				}

				So(find, ShouldBeTrue)
				if find {
					So(rstServiceInfos[index].ServiceValue.APIVersion, ShouldEqual, serviceInfo.ServiceValue.APIVersion)
					So(rstServiceInfos[index].ServiceValue.Kind, ShouldEqual, serviceInfo.ServiceValue.Kind)
					So(rstServiceInfos[index].ServiceValue.Status, ShouldEqual, serviceInfo.ServiceValue.Status)
					So(rstServiceInfos[index].ServiceValue.MetaData.Name, ShouldEqual, serviceInfo.ServiceValue.MetaData.Name)
					So(rstServiceInfos[index].ServiceValue.MetaData.Namespace, ShouldEqual, serviceInfo.ServiceValue.MetaData.Namespace)
					//					So(len(rstServiceInfos[index].ServiceValue.MetaData.Labels), ShouldEqual, len(serviceInfo.ServiceValue.MetaData.Labels))
					So(rstServiceInfos[index].ServiceValue.Spec.Custom, ShouldEqual, serviceInfo.ServiceValue.Spec.Custom)
					So(rstServiceInfos[index].ServiceValue.Spec.EnableReferMatch, ShouldEqual, serviceInfo.ServiceValue.Spec.EnableReferMatch)
					So(rstServiceInfos[index].ServiceValue.Spec.EnableSSL, ShouldEqual, serviceInfo.ServiceValue.Spec.EnableSSL)
					So(rstServiceInfos[index].ServiceValue.Spec.EnableTLS, ShouldEqual, serviceInfo.ServiceValue.Spec.EnableTLS)
					So(rstServiceInfos[index].ServiceValue.Spec.Host, ShouldEqual, serviceInfo.ServiceValue.Spec.Host)
					So(rstServiceInfos[index].ServiceValue.Spec.LBPolicy, ShouldEqual, serviceInfo.ServiceValue.Spec.LBPolicy)
					//					So(len(rstServiceInfos[index].ServiceValue.Spec.Nodes), ShouldEqual, len(serviceInfo.ServiceValue.Spec.Nodes))
					So(rstServiceInfos[index].ServiceValue.Spec.Path, ShouldEqual, serviceInfo.ServiceValue.Spec.Path)
					So(rstServiceInfos[index].ServiceValue.Spec.Protocol, ShouldEqual, serviceInfo.ServiceValue.Spec.Protocol)
					So(rstServiceInfos[index].ServiceValue.Spec.ProxyRule.HTTPProxy.ReadTimeout, ShouldEqual, serviceInfo.ServiceValue.Spec.ProxyRule.HTTPProxy.ReadTimeout)
					So(rstServiceInfos[index].ServiceValue.Spec.ProxyRule.HTTPProxy.SendTimeout, ShouldEqual, serviceInfo.ServiceValue.Spec.ProxyRule.HTTPProxy.SendTimeout)
					So(rstServiceInfos[index].ServiceValue.Spec.ProxyRule.StreamProxy.ProxyResponse, ShouldEqual, serviceInfo.ServiceValue.Spec.ProxyRule.StreamProxy.ProxyResponse)
					So(rstServiceInfos[index].ServiceValue.Spec.ProxyRule.StreamProxy.ProxyTimeout, ShouldEqual, serviceInfo.ServiceValue.Spec.ProxyRule.StreamProxy.ProxyTimeout)
					So(rstServiceInfos[index].ServiceValue.Spec.PublishPort, ShouldEqual, serviceInfo.ServiceValue.Spec.PublishPort)
					So(rstServiceInfos[index].ServiceValue.Spec.SwaggerURL, ShouldEqual, serviceInfo.ServiceValue.Spec.SwaggerURL)
					So(rstServiceInfos[index].ServiceValue.Spec.URL, ShouldEqual, serviceInfo.ServiceValue.Spec.URL)
					So(rstServiceInfos[index].ServiceValue.Spec.VisualRange, ShouldEqual, serviceInfo.ServiceValue.Spec.VisualRange)
				}
			}
		})
	})
}

func TestCovertServiceKey(t *testing.T) {

}
