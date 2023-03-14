package plugin

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"math/rand"
	"strings"
	"time"
)

// LoadBalancePlugin load balance
type LoadBalancePlugin struct {
	ServerName    string
	Location      string
	ProxyPassList []string
	PolicyFunc    func(proxyPassList []string) string
}

// Config config
func (p *LoadBalancePlugin) Config(config *ng.PluginConfig) {
	config.Name("ng_load_balance_plugin")
	if p.ServerName != "" {
		config.Host(p.ServerName)
	}
	if p.PolicyFunc == nil {
		p.PolicyFunc = DefaultPolicyFunc
	}
	config.ProxyPass(p.Location, "")
}

// Interceptor interceptor
func (p *LoadBalancePlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	proxyPass := p.PolicyFunc(p.ProxyPassList)
	request.Url = fmt.Sprintf("%s%s", strings.TrimSuffix(proxyPass, "/"), request.HttpRequest.RequestURI)
	return ng.Invoke(request, response)
}

// DefaultPolicyFunc default policy
func DefaultPolicyFunc(proxyPassList []string) string {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(proxyPassList))
	proxyPass := proxyPassList[randIndex]
	return proxyPass
}
