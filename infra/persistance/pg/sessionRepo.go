package pg

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/miliya612/webauthn-demo/domain/model"
	"github.com/miliya612/webauthn-demo/domain/repo"
	"github.com/pkg/errors"
	"time"
)

type sessionRepo struct {
	db *sql.DB
}

var sessions []*model.Session

func NewSessionRepo(db *sql.DB) repo.SessionRepo {
	return sessionRepo{db: db}
}

func (repo sessionRepo) GetByID(id string) (*model.Session, error) {
	for _, s := range sessions {
		if id == s.ID {
			s.LastAccessed = time.Now()
			return s, nil
		}
	}
	return nil, errors.New("session not found")
}

func (repo sessionRepo) Create(session model.Session) (*model.Session, error) {
	fmt.Println("session: ", session)
	sessions = append(sessions, &session)
	return &session, nil
}

func (repo sessionRepo) Delete(id string) (int, error) {
	for i, s := range sessions {
		if id == s.ID {
			sessions = append(sessions[:i], sessions[i+1:]...)
			return 1, nil
		}
	}
	return 0, errors.New("session not found")
}
