package repo

import "github.com/miliya612/webauthn-demo/domain/model"

type UserRepo interface {
	GetByID(id []byte) (*model.User, error)
	Create(user model.User, chal []byte) (*model.User, error)
	Update(user model.User) (*model.User, error)
}
