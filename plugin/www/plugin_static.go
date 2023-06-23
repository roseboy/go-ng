package www

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"net/http"
)

// PluginStatic static file server
type PluginStatic struct {
	WebRoot string
}

// Config config
func (p *PluginStatic) Config(config *ng.PluginConfig) {
	if p.WebRoot == "" {
		p.WebRoot = "www"
	}
	config.Name("ng_file_server_plugin")
	config.ProxyPass("/", "")
}

// Interceptor interceptor
func (p *PluginStatic) Interceptor(request *ng.Request, response *ng.Response) error {
	fileName := fmt.Sprintf("%s/%s", p.WebRoot, request.HttpRequest.URL.Path)
	http.ServeFile(response.ResponseWriter, request.HttpRequest, fileName)
	return nil
}
