package httputil

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

var KB int64 = 1024

func ParseBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*KB)) // 1MB
	if err != nil {
		return nil, errors.New("request body is too large")
	}
	defer r.Body.Close()

	//fmt.Println(r.GetBody())

	if len(body) == 0 {
		return nil, errors.New("empty body is not acceptable")
	}

	return body, nil
}
