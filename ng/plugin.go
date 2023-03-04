package ng

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
)

var (
	pluginList            = make([]*plugin, 0)
	pluginURIMappingCache = sync.Map{}
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
func RegisterPlugins(plugins ...pluginInterface) {
	for _, pg := range plugins {
		config := &PluginConfig{}
		pg.Config(config)
		log.Printf("RegisterPlugin: %s %v", config.name, config.locations)
		for location, proxyPass := range config.locations {
			pluginList = append(pluginList,
				&plugin{name: config.name, proxyPass: proxyPass, location: location,
					locationRegexp: regexp.MustCompile(location),
					interceptor:    pg.Interceptor})
		}
	}
	return
}

// getPluginByRequest 根据请求获取插件
func getPluginByRequest(request *http.Request) []*plugin {
	url := fmt.Sprintf("%s%s", request.Host, request.URL.Path)

	if plugins, ok := pluginURIMappingCache.Load(url); ok {
		return plugins.([]*plugin)
	}
	plugins := make([]*plugin, 0)
	for _, plugin := range pluginList {
		if plugin.locationRegexp.MatchString(url) {
			plugins = append(plugins, plugin)
		}
	}
	if len(plugins) > 0 {
		pluginURIMappingCache.Store(url, plugins)
	}
	return plugins
}
