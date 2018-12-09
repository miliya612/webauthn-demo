package usecase

import (
	"crypto/sha256"
	"github.com/miliya612/webauthn-demo/domain/service"
	"github.com/miliya612/webauthn-demo/presentation/usecase/input"
	"github.com/miliya612/webauthn-demo/presentation/usecase/output"
)

type RegistrationUseCase interface {
	Registration(input input.Registration) (*output.Registration, error)
}

type registrationUseCase struct {
	service service.RegistrationService
}

func NewRegistrationUseCase(service service.RegistrationService) (RegistrationUseCase)  {
	return &registrationUseCase{service: service}
}

// 7.1.
// Registering a new credential
// When registering a new credential, represented by an AuthenticatorAttestationResponse structure response and an
// AuthenticationExtensionsClientOutputs structure clientExtensionResults, as part of a registration ceremony, a Relying
// Party MUST proceed as follows:
func (uc registrationUseCase)Registration(input input.Registration) (*output.Registration, error) {
	d, err := uc.service.ParseClientData(input.Body.Response)
	if err != nil {
		return nil, err
	}

	c := d.ClientData
	err = uc.service.ValidateClientData(c)
	if err != nil {
		return nil, err
	}

	// 7. Compute the hash of response.clientDataJSON using SHA-256.
	hashedClientDataJSON := sha256.Sum256(input.Body.Response.ClientDataJSON)

	d, err = uc.service.ParseAttestationObj(input.Body.Response.AttestationObject, d)
	if err != nil {
		return nil, err
	}

	err = uc.service.ValidateAuthenticatorData(d.DecodedAttestationObject.AuthData)
	if err != nil {
		return nil, err
	}

	err = uc.service.ValidateClientExtensionOutputs(input.Body.AuthenticationExtensionsClientOutputs)
	if err != nil {
		return nil, err
	}

	err = uc.service.ValidateAttestationResponse(d.DecodedAttestationObject, hashedClientDataJSON)
	if err != nil {
		return nil, err
	}

	err = uc.service.Register([]byte{}, d.DecodedAttestationObject.AuthData.AttestedCredentialData)
	if err != nil {
		return nil, err
	}

	return nil, nil
}