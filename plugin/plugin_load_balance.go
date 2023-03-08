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
}

// Config config
func (p *LoadBalancePlugin) Config(config *ng.PluginConfig) {
	config.Name("ng_load_balance_plugin")
	if p.ServerName != "" {
		config.Host(p.ServerName)
	}
	config.ProxyPass(p.Location, "")
}

// Interceptor interceptor
func (p *LoadBalancePlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(p.ProxyPassList))
	proxyPass := p.ProxyPassList[randIndex]
	request.Url = fmt.Sprintf("%s%s", strings.TrimSuffix(proxyPass, "/"), request.HttpRequest.RequestURI)
	fmt.Println("===>", request.Url)
	return ng.Invoke(request, response)
}
