package ng

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	HttpCodeServerError = 500
	HttpCodeNotFound    = 404
	HttpCodeNormal      = 200

	HttpCodeNotFoundText = "404 page not found"
)

type server struct {
	Port                  int
	CertFile, KeyFile     string
	pluginList            []*plugin
	cacheMappingPluginURL sync.Map
}

// NewServer new ng server
func NewServer(port int) *server {
	return &server{
		Port:                  port,
		pluginList:            make([]*plugin, 0),
		cacheMappingPluginURL: sync.Map{},
	}
}

// Start start a ng server
func (s *server) Start() (err error) {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", s.Port)}
	http.HandleFunc("/", s.httpHandler)
	log.Printf("ng server started on port: %d", s.Port)
	if s.CertFile != "" && s.KeyFile != "" {
		err = srv.ListenAndServeTLS(s.CertFile, s.KeyFile)
	} else {
		err = srv.ListenAndServe()
	}
	return err
}

// httpHandler http handle
func (s *server) httpHandler(rw http.ResponseWriter, request *http.Request) {
	plugins := s.getPluginByRequest(request)
	if len(plugins) == 0 {
		rw.WriteHeader(HttpCodeNotFound)
		_, err := fmt.Fprint(rw, HttpCodeNotFoundText)
		if err != nil {
			panic(err)
		}
		return
	}

	resp := &Response{ResponseWriter: rw}
	req := newRequest().HttpRequest(request).GetRequest()
	req.plugins = plugins

	err := doInterceptor(req, resp, req.plugins[req.pluginPos])
	if err != nil {
		rw.WriteHeader(HttpCodeServerError)
		_, err = fmt.Fprint(rw, err.Error())
		if err != nil {
			panic(err)
		}
		return
	}

	for k, v := range resp.Headers {
		rw.Header().Set(k, v)
	}
	if resp.Status > 0 {
		rw.WriteHeader(resp.Status)
	}
	_, err = fmt.Fprint(rw, resp.Body)
	if err != nil {
		panic(err)
	}
}

func doInterceptor(req *Request, resp *Response, plg *pluginWrapper) error {
	if plg.proxyPass != "" {
		req.Url = fmt.Sprintf("%s%s", strings.TrimSuffix(plg.proxyPass, "/"), req.HttpRequest.RequestURI)
	}
	return plg.plugin.interceptor(req, resp)
}

// Invoke invoke
func Invoke(req *Request, resp *Response) error {
	req.pluginPos++
	if req.pluginPos < len(req.plugins) {
		return doInterceptor(req, resp, req.plugins[req.pluginPos])
	}

	if len(req.Url) == 0 {
		resp.Status = HttpCodeNotFound
		resp.Body = HttpCodeNotFoundText
		return nil
	}

	response, err := newRequest().SendRequest(req)
	if err != nil {
		return err
	}

	header := make(map[string]string)
	for k := range response.Header {
		v := response.Header.Get(k)
		header[k] = v
	}
	for k, v := range resp.Headers {
		header[k] = v
	}

	defer func() {
		_ = response.Body.Close()
	}()
	body := response.Body
	if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
		body, err = gzip.NewReader(body)
		if err != nil {
			return err
		}
	}

	bodyData, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	resp.HttpResponse = response
	resp.Headers = header
	resp.Body = string(bodyData)
	resp.Status = response.StatusCode

	return nil
}
