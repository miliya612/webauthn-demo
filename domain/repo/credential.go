package repo

import "github.com/miliya612/webauthn-demo/domain/model"

type CredentialRepo interface {
	GetByCredentialID(id []byte) (*model.Credential, error)
	Create(credential model.Credential) (*model.Credential, error)
	Update(credential model.Credential) (*model.Credential, error)
	Delete(id []byte) ([]byte,error)
	// GetCount(id []byte) (int, error)
	// UpdateCount(id []byte, count id) (*model.Credential, error)
}
