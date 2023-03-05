package ng

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type plugin struct {
	name           string
	proxyPass      string
	location       string
	locationRegexp *regexp.Regexp
	interceptor    func(request *Request, response *Response) error
}

// pluginInterface plugin interface
type pluginInterface interface {
	Interceptor(request *Request, response *Response) error
	Config(config *PluginConfig)
}

// PluginConfig plugin config
type PluginConfig struct {
	name      string
	locations map[string]string
}

// SetName plugin name
func (c *PluginConfig) SetName(name string) {
	c.name = name
}

// AddProxyPass 设置转发
func (c *PluginConfig) AddProxyPass(location, proxyPass string) {
	if c.locations == nil {
		c.locations = make(map[string]string)
	}
	c.locations[location] = proxyPass
}

// AddLocation  设置转发
func (c *PluginConfig) AddLocation(location string) {
	c.AddProxyPass(location, "")
}

// RegisterPlugins 注册插件
func (s *server) RegisterPlugins(plugins ...pluginInterface) *server {
	for _, pg := range plugins {
		config := &PluginConfig{}
		pg.Config(config)
		log.Printf("RegisterPlugin: %s %v", config.name, config.locations)
		for location, proxyPass := range config.locations {
			s.pluginList = append(s.pluginList,
				&plugin{name: config.name, proxyPass: proxyPass, location: location,
					locationRegexp: regexp.MustCompile(location),
					interceptor:    pg.Interceptor})
		}
	}
	return s
}

// getPluginByRequest 根据请求获取插件
func (s *server) getPluginByRequest(request *http.Request) []*plugin {
	url := fmt.Sprintf("%s%s", request.Host, request.URL.Path)

	if plugins, ok := s.pluginURIMappingCache.Load(url); ok {
		return plugins.([]*plugin)
	}
	plugins := make([]*plugin, 0)
	for _, plugin := range s.pluginList {
		if plugin.locationRegexp.MatchString(url) {
			plugins = append(plugins, plugin)
		}
	}
	if len(plugins) > 0 {
		s.pluginURIMappingCache.Store(url, plugins)
	}
	return plugins
}
