package registry

import (
	"database/sql"
	"github.com/miliya612/webauthn-demo/domain/repo"
	"github.com/miliya612/webauthn-demo/domain/service"
	"github.com/miliya612/webauthn-demo/infra/persistance/pg"
	"github.com/miliya612/webauthn-demo/presentation/handler"
	"github.com/miliya612/webauthn-demo/presentation/usecase"
)

type Registration struct{}

type Registerer interface {
	InjectDB() *sql.DB
	InjectTodoRepo() repo.TodoRepo
	InjectTodoService() service.TodoService
	InjectTodoHandler() handler.CredentialHandler
}

func (r *Registration) RegisterDB() *sql.DB {
	db, err := sql.Open("postgres", "user=todoapp dbname=todoapp password=todopass sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

func (r *Registration) RegisterCredentialRepo() repo.CredentialRepo {
	return pg.NewCredentialRepo(r.RegisterDB())
}

func (r *Registration) RegisterUserRepo() repo.UserRepo {
	return pg.NewUserRepo(r.RegisterDB())
}

func (r *Registration) RegisterSessionRepo() repo.SessionRepo {
	return pg.NewSessionRepo(r.RegisterDB())
}

func (r *Registration) RegisterCredentialService() service.RegistrationService {
	return service.NewRegistrationService(
		r.RegisterCredentialRepo(),
		r.RegisterUserRepo(),
		r.RegisterSessionRepo(),
	)
}

func (r *Registration) RegisterSessionService() service.SessionService {
	return service.NewSessionService(r.RegisterSessionRepo())
}

func (r *Registration) RegisterCredentialInitUsecase() usecase.RegistrationInitUseCase {
	return usecase.NewRegistrationInitUseCase(
		r.RegisterCredentialService(),
		r.RegisterSessionService(),
	)
}

func (r *Registration) RegisterCredentialRegisterUsecase() usecase.RegistrationUseCase {
	return usecase.NewRegistrationUseCase(
		r.RegisterCredentialService(),
		r.RegisterSessionService(),
	)
}

func (r *Registration) RegisterCredentialHandler() handler.AppHandler {
	return handler.NewCredentialHandler(
		r.RegisterCredentialInitUsecase(),
		r.RegisterCredentialRegisterUsecase(),
	)
}
