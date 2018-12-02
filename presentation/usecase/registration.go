package usecase

import (
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


func (uc registrationUseCase)Registration(input input.Registration) (*output.Registration, error){
	//options, err := uc.service.Register()
	//if err != nil {
	//	return nil, err
	//}
	return &output.Registration{
		//Options: *options,
	}, nil
}