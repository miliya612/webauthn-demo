package output

import "github.com/miliya612/webauthn-demo/webauthnif"

type RegistrationInit struct {
	Options webauthnif.CredentialCreationOptions `json:"options"`
}
