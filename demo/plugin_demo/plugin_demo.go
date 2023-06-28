package main

import (
	"log"

	"github.com/roseboy/go-ng/ng"
)

// DemoPlugin demo
type DemoPlugin struct{}

// Config config
func (v *DemoPlugin) Config(config *ng.PluginConfig) {
	config.Name("ng_demo_plugin")
	config.ProxyPass("/", "https://news.baidu.com/")
}

// Interceptor interceptor
func (v *DemoPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	log.Println("===> DemoPlugin before...")
	request.Headers["ProxyPass"] = "NG"
	log.Println(request.HttpRequest.RequestURI)
	err := ng.Invoke(request, response)
	log.Println("===> DemoPlugin after...")
	return err
}
