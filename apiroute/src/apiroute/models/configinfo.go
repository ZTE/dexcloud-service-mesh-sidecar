package models

import (
	"github.com/astaxie/beego"
)

//basic cfg
type Config struct {
	EnableTest        bool
	BasicConfig       *beego.Config
	Listenport        ListenPortInfo
	Apigatewaycfg     ApigatewayCfgInfo
	Discover          DiscoverInfo
	ApigatewayMetrics InternalApigatewayMetricsInfo
	Redis             RedisInfo
}

type ListenPortInfo struct {
	Httpdefaultport  int `env:"HTTP_OVERWRITE_PORT"`
	Httpsdefaultport int `env:"HTTPS_OVERWRITE_PORT"`
	Redisdefaultport int `env:"APIGATEWAY_REDIS_PORT"`
}

type ApigatewayCfgInfo struct {
	Namespace          string `env:"NAMESPACE"`
	Lables             string `env:"ROUTE_LABELS"`
	CustomFilterConfig string `env:"CUSTOM_FILTER_CONFIG"`
	RouteWay           string `env:"ROUTE_WAY"`
	RouteSubdomain     string `env:"ROUTER_SUBDOMAIN"`
	ServiceIP          string `env:"SERVICE_IP"`
	MetricsIP          string `env:"METRICS_IP"`
}

type DiscoverInfo struct {
	Enabled bool
	IP      string `env:"SDCLIENT_IP"`
	Port    int64
}

type InternalApigatewayMetricsInfo struct {
	IP   string `env:"APIGATEWAY_METRICS_IP"`
	Port int64
}

type RedisPoolInfo struct {
	MaxTotal      int
	MaxIdle       int
	MaxWaitMillis int
	TestOnBorrow  bool
	TestOnReturn  bool
}
type RedisInfo struct {
	Host              string
	Port              int `env:"APIGATEWAY_REDIS_PORT"`
	ConnectionTimeout int
	DBIndexRoute      int
	DBIndexService    int
	Pool              RedisPoolInfo
}

//logger config
type Logger struct {
	Console ConsoleOutput
	File    FileOutput
}

type ConsoleOutput struct {
	Level string
}

type FileOutput struct {
	Filename string
	Level    string
	Maxlines int
	Maxsize  int
	Daily    bool
	Maxdays  int64
	Rotate   bool
}
