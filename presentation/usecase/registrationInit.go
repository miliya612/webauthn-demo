package usecase

import (
	"context"
	"github.com/miliya612/webauthn-demo/domain/service"
	"github.com/miliya612/webauthn-demo/presentation/httputil"
	"github.com/miliya612/webauthn-demo/presentation/usecase/input"
	"github.com/miliya612/webauthn-demo/presentation/usecase/output"
	"github.com/miliya612/webauthn-demo/webauthnif"
)

type RegistrationInitUseCase interface {
	RegistrationInit(ctx context.Context, input input.RegistrationInit) (*output.RegistrationInit, error)
}

type registrationInitUseCase struct {
	registration service.RegistrationService
	session      service.SessionService
}

func NewRegistrationInitUseCase(registration service.RegistrationService, session service.SessionService,
) RegistrationInitUseCase {
	return &registrationInitUseCase{
		registration: registration,
		session:      session,
	}
}

func (uc registrationInitUseCase) RegistrationInit(ctx context.Context, input input.RegistrationInit,
) (*output.RegistrationInit, error) {
	options, err := uc.registration.GetOptions(input.ID, input.DisplayName)
	if err != nil {
		return nil, err
	}
	u := options.PublicKey.User
	err = uc.registration.ReserveClientInfo(u.ID, u.Name, u.DisplayName, u.Icon)
	if err != nil {
		return nil, err
	}

	rawSid := ctx.Value(httputil.KeySessionID)
	sid := rawSid.(string)
	err = uc.session.Store(sid, u.ID, options.PublicKey.Challenge)
	if err != nil {
		return nil, err
	}

	return &output.RegistrationInit{
		CredentialCreationOptions: webauthnif.CredentialCreationOptions{
			PublicKey: options.PublicKey,
		},
	}, nil
}
