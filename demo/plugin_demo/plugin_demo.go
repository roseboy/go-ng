package main

import (
	"context"
	"log"

	"github.com/roseboy/go-ng/ng"
)

// DemoPlugin demo
type DemoPlugin struct{}

// Config config
func (v *DemoPlugin) Config(config *ng.PluginConfig) {
	config.Name("ng_demo_plugin")
	config.Location("/qq", "https://news.baidu.com/")
	config.Host("localhost")
}

// Interceptor interceptor
func (v *DemoPlugin) Interceptor(ctx context.Context, request *ng.Request, response *ng.Response) error {
	log.Println("===> DemoPlugin before...")
	request.Headers["ProxyPass"] = "NG"
	log.Println(request.HttpRequest.RequestURI)
	request.AllowRedirect = true
	err := ng.Invoke(ctx, request, response)
	log.Println("===> DemoPlugin after...")

	return err
}
