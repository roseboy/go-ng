package main

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"log"
)

// DemoPlugin demo
type DemoPlugin struct{}

// Config config
func (v *DemoPlugin) Config(config *ng.PluginConfig) {
	config.SetName("ng_plugin_demo")
	config.AddProxyPass("/", "https://news.baidu.com/")
}

// Interceptor interceptor
func (v *DemoPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	log.Println("===> DemoPlugin before...")
	request.Headers["ProxyPass"] = "NG"
	fmt.Println(request.HttpRequest.RequestURI)
	err := ng.Invoke(request, response)
	log.Println("===> DemoPlugin after...")
	return err
}
