package plugin

import (
	"github.com/roseboy/go-ng/ng"
	"log"
)

// DemoPlugin demo
type DemoPlugin struct{}

// Config config
func (v *DemoPlugin) Config(config *ng.PluginConfig) {
	config.SetName("ng_plugin_demo")
	config.AddProxyPass("/", "http://localhost")
}

// Interceptor interceptor
func (v *DemoPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	log.Println("===> DemoPlugin before...")
	err := ng.Invoke(request, response)
	log.Println("===> DemoPlugin after...")
	return err
}
