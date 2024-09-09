package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"net/http"
	"time"
)

// PluginBasicAuth auth
type PluginBasicAuth struct {
	AuthLocation string
	GetAuthInfo  func(string) (string, string, error)
}

// Config config
func (p *PluginBasicAuth) Config(config *ng.PluginConfig) {
	config.Name("ng_basic_auth_plugin")
	if p.AuthLocation == "" {
		p.AuthLocation = "/"
	}
	config.Location(p.AuthLocation, "")
}

// Interceptor interceptor
func (p *PluginBasicAuth) Interceptor(ctx context.Context, request *ng.Request, response *ng.Response) error {
	var err error
	err = p.checkAuthorization(ctx, request)
	if err == nil {
		return ng.Invoke(ctx, request, response)
	}

	contentType := request.Headers["Content-Type"]
	requestId := fmt.Sprintf("s%d", time.Now().UnixNano())
	body := ""
	if contentType == "application/json" {
		resp := map[string]any{
			"Error": map[string]any{
				"Code": http.StatusUnauthorized,
				"Msg":  err.Error(),
			},
			"RequestId": requestId,
		}
		data, _ := json.Marshal(resp)
		body = string(data)
	} else {
		body = http.StatusText(http.StatusUnauthorized)
	}

	response.SetHeader("Content-Type", contentType)
	response.SetHeader("X-Request-Id", requestId)
	response.SetHeader("WWW-Authenticate", `Basic realm="Restricted"`)
	response.Body, response.Status = body, http.StatusUnauthorized
	return nil
}

func (p *PluginBasicAuth) checkAuthorization(ctx context.Context, request *ng.Request) error {
	username, password, ok := request.HttpRequest.BasicAuth()
	if !ok {
		return errors.New("AuthFailure")
	}
	if p.GetAuthInfo == nil {
		return errors.New("AuthFailure")
	}
	_, pass, err := p.GetAuthInfo(username)
	if err != nil {
		return err
	}
	if pass != password {
		return errors.New("AuthFailure")
	}
	return nil
}
