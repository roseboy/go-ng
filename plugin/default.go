package plugin

import (
	"github.com/roseboy/go-ng/plugin/action"
	"github.com/roseboy/go-ng/plugin/lb"
	"github.com/roseboy/go-ng/plugin/ssl"
	"github.com/roseboy/go-ng/plugin/www"
)

// NewStaticFilePlugin www
func NewStaticFilePlugin(webRoot string) *www.PluginStatic {
	return &www.PluginStatic{WebRoot: webRoot}
}

// NewActionPlugin action
func NewActionPlugin(endpoint string) *action.PluginAction {
	actionPlg := &action.PluginAction{
		Endpoint: endpoint,
	}
	action.RegisterInitAction(actionPlg)
	return actionPlg
}

// NewActionParamsPlugin action
func NewActionParamsPlugin(endpoint string) *action.PluginActionParams {
	return &action.PluginActionParams{
		Endpoint: endpoint,
	}
}

// NewLoadBalancePlugin bl
func NewLoadBalancePlugin(serverName string, location string, proxyPassList []string) *lb.PluginLoadBalance {
	return &lb.PluginLoadBalance{ServerName: serverName, Location: location, ProxyPassList: proxyPassList}
}

// NewLoadBalancePluginWithPolicy bl
func NewLoadBalancePluginWithPolicy(serverName string, location string, proxyPassList []string,
	policyFunc func(proxyPassList []string) string) *lb.PluginLoadBalance {
	return &lb.PluginLoadBalance{
		ServerName:    serverName,
		Location:      location,
		ProxyPassList: proxyPassList,
		PolicyFunc:    policyFunc,
	}
}

// NewSSLPlugin ssl
func NewSSLPlugin(certFile, keyFile string) *ssl.PluginSSL {
	return &ssl.PluginSSL{CertFile: certFile, KeyFile: keyFile}
}

// NewSSLPluginWithAutoRedirect ssl
func NewSSLPluginWithAutoRedirect(certFile, keyFile string, httpServerPort int) *ssl.PluginSSL {
	return &ssl.PluginSSL{CertFile: certFile, KeyFile: keyFile, HttpServerPort: httpServerPort, AutoRedirect: true}
}
