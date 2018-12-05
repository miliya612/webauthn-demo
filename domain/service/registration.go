package service

import "C"
import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/miliya612/webauthn-demo/domain/repo"
	"github.com/miliya612/webauthn-demo/webauthnif"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
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
	CLIENTDATATYPE string = "webauthn.create"
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

// 7.1.
// Registering a new credential
// When registering a new credential, represented by an AuthenticatorAttestationResponse structure response and an
// AuthenticationExtensionsClientOutputs structure clientExtensionResults, as part of a registration ceremony, a Relying
// Party MUST proceed as follows:
func (s registrationService) Register(req webauthnif.AuthenticatorAttestationResponse) (error) {
	decodedReq, err := parseAttestationResponse(req)
	if err != nil {
		return err
	}


	c := decodedReq.ClientData

	// 3. Verify that the value of C.type is webauthn.create.
	if c.Type != CLIENTDATATYPE {
		errMsg := fmt.Sprintf("got %q, but %q is required", c.Type, CLIENTDATATYPE)
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 4. Verify that the value of C.challenge matches the challenge that was sent to the authenticator in the create()
	// call.
	var rawChallenge interface{} = "TODO: GET FROM USERS SESSION"
	bytesArrayChallenge, err := rawChallenge.([]byte)
	if err != nil {
		return err
	}
	orgnChallenge := (webauthnif.BufferSource)(bytesArrayChallenge)

	if !c.Challenge.Equals(orgnChallenge) {
		errMsg := "invalid challenge"
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 5. Verify that the value of C.origin matches the Relying Party's origin.
	// TODO: RPIDマッチのロジック必要そう。subDomainとrootDomainの判定処理も書く
	if c.Origin != RPID {
		errMsg := "invalid origin"
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 6. Verify that the value of C.tokenBinding.status matches the state of Token Binding for the TLS connection over
	// which the assertion was obtained. If Token Binding was used on that TLS connection, also verify that
	// C.tokenBinding.id matches the base64url encoding of the Token Binding ID for the connection.
	// TODO: やる

	// 7. Compute the hash of response.clientDataJSON using SHA-256.
	hashedClientDataJSON := sha256.Sum256(req.ClientDataJSON)

	// 8. Perform CBOR decoding on the attestationObject field of the AuthenticatorAttestationResponse structure to
	// obtain the attestation statement format fmt, the authenticator data authData, and the attestation statement
	// attStmt.
	handle := codec.CborHandle{}
	decoder := codec.NewDecoder(bytes.NewReader(req.AttestationObject), &handle).Decode()

}

func parseAttestationResponse(req webauthnif.AuthenticatorAttestationResponse) (
	*webauthnif.DecodedAuthenticatorResponse, error) {

	// 1. Let JSONtext be the result of running UTF-8 decode on the value of response.clientDataJSON.
	// 2. Let C, the client data claimed as collected during the credential creation, be the result of running an
	// implementation-specific JSON parser on JSONtext.
	decodedAttestationRes := webauthnif.DecodedAuthenticatorResponse{}
	c := decodedAttestationRes.ClientData
	if err := json.Unmarshal(req.ClientDataJSON, &c); err != nil {
		return nil, errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", err))
	}

	return &decodedAttestationRes, nil
}
