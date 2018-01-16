package models

import (
	"apiroute/util"
)

type ServiceInfo struct {
	ServicePrefix string
	ServiceKey    ServiceKey
	ServiceValue  ServiceDetailInfo
}

type ServiceNodeObject struct {
	IP             string `json:"ip"`
	Port           string `json:"port"`
	TTL            int    `json:"ttl"`
	LBServerParams string `json:"lb_server_params,omitempty"`
}

type ServiceSpecObject struct {
	VisualRange      string              `json:"visualRange"`
	URL              string              `json:"url"`
	Path             string              `json:"path"`
	PublishPort      string              `json:"publish_port"`
	Host             string              `json:"host"`
	Protocol         string              `json:"protocol"`
	Custom           string              `json:"custom"`
	LBPolicy         string              `json:"lb_policy"`
	EnableSSL        bool                `json:"enable_ssl"`
	EnableTLS        bool                `json:"enable_tls"`
	SwaggerURL       string              `json:"swagger_url,omitempty"`
	EnableReferMatch string              `json:"enable_refer_match"`
	ProxyRule        util.Rules          `json:"proxy_rule,omitempty"`
	Nodes            []ServiceNodeObject `json:"nodes"`
}

type ServiceMetaDataObjet struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Namespace       string            `json:"namespace"`
	UID             string            `json:"uid"`
	UpdateTimestamp string            `json:"updateTimestamp"`
	Labels          map[string]string `json:"labels"`
	Annotations     []string          `json:"annotations"`
}

type ServiceDetailInfo struct {
	Kind       string               `json:"kind"`
	APIVersion string               `json:"apiVersion"`
	Status     string               `json:"status"`
	MetaData   ServiceMetaDataObjet `json:"metadata"`
	Spec       ServiceSpecObject    `json:"spec"`
}

type ServiceKey struct {
	Namespace      string
	ServiceName    string
	ServiceVersion string
}

type PublishInfo struct {
	Protocol    string
	PublishPort string
	NumOfNode   int
}
