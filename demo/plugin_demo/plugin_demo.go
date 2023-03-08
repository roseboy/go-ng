package main

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"log"
	"strings"
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
	fmt.Println(request.HttpRequest.RequestURI)
	if strings.Contains(request.Url, "?") {
		request.Url = request.Url[0:strings.Index(request.Url, "?")]
	}
	fmt.Println(request.Url)
	err := ng.Invoke(request, response)
	log.Println("===> DemoPlugin after...")
	return err
}
