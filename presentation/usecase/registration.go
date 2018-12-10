package usecase

import (
	"context"
	"crypto/sha256"
	"github.com/miliya612/webauthn-demo/domain/service"
	"github.com/miliya612/webauthn-demo/presentation/httputil"
	"github.com/miliya612/webauthn-demo/presentation/usecase/input"
	"github.com/miliya612/webauthn-demo/presentation/usecase/output"
	"github.com/miliya612/webauthn-demo/webauthnif"
	"github.com/pkg/errors"
)

type RegistrationUseCase interface {
	Registration(ctx context.Context, input input.Registration) (*output.Registration, error)
}

type registrationUseCase struct {
	registration service.RegistrationService
	session      service.SessionService
}

func NewRegistrationUseCase(registration service.RegistrationService, session service.SessionService,
) RegistrationUseCase {
	return &registrationUseCase{
		registration: registration,
		session:      session,
	}
}

// 7.1.
// Registering a new credential
// When registering a new credential, represented by an AuthenticatorAttestationResponse structure response and an
// AuthenticationExtensionsClientOutputs structure clientExtensionResults, as part of a registration ceremony, a Relying
// Party MUST proceed as follows:
func (uc registrationUseCase) Registration(ctx context.Context, input input.Registration) (*output.Registration, error) {

	rawSid := ctx.Value(httputil.KeySessionID)
	sid := rawSid.(string)
	session, ok := uc.session.Get(sid)
	if !ok {
		return nil, errors.New("session is not set")
	}

	d := &webauthnif.DecodedAuthenticatorAttestationResponse{}

	c, err := uc.registration.ParseClientData(input.Response)
	if err != nil {
		return nil, err
	}

	d.ClientData = *c

	err = uc.registration.ValidateClientData(session.Challenge, d.ClientData)
	if err != nil {
		return nil, err
	}

	// 7. Compute the hash of response.clientDataJSON using SHA-256.
	hashedClientDataJSON := sha256.Sum256(input.Response.ClientDataJSON)

	d, err = uc.registration.ParseAttestationObj(input.Response.AttestationObject, d)
	if err != nil {
		return nil, err
	}

	err = uc.registration.ValidateAuthenticatorData(d.DecodedAttestationObject.AuthData)
	if err != nil {
		return nil, err
	}

	err = uc.registration.ValidateClientExtensionOutputs(input.AuthenticationExtensionsClientOutputs)
	if err != nil {
		return nil, err
	}

	err = uc.registration.ValidateAttestationResponse(d.DecodedAttestationObject, hashedClientDataJSON)
	if err != nil {
		return nil, err
	}

	err = uc.registration.Register(session.UserID, d.DecodedAttestationObject.AuthData.AttestedCredentialData)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
