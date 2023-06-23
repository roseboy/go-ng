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
func NewActionPlugin(endpoint string, signatureCheck bool, authInfoFunc func(string) (uint64, string),
) *action.PluginAction {
	return &action.PluginAction{Endpoint: endpoint, SignatureCheck: signatureCheck, AuthInfoFunc: authInfoFunc}
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
