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

	ssl := &plugin.SSLPlugin{
		CertFile:       "./demo/ssl_demo/test.crt",
		KeyFile:        "./demo/ssl_demo/test.key",
		AutoRedirect:   true,
		HttpServerPort: 80,
	}

	lb := &plugin.LoadBalancePlugin{
		ServerName: "test.com",
		Location:   "/",
		ProxyPassList: []string{
			"http://127.0.0.1:8080",
			"http://127.0.0.1:9090",
			"http://127.0.0.1:8081",
			"http://127.0.0.1:9091",
		},
	}

	err := ng.NewServer(443).RegisterPlugins(lb, ssl).Start()
	if err != nil {
		panic(err)
	}
}
