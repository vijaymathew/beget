package beget

import (
	"fmt"
)

type HttpResponse struct {
	Status string // e.g. 200 OK
	StatusCode int // e.g. 200
	Header map[string]string
	Body string
}

func Get(url string) (response HttpResponse) {
	fmt.Printf("GET %s\n", url)
	response.Status = "200 OK"
	response.StatusCode = 200
	// TODO: fetch the contents of url, put that into
	// and HttpResponse and return it.
	return
}
