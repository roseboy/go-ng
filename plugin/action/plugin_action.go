package action

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/util"
)

const actionMetaKey = "CTX_ACTION_META_KEY"

var initActionMap = make(map[string]map[string]actionFunc)
var initActionReqMap = make(map[string]map[string]reflect.Type)
var initActionRespMap = make(map[string]map[string]reflect.Type)

// PluginAction action
type PluginAction struct {
	Endpoint  string
	ActionMap map[string]Action
	actionMap sync.Map
}

// Config config
func (p *PluginAction) Config(config *ng.PluginConfig) {
	path := util.If(strings.HasPrefix(p.Endpoint, "/"), p.Endpoint, "/"+p.Endpoint)
	config.Name("ng_action_plugin")
	config.Location(path, "")
	for k, v := range p.ActionMap {
		p.actionMap.Store(k, v)
	}
}

// Interceptor interceptor
func (p *PluginAction) Interceptor(ctx context.Context, request *ng.Request, response *ng.Response) error {
	var (
		err       error
		requestId string
	)

	defer func() {
		if err == nil {
			return
		}
		actResp := &actionResponse{RequestId: requestId}
		var e = &actionError{}
		if errors.As(err, &e) {
			actResp.Error = e
		} else {
			actResp.Error = &actionError{Code: -1, Msg: err.Error()}
		}
		data, _ := json.Marshal(actResp)
		response.Body, response.Status = string(data), http.StatusOK
	}()

	meta := ExtractMeta(ctx)
	if meta == nil {
		meta = &Meta{}
		err = json.Unmarshal([]byte(request.Body), meta)
		if err != nil {
			return nil
		}
		ctx = context.WithValue(ctx, actionMetaKey, meta)
	}

	requestId = meta.RequestId
	requestId = util.If(requestId == "", fmt.Sprintf("s%d", time.Now().UnixNano()), requestId)
	meta.RequestId = requestId
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("X-Request-Id", requestId)
	err = p.doAction(ctx, request, response)
	return nil
}

func (p *PluginAction) doAction(ctx context.Context, request *ng.Request, response *ng.Response) error {
	var meta = ExtractMeta(ctx)
	meta.Headers = make(map[string]string)
	for k, v := range request.Headers {
		meta.Headers[k] = v
	}

	actFunc, ok := p.actionMap.Load(meta.Action)
	if !ok {
		return errors.New("action not found")
	}

	fun, req, resp := actFunc.(Action)()
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return err
	}
	metaData, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	err = json.Unmarshal(metaData, req)
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

	actResp := &actionResponse{RequestId: meta.RequestId, Response: resp}
	data, _ := json.Marshal(actResp)
	response.Body, response.Status = string(data), http.StatusOK
	return nil
}

// RegisterAction register action
func (p *PluginAction) RegisterAction(actionFunc actionFunc, reqType, respType reflect.Type) {
	actionName := runtime.FuncForPC(reflect.ValueOf(actionFunc).Pointer()).Name()
	actionName = actionName[strings.LastIndex(actionName, ".")+1:]
	p.RegisterActionWithName(actionName, actionFunc, reqType, respType)
}

// RegisterActionWithName register action
func (p *PluginAction) RegisterActionWithName(actionName string, actionFunc actionFunc, reqType, respType reflect.Type) {
	if p.ActionMap == nil {
		p.ActionMap = map[string]Action{}
	}
	p.ActionMap[actionName] = NewAction(actionFunc, reqType, respType)
}

// Meta meta
type Meta struct {
	Headers map[string]string

	Action    string
	AppId     uint64
	RequestId string
	Nonce     string
	SecretId  string
	Timestamp int64
	//Version   string
	//Region    string
	//Language  string
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

type actionFunc func(context.Context, any, any) error
type Action func() (actionFunc, any, any)

// NewAction new
func NewAction(fun actionFunc, reqType, respType reflect.Type) Action {
	return func() (actionFunc, any, any) {
		return fun, util.NewInstanceByType(reqType), util.NewInstanceByType(respType)
	}
}

// NewError new error
func NewError(code int, msg string) error {
	return &actionError{
		Code: code,
		Msg:  msg,
	}
}

// ExtractMeta extract meta
func ExtractMeta(ctx context.Context) *Meta {
	meta, ok := ctx.Value(actionMetaKey).(*Meta)
	if ok {
		return meta
	}
	return nil
}

// RegisterAction register action
func RegisterAction(endpoint string, actionFunc actionFunc, reqType, respType reflect.Type) {
	actionName := runtime.FuncForPC(reflect.ValueOf(actionFunc).Pointer()).Name()
	actionName = actionName[strings.LastIndex(actionName, ".")+1:]
	RegisterActionWithName(endpoint, actionName, actionFunc, reqType, respType)
}

// RegisterActionWithName register action
func RegisterActionWithName(endpoint, actionName string, actFunc actionFunc, reqType, respType reflect.Type) {
	path := util.If(strings.HasPrefix(endpoint, "/"), endpoint, "/"+endpoint)
	if _, ok := initActionMap[path]; !ok {
		initActionMap[path] = map[string]actionFunc{}
		initActionReqMap[path] = map[string]reflect.Type{}
		initActionRespMap[path] = map[string]reflect.Type{}
	}
	initActionMap[path][actionName] = actFunc
	initActionReqMap[path][actionName] = reqType
	initActionRespMap[path][actionName] = respType
}

// RegisterInitAction .
func RegisterInitAction(action *PluginAction) {
	for path, actions := range initActionMap {
		if action.Endpoint != path {
			continue
		}
		for name, act := range actions {
			action.RegisterActionWithName(name, act, initActionReqMap[path][name], initActionRespMap[path][name])
		}
	}
}
