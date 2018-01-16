package configmanager

import (
	"apiroute/logs"
	"apiroute/models"
	"apiroute/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/jinzhu/configor"
)

var (
	CfgM *cfgmanager
)

const (
	yamlSuffix string = ".yml"
)

func init() {
	if CfgM == nil {
		CfgM = newCfgManager()
		CfgM.readConfigFiles()
		//extention
		CfgM.ParseRouteWay()
		CfgM.ParseLables()
		CfgM.ParseVisualRange()
	}

	CfgM.printConfigInfo()

}

type cfgmanager struct {
	config      *models.Config
	routeWay    []string
	lables      map[string]string
	visualRange []string
}

func (cm *cfgmanager) GetConfigInfo() *models.Config {
	return cm.config
}

func (cm *cfgmanager) readConfigFiles() {
	files := cm.scanYMLFiles()

	if files == nil || len(files) == 0 {
		fmt.Println("no yml config files")
		return
	}

	err := configor.Load(cm.config, files[:]...)

	if err != nil {
		fmt.Println("read config file failed:%s", err.Error())
		return
	}

	cm.config.BasicConfig = beego.BConfig
}

func (cm *cfgmanager) scanYMLFiles() (files []string) {
	confdir := util.GetCfgFilePath()

	files = make([]string, 0, 20)

	dir, err := ioutil.ReadDir(confdir)

	if err != nil {
		fmt.Println("read dir failed:%s", err.Error())
		return nil
	}

	pthSep := string(os.PathSeparator)
	suffix := strings.ToUpper(yamlSuffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		//		fmt.Println(fi.Name())

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, confdir+pthSep+fi.Name())
		}
	}
	return files
}

func (cm *cfgmanager) printConfigInfo() {
	//
	logs.Log.Info("EnableTest:%t", cm.config.EnableTest)

	logs.Log.Info("---------listenport-------")
	logs.Log.Info(string(changeToJSONStr(cm.config.Listenport)))

	logs.Log.Info("---------apigatewaycfg-------")
	logs.Log.Info(string(changeToJSONStr(cm.config.Apigatewaycfg)))

	logs.Log.Info("---------discover-------")
	logs.Log.Info(string(changeToJSONStr(cm.config.Discover)))

	logs.Log.Info("---------apigatewaymetrics-------")
	logs.Log.Info(string(changeToJSONStr(cm.config.ApigatewayMetrics)))

	logs.Log.Info("-----------redis-----------")
	logs.Log.Info(string(changeToJSONStr(cm.config.Redis)))
}

func changeToJSONStr(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		logs.Log.Info("change json failed:%t", err.Error())
		var byteArr []byte
		return byteArr
	}
	return b
}

func newCfgManager() *cfgmanager {
	cfgm := &cfgmanager{}
	cfgm.config = &models.Config{}
	return cfgm
}

////////////////////////extension function ////////
func (cm *cfgmanager) ParseRouteWay() {
	if strings.Compare(cm.config.Apigatewaycfg.RouteWay, "") == 0 {
		cm.routeWay = make([]string, 1)
		cm.routeWay[0] = "ip"
		return
	}
	cm.routeWay = strings.Split(cm.config.Apigatewaycfg.RouteWay, "|")
	return
}

func (cm *cfgmanager) ParseLables() {
	cm.lables = make(map[string]string)

	lbs := strings.Split(cm.config.Apigatewaycfg.Lables, ",")

	for _, v := range lbs {
		kv := strings.Split(v, ":")
		if len(kv) == 2 {
			cm.lables[kv[0]] = kv[1]
		}
	}
}

func (cm *cfgmanager) ParseVisualRange() {
	if vs, ok := cm.lables["visualRange"]; ok {
		if vs == "0" || vs == "1" {
			cm.visualRange = make([]string, 1)
			cm.visualRange[0] = vs
			return
		}
	}

	//other
	cm.visualRange = make([]string, 2)
	cm.visualRange[0] = "0"
	cm.visualRange[1] = "1"
	return
}

func (cm *cfgmanager) GetRouteWay() []string {
	return cm.routeWay
}

func (cm *cfgmanager) GetVisualRange() []string {
	return cm.visualRange
}
