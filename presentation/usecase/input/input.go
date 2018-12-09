package input

import "github.com/miliya612/webauthn-demo/webauthnif"

type RegistrationInit struct {
	// ID is a identifier
	ID string `json:"id"`
	// DisplayName is a human-palatable name for the user account, intended only for display.
	DisplayName string `json:"displayName"`
}

type Registration struct {
	Body webauthnif.PublicKeyCredential
	Challenge []byte
}
