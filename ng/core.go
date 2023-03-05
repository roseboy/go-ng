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

type server struct {
	pluginList            []*plugin
	pluginURIMappingCache sync.Map
}

// NewServer new ng server
func NewServer() *server {
	return &server{
		pluginList:            make([]*plugin, 0),
		pluginURIMappingCache: sync.Map{},
	}
}

// Start start a ng server
func (s *server) Start(port int) error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	http.HandleFunc("/", s.httpHandler)
	log.Printf("ng server started on port: %d", port)
	err := srv.ListenAndServe()
	return err
}

// httpHandler http handle
func (s *server) httpHandler(rw http.ResponseWriter, request *http.Request) {
	plugins := s.getPluginByRequest(request)
	if len(plugins) == 0 {
		rw.WriteHeader(404)
		_, err := fmt.Fprint(rw, "404 page not found")
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
		rw.WriteHeader(500)
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

func doInterceptor(req *Request, resp *Response, plg *plugin) error {
	req.pluginPos++
	if plg.proxyPass != "" {
		req.Url = fmt.Sprintf("%s%s", strings.TrimSuffix(plg.proxyPass, "/"), req.HttpRequest.RequestURI)
	}
	return plg.interceptor(req, resp)
}

// Invoke invoke
func Invoke(req *Request, resp *Response) error {
	if req.pluginPos < len(req.plugins) {
		return doInterceptor(req, resp, req.plugins[req.pluginPos])
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
