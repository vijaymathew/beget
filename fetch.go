package beget

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
)

type BaseResponse struct {
	StatusCode int // e.g -1, 200
	Body string
}

type HTTPResponse struct {
	BaseResponse
	Status string // e.g. 200 OK
	Header map[string][]string
}

// Get fetches a document via an HTTP GET request.
// The status, header and body information of the HTTP response
// will be returned in an HTTPResponse object.
// In case of an IO or system error, response.StatusCode will be set to -1
// and response.Body will contain a description of the error.
func Get(url string) (response HTTPResponse) {
	res, err := http.Get(url)
	response.StatusCode = -1
	if err != nil {
		log.Fatal(err)
		response.Body = fmt.Sprintf("%s", err)
		return
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	response.Body = string(bytes)
	if err != nil {
		log.Fatal(err)
		response.Body = fmt.Sprintf("%s", err)
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
