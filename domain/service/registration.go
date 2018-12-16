package service

import "C"
import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/miliya612/webauthn-demo/domain/model"
	"github.com/miliya612/webauthn-demo/domain/repo"
	"github.com/miliya612/webauthn-demo/domain/service/attestation"
	"github.com/miliya612/webauthn-demo/webauthnif"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

type RegistrationService interface {
	GetOptions(id, displayName string) (*webauthnif.CredentialCreationOptions, error)
	ReserveClientInfo(userId []byte, name, displayName, icon string) error
	Register(userId []byte, data webauthnif.AuthenticatorData) error
	ParseClientData(req webauthnif.AuthenticatorAttestationResponse) (
		*webauthnif.CollectedClientData, error)
	ValidateClientData(rawChal []byte, c webauthnif.CollectedClientData) error
	ParseAttestationObj(
		req []byte,
		d *webauthnif.DecodedAuthenticatorAttestationResponse,
	) (*webauthnif.DecodedAuthenticatorAttestationResponse, error)
	ValidateClientExtensionOutputs(outputs webauthnif.AuthenticationExtensionsClientOutputs) error
	ValidateAttestationResponse(attObj webauthnif.DecodedAttestationObject, hashedClientData [32]byte) error
	ValidateAuthenticatorData(data webauthnif.AuthenticatorData) error
}

type registrationService struct {
	credentialRepo repo.CredentialRepo
	userRepo       repo.UserRepo
}

func NewRegistrationService(
	credential repo.CredentialRepo, user repo.UserRepo, session repo.SessionRepo) RegistrationService {
	return &registrationService{
		credentialRepo: credential,
		userRepo:       user,
	}
}

