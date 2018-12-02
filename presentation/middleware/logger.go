package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(f http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		f.ServeHTTP(w, r)
		t2 := time.Now()
		t := t2.Sub(t1)
		log.Printf("[%s] %s %s", r.Method, r.URL, t.String())
	})
}