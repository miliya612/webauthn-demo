package routes

import (
	"github.com/gorilla/mux"
	"github.com/miliya612/webauthn-demo/presentation/handler"
	mw "github.com/miliya612/webauthn-demo/presentation/middleware"
)

func NewRouter(app handler.AppHandler) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range getRoutes(app) {
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			HandlerFunc(mw.AccessControl(mw.Logging(route.HandlerFunc)))
	}
	return router
}