const (
	RPID                       string = "localhost"
	RPNAME                     string = "miliya612 - webauthn demo"
	TIMEOUTMILLSEC             uint32 = 6000
	CLIENTDATATYPE             string = "webauthn.create"
	IsUserVerificationRequired bool   = false
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

	acp := webauthnif.AttestationConveyancePreferenceDirect

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

func (s registrationService) ReserveClientInfo(userId []byte, name, displayName, icon string) error {
	u := &model.User{
		ID:          userId,
		Name:        name,
		DisplayName: displayName,
		Icon:        icon,
	}
	_, err := s.userRepo.Create(*u)
	if err != nil {
		return err
	}

	return nil
}

func (s registrationService) ParseClientData(req webauthnif.AuthenticatorAttestationResponse) (
	*webauthnif.CollectedClientData, error) {

	// 1. Let JSONtext be the result of running UTF-8 decode on the value of response.clientDataJSON.
	// 2. Let C, the client data claimed as collected during the credential creation, be the result of running an
	// implementation-specific JSON parser on JSONtext.
	c := webauthnif.CollectedClientData{}
	if err := json.Unmarshal(req.ClientDataJSON, &c); err != nil {
		return nil, errors.New(fmt.Sprintf("invalidRegistrationRequest: parsingClientData: %v", err))
	}

	return &c, nil
}

func (s registrationService) ValidateClientData(rawChal []byte, c webauthnif.CollectedClientData) error {
	// 3. Verify that the value of C.type is webauthn.create.
	if c.Type != CLIENTDATATYPE {
		errMsg := fmt.Sprintf("got %q, but %q is required", c.Type, CLIENTDATATYPE)
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 4. Verify that the value of C.challenge matches the challenge that was sent to the authenticator in the create()
	// call.
	orgnChallenge := (webauthnif.BufferSource)(rawChal)

	byteChal, err := base64.RawURLEncoding.DecodeString(c.Challenge) // This is raw URL encoding, so the JSON parser does not handle it
	if err != nil {
		return err
	}
	challenge := (webauthnif.BufferSource)(byteChal)

	if !challenge.Equals(orgnChallenge) {
		fmt.Println("want: ", orgnChallenge)
		fmt.Println("got: ", challenge)
		errMsg := "invalid challenge"
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 5. Verify that the value of C.origin matches the Relying Party's origin.
	// TODO: RPID, subdomainマッチのロジック必要そう。subDomainとrootDomainの判定処理も書く
	if c.Origin != "http://" + RPID + ":8080" {
		errMsg := "invalid origin"
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 6. Verify that the value of C.tokenBinding.status matches the state of Token Binding for the TLS connection over
	// which the assertion was obtained. If Token Binding was used on that TLS connection, also verify that
	// C.tokenBinding.id matches the base64url encoding of the Token Binding ID for the connection.
	// TODO: やる

	return nil
}

func (s registrationService) ParseAttestationObj(
	req []byte,
	d *webauthnif.DecodedAuthenticatorAttestationResponse,
) (*webauthnif.DecodedAuthenticatorAttestationResponse, error) {
	// 8. Perform CBOR decoding on the attestationObject field of the AuthenticatorAttestationResponse structure to
	// obtain the attestation statement format fmt, the authenticator data authData, and the attestation statement
	// attStmt.
	handle := codec.CborHandle{}
	err := codec.NewDecoder(bytes.NewReader(req), &handle).Decode(&d.DecodedAttestationObject)
	if err != nil {
		return nil, err
	}
	fmt.Println("-----")
	fmt.Println(d.DecodedAttestationObject.RawAuthData)

	err = d.DecodedAttestationObject.UnmarshalBinary()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s registrationService) ValidateAuthenticatorData(data webauthnif.AuthenticatorData) error {
	// 9. Verify that the RP ID hash in authData is indeed the SHA-256 hash of the RP ID expected by the RP.
	wantRpIdHash := sha256.Sum256([]byte(RPID))
	gotRpIdHash := data.RPIDHash
	if !bytes.Equal(wantRpIdHash[:], gotRpIdHash) {
		errMsg := "invalid origin"
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 10. Verify that the User Present bit of the flags in authData is set.
	if !data.Flags.UserPresent() {
		errMsg := "no user presentation"
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 11. If user verification is required for this registration, verify that the User Verified bit of the flags in
	// authData is set.
	if IsUserVerificationRequired {
		if !data.Flags.UserVerified() {
			errMsg := "no user verification"
			return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
		}
	}
	return nil
}

func (s registrationService) ValidateClientExtensionOutputs(
	outputs webauthnif.AuthenticationExtensionsClientOutputs) error {
	// 12. Verify that the values of the client extension outputs in clientExtensionResults and the authenticator
	// extension outputs in the extensions in authData are as expected, considering the client extension input values
	// that were given as the extensions option in the create() call. In particular, any extension identifier values in
	// the clientExtensionResults and the extensions in authData MUST be also be present as extension identifier values
	// in the extensions member of options, i.e., no extensions are present that were not requested. In the general
	// case, the meaning of "are as expected" is specific to the Relying Party and which extensions are in use.
	return nil
}

func (s registrationService) ValidateAttestationResponse(
	attObj webauthnif.DecodedAttestationObject, hashedClientData [32]byte) error {
	// 13. Determine the attestation statement format by performing a USASCII case-sensitive match on fmt against the
	// set of supported WebAuthn Attestation Statement Format Identifier values. The up-to-date list of registered
	// WebAuthn Attestation Statement Format Identifier values is maintained in the in the IANA registry of the same
	// name.

	verifier, ok := attestation.AttVerifiers[attObj.Fmt]
	if !ok {
		errMsg := "Unsupported attestation statement format identifier"
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 14. Verify that attStmt is a correct attestation statement, conveying a valid attestation signature, by using the
	// attestation statement format fmt’s verification procedure given attStmt, authData and the hash of the serialized
	// client data computed in step 7.
	err := verifier(attObj, hashedClientData)
	if err != nil {
		errMsg := fmt.Sprintf("attestation statement is not matched with its format: %v", attObj.Fmt)
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 15. If validation is successful, obtain a list of acceptable trust anchors (attestation root certificates or
	// ECDAA-Issuer public keys) for that attestation type and attestation statement format fmt, from a trusted source
	// or from policy. For example, the FIDO Metadata Service  [FIDOMetadataService] provides one way to obtain such
	// information, using the aaguid in the attestedCredentialData in authData.

	// 16. Assess the attestation trustworthiness using the outputs of the verification procedure in step 14, as follows:
	//     - If self attestation was used, check if self attestation is acceptable under Relying Party policy.
	//     - If ECDAA was used, verify that the identifier of the ECDAA-Issuer public key used is included in the set of
	//     acceptable trust anchors obtained in step 15.
	//     - Otherwise, use the X.509 certificates returned by the verification procedure to verify that the attestation
	//     public key correctly chains up to an acceptable root certificate.

	return nil
}

func (s registrationService) Register(userId []byte, data webauthnif.AuthenticatorData) error {
	// 17. Check that the credentialId is not yet registered to any other user. If registration is requested for a
	// credential that is already registered to a different user, the Relying Party SHOULD fail this registration
	// ceremony, or it MAY decide to accept the registration, e.g. while deleting the older registration.
	cred, err := s.credentialRepo.GetByCredentialID(data.AttestedCredentialData.CredentialID)
	if err != nil {
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", err))
	}
	if cred != nil {
		errMsg := "credentialId has already been registered"
		return errors.New(fmt.Sprintf("invalidRegistrationRequest: %v", errMsg))
	}

	// 18. If the attestation statement attStmt verified successfully and is found to be trustworthy, then register the
	// new credential with the account that was denoted in the options.user passed to create(), by associating it with
	// the credentialId and credentialPublicKey in the attestedCredentialData in authData, as appropriate for the
	// Relying Party's system.
	newCred := model.Credential{
		CredentialID: data.AttestedCredentialData.CredentialID,
		UserID:       userId,
		PublicKey:    data.AttestedCredentialData.CredentialPublicKey,
		SignCount:    data.SignCount,
	}
	_, err = s.credentialRepo.Create(newCred)
	if err != nil {
		return err
	}

	return nil
}
