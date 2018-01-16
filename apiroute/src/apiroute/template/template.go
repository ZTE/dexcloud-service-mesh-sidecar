package template

const (
	httpTemplate   = "{{range .}}\nserver {\n\tlisten {{.}};\n\tinclude ../msb-enabled/location-default/msblocations.conf;\n}\n {{end}}"
	httpsTemplate  = "{{range .}}\nserver {\n\tlisten {{.}} ssl;\n\tssl_certificate ../ssl/cert/cert.crt;\n\tssl_certificate_key ../ssl/cert/cert.key;\n\tssl_protocols TLSv1.2;\n\tssl_dhparam ../ssl/dh-pubkey/dhparams.pem;\n\tinclude ../msb-enabled/location-default/msblocations.conf;\n}\n {{end}}"
	streamTemplate = "{{range .}}\nserver {\n\tlisten {{.Data.PublishPort}}{{if .Data.EnableTLS}} ssl{{end}}{{if .Data.IsUDP}} udp{{end}};{{range .Data.ProxyRules}}\n\t{{.Name}} {{.Value}};{{end}}\n\tproxy_pass {{.Name}};{{if .Data.EnableTLS}}\n\tssl_certificate ../ssl/cert/cert.crt;\n\tssl_certificate_key ../ssl/cert/cert.key;\n\tssl_protocols TLSv1.2;\n\tssl_dhparam ../ssl/dh-pubkey/dhparams.pem;{{end}}\n}\nupstream {{.Name}} {\n\t{{if .Data.LBPolicy}}{{.Data.LBPolicy}};{{end}}{{range .Data.ServerList}}\n\tserver {{.Addr}}{{if .Params}} {{.Params}}{{end}};{{end}}\n}{{end}}"
)
