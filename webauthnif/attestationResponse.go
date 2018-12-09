package webauthnif

// 5.2.1
// AuthenticatorAttestationResponse represents the authenticator's response to a client’s request for the creation of a
// new public key credential. It contains information about the new credential that can be used to identify it for later
// use, and metadata that can be used by the WebAuthn Relying Party to assess the characteristics of the credential
// during registration.
// See https://www.w3.org/TR/webauthn/#iface-authenticatorattestationresponse
type AuthenticatorAttestationResponse struct {
	// AuthenticatorResponse contains the JSON-serialized client data (see §6.4 Attestation) passed to the authenticator
	// by the client in order to generate this credential. The exact JSON serialization MUST be preserved, as the hash
	// of the serialized client data has been computed over it.
	// the hash of the serialized client data is the hash (computed using SHA-256) of the JSON-serialized client data,
	// as constructed by the client.
	AuthenticatorResponse
	// AttestationObject contains an attestation object, which is opaque to, and cryptographically protected against tampering by, the
	// client. The attestation object contains both authenticator data and an attestation statement. The former contains
	// the AAGUID, a unique credential ID, and the credential public key. The contents of the attestation statement are
	// determined by the attestation statement format used by the authenticator. It also contains any additional
	// information that the Relying Party's server requires to validate the attestation statement, as well as to decode
	// and validate the authenticator data along with the JSON-serialized client data. For more details,
	// see §6.4 Attestation, §6.4.4 Generating an Attestation Object, and Figure 5.
	AttestationObject []byte `json:"attestationObject"`
}

// 5.2
// AuthenticatorResponse represents Authenticators response to which Relying Party requests
// See https://www.w3.org/TR/webauthn/#authenticatorresponse
type AuthenticatorResponse struct {
	// ClientDataJSON contains a JSON serialization of the client data passed to the authenticator by the client in its
	// call to either create() or get().
	ClientDataJSON []byte `json:"clientDataJson"`
}

// DecodedAuthenticatorAttestationResponse represents the result of parsing the AuthenticatorAttestationResponse
type DecodedAuthenticatorAttestationResponse struct {
	DecodedAuthenticatorResponse
	DecodedAttestationObject DecodedAttestationObject
}

// DecodedAuthenticatorResponse represents the result of running UTF-8 decode on the value of
// AuthenticatorResponse.clientDataJSON
type DecodedAuthenticatorResponse struct {
	ClientData CollectedClientData `json:"clientData"`
}

// CollectedClientData represents
type CollectedClientData struct {
	Type string `json:"type"`
	Challenge BufferSource `json:"challenge"`
	Origin string `json:"origin"`
}


// DecodedAttestationObject
type DecodedAttestationObject struct {
	Fmt AttestationStatementFormatIdentifier `json:"fmt"`
	AuthData AuthenticatorData `json:"authenticatorData"`
	AttStmt interface{} `json:"attStmt"`
}

// 8.1
// AttestationStatementFormatIdentifier defines a data format which represents a cryptographic signature by an authenticator over
// a set of contextual bindings
// AttestationStatementFormats are identified by a string, called an attestation statement format identifier, chosen
// by the author of the attestation statement format.
// See https://www.w3.org/TR/webauthn/#attestation-statement-format-identifier
type AttestationStatementFormatIdentifier string

const (
	AttestationStatementFormatPacked           AttestationStatementFormatIdentifier = "packed"
	AttestationStatementFormatTPM              AttestationStatementFormatIdentifier = "tpm"
	AttestationStatementFormatAndroidKey       AttestationStatementFormatIdentifier = "android-key"
	AttestationStatementFormatAndroidSafetyNet AttestationStatementFormatIdentifier = "android-safetynet"
	AttestationStatementFormatFIDOU2F          AttestationStatementFormatIdentifier = "fido-u2f"
	AttestationStatementFormatNone             AttestationStatementFormatIdentifier = "none"
)

