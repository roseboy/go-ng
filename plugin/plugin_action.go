package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"reflect"
	"strings"
	"sync"
	"time"
)

var ctxMetaKey = &ActionMeta{}

// ActionPlugin action
type ActionPlugin struct {
	Endpoint  string
	ActionMap map[string]Action
	actionMap sync.Map
}

// Config config
func (p *ActionPlugin) Config(config *ng.PluginConfig) {
	config.SetName("ng_action_plugin")
	config.AddLocation(ng.If(strings.HasPrefix(p.Endpoint, "/"), p.Endpoint, "/"+p.Endpoint))
	for k, v := range p.ActionMap {
		p.actionMap.Store(k, v)
	}
}

// Interceptor interceptor
func (p *ActionPlugin) Interceptor(request *ng.Request, response *ng.Response) error {
	var (
		requestId = fmt.Sprintf("s%d", time.Now().UnixNano())
		meta      = &ActionMeta{RequestId: requestId, Headers: map[string]string{}}
		ctx       = context.WithValue(context.Background(), ctxMetaKey, meta)
	)

	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("X-Request-Id", requestId)
	for k, v := range request.Headers {
		meta.Headers[k] = v
	}

	err := p.doAction(ctx, request, response)
	if err != nil {
		actionResponse := &actionResponse{RequestId: requestId}
		if e, ok := err.(*actionError); ok {
			actionResponse.Error = e
		} else {
			actionResponse.Error = &actionError{Code: -1, Msg: err.Error()}
		}
		data, _ := json.Marshal(actionResponse)
		response.Body, response.Status = string(data), 200
	}
	return nil
}

func (p *ActionPlugin) doAction(ctx context.Context, request *ng.Request, response *ng.Response) error {
	var meta = ctx.Value(ctxMetaKey).(*ActionMeta)
	actionRequest := actionRequest{}
	err := json.Unmarshal([]byte(request.Body), &actionRequest)
	if err != nil {
		return err
	}

	actionFunc, ok := p.actionMap.Load(actionRequest.Action)
	if !ok {
		return errors.New("action not found")
	}

	fun, req, resp := actionFunc.(Action)()
	err = json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return err
	}

	err = fun(ctx, req, resp)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("action response is null")
	}

	actionResponse := &actionResponse{RequestId: meta.RequestId, Response: resp}
	data, _ := json.Marshal(actionResponse)
	response.Body, response.Status = string(data), 200
	return nil
}

func (p *ActionPlugin) RegisterAction(actionName string, actionFunc Action) *ActionPlugin {
	p.actionMap.Store(actionName, actionFunc)
	return p
}

// ActionMeta meta
type ActionMeta struct {
	RequestId string
	Headers   map[string]string
}

// actionRequest request
type actionRequest struct {
	Action string
}

// actionResponse response
type actionResponse struct {
	Response  any          `json:"Response,omitempty"`
	RequestId string       `json:"RequestId,omitempty"`
	Error     *actionError `json:"Error,omitempty"`
}

// ActionError error
type actionError struct {
	Code int
	Msg  string
}

// Error error
func (e *actionError) Error() string {
	return e.Msg
}

// NewError error
func NewError(code int, msg string) error {
	return &actionError{
		Code: code,
		Msg:  msg,
	}
}

type actionFunc func(context.Context, any, any) error
type Action func() (actionFunc, any, any)

// NewAction new
func NewAction(fun actionFunc, req, resp any) Action {
	return func() (actionFunc, any, any) {
		return fun, ng.NewInstanceByType(reflect.TypeOf(req)), ng.NewInstanceByType(reflect.TypeOf(resp))
	}
}
