package main

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("open 'http://localhost:8000/test' in browser")
	fmt.Println()

	lb1 := plugin.NewLoadBalancePluginWithPolicy("localhost", "/test", []string{
		"http://127.0.0.1:8080",
		"http://127.0.0.1:9090",
		"http://127.0.0.1:8081",
		"http://127.0.0.1:9091",
	}, RandPolicyFunc)

	lb2 := plugin.NewLoadBalancePlugin("test.com", "/", []string{
		"http://127.0.0.1:18080",
		"http://127.0.0.1:19090",
		"http://127.0.0.1:18081",
		"http://127.0.0.1:19091",
	})

	err := ng.NewServer(8000).RegisterPlugins(lb1, lb2).Start()
	if err != nil {
		panic(err)
	}
}

func RandPolicyFunc(proxyPassList []string) string {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(proxyPassList))
	proxyPass := proxyPassList[randIndex]
	fmt.Println(proxyPass)
	return proxyPass
}
