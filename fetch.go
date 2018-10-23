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
	"io/ioutil"
	"net/http"
)

type HTTPResponse struct {
	StatusCode int // e.g -1, 200
	Status string // e.g. 200 OK
	Header map[string][]string
	Body string
}

// Get fetches a document via an HTTP GET request.
// The status, header and body information of the HTTP response
// will be returned in an HTTPResponse object.
func Get(url string) (response HTTPResponse, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
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
