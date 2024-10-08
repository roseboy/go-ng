package www

import (
	"context"
	"fmt"
	"net/http"

	"github.com/roseboy/go-ng/ng"
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
	config.Location("/", "")
}

// Interceptor interceptor
func (p *PluginStatic) Interceptor(ctx context.Context, request *ng.Request, response *ng.Response) error {
	fileName := fmt.Sprintf("%s/%s", p.WebRoot, request.HttpRequest.URL.Path)
	http.ServeFile(response.ResponseWriter, request.HttpRequest, fileName)
	return nil
}
