package routes

import (
	"github.com/miliya612/webauthn-demo/presentation/handler"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func getRoutes(app handler.AppHandler) []Route {
	return Routes{
		Route{"RegistrationInit", "POST", "/credential/register/init", app.RegistrationInit},
		//Route{"Registration", "POST", "/credential/register", app.Registration},
		//Route{"AuthenticationInit", "POST", "/credential/authenticate/init", app.TodoCreate},
		//Route{"Authentication", "POST", "/credential/authenticate", app.TodoDelete},
	}
}