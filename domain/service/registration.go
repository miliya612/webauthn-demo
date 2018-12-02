package service

import (
	"github.com/miliya612/webauthn-demo/domain/repo"
	"github.com/miliya612/webauthn-demo/webauthnif"
)

type RegistrationService interface {
	GetOptions(id, displayName string) (*webauthnif.CredentialCreationOptions, error)
	Register()
}

type registrationService struct {
	repo repo.CredentialRepo
}

func NewRegistrationService(repository repo.CredentialRepo) RegistrationService {
	return &registrationService{repo: repository}
}

const (
	RPID           string = "miliya.tech"
	RPNAME         string = "miliya612 - webauthn demo"
	TIMEOUTMILLSEC uint32 = 6000
)

// GetOptions returns CredentialCreationOptions to client. It will be used when calling navigator.credentials.create().
// Parameters:
//   - id:          REQUIRED. This param identifies user who will register a credential to RP.
//   - displayName: OPTIONAL. This params is intended to be shown to users.
func (s registrationService) GetOptions(id, displayName string) (*webauthnif.CredentialCreationOptions, error) {

	rp := &webauthnif.PublicKeyCredentialRpEntity{
		// rpidのscopeを指定した場合はここで指定
		// defaultではsubdomainつきのFQDNとか?
		ID: RPID,
		PublicKeyCredentialEntity: webauthnif.PublicKeyCredentialEntity{
			// TODO: これheaderのtitleから取った方が良い?
			Name: RPNAME,
		},
	}

	user := &webauthnif.PublicKeyCredentialUserEntity{
		ID:          webauthnif.ToBufferSource(id),
		DisplayName: displayName,
		PublicKeyCredentialEntity: webauthnif.PublicKeyCredentialEntity{
			// TODO: name何入れよう
			Name: id,
		},
	}

	challenge, err := webauthnif.GenChallenge()
	if err != nil {
		return nil, err
	}

	credentialParams := &webauthnif.PublicKeyCredentialParameters{
		webauthnif.PublicKeyCredentialParameter{
			Type: webauthnif.PublicKeyCredentialTypePublicKey,
			Alg:  webauthnif.COSEAlgorithmIdentifierES256,
		},
		webauthnif.PublicKeyCredentialParameter{
			Type: webauthnif.PublicKeyCredentialTypePublicKey,
			Alg:  webauthnif.COSEAlgorithmIdentifierRS256,
		},
	}

	excludeCredentials := webauthnif.PublicKeyCredentialDescriptors{
		webauthnif.PublicKeyCredentialDescriptor{
			Type: webauthnif.PublicKeyCredentialTypePublicKey,
			// TODO: ここ何いれる？
			// authenticator内で二重にcredentialIDを登録しないために使う
			// https://www.w3.org/TR/webauthn/#op-make-cred
			ID: webauthnif.ToBufferSource(id),
			Transports: webauthnif.AuthenticatorTransports{
				webauthnif.AuthenticatorTransportUSB,
				webauthnif.AuthenticatorTransportInternal,
			},
		},
	}

	authenticatorSelection := webauthnif.AuthenticatorSelectionCriteria{
		AuthenticatorAttachment: webauthnif.AuthenticatorAttachmentEmpty,
		RequireResidentKey:      false, // default
		UserVerification:        webauthnif.UserVerificationRequirementEmpty,
	}

	acp := webauthnif.AttestationConveyancePreferenceEmpty

	extensions := webauthnif.AuthenticationExtensionsClientInputs{}

	pkoptions := &webauthnif.PublicKeyCredentialCreationOptions{
		RP:                     *rp,
		User:                   *user,
		Challenge:              challenge,
		PubKeyCredParams:       *credentialParams,
		Timeout:                TIMEOUTMILLSEC,
		ExcludeCredentials:     excludeCredentials,
		AuthenticatorSelection: authenticatorSelection,
		Attestation:            acp,
		Extensions:             extensions,
	}

	options := &webauthnif.CredentialCreationOptions{
		PublicKey: *pkoptions,
	}
	return options, nil
}

func (s registrationService) Register() {}
