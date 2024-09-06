package ng

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type server struct {
	Port                  int
	httpServer            *http.Server
	certFile, keyFile     string
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

// NewServerWithHTTPServer new
func NewServerWithHTTPServer(httpServer *http.Server, port int) *server {
	return &server{
		Port:                  port,
		httpServer:            httpServer,
		pluginList:            make([]*plugin, 0),
		cacheMappingPluginURL: sync.Map{},
	}
}

// Start a ng server
func (s *server) Start() (err error) {
	if s.httpServer != nil {
		s.httpServer.Addr = fmt.Sprintf(":%d", s.Port)
		s.httpServer.Handler = http.HandlerFunc(s.httpHandler)
		log.Printf("ng server started on port: %d", s.Port)
		if s.certFile != "" && s.keyFile != "" {
			err = s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile)
		} else {
			err = s.httpServer.ListenAndServe()
		}
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpHandler)
	addr := fmt.Sprintf(":%d", s.Port)
	log.Printf("ng server started on port: %d", s.Port)
	if s.certFile != "" && s.keyFile != "" {
		err = http.ListenAndServeTLS(addr, s.certFile, s.keyFile, mux)
	} else {
		err = http.ListenAndServe(addr, mux)
	}
	return err
}

// WithTLS set tls
func (s *server) WithTLS(certFile, keyFile string) *server {
	s.certFile, s.keyFile = certFile, keyFile
	return s
}

// httpHandler http handle
func (s *server) httpHandler(rw http.ResponseWriter, request *http.Request) {
	plugins := s.getPluginByRequest(request)
	if len(plugins) == 0 {
		http.NotFound(rw, request)
		return
	}

	resp := &Response{ResponseWriter: rw, Headers: map[string]string{}}
	req := newRequest().HttpRequest(request).GetRequest()
	req.plugins = plugins

	err := doInterceptor(req, resp, req.plugins[req.pluginPos])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Headers {
		rw.Header().Set(k, v)
	}
	if resp.Status > 0 {
		rw.WriteHeader(resp.Status)
	}
	_, err = rw.Write([]byte(resp.Body))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func doInterceptor(req *Request, resp *Response, plg *pluginWrapper) error {
	if plg.proxyPass == "" {
		return plg.plugin.interceptor(req, resp)
	}

	if strings.HasSuffix(plg.proxyPass, "/") {
		req.Url = fmt.Sprintf("%s%s", plg.proxyPass,
			strings.TrimLeft(req.HttpRequest.RequestURI, plg.location))
	} else {
		req.Url = fmt.Sprintf("%s%s", plg.proxyPass, req.HttpRequest.RequestURI)
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
		return errors.New("proxy pass url not found")
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
