package ng

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/roseboy/go-ng/util"
)

type plugin struct {
	name                string
	hosts               []string
	locationProxyPasses map[string]string
	locationRegexps     map[string]*regexp.Regexp
	locationWeights     map[string]int
	interceptor         func(ctx context.Context, request *Request, response *Response) error
}

type pluginWrapper struct {
	location  string
	proxyPass string
	plugin    *plugin
}

// pluginInterface plugin interface
type pluginInterface interface {
	Interceptor(ctx context.Context, request *Request, response *Response) error
	Config(config *PluginConfig)
}

// PluginConfig plugin config
type PluginConfig struct {
	Server    *server
	name      string
	hosts     []string
	locations map[string]string
}

// Name plugin name
func (c *PluginConfig) Name(name string) {
	c.name = name
}

// Location  add location proxy pass
func (c *PluginConfig) Location(location, proxyPass string) {
	if !strings.HasPrefix(location, "/") {
		panic("location must start with /")
	}
	c.locations[location] = proxyPass
}

// Host add intercept host
func (c *PluginConfig) Host(hosts ...string) {
	c.hosts = append(c.hosts, hosts...)
}

// RegisterPlugins 注册插件
func (s *server) RegisterPlugins(plugins ...pluginInterface) *server {
	for _, pg := range plugins {
		config := &PluginConfig{
			locations: map[string]string{},
			Server:    s,
		}
		pg.Config(config)

		if config.name == "" {
			panic("plugin name can not be empty")
		}
		log.Printf("register plugin: %s", config.name)

		plg := &plugin{
			name:                config.name,
			hosts:               config.hosts,
			locationProxyPasses: config.locations,
			locationRegexps:     map[string]*regexp.Regexp{},
			locationWeights:     map[string]int{},
			interceptor:         pg.Interceptor,
		}
		for location := range config.locations {
			plg.locationRegexps[location] = regexp.MustCompile(location)
			plg.locationWeights[location] = len(location)
		}
		s.pluginList = append(s.pluginList, plg)
	}
	return s
}

// getPluginByRequest 根据请求获取插件
func (s *server) getPluginByRequest(request *http.Request) []*pluginWrapper {
	url := fmt.Sprintf("%s%s", request.Host, request.URL.Path)
	if plugins, ok := s.cacheMappingPluginURL.Load(url); ok {
		return plugins.([]*pluginWrapper)
	}

	host := strings.Split(request.Host, ":")[0]
	plgWrappers := make([]*pluginWrapper, 0)
	for _, plg := range s.pluginList {
		if len(plg.hosts) > 0 && !util.In(plg.hosts, host) {
			continue
		}

		var (
			weight      int
			maxPlg      *plugin
			proxyPass   string
			locationStr string
		)

		for location, reg := range plg.locationRegexps {
			if !reg.MatchString(request.URL.Path) {
				continue
			}
			if plg.locationWeights[location] > weight {
				weight = plg.locationWeights[location]
				maxPlg = plg
				proxyPass = plg.locationProxyPasses[location]
				locationStr = location
			}
		}
		if maxPlg != nil {
			plgWrappers = append(plgWrappers, &pluginWrapper{
				location:  locationStr,
				proxyPass: proxyPass,
				plugin:    maxPlg,
			})
		}
	}

	s.cacheMappingPluginURL.Store(url, plgWrappers)
	return plgWrappers
}
