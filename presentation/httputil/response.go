package httputil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Responder interface {
	Write(w http.ResponseWriter)
	Status() int
	Body() []byte
}

type Response struct {
	status int
	header http.Header
	body   []byte
}

func (r *Response) Write(w http.ResponseWriter) {
	header := w.Header()
	for k, v := range r.header {
		header[k] = v
	}
	w.WriteHeader(r.status)
	w.Write(r.body)
}

func (r Response) Status() int {
	return r.status
}

func (r Response) Body() []byte {
	return r.body
}

func (r Response) Header(k, v string) *Response {
	r.header.Set(k, v)
	return &r
}

func Empty(status int) *Response {
	if isStatusCodeErrors(status) {
		panic("status code is not 4xx or 5xx: Error() should be used for returning error response")
	}
	return respond(status, nil)
}

func jsonBody(status int, body interface{}) *Response {
	return respond(status, body).Header("Content-Type", "application/json; charset=UTF-8")
}

func Ok(body interface{}) *Response {
	return jsonBody(http.StatusOK, body)
}

func Created(body interface{}, location string) *Response {
	return jsonBody(http.StatusCreated, body).Header("Location", location)
}

func Error(status int, message string, err error) *Response {
	log.Printf("[ERROR]\t%s, %s", message, err)
	msg := createErrMsg(status, fmt.Sprintf("%s: %s", message, err))
	if !isStatusCodeErrors(status) {
		panic("status code is not 4xx or 5xx")
	}
	return jsonBody(status, msg)
}

func respond(status int, body interface{}) *Response {
	var b []byte
	var err error
	switch t := body.(type) {
	case string:
		b = []byte(t)
	default:
		if b, err = json.Marshal(body); err != nil {
			return Error(http.StatusInternalServerError, "failed marshalling json", err)
		}
	}

	return &Response{
		status: status,
		body:   b,
		header: make(http.Header),
	}
}

type errResp struct {
	Error string `json:"error"`
	Message string `json:"message"`
}

 func createErrMsg(code int, msg string) errResp {
 	return errResp{
 		Error: http.StatusText(code),
 		Message: msg,
	}
 }

 func isStatusCodeErrors(code int) bool {
 	return code / 100 == 4 || code / 100 == 5
 }