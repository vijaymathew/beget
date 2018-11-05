// Copyright 2018, 2019 Vijay Mathew Pandyalakal<vijay.the.lisper@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beget

import (
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"net/url"
)

const defaultRedirects = 10

type HTTPRequestContext struct {
	ProxyURL string `json:"proxyURL"`
	MaxRedirects int `json:"maxRedirects"` // defaults to 10
	Jar http.CookieJar `json:"jar"`
	TimeoutSecs time.Duration `json:"timeoutSecs"` // defaults to 5secs
	Header map[string]string `json:"header"`
}

type HTTPResponse struct {
	StatusCode int `json:"statusCode"` // e.g -1, 200
	Status string `json:"status"` // e.g. 200 OK
	Header http.Header `json:"headers"`
	Jar http.CookieJar `json:"-"`
	Body string `json:"body"`
}

func NewHTTPRequestContext() (ctx HTTPRequestContext) {
	ctx.MaxRedirects = 10
	ctx.TimeoutSecs = time.Duration(5 * time.Second)
	return
}

// Get fetches a document via an HTTP GET request.
// The status, header and body information of the HTTP response
// will be returned in an HTTPResponse object.
func Get(url string, ctx *HTTPRequestContext) (response HTTPResponse, err error) {
	if ctx == nil {
		return simpleGet(url)
	} else {
		return customGet(url, ctx)
	}
}

func simpleGet(url string) (response HTTPResponse, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	return parseResponse(res)
}

func parseResponse(res *http.Response) (response HTTPResponse, err error) {
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	response.Body = string(bytes)
	if err != nil {
		err = fmt.Errorf("Get: %v", err)
		return
	}
	response.Status = res.Status
	response.StatusCode = res.StatusCode
	response.Header = make(map[string][]string)
	for key, values := range res.Header {
		newValues := make([]string, len(values))
		for i, v := range values {
			newValues[i] = v
		}
		response.Header[key] = newValues
	}
	return
}

func makeRedirectCounter(maxRedirects int) (func(req *http.Request, via []*http.Request) error) {
	counter := func(_ *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return fmt.Errorf("stopped after %d redirects", maxRedirects)
		}
		return nil
	}
	return counter
}

func noRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

func customGet(urlStr string, ctx *HTTPRequestContext) (response HTTPResponse, err error) {
	client := &http.Client{}
	if ctx.MaxRedirects != defaultRedirects {
		client.CheckRedirect = makeRedirectCounter(ctx.MaxRedirects)
	} else if ctx.MaxRedirects <= 0 { // No redirects
		client.CheckRedirect = noRedirect
	}
	if ctx.TimeoutSecs >= 0 {
		client.Timeout = time.Duration(ctx.TimeoutSecs * time.Second)
	}
	if ctx.ProxyURL != "" {
		proxyURL, e := url.Parse(ctx.ProxyURL)
		if e != nil {
			err = e
			return
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		client.Transport = transport
	}
	request, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}
	for key, value := range ctx.Header {
		request.Header.Add(key, value)
	}
	res, err := client.Do(request)
	if err != nil {
		return
	}
	return parseResponse(res)
}
