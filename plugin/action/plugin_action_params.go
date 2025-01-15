package action

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/util"
)

// PluginActionParams params
type PluginActionParams struct {
	Endpoint string
}

// Config config
func (p *PluginActionParams) Config(config *ng.PluginConfig) {
	path := util.If(strings.HasPrefix(p.Endpoint, "/"), p.Endpoint, "/"+p.Endpoint)
	config.Name("ng_action_params_plugin")
	config.Location(path, "")
}

// Interceptor interceptor
func (p *PluginActionParams) Interceptor(ctx context.Context, request *ng.Request, response *ng.Response) error {
	_ = request.HttpRequest.ParseForm()
	bodyMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(request.Body), &bodyMap)
	for k, v := range request.HttpRequest.Form {
		if len(v) == 0 {
			bodyMap[k] = ""
		} else if len(v) == 1 {
			bodyMap[k] = v[0]
		} else {
			bodyMap[k] = v
		}
	}
	formBody, _ := json.Marshal(bodyMap)
	request.Body = string(formBody)
	err := ng.Invoke(ctx, request, response)
	return err
}
