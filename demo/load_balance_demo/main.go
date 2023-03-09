package main

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
)

func main() {
	fmt.Println("open 'http://localhost:8000/test' in browser")
	fmt.Println()

	lb := &plugin.LoadBalancePlugin{
		ServerName: "localhost:8000",
		Location:   "/test",
		ProxyPassList: []string{
			"http://127.0.0.1:8080",
			"http://127.0.0.1:9090",
			"http://127.0.0.1:8081",
			"http://127.0.0.1:9091",
		},
	}

	lb2 := &plugin.LoadBalancePlugin{
		ServerName: "test.com:8000",
		Location:   "/",
		ProxyPassList: []string{
			"http://127.0.0.1:18080",
			"http://127.0.0.1:19090",
			"http://127.0.0.1:18081",
			"http://127.0.0.1:19091",
		},
	}

	err := ng.NewServer(8000).RegisterPlugins(lb, lb2).Start()
	if err != nil {
		panic(err)
	}
}
