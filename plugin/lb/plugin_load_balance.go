package lb

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/roseboy/go-ng/ng"
)

// PluginLoadBalance load balance
type PluginLoadBalance struct {
	ServerName    string
	Location      string
	ProxyPassList []string
	PolicyFunc    func(proxyPassList []string) string
}

// Config config
func (p *PluginLoadBalance) Config(config *ng.PluginConfig) {
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
func (p *PluginLoadBalance) Interceptor(request *ng.Request, response *ng.Response) error {
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
