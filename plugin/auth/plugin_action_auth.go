package auth

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/roseboy/go-ng/ng"
	"github.com/roseboy/go-ng/plugin/action"
	"github.com/roseboy/go-ng/util"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	expireSecond  = 10
	actionMetaKey = "CTX_ACTION_META_KEY"
)

// PluginActionAuth auth
type PluginActionAuth struct {
	AuthLocation  string
	SignatureFunc func(args *SignatureArgs) string
	AuthInfoFunc  func(string) (uint64, string)
}

// Config config
func (p *PluginActionAuth) Config(config *ng.PluginConfig) {
	config.Name("ng_action_auth_plugin")
	if p.AuthLocation == "" {
		p.AuthLocation = "/"
	}
	config.Location(p.AuthLocation, "")
	if p.SignatureFunc == nil {
		p.SignatureFunc = CalcSignatureV1
	}
}

// Interceptor interceptor
func (p *PluginActionAuth) Interceptor(ctx context.Context, request *ng.Request, response *ng.Response) error {
	var (
		err error
	)
	meta := action.ExtractMeta(ctx)
	if meta == nil {
		meta = &action.Meta{}
		err = json.Unmarshal([]byte(request.Body), meta)
		if err != nil {
			return nil
		}
		ctx = context.WithValue(ctx, actionMetaKey, meta)
	}

	err = p.checkAuthorization(meta, request)
	if err == nil {
		return ng.Invoke(ctx, request, response)
	}

	requestId := util.If(meta.RequestId == "", fmt.Sprintf("s%d", time.Now().UnixNano()), meta.RequestId)
	resp := map[string]any{
		"Error": map[string]any{
			"Code": http.StatusUnauthorized,
			"Msg":  err.Error(),
		},
		"RequestId": requestId,
	}
	data, _ := json.Marshal(resp)
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("X-Request-Id", requestId)
	response.Body, response.Status = string(data), http.StatusUnauthorized
	return nil
}

func (p *PluginActionAuth) checkAuthorization(meta *action.Meta, request *ng.Request) error {
	authorization := strings.Split(request.Headers["Authorization"], ";")
	if len(authorization) < 2 {
		return errors.New("AuthFailure.SignatureFailure")
	}
	secretId := authorization[0]
	reqSign := authorization[1]
	appId, secretKey := p.AuthInfoFunc(secretId)
	if appId == 0 {
		return errors.New("AuthFailure.AppNotExist")
	}

	nowTimestamp := time.Now().Unix()
	timestamp := meta.Timestamp
	if math.Abs(float64(timestamp-nowTimestamp)) > expireSecond {
		timestamp = nowTimestamp
	}
	sign := p.SignatureFunc(&SignatureArgs{
		Service:   p.AuthLocation,
		Timestamp: timestamp,
		Method:    request.HttpRequest.Method,
		Host:      request.HttpRequest.Host,
		URI:       request.HttpRequest.RequestURI,
		Payload:   request.Body,
		SecretKey: secretKey,
	})

	if reqSign != sign {
		return errors.New("AuthFailure.SignatureFailure")
	}

	meta.AppId = appId
	meta.SecretId = secretId

	return nil
}

// SignatureArgs signature args
type SignatureArgs struct {
	Service   string
	Timestamp int64
	Method    string
	Host      string
	URI       string
	Payload   string
	SecretKey string
}

// CalcSignatureV1 calc signature
func CalcSignatureV1(args *SignatureArgs) string {
	hashedPayload := util.SHA256Hex(args.Payload)
	canonicalRequest := fmt.Sprintf("%d;%s", args.Timestamp, hashedPayload)
	signatureString := util.HMacSHA256(canonicalRequest, args.SecretKey)
	signature := hex.EncodeToString([]byte(signatureString))
	return signature
}

// CalcSignatureV2 calc signature
func CalcSignatureV2(args *SignatureArgs) string {
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
