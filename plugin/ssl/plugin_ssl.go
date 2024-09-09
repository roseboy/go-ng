package ssl

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/roseboy/go-ng/ng"
)

// PluginSSL ssl
type PluginSSL struct {
	CertFile, KeyFile string
	AutoRedirect      bool
	HttpServerPort    int
}

// Config config
func (p *PluginSSL) Config(config *ng.PluginConfig) {
	config.Name("ng_ssl_plugin")
	config.Server.WithTLS(p.CertFile, p.KeyFile)
	if p.AutoRedirect && p.HttpServerPort <= 0 {
		panic("SSLPlugin.HttpServerPort error")
	}
	if p.AutoRedirect {
		go func() {
			err := ng.NewServer(p.HttpServerPort).
				RegisterPlugins(&pluginRedirectSSL{httpsPort: config.Server.Port}).Start()
			if err != nil {
				panic(err)
			}
		}()
	}
}

// Interceptor interceptor
func (p *PluginSSL) Interceptor(ctx context.Context, request *ng.Request, response *ng.Response) error {
	return nil
}

// pluginRedirectSSL redirect
type pluginRedirectSSL struct {
	httpsPort int
}

// Config config
func (p *pluginRedirectSSL) Config(config *ng.PluginConfig) {
	config.Name("ng_ssl_plugin_redirect")
	config.Location("/", "")
}

// Interceptor interceptor
func (p *pluginRedirectSSL) Interceptor(ctx context.Context, request *ng.Request, response *ng.Response) error {
	if !strings.HasPrefix(strings.ToLower(request.HttpRequest.Proto), "https") {
		host := strings.Split(request.HttpRequest.Host, ":")[0]
		uri := request.HttpRequest.RequestURI

		response.Status = http.StatusFound
		response.Headers["Location"] = fmt.Sprintf("https://%s:%d%s", host, p.httpsPort, uri)
	}
	return nil
}
