package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/miliya612/webauthn-demo/presentation/routes"
	"github.com/miliya612/webauthn-demo/registry"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	f, err := os.OpenFile("tmp/application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening file :", err.Error())
	}
	log.SetOutput(f)
}

func main() {
	r := registry.Registration{}
	router := routes.NewRouter(r.RegisterCredentialHandler())
	corsMw := mux.CORSMethodMiddleware(router)
	router.Use(corsMw)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})


	log.Printf("server started at: %v", time.Now())
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
