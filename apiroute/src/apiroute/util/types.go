package util

type ServiceUnit struct {
	Name             string         `json:"serviceName,omitempty"`
	Version          string         `json:"version"`
	URL              string         `json:"url"`
	Protocol         string         `json:"protocol"`
	VisualRange      string         `json:"visualRange"`
	LBPolicy         string         `json:"lb_policy"`
	PublishPort      string         `json:"publish_port"`
	Namespace        string         `json:"namespace"`
	NWPlaneType      string         `json:"network_plane_type"`
	Host             string         `json:"host"`
	Path             string         `json:"path"`
	Instances        []InstanceUnit `json:"nodes"`
	Metadata         []MetaUnit     `json:"metadata"`
	Labels           []string       `json:"labels"`
	SwaggerURL       string         `json:"swagger_url,omitempty"`
	IsManual         bool           `json:"is_manual"`
	EnableSSL        bool           `json:"enable_ssl"`
	EnableTLS        bool           `json:"enable_tls"`
	EnableReferMatch string         `json:"enable_refer_match"`
	ProxyRule        Rules          `json:"proxy_rule,omitempty"`
}

type InstanceUnit struct {
	ServiceIP      string `json:"ip,omitempty"`
	ServicePort    string `json:"port,omitempty"`
	LBServerParams string `json:"lb_server_params,omitempty"`
	CheckType      string `json:"checkType,omitempty"`
	CheckURL       string `json:"checkUrl,omitempty"`
	CheckInterval  string `json:"checkInterval,omitempty"`
	CheckTTL       string `json:"ttl,omitempty"`
	CheckTimeOut   string `json:"checkTimeOut,omitempty"`
	HaRole         string `json:"ha_role,omitempty"`
	ServiceID      string `json:"nodeId,omitempty"`
	ServiceStatus  string `json:"status,omitempty"`
}

type MetaUnit struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Rules struct {
	HTTPProxy   HTTPProxyRule   `json:"http_proxy,omitempty"`
	StreamProxy StreamProxyRule `json:"stream_proxy,omitempty"`
}

type HTTPProxyRule struct {
	SendTimeout string `json:"send_timeout,omitempty"`
	ReadTimeout string `json:"read_timeout,omitempty"`
}

type StreamProxyRule struct {
	ProxyTimeout  string `json:"proxy_timeout,omitempty"`
	ProxyResponse string `json:"proxy_responses,omitempty"`
}

type NameAndNamespace struct {
	Name      string
	Namespace string
}

type DigestUnit struct {
	Namespace string
	NumOfNode int
}