// 6.1
// AuthenticatorData encodes contextual bindings made by the authenticator. These bindings are controlled by the
// authenticator itself, and derive their trust from the WebAuthn Relying Party's assessment of the security properties
// of the authenticator. In one extreme case, the authenticator may be embedded in the client, and its bindings may be
// no more trustworthy than the client data. At the other extreme, the authenticator may be a discrete entity with
// high-security hardware and software, connected to the client over a secure channel. In both cases, the Relying Party
// receives the authenticator data in the same format, and uses its knowledge of the authenticator to make trust
// decisions.
// AuthenticatorData has a compact but extensible encoding. This is desired since authenticators can be devices
// with limited capabilities and low power requirements, with much simpler software stacks than the client platform.
// See https://www.w3.org/TR/webauthn/#authenticator-data
type AuthenticatorData struct {
	// RpIdHash is SHA-256 hash of the RP ID associated with the credential
	RpIdHash []byte
	// Flags
	Flags AuthenticatorDataFlags
	// SignCount is Signature counter, 32-bit unsigned big-endian integer
	SignCount uint32
	// AttestedCredentialData
	AttestedCredentialData AttestedCredentialData
}

type AuthenticatorDataFlags byte

const (
	// AuthenticatorDataFlagUserPresent indicates the UP flag.
	AuthenticatorDataFlagUserPresent = 0x001 // 0000 0001
	// AuthenticatorDataFlagUserVerified indicates the UV flag.
	AuthenticatorDataFlagUserVerified = 0x002 // 0000 0010
	// AuthenticatorDataFlagHasCredentialData indicates the AT flag.
	AuthenticatorDataFlagHasCredentialData = 0x040 // 0100 0000
	// AuthenticatorDataFlagHasExtension indicates the ED flag.
	AuthenticatorDataFlagHasExtension = 0x080 // 1000 0000
)

// UserPresent returns whether the UP flag is set.
func (f AuthenticatorDataFlags) UserPresent() bool {
	return (f & AuthenticatorDataFlagUserPresent) == AuthenticatorDataFlagUserPresent
}

// UserVerified returns whether the UV flag is set.
func (f AuthenticatorDataFlags) UserVerified() bool {
	return (f & AuthenticatorDataFlagUserVerified) == AuthenticatorDataFlagUserVerified
}

// HasAttestedCredentialData returns whether the AT flag is set.
func (f AuthenticatorDataFlags) HasAttestedCredentialData() bool {
	return (f & AuthenticatorDataFlagHasCredentialData) == AuthenticatorDataFlagHasCredentialData
}

// HasExtensions returns whether the ED flag is set.
func (f AuthenticatorDataFlags) HasExtensions() bool {
	return (f & AuthenticatorDataFlagHasExtension) == AuthenticatorDataFlagHasExtension
}

// 6.4.1
// AttestedCredentialData a variable-length byte array added to the authenticator data when generating an attestation
// object for a given credential.
// See https://www.w3.org/TR/webauthn/#sec-attested-credential-data
type AttestedCredentialData struct {
	// AAGUID is the one of the authenticator. its length is 16bit
	AAGUID []byte
	// CredentialIdLength is byte length L of Credential ID, 16-bit unsigned big-endian integer. its length is 2bit
	CredentialIdLength []byte
	// CredentialID is a probabilistically-unique byte sequence identifying a public key credential source and its
	// authentication assertions.
	// Credential IDs are generated by authenticators in two forms:
	//     1. At least 16 bytes that include at least 100 bits of entropy, or
	//     2. The public key credential source, without its Credential ID, encrypted so only its managing authenticator
	//        can decrypt it. This form allows the authenticator to be nearly stateless, by having the Relying Party
	//        store any necessary state.
	// Relying Parties do not need to distinguish these two Credential ID forms.
	// See https://www.w3.org/TR/webauthn/#credential-id
	CredentialID []byte

	// CredentialPublicKey is the public key portion of a Relying Party-specific credential key pair, generated by an
	// authenticator and returned to a Relying Party at registration time (see also public key credential). The private
	// key portion of the credential key pair is known as the credential private key. Note that in the case of self
	// attestation, the credential key pair is also used as the attestation key pair, see self attestation for details.
	CredentialPublicKey []byte
	// TODO: しらべる
	// credential public keyの要素

	DecodedCredentialPublicKey interface{}
	// 完成したpublic key
}

type AttestationStatement struct {

}