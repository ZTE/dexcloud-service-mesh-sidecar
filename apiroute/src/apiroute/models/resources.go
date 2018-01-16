package models

type HTTPHTTPSDefaultPorts struct {
	HTTPPort  int `json:"httpPort"`
	HTTPSPort int `json:"httpsPort"`
}

type RouteAbstractInfo struct {
	Namespace       string `json:"namespace"`
	ServiceName     string `json:"serviceName"`
	ServiceVersion  string `json:"serviceVersion"`
	PublishProtocol string `json:"publishProtocol"`
	PublishPort     int    `json:"publishPort"`
	RouteType       string `json:"routeType"`
	RouterName      string `json:"routeName"`
	RouterVersion   string `json:"routeVersion"`
	PublishURL      string `json:"publishUrl"`
}

type TCPUDPRouteSpecObject struct {
	VisualRange     string            `json:"visualRange"`
	PublishPort     string            `json:"publish_port"`
	PublishProtocol string            `json:"publish_protocol"`
	Host            string            `json:"host"`
	EnableTLS       bool              `json:"enable_tls"`
	LBPolicy        string            `json:"lb_policy"`
	ConnectTimeOut  string            `json:"connect_timeout"`
	ProxyTimeout    string            `json:"proxy_timeout"`
	ProxyResponse   string            `json:"proxy_responses"`
	Nodes           []RouteNodeObject `json:"nodes"`
}

type TCPUDPRouteMetaDataObjet struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Namespace       string `json:"namespace"`
	ServiceName     string `json:"serviceName"`
	ServiceVersion  string `json:"serviceVersion"`
	UpdateTimestamp string `json:"updateTimestamp"`
}

type TCPUDPRouteDetailInfo struct {
	Kind       string                   `json:"kind"`
	APIVersion string                   `json:"apiVersion"`
	Status     string                   `json:"status"`
	MetaData   TCPUDPRouteMetaDataObjet `json:"metadata"`
	Spec       TCPUDPRouteSpecObject    `json:"spec"`
}
