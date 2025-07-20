package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/dika22/news-service/package/config"

	"github.com/segmentio/encoding/json"
)

const (
	HTTPClientTest = iota
	HTTPClientMailgun
	HTTPClientCloudflare
	HTTPClientSlack
)

type (
	HTTPResponseFailure struct {
		HTTPResponse
		Message    string
		StatusCode int
	}
	HTTPHeader struct {
		Key   string
		Value string
	}
	HTTPResponse struct {
		Body       interface{}
		Cookie     []*http.Cookie
		Header     http.Header
		StatusCode int
	}
	HTTPClient struct {
		baseURL      string
		Client       http.Client
		headers      []HTTPHeader
		Req          *http.Request
		WithDumpHTTP bool
	}
)

func (f HTTPResponseFailure) Error() string {
	return fmt.Sprintf("Failure Response (%v), Body: %v", f.StatusCode, f.Body)
}

func NewHTTPClient(httpClientType int, c *config.Config) HTTPClient {
	httpClient := HTTPClient{
		WithDumpHTTP: c.DebugHTTP == "true",
	}
	switch httpClientType {
	case HTTPClientTest:
		httpClient.headers = []HTTPHeader{}
	default:
		panic("unknown http Client type")
	}
	httpClient.Client = http.Client{
		Timeout: time.Second * 5,
	}
	return httpClient
}

func (h HTTPClient) WithHeader(headers []HTTPHeader) HTTPClient {
	if h.Req != nil {
		for _, header := range headers {
			h.Req.Header.Set(header.Key, header.Value)
		}
	}
	h.headers = append(h.headers, headers...)
	return h
}

func (h HTTPClient) WithCookies(cookies []*http.Cookie) HTTPClient {
	for _, cookie := range cookies {
		h.Req.AddCookie(cookie)
	}
	return h
}

func (h HTTPClient) PrepareRequestFormData(ctx context.Context, body interface{}, method, uri string) HTTPClient {
	v := reflect.Indirect(reflect.ValueOf(body))
	t := v.Type()
	formData := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		formData[t.Field(i).Tag.Get("form")] = []string{v.Field(i).String()}
	}
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%v%v", h.baseURL, uri),
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		log.Error(err)
	}
	h.Req = req
	return h
}

func (h HTTPClient) PrepareRequestJSON(ctx context.Context, body interface{}, method, uri string) HTTPClient {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return HTTPClient{}
	}
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%v%v", h.baseURL, uri),
		bytes.NewReader(bodyByte),
	)
	if err != nil {
		log.Error(err)
	}
	h.Req = req
	return h
}

func (h HTTPClient) Do(dest interface{}) (HTTPResponse, error) {
	httpResp := HTTPResponse{}
	for _, header := range h.headers {
		h.Req.Header.Set(header.Key, header.Value)
	}
	resp, err := h.Client.Do(h.Req)
	if err != nil {
		return httpResp, err
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return httpResp, err
	}
	if h.WithDumpHTTP {
		dumpReq, err := httputil.DumpRequest(h.Req, true)
		log.Printf("Req: %+v err:%v", string(dumpReq), err)
		dumpRes, err := httputil.DumpResponse(resp, true)
		log.Printf("Res: %+v err:%v", string(dumpRes), err)
	}
	if dest == nil {
		return httpResp, nil
	}
	if err = json.Unmarshal(respBody, dest); err != nil {
		return httpResp, err
	}
	httpResp.Body = dest
	httpResp.Cookie = resp.Cookies()
	httpResp.Header = resp.Header
	httpResp.StatusCode = resp.StatusCode
	if resp.StatusCode >= 400 {
		return httpResp, &HTTPResponseFailure{
			HTTPResponse: httpResp,
			StatusCode:   httpResp.StatusCode,
		}
	}
	return httpResp, nil
}
