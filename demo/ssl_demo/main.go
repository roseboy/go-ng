package main

import (
	"fmt"

	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
)

func main() {
	fmt.Println("edit host file: 127.0.0.1 test.com")
	fmt.Println("open 'https://test.com' in browser")
	fmt.Println()

	ssl := plugin.NewSSLPluginWithAutoRedirect(
		"./demo/ssl_demo/test.crt",
		"./demo/ssl_demo/test.key",
		80)

	lb := plugin.NewLoadBalancePlugin("test.com", "/", []string{
		"http://127.0.0.1:18080",
		"http://127.0.0.1:19090",
		"http://127.0.0.1:18081",
		"http://127.0.0.1:19091",
	})

	err := ng.NewServer(443).RegisterPlugins(lb, ssl).Start()
	if err != nil {
		panic(err)
	}
}
