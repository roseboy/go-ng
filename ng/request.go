package ng

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// requestBuild 请求build
type requestBuild struct {
	request *Request
}

// Request 请求结构体
type Request struct {
	HttpRequest   *http.Request
	Method        string
	Url           string
	Body          string
	Params        map[string]string
	Headers       map[string]string
	AllowRedirect bool
	plugins       []*pluginWrapper
	pluginPos     int
}

// Response 响应结构体
type Response struct {
	HttpResponse   *http.Response
	ResponseWriter http.ResponseWriter
	Body           string
	Headers        map[string]string
	Status         int
}

// newRequest 创建新的请求
func newRequest() *requestBuild {
	return &requestBuild{request: &Request{Headers: make(map[string]string), Params: make(map[string]string)}}
}

// HttpRequest 设置请求
func (rb *requestBuild) HttpRequest(request *http.Request) *requestBuild {
	defer request.Body.Close()
	body, _ := io.ReadAll(request.Body)

	header := make(map[string]string)
	for k := range request.Header {
		v := request.Header.Get(k)
		header[k] = v
	}

	rb.request.HttpRequest = request
	rb.request.Method = request.Method
	rb.request.Headers = header
	rb.request.Body = string(body)
	return rb
}

// Url 设置url
func (rb *requestBuild) Url(url string) *requestBuild {
	rb.request.Url = url
	return rb
}

// Method 设置请求方法
func (rb *requestBuild) Method(method string) *requestBuild {
	rb.request.Method = method
	return rb
}

// Get 请求
func (rb *requestBuild) Get(url string) *requestBuild {
	rb.request.Method = "GET"
	rb.request.Url = url
	return rb
}

// Post 请求
func (rb *requestBuild) Post(url string) *requestBuild {
	rb.request.Method = "POST"
	rb.request.Url = url
	return rb
}

// Body 设置请求体
func (rb *requestBuild) Body(body string) *requestBuild {
	rb.request.Body = body
	return rb
}

// Param 设置请求参数
func (rb *requestBuild) Param(key string, value string) *requestBuild {
	rb.request.Params[key] = value
	return rb
}

// Params 批量设置请求参数
func (rb *requestBuild) Params(params map[string]string) *requestBuild {
	rb.request.Params = params
	return rb
}

// Header 设置请求头
func (rb *requestBuild) Header(key string, value string) *requestBuild {
	rb.request.Headers[key] = value
	return rb
}

// Headers 批量设置请求头
func (rb *requestBuild) Headers(headers map[string]string) *requestBuild {
	rb.request.Headers = headers
	return rb
}

// AllowRedirect 允许重定向
func (rb *requestBuild) AllowRedirect(allow bool) *Request {
	rb.request.AllowRedirect = allow
	return rb.request
}

// GetRequest 获取请求
func (rb *requestBuild) GetRequest() *Request {
	return rb.request
}

// Send 发送请求
func (rb *requestBuild) Send() (*http.Response, error) {
	return rb.SendRequest(rb.request)
}

// SendRequest 发送请求
func (rb *requestBuild) SendRequest(req *Request) (*http.Response, error) {
	var (
		resp *http.Response
	)

	params := ""
	if req.Params != nil {
		for k, v := range req.Params {
			params = fmt.Sprintf("%s&%s=%s", params, k, v)
		}
	}

	if req.Method == "POST" && req.Body != "" {
		params = req.Body
		if _, ok := req.Headers["Content-Type"]; !ok {
			req.Headers["Content-Type"] = "application/json"
		}
		if _, ok := req.Headers["Accept"]; !ok {
			req.Headers["Accept"] = "application/json"
		}
	} else if req.Method == "POST" {
		if _, ok := req.Headers["Content-Type"]; !ok {
			req.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	request, err := http.NewRequest(req.Method, req.Url, strings.NewReader(params))
	if err != nil {
		return resp, err
	}

	log.Println(request.Method, request.URL)
	for k, v := range req.Headers {
		request.Header.Set(k, v)
	}

	if req.AllowRedirect {
		c := http.Transport{}
		resp, err = c.RoundTrip(request)
	} else {
		c := http.Client{}
		resp, err = c.Do(request)
	}

	return resp, err
}

// SetHeader 设置响应的请求头
func (res *Response) SetHeader(key string, value string) {
	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}
	res.Headers[key] = value
}
