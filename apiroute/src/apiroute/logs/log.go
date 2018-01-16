package logs

import (
	"apiroute/models"
	"apiroute/util"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/configor"
)

var (
	Log         *logs.BeeLogger
	LoggerLevel = map[string]int{"Emergency": 0, "Alert": 1, "Critical": 2, "Error": 3,
		"Warn": 4, "Notice": 5, "Info": 6, "Debug": 7}
)

const (
	logFileName string = "log.yml"
)

func init() {
	Log = logs.NewLogger()
	Log.EnableFuncCallDepth(true)

	loggerCfg := readLoggerCfg()

	if loggerCfg != nil {
		if !setCustomConsoleLogger(loggerCfg) {
			setDefaultConsoleLogger()
		}
		if !setCustomFileLogger(loggerCfg) {
			setDefaultFileLogger()
		}
		Log.Info("---------custom logger conf file-------")
		printLoggerCfg(loggerCfg)
	} else {
		setDefaultConsoleLogger()
		setDefaultFileLogger()
	}
}

func readLoggerCfg() *models.Logger {
	loggerCfg := &models.Logger{}

	confdir := util.GetCfgFilePath()
	pthSep := string(os.PathSeparator)
	filPath := confdir + pthSep + logFileName

	err := configor.Load(loggerCfg, filPath)

	if err != nil {
		fmt.Printf("read config file failed:%s", err.Error())
		return nil
	}

	return loggerCfg
}

func setCustomConsoleLogger(lc *models.Logger) bool {
	//console
	consolecfg := &ConsoleCfg{}
	consolecfg.Level = LoggerLevel[lc.Console.Level]
	if !setConsoleLogger(consolecfg) {
		return false
	}

	return true
}

func setDefaultConsoleLogger() {
	consolecfg := &ConsoleCfg{}
	consolecfg.Level = LoggerLevel["Warn"]
	setConsoleLogger(consolecfg)
}

func setCustomFileLogger(lc *models.Logger) bool {

	filecfg := &FileCfg{}
	filecfg.Filename = lc.File.Filename
	checkAndCreateLogDir(lc.File.Filename)
	filecfg.Level = LoggerLevel[lc.File.Level]
	filecfg.MaxLines = lc.File.Maxlines
	filecfg.MaxSize = lc.File.Maxsize * 1024 * 1024
	filecfg.Daily = lc.File.Daily
	filecfg.MaxDays = lc.File.Maxdays
	filecfg.Rotate = lc.File.Rotate

	if !setFileLogger(filecfg) {
		return false
	}

	return true
}

var checkAndCreateLogDir = func(fileName string) {
	if fileName == "" {
		return
	}

	// no,/, ./ ,../
	var index int
	if index = strings.LastIndex(fileName, "/"); index <= 2 {
		return
	}
	perm, _ := strconv.ParseInt("0660", 8, 64)

	if mkerr := os.MkdirAll(fileName[0:index], os.FileMode(perm)); mkerr != nil {
		return
	}
}

func setDefaultFileLogger() {
	filecfg := &FileCfg{}
	filecfg.Filename = "apiroute.log"
	filecfg.Level = LoggerLevel["Info"]
	filecfg.MaxLines = 100000
	filecfg.MaxSize = 30 * 1024 * 1024
	filecfg.Daily = true
	filecfg.MaxDays = 10
	filecfg.Rotate = true
	setFileLogger(filecfg)
}

func printLoggerCfg(lc *models.Logger) {
	Log.Info("---------console-------")
	Log.Info("level:%s", lc.Console.Level)

	Log.Info("---------file-------")
	Log.Info("filename:%s", lc.File.Filename)
	Log.Info("level:%s", lc.File.Level)
	Log.Info("maxlines:%d", lc.File.Maxlines)
	Log.Info("maxsize:%d", lc.File.Maxsize)
	Log.Info("daily:%t", lc.File.Daily)
	Log.Info("maxdays:%d", lc.File.Maxdays)
	Log.Info("rotate:%t", lc.File.Rotate)
}

type ConsoleCfg struct {
	Level int `json:"level"`
}

type FileCfg struct {
	Filename string `json:"filename"`
	Level    int    `json:"level"`
	MaxLines int    `json:"maxlines"`
	MaxSize  int    `json:"maxsize"`
	Daily    bool   `json:"daily"`
	MaxDays  int64  `json:"maxdays"`
	Rotate   bool   `json:"rotate"`
}

//set console
//Level    int  `json:"level"`
//Colorful bool `json:"color"`
func setConsoleLogger(cc *ConsoleCfg) bool {
	byteconfig, err := json.Marshal(cc)
	if err != nil {
		fmt.Printf("set console logger,change to json failed:%s", err.Error())
		return false
	}
	//	fmt.Println(string(byteconfig))
	seterr := Log.SetLogger(logs.AdapterConsole, string(byteconfig))
	if seterr != nil {
		fmt.Printf("set console logger failed:", seterr.Error())
		return false
	}
	return true
}

//set file
// config need to be correct JSON as string: {"interval":360}.
// It writes messages by lines limit, file size limit, or time frequency.
//(1-999) filename.date.num.log
//Filename   string `json:"filename"`
//MaxLines         int `json:"maxlines"` 1000000
//MaxSize        int `json:"maxsize"` 1 << 28 256M  length in bytes
//Daily         bool  `json:"daily"`
//MaxDays       int64 `json:"maxdays"` 7
//Rotate bool `json:"rotate"`
//Level int `json:"level"`
//Perm string `json:"perm"`
func setFileLogger(fileCfg *FileCfg) bool {
	byteconfig, err := json.Marshal(fileCfg)
	if err != nil {
		fmt.Printf("set file logger,change to json failed:%s", err.Error())
		return false
	}

	seterr := Log.SetLogger(logs.AdapterFile, string(byteconfig))
	if seterr != nil {
		fmt.Printf("set file logger failed:", seterr.Error())
		return false
	}
	return true
}
