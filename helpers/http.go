package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpApi struct {
}

type HttpRequest struct {
	Url         string
	Data        string
	ContentType string
	TimeOut     int
}

type NewHttpRequest struct {
	HttpRequest
	Method     string
	Headers    map[string]string
	RetryTimes int
	GetParams  map[string]string
}

//发送post请求 - 无重试
func (httpApi *HttpApi) SendPostRequest(httpRequest *HttpRequest) ([]byte, error) {
	var client = new(http.Client)
	client.Timeout = 10 * time.Second
	if httpRequest.TimeOut > 0 {
		client.Timeout = time.Duration(1000000000 * httpRequest.TimeOut)
	}
	var req, err = http.NewRequest("POST", httpRequest.Url, strings.NewReader(httpRequest.Data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if httpRequest.ContentType != "" {
		req.Header.Set("Content-Type", httpRequest.ContentType)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//发送post请求 - 含EOF错误重试机制
func (httpApi *HttpApi) SendPostRequestRetry(httpRequest *HttpRequest) ([]byte, error) {
	//重试次数为3
	var respBytes []byte
	var err error
	for i := 1; i <= 3; i++ {
		respBytes, err = httpApi.SendPostRequest(httpRequest)
		if err != nil {
			b := strings.Contains(err.Error(), "EOF")
			if !b {
				break
			}
		} else {
			break
		}
	}
	return respBytes, err
}

//发送get请求 - 含EOF错误重试机制
func (httpApi *HttpApi) SendGetRequestRetry(httpRequest *HttpRequest) ([]byte, error) {
	//重试次数为3
	var respBytes []byte
	var err error
	for i := 1; i <= 3; i++ {
		resp, err := http.Get(httpRequest.Url)
		if err != nil {
			b := strings.Contains(err.Error(), "EOF")
			if !b {
				break
			}
		}
		defer resp.Body.Close()
		respBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			b := strings.Contains(err.Error(), "EOF")
			if !b {
				break
			}
		} else {
			break
		}
	}
	return respBytes, err
}

/**
 * 发送http请求，包含重试机制
 */
func (httpApi *HttpApi) SendRequest (request *NewHttpRequest) ([]byte, error) {
	var err error
	var req *http.Request
	var resultBytes []byte
	for request.RetryTimes >= 0 {
		if request.Method == http.MethodGet {
			req, err = http.NewRequest(request.Method, request.Url, nil)
		} else {
			req, err = http.NewRequest(request.Method, request.Url, bytes.NewBuffer([]byte(request.Data)))
		}
		if err != nil {
			request.RetryTimes--
			if request.RetryTimes < 0 {
				return resultBytes, err
			}
			continue
		}
		if request.Method == http.MethodGet && len(request.GetParams) > 0{
			q := req.URL.Query()
			for getKey, getValue := range request.GetParams {
				q.Add(getKey, getValue)
			}
			req.URL.RawQuery = q.Encode()
		}
		for headerKey, headerValue := range request.Headers {
			req.Header.Add(headerKey, headerValue)
		}
		var timeout = 10
		if request.TimeOut > 0 {
			timeout = request.TimeOut
		}
		client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			request.RetryTimes--
			if request.RetryTimes < 0 {
				if req.Body != nil {
					_ = req.Body.Close()
				}
				return resultBytes, err
			}
			continue
		}
		resultBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			request.RetryTimes--
			if request.RetryTimes < 0 {
				if req.Body != nil {
					_ = req.Body.Close()
				}
				if resp.Body != nil {
					_ = resp.Body.Close()
				}
				return resultBytes, err
			}
			continue
		}
		if req.Body != nil {
			_ = req.Body.Close()
		}
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
		return resultBytes, nil
	}
	return resultBytes, nil
}
