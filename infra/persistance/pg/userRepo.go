package pg

import (
	"bytes"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/miliya612/webauthn-demo/domain/model"
	"github.com/miliya612/webauthn-demo/domain/repo"
	"github.com/pkg/errors"
)

type userRepo struct {
	db *sql.DB
}

var users []*model.User

func NewUserRepo(db *sql.DB) repo.UserRepo {
	return userRepo{db: db}
}

func (repo userRepo) GetByID(id []byte) (*model.User, error) {
	for _, u := range users {
		if bytes.Equal(id, u.ID) {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (repo userRepo) Create(user model.User) (*model.User, error) {
	users = append(users, &user)
	return &user, nil
}

func (repo userRepo) Update(user model.User) (*model.User, error) {
	return nil, nil
}
