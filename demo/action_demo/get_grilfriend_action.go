package main

import (
	"context"
	"github.com/roseboy/go-ng/plugin"
)

// GetGirlfriendRequest Request
type GetGirlfriendRequest struct {
	Age int
}

// GetGirlfriendResponse Response
type GetGirlfriendResponse struct {
	Result string
}

// GetGirlfriend action
func GetGirlfriend(ctx context.Context, request, response any) error {
	var req, resp = request.(*GetGirlfriendRequest), response.(*GetGirlfriendResponse)
	if req.Age <= 18 {
		return plugin.NewError(10001, "get girlfriend error")
	}
	resp.Result = "congratulations, you got a girlfriend!!!"
	return nil
}
