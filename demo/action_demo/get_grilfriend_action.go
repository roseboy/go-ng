package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/roseboy/go-ng/plugin/action"
)

func init() {
	action.RegisterAction("/api", GetGirlfriend,
		reflect.TypeOf(new(GetGirlfriendRequest)), reflect.TypeOf(new(GetGirlfriendResponse)))
}

// GetGirlfriendRequest Request
type GetGirlfriendRequest struct {
	action.Meta
	Age int
}

// GetGirlfriendResponse Response
type GetGirlfriendResponse struct {
	Result string
}

// GetGirlfriend action
func GetGirlfriend(ctx context.Context, request, response any) error {
	var req, resp = request.(*GetGirlfriendRequest), response.(*GetGirlfriendResponse)
	meta := action.ExtractMeta(ctx)
	fmt.Println(req.AppId)
	mb, _ := json.Marshal(meta)
	fmt.Println(string(mb))
	if req.Age <= 18 {
		return action.NewError(10001, "get girlfriend error")
	}
	resp.Result = "congratulations, you got a girlfriend!!!"
	return nil
}
