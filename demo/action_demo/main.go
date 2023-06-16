package main

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
	"time"
)

func main() {
	timestamp := time.Now().Unix() + 5
	body := fmt.Sprintf(`{"Action":"GetGirlfriend","Age":20,"Timestamp":%d,"AppId":111}`, timestamp)
	sign := plugin.CalcSignature(&plugin.CalcSignatureArgs{
		Service:   "/api",
		Timestamp: timestamp,
		Method:    "POST",
		Host:      "localhost:8000",
		URI:       "/api",
		Payload:   body,
		SecretKey: secretKey,
	})
	authorization := fmt.Sprintf("%s;%s", secretId, sign)

	fmt.Println("execute the following command in 10s:")
	fmt.Printf("curl localhost:8000/api -d'%s' -H'Authorization:%s'\n", body, authorization)
	fmt.Println()

	plg := &plugin.ActionPlugin{Endpoint: "/api", SignatureCheck: true, AuthInfoFunc: GetAuthInfo}
	plg.RegisterAction("GetGirlfriend", GetGirlfriend, &GetGirlfriendRequest{}, &GetGirlfriendResponse{})

	err := ng.NewServer(8000).RegisterPlugins(plg).Start()
	if err != nil {
		panic(err)
	}
}
