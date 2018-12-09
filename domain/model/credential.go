package model

type Credential struct {
	CredentialID []byte
	UserID []byte
	PublicKey []byte
	SignCount uint32
}
