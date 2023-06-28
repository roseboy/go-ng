package main

import (
	"fmt"
	"time"

	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
	"github.com/roseboy/go-ng/plugin/action"
)

func main() {
	timestamp := time.Now().Unix() + 5
	body := fmt.Sprintf(`{"Action":"GetGirlfriend","Age":20,"Timestamp":%d,"AppId":111}`, timestamp)
	sign := action.CalcSignature(&action.CalcSignatureArgs{
		Service:   "/api",
		Timestamp: timestamp,
		Method:    "POST",
		Host:      "localhost:8000",
		URI:       "/api?a=1",
		Payload:   body,
		SecretKey: secretKey,
	})
	authorization := fmt.Sprintf("%s;%s", secretId, sign)

	fmt.Println(" execute the following command in 10s:")
	fmt.Printf(" curl 'http://localhost:8000/api?a=1' -d'%s' -H'Authorization:%s'\n", body, authorization)
	fmt.Println()

	plg := plugin.NewActionPlugin("/api", true, GetAuthInfo)
	plg.RegisterAction("GetGirlfriend", GetGirlfriend, &GetGirlfriendRequest{}, &GetGirlfriendResponse{})

	err := ng.NewServer(8000).RegisterPlugins(plg).Start()
	if err != nil {
		panic(err)
	}
}
