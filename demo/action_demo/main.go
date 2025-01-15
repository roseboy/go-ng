package main

import (
	"fmt"
	"github.com/roseboy/go-ng/plugin/auth"
	"time"

	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
)

func main() {
	printTestInfo()
	plg := plugin.NewActionPlugin("/api")
	plgParams := plugin.NewActionParamsPlugin("/api")

	//actionAuthPlg := &auth.PluginActionAuth{AuthLocation: "/api", AuthInfoFunc: GetAuthInfo}
	//actionAuthPlg := &auth.PluginBasicAuth{AuthLocation: "/api", GetAuthInfo: GetBasicAuthInfo}

	err := ng.NewServer(8000).RegisterPlugins(plgParams, plg).Start()
	if err != nil {
		panic(err)
	}
}

func printTestInfo() {
	timestamp := time.Now().Unix() + 5
	body := fmt.Sprintf(`{"Action":"GetGirlfriend","Age":20,"Timestamp":%d,"AppId":111}`, timestamp)
	sign := auth.CalcSignatureV1(&auth.SignatureArgs{
		Timestamp: timestamp, Payload: body, SecretKey: secretKey,
	})
	authorization := fmt.Sprintf("%s;%s", secretId, sign)

	fmt.Println(" execute the following command in 10s:")
	fmt.Printf(" curl 'http://localhost:8000/api?a=1' -d'%s' -H'Authorization:%s'\n", body, authorization)
	fmt.Println()
}
