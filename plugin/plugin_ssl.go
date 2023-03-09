package plugin

import (
	"github.com/roseboy/go-ng/ng"
)

// SSLPlugin ssl
type SSLPlugin struct {
	CertFile, KeyFile string
}

// Config config
func (p *SSLPlugin) Config(config *ng.PluginConfig) {
	config.Name("ng_ssl_plugin")
	config.Server.CertFile = p.CertFile
	config.Server.KeyFile = p.KeyFile
}

// Interceptor interceptor
func (p *SSLPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	return nil
}
