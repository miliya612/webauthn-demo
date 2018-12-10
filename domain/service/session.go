package service

import (
	"github.com/miliya612/webauthn-demo/domain/model"
	"github.com/miliya612/webauthn-demo/domain/repo"
	"time"
)

type SessionService interface {
	Get(sid string) (*model.Session, bool)
	Store(sid string, uid, chal []byte) error
}

type sessionService struct {
	repo repo.SessionRepo
}

func NewSessionService(session repo.SessionRepo) SessionService {
	return &sessionService{
		repo: session,
	}
}

func (s *sessionService) Get(sid string) (*model.Session, bool) {
	session, err := s.repo.GetByID(sid)
	if err != nil {
		return nil, false
	}
	return session, true
}

func (s *sessionService) Store(sid string, uid, chal []byte) error {
	session := &model.Session{
		ID:           sid,
		UserID:       uid,
		Challenge:    chal,
		LastAccessed: time.Now(),
	}
	_, err := s.repo.Create(*session)
	if err != nil {
		return err
	}
	return nil
}
