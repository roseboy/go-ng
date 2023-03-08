package plugin

import (
	"github.com/roseboy/go-ng/ng"
)

// SSLPlugin ssl
type SSLPlugin struct {
}

// Config config
func (p *SSLPlugin) Config(config *ng.PluginConfig) {
	config.Name("ng_ssl_plugin")
	config.ProxyPass("/", "")
}

// Interceptor interceptor
func (p *SSLPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	return ng.Invoke(request, response)
}
