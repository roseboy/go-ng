package action

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/roseboy/go-ng/util"
	"math"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/roseboy/go-ng/ng"
)

const (
	expireSecond = 10
)

var ctxMetaKey = &Meta{}

// PluginAction action
type PluginAction struct {
	Endpoint       string
	ActionMap      map[string]Action
	SignatureCheck bool
	AuthInfoFunc   func(string) (uint64, string)
	actionMap      sync.Map
}

// Config config
func (p *PluginAction) Config(config *ng.PluginConfig) {
	path := util.If(strings.HasPrefix(p.Endpoint, "/"), p.Endpoint, "/"+p.Endpoint)
	config.Name("ng_action_plugin")
	config.ProxyPass(path, "")
	for k, v := range p.ActionMap {
		p.actionMap.Store(k, v)
	}
}

// Interceptor interceptor
func (p *PluginAction) Interceptor(request *ng.Request, response *ng.Response) error {
	var (
		err       error
		requestId string
	)

	defer func() {
		if err == nil {
			return
		}
		actionResponse := &actionResponse{RequestId: requestId}
		if e, ok := err.(*actionError); ok {
			actionResponse.Error = e
		} else {
			actionResponse.Error = &actionError{Code: -1, Msg: err.Error()}
		}
		data, _ := json.Marshal(actionResponse)
		response.Body, response.Status = string(data), http.StatusOK
	}()

	var meta Meta
	err = json.Unmarshal([]byte(request.Body), &meta)
	if err != nil {
		return nil
	}

	requestId = meta.RequestId
	requestId = util.If(requestId == "", fmt.Sprintf("s%d", time.Now().UnixNano()), requestId)
	meta.RequestId = requestId
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("X-Request-Id", requestId)
	ctx := context.WithValue(context.Background(), ctxMetaKey, &meta)

	if p.SignatureCheck {
		err = p.checkSignature(ctx, request)
	}
	if err != nil {
		return nil
	}

	err = p.doAction(ctx, request, response)
	return nil
}

func (p *PluginAction) doAction(ctx context.Context, request *ng.Request, response *ng.Response) error {
	var meta = ctx.Value(ctxMetaKey).(*Meta)
	meta.Headers = make(map[string]string)
	for k, v := range request.Headers {
		meta.Headers[k] = v
	}

	actionFunc, ok := p.actionMap.Load(meta.Action)
	if !ok {
		return errors.New("action not found")
	}

	fun, req, resp := actionFunc.(Action)()
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

	actionResponse := &actionResponse{RequestId: meta.RequestId, Response: resp}
	data, _ := json.Marshal(actionResponse)
	response.Body, response.Status = string(data), http.StatusOK
	return nil
}

// RegisterAction register action
func (p *PluginAction) RegisterAction(actionName string, actionFunc actionFunc, request, response any) {
	if p.ActionMap == nil {
		p.ActionMap = map[string]Action{}
	}
	p.ActionMap[actionName] = NewAction(actionFunc, reflect.TypeOf(request), reflect.TypeOf(response))
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
	meta, ok := ctx.Value(ctxMetaKey).(*Meta)
	if ok {
		return meta
	}
	return nil
}

func (p *PluginAction) checkSignature(ctx context.Context, request *ng.Request) error {
	if p.AuthInfoFunc == nil {
		return NewError(http.StatusUnauthorized, "AuthFailure.AuthInfoError")
	}

	var meta = ctx.Value(ctxMetaKey).(*Meta)
	authorization := strings.Split(request.Headers["Authorization"], ";")
	if len(authorization) < 2 {
		return NewError(http.StatusUnauthorized, "AuthFailure.SignatureFailure")
	}
	secretId := authorization[0]
	reqSign := authorization[1]
	appId, secretKey := p.AuthInfoFunc(secretId)
	if appId == 0 {
		return NewError(http.StatusUnauthorized, "AuthFailure.AppNotExist")
	}

	nowTimestamp := time.Now().Unix()
	timestamp := meta.Timestamp
	if math.Abs(float64(timestamp-nowTimestamp)) > expireSecond {
		timestamp = nowTimestamp
	}
	sign := CalcSignature(&CalcSignatureArgs{
		Service:   p.Endpoint,
		Timestamp: timestamp,
		Method:    request.HttpRequest.Method,
		Host:      request.HttpRequest.Host,
		URI:       request.HttpRequest.RequestURI,
		Payload:   request.Body,
		SecretKey: secretKey,
	})

	if reqSign != sign {
		return NewError(http.StatusUnauthorized, "AuthFailure.SignatureFailure")
	}

	meta.AppId = appId
	meta.SecretId = secretId
	return nil
}

// CalcSignatureArgs signature args
type CalcSignatureArgs struct {
	Service   string
	Timestamp int64
	Method    string
	Host      string
	URI       string
	Payload   string
	SecretKey string
}

// CalcSignature calc signature
func CalcSignature(args *CalcSignatureArgs) string {
	hashedPayload := util.SHA256Hex(args.Payload)
	canonicalRequest := fmt.Sprintf("%s;%s;%s;%d;%s",
		args.Method,
		args.Host,
		args.URI,
		args.Timestamp,
		hashedPayload)
	date := time.Unix(args.Timestamp, 0).UTC().Format("2006-01-02")
	secretDate := util.HMacSHA256(date, args.SecretKey)
	secretService := util.HMacSHA256(args.Service, secretDate)
	signatureString := util.HMacSHA256(canonicalRequest, secretService)
	signature := hex.EncodeToString([]byte(signatureString))
	return signature
}
