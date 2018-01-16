package template

import (
	"apiroute/cache"
	"apiroute/logs"
	"os"
	"text/template"
)

const (
	httpFile   = "../openresty/nginx/sites-enabled/server_http.conf"
	httpsFile  = "../openresty/nginx/sites-enabled/server_https.conf"
	streamFile = "../openresty/nginx/stream-enabled/service.conf"
)

type StreamRenderUnit struct {
	Name string
	Data *cache.StreamMetaData
}

func RenderHTTP(ports []string) error {
	f, err := os.OpenFile(httpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logs.Log.Error("Failed to open %s:%v", httpFile, err)
		return err
	}

	defer f.Close()
	t := template.New("HTTP Render")
	t = template.Must(t.Parse(httpTemplate))
	t.Execute(f, ports)

	logs.Log.Info("Render HTTP file completes")
	return nil
}

func RenderHTTPS(ports []string) error {
	f, err := os.OpenFile(httpsFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logs.Log.Error("Failed to open %s:%v", httpsFile, err)
		return err
	}

	defer f.Close()
	t := template.New("HTTPS Render")
	t = template.Must(t.Parse(httpsTemplate))
	t.Execute(f, ports)

	logs.Log.Info("Render HTTPS file completes")
	return nil
}

func RenderStream(streams map[string]*cache.StreamMetaData) error {
	f, err := os.OpenFile(streamFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logs.Log.Error("Failed to open %s:%v", streamFile, err)
		return err
	}

	defer f.Close()

	renderData := []*StreamRenderUnit{}
	for k, v := range streams {
		renderData = append(renderData, &StreamRenderUnit{
			Name: k,
			Data: v,
		})
	}

	t := template.New("Stream Render")
	t = template.Must(t.Parse(streamTemplate))
	t.Execute(f, renderData)

	logs.Log.Info("Render Stream file completes")
	return nil
}
