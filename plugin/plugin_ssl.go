package plugin

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"net/http"
	"strings"
)

// SSLPlugin ssl
type SSLPlugin struct {
	CertFile, KeyFile string
	AutoRedirect      bool
	HttpServerPort    int
}

// Config config
func (p *SSLPlugin) Config(config *ng.PluginConfig) {
	config.Name("ng_ssl_plugin")
	config.Server.CertFile = p.CertFile
	config.Server.KeyFile = p.KeyFile
	if p.AutoRedirect && p.HttpServerPort <= 0 {
		panic("SSLPlugin.HttpServerPort error")
	}
	if p.AutoRedirect {
		go ng.NewServer(p.HttpServerPort).RegisterPlugins(&redirectSSLPlugin{httpsPort: config.Server.Port}).Start()
	}
}

// Interceptor interceptor
func (p *SSLPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	return nil
}

// redirectSSLPlugin redirect
type redirectSSLPlugin struct {
	httpsPort int
}

// Config config
func (p *redirectSSLPlugin) Config(config *ng.PluginConfig) {
	config.Name("ng_ssl_plugin_redirect")
	config.ProxyPass("/", "")
}

// Interceptor interceptor
func (p *redirectSSLPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	if !strings.HasPrefix(strings.ToLower(request.HttpRequest.Proto), "https") {
		host := strings.Split(request.HttpRequest.Host, ":")[0]
		uri := request.HttpRequest.RequestURI

		response.Status = http.StatusFound
		response.Headers["Location"] = fmt.Sprintf("https://%s:%d%s", host, p.httpsPort, uri)
	}
	return nil
}
