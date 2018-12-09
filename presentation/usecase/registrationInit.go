package usecase

import (
	"github.com/miliya612/webauthn-demo/domain/service"
	"github.com/miliya612/webauthn-demo/presentation/usecase/input"
	"github.com/miliya612/webauthn-demo/presentation/usecase/output"
	"github.com/miliya612/webauthn-demo/webauthnif"
)

type RegistrationInitUseCase interface {
	RegistrationInit(input input.RegistrationInit) (*output.RegistrationInit, error)
}

type registrationInitUseCase struct {
	service service.RegistrationService
}

func NewRegistrationInitUseCase(service service.RegistrationService) (RegistrationInitUseCase) {
	return &registrationInitUseCase{service: service}
}

func (uc registrationInitUseCase) RegistrationInit(input input.RegistrationInit) (*output.RegistrationInit, error) {
	options, err := uc.service.GetOptions(input.ID, input.DisplayName)
	if err != nil {
		return nil, err
	}
	u := options.PublicKey.User
	err = uc.service.ReserveClientInfo(u.ID, u.Name, u.DisplayName, u.Icon)
	if err != nil {
		return nil, err
	}
	return &output.RegistrationInit{
		CredentialCreationOptions: webauthnif.CredentialCreationOptions{
			PublicKey: options.PublicKey,
		},
	}, nil
}
