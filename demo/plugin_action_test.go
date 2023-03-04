package demo

import (
	"context"
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin"
	"testing"
	"time"
)

func TestStartActionPlugin(t *testing.T) {

	fmt.Println("execute the following command:")
	fmt.Println("curl \"localhost:8000/api\" -d'{\"Action\":\"GetGirlFriend\",\"Age\":20}'")
	fmt.Println()

	var actions = map[string]plugin.Action{
		"GetGirlFriend": plugin.NewAction(GetGirlFriend, &GetGirlFriendRequest{}, &GetGirlFriendResponse{}),
	}

	ng.RegisterPlugins(
		&plugin.ActionPlugin{
			Endpoint:  "/api",
			ActionMap: actions,
		},
	)
	ng.Start(8000)

	time.Sleep(24 * time.Hour)
}

// GetGirlFriendRequest Request
type GetGirlFriendRequest struct {
	Age int
}

// GetGirlFriendResponse Response
type GetGirlFriendResponse struct {
	Result string
}

// GetGirlFriend action
func GetGirlFriend(ctx context.Context, request, response any) error {
	var req, resp = request.(*GetGirlFriendRequest), response.(*GetGirlFriendResponse)
	if req.Age <= 18 {
		return plugin.NewError(10001, "get girlfriend error")
	}
	resp.Result = "congratulations, you got a girlfriend!!!"
	return nil
}
