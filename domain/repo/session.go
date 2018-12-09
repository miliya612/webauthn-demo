package repo

import "github.com/miliya612/webauthn-demo/domain/model"

type SessionRepo interface {
	GetByID(id string) (*model.Session, error)
	Create(session model.Session) (*model.Session, error)
	Delete(id string) (int, error)
}
