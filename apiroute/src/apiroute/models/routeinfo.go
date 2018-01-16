package models

type RouteInfo struct {
	RoutePrefix string
	RouteKey    RouteKey
	RouteValue  RouteDetailInfo
}

type RouteNodeObject struct {
	IP             string `json:"ip"`
	Port           int    `json:"port"`
	Weight         int    `json:"weight"`
	LBServerParams string `json:"lb_server_params,omitempty"`
}

type RouteSpecObject struct {
	VisualRange       string            `json:"visualRange"`
	URL               string            `json:"url"`
	PublishPort       string            `json:"publish_port"`
	PublishProtocol   string            `json:"publish_protocol"`
	Host              string            `json:"host"`
	Apijson           string            `json:"apijson"`
	Apijsontype       string            `json:"apijsontype"`
	MetricsURL        string            `json:"metricsUrl"`
	LBPolicy          string            `json:"lb_policy"`
	ConsulServiceName string            `json:"consulServiceName"`
	UseOwnUpstream    string            `json:"useOwnUpstream"`
	Enablessl         bool              `json:"enable_ssl"`
	Control           string            `json:"control"`
	Scenario          int               `json:"scenario"`
	EnableReferMatch  string            `json:"enable_refer_match"`
	ConnectTimeOut    string            `json:"connect_timeout"`
	SendTimeout       string            `json:"send_timeout"`
	ReadTimeout       string            `json:"read_timeout"`
	Nodes             []RouteNodeObject `json:"nodes"`
}

type RouteMetaDataObjet struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Namespace       string            `json:"namespace"`
	ServiceName     string            `json:"serviceName"`
	ServiceVersion  string            `json:"serviceVersion"`
	UID             string            `json:"uid"`
	UpdateTimestamp string            `json:"updateTimestamp"`
	Labels          map[string]string `json:"labels"`
	Annotations     []string          `json:"annotations"`
}

type RouteDetailInfo struct {
	Kind       string             `json:"kind"`
	APIVersion string             `json:"apiVersion"`
	Status     string             `json:"status"`
	MetaData   RouteMetaDataObjet `json:"metadata"`
	Spec       RouteSpecObject    `json:"spec"`
}

type RouteKey struct {
	PublishPort  string //"host|publishport|routing"
	RouteType    string
	RouteName    string
	RouteVersion string
}
