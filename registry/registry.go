package registry

import (
	"database/sql"
	"github.com/miliya612/webauthn-demo/domain/repo"
	"github.com/miliya612/webauthn-demo/domain/service"
	"github.com/miliya612/webauthn-demo/infra/persistance/datastore"
	"github.com/miliya612/webauthn-demo/presentation/handler"
)

type Registration struct {}

type Registerer interface {
	InjectDB() *sql.DB
	InjectTodoRepo() repo.TodoRepo
	InjectTodoService() service.TodoService
	InjectTodoHandler() handler.TodoHandler
}

func (r *Registration) RegisterDB() *sql.DB {
	db, err := sql.Open("postgres", "user=todoapp dbname=todoapp password=todopass sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

func (r *Registration) RegisterTodoRepo() repo.TodoRepo {
	return datastore.NewTodoRepo(r.RegisterDB())
}

func (r *Registration) RegisterTodoService() service.TodoService {
	return service.NewTodoService(r.RegisterTodoRepo())
}

func (r *Registration) RegisterTodoHandler() handler.AppHandler {
	return handler.NewTodoHandler(r.RegisterTodoService())
}