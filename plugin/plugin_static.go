package plugin

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"net/http"
)

// StaticPlugin static file server
type StaticPlugin struct {
	WebRoot string
}

// Config config
func (p *StaticPlugin) Config(config *ng.PluginConfig) {
	if p.WebRoot == "" {
		p.WebRoot = "www"
	}
	config.Name("ng_file_server_plugin")
	config.ProxyPass("/", "")
}

// Interceptor interceptor
func (p *StaticPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	fileName := fmt.Sprintf("%s/%s", p.WebRoot, request.HttpRequest.URL.Path)
	http.ServeFile(response.ResponseWriter, request.HttpRequest, fileName)
	return nil
}
