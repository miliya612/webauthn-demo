package middleware

import "net/http"

func AccessControl(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		//w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		//w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
