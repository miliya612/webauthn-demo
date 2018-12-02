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

type response struct {
	status int
	header http.Header
	body   []byte
}

func (r *response) Write(w http.ResponseWriter) {
	header := w.Header()
	for k, v := range r.header {
		header[k] = v
	}
	w.WriteHeader(r.status)
	if _, err := w.Write(r.body); err != nil {
		panic(err)
	}
}

func (r response) Status() int {
	return r.status
}

func (r response) Body() []byte {
	return r.body
}

func (r *response) Header(k, v string) *response {
	r.header.Set(k, v)
	return r
}

func Empty(w http.ResponseWriter, status int) {
	if isStatusCodeErrors(status) {
		panic("status code is not 4xx or 5xx: Error() should be used for returning error response")
	}
	res := &response{
		status: status,
		header: make(http.Header),
	}
	res.respond(w, nil)
}

func Ok(w http.ResponseWriter, body interface{}) {
	res := &response{
		status: http.StatusOK,
		header: make(http.Header),
	}
	res.jsonBody(w, body)
}

func Created(w http.ResponseWriter, body interface{}) {
	r := response{
		status: http.StatusCreated,
		header: make(http.Header),
	}
	r.jsonBody(w, body)
}

func Accepted(w http.ResponseWriter, body interface{}) {
	r := response{
		status: http.StatusAccepted,
		header: make(http.Header),
	}
	r.jsonBody(w, body)
}

func Error(w http.ResponseWriter, status int, message string, err error) {
	log.Printf("[ERROR]\t%s, %s", message, err)
	msg := createErrMsg(status, fmt.Sprintf("%s: %s", message, err))
	if !isStatusCodeErrors(status) {
		panic("status code is not 4xx or 5xx")
	}
	res := &response{
		status: status,
		header: make(http.Header),
	}
	res.jsonBody(w, msg)
}

func (r *response) jsonBody(w http.ResponseWriter, body interface{}) {
	r.Header("Content-Type", "application/json; charset=UTF-8")
	r.respond(w, body)
}

func (r *response) respond(w http.ResponseWriter, body interface{}) {
	var b []byte
	var err error
	switch t := body.(type) {
	case string:
		b = []byte(t)
	default:
		if b, err = json.Marshal(body); err != nil {
			Error(w, http.StatusInternalServerError, "failed marshalling json", err)
		}
	}

	r.body = b
	r.Write(w)
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