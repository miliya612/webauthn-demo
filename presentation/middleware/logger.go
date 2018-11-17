package middleware

import (
	"github.com/miliya612/webauthn-demo/presentation/handler"
	"github.com/miliya612/webauthn-demo/presentation/httputil"
	"log"
	"net/http"
	"time"
)

func logger(status int, method, uri, name string, start time.Time) {
	log.Printf("%d\t%s\t%s\t%s\t%s", status, method, uri, name, time.Since(start))
}

func logging(h handler.APIHandleFunc, name string) handler.APIHandleFunc {
	return func(r *http.Request) httputil.Responder {
		start := time.Now()
		result := h(r)
		logger(result.Status(), r.Method, r.URL.EscapedPath(), name, start)
		return result
	}
}

func Logging(f handler.APIHandleFunc, name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fun := logging(f, name)
		result := fun(r)
		result.Write(w)
	}
}