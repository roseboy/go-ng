package main

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
)

func main() {
	fmt.Println("execute the following command:")
	fmt.Println("curl \"localhost:8000/api\" -d'{\"Action\":\"GetGirlfriend\",\"Age\":20}'")
	fmt.Println()

	plg := &plugin.ActionPlugin{Endpoint: "/api"}
	plg.RegisterAction("GetGirlfriend", GetGirlfriend, &GetGirlfriendRequest{}, &GetGirlfriendResponse{})

	err := ng.NewServer().RegisterPlugins(plg).Start(8000)
	if err != nil {
		panic(err)
	}
}
