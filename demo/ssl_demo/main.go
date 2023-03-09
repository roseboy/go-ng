package main

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
)

func main() {
	fmt.Println("edit host file: 127.0.0.1 test.com")
	fmt.Println("open 'https://test.com:8000/' in browser")
	fmt.Println()

	ssl := &plugin.SSLPlugin{
		CertFile: "/data/test.crt",
		KeyFile:  "/data/test.key",
	}

	lb := &plugin.LoadBalancePlugin{
		ServerName: "test.com:8000",
		Location:   "/",
		ProxyPassList: []string{
			"http://127.0.0.1:8080",
			"http://127.0.0.1:9090",
			"http://127.0.0.1:8081",
			"http://127.0.0.1:9091",
		},
	}

	err := ng.NewServer(8000).RegisterPlugins(lb, ssl).Start()
	if err != nil {
		panic(err)
	}
}
