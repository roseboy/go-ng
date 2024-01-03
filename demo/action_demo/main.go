package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
	"github.com/roseboy/go-ng/plugin/action"
)

func main() {
	printTestInfo()
	plg := plugin.NewActionPlugin("/api", true, GetAuthInfo)
	plg.RegisterAction(GetGirlfriend, reflect.TypeOf(new(GetGirlfriendRequest)), reflect.TypeOf(new(GetGirlfriendResponse)))

	err := ng.NewServer(8000).RegisterPlugins(plg).Start()
	if err != nil {
		panic(err)
	}
}

func printTestInfo() {
	timestamp := time.Now().Unix() + 5
	body := fmt.Sprintf(`{"Action":"GetGirlfriend","Age":20,"Timestamp":%d,"AppId":111}`, timestamp)
	sign := action.CalcSignatureV1(&action.SignatureArgs{
		Timestamp: timestamp, Payload: body, SecretKey: secretKey,
	})
	authorization := fmt.Sprintf("%s;%s", secretId, sign)

	fmt.Println(" execute the following command in 10s:")
	fmt.Printf(" curl 'http://localhost:8000/api?a=1' -d'%s' -H'Authorization:%s'\n", body, authorization)
	fmt.Println()
}
