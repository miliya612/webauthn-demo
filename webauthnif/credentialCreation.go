package webauthnif

// 2.2
// Credential represents an interface of credentials which are returned in the API `navigator.credentials.* `
// See https://w3c.github.io/webappsec-credential-management/#credential
type Credential struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// 5.1
// PublicKeyCredential interface inherits from Credential [CREDENTIAL-MANAGEMENT-1], and contains the attributes that
// are returned to the caller when a new credential is created, or a new assertion is requested.
type PublicKeyCredential struct {
	Credential
	// RawID returns the ArrayBuffer contained in the [[identifier]] internal slot.
	RawID []byte `json:"rawId"`
	// Response (AuthenticatorAttestationResponse) contains the authenticator's response to the client’s request to
	// either create a public key credential, or generate an authentication assertion. If the PublicKeyCredential is
	// created in response to create(), this attribute’s value will be an AuthenticatorAttestationResponse, otherwise,
	// the PublicKeyCredential was created in response to get(), and this attribute’s value will be an
	// AuthenticatorAssertionResponse.
	Response AuthenticatorAttestationResponse `json:"response"`
	// AuthenticationExtensionsClientOutputs returns the value of [[clientExtensionsResults]], which is a map containing
	// extension identifier → client extension output entries produced by the extension’s client extension processing.
	AuthenticationExtensionsClientOutputs AuthenticationExtensionsClientOutputs `json:"authenticationExtensionsClientOutputs"`
}

// 5.9
// AuthenticationExtensionsClientOutputs dictionary containing the authenticator extension input values for zero or more
// WebAuthn extensions, as defined in §9 WebAuthn Extensions.
type AuthenticationExtensionsClientOutputs struct {
}

// 5.1.1
// CredentialCreationOptions is an extension of the CredentialCreationOptions dictionary in order to support
// registration via navigator.credentials.create()
// See https://www.w3.org/TR/webauthn/#credentialcreationoptions-extension
type CredentialCreationOptions struct {
	PublicKey PublicKeyCredentialCreationOptions `json:"publicKey"`
}

// 5.4
// PublicKeyCredentialCreationOptions is an entity of the options
// for the credential creation on authenticator.
// See https://www.w3.org/TR/webauthn/#dictionary-makecredentialoptions
type PublicKeyCredentialCreationOptions struct {
	// RP contains data about the Relying Party responsible for the request.
	RP PublicKeyCredentialRpEntity `json:"rp"`
	// User contains data about the user account
	// for which the Relying Party is requesting attestation.
	User PublicKeyCredentialUserEntity `json:"user"`
	// Challenge contains a challenge intended to be used for
	// generating the newly created credential’s attestation object.
	Challenge BufferSource `json:"challenge"`
	// PubKeyCredParams contains information about the desired properties of the credential to be created.
	// The sequence is ordered from most preferred to least preferred.
	// The client makes a best-effort to create the most preferred credential that it can.
	PubKeyCredParams PublicKeyCredentialParameters `json:"pubKeyCredParams"`
	// Timeout specifies a time, in milliseconds, that the caller is willing to wait for the call to complete. This is
	// treated as a hint, and MAY be overridden by the client.
	Timeout uint32 `json:"timeout,omitempty"`
	// ExcludeCredentials is intended for use by Relying Parties that wish to limit the creation of multiple credentials
	// for the same account on a single authenticator. The client is requested to return an error if the new credential
	// would be created on an authenticator that also contains one of the credentials enumerated in this parameter.
	// Default: None
	ExcludeCredentials PublicKeyCredentialDescriptors `json:"excludeCredentials,omitempty"`
	// AuthenticatorSelection is intended for use by Relying Parties that wish to select the appropriate authenticators
	// to participate in the create() operation.
	AuthenticatorSelection AuthenticatorSelectionCriteria `json:"authenticatorSelection,omitempty"`
	// Attestation is intended for use by Relying Parties that wish to express their preference for attestation
	// conveyance.
	// Default: none
	Attestation AttestationConveyancePreference `json:"attestation,omitempty"`
	// Extensions contains additional parameters requesting additional processing by the client and authenticator.
	// For example, the caller may request that only authenticators with certain capabilities be used to create the
	// credential, or that particular information be returned in the attestation object. Some extensions are defined
	// in §9 WebAuthn Extensions; consult the IANA "WebAuthn Extension Identifier" registry established by
	// [WebAuthn-Registries] for an up-to-date list of registered WebAuthn Extensions.
	Extensions AuthenticationExtensionsClientInputs `json:"extensions,omitempty"`
}

// 5.4.2
// PublicKeyCredentialRpEntity is used to supply additional Relying Party attributes when creating a new credential.
// See https://www.w3.org/TR/webauthn/#dictdef-publickeycredentialrpentity
type PublicKeyCredentialRpEntity struct {
	// ID is a unique identifier for the Relying Party entity, which sets the RP ID
	ID string `json:"id,omitempty"`
	PublicKeyCredentialEntity
}

// 5.4.3
// PublicKeyCredentialUserEntity contains data about the user account for which the Relying Party is requesting
// attestation.
// See https://www.w3.org/TR/webauthn/#dictdef-publickeycredentialuserentity
type PublicKeyCredentialUserEntity struct {
	// ID is the user handle of the user account entity. To ensure secure operation, authentication and authorization
	// decisions MUST be made on the basis of this id member, not the displayName nor name members.
	ID BufferSource `json:"id"`
	// DisplayName is a human-palatable name for the user account, intended only for display. For example,
	// "Alex P. Müller" or "田中 倫". The Relying Party SHOULD let the user choose this, and SHOULD NOT restrict the
	// choice more than necessary.
	//   - Relying Parties SHOULD perform enforcement, as prescribed in Section 2.3 of [RFC8266] for the Nickname
	//     Profile of the PRECIS FreeformClass [RFC8264], when setting displayName's value, or displaying the value to
	//     the user.
	//   - Clients SHOULD perform enforcement, as prescribed in Section 2.3 of [RFC8266] for the Nickname Profile of the
	//     PRECIS FreeformClass [RFC8264], on displayName's value prior to displaying the value to the user or including
	//     the value as a parameter of the authenticatorMakeCredential operation.
	// When clients, client platforms, or authenticators display a displayName's value, they should always use UI
	// elements to provide a clear boundary around the displayed value, and not allow overflow into other elements
	// [css-overflow-3].
	// Authenticators MUST accept and store a 64-byte minimum length for a displayName member’s value. Authenticators
	// MAY truncate a displayName member’s value to a length equal to or greater than 64 bytes.
	DisplayName string `json:"displayName"`
	PublicKeyCredentialEntity
}

// 5.3
// PublicKeyCredentialParameter specifies the type of credential to be created.
// See https://www.w3.org/TR/webauthn/#credential-params
type PublicKeyCredentialParameter struct {
	// PublicKeyCredentialType specifies the type of credential to be created.
	Type PublicKeyCredentialType `json:"type"`

	// COSEAlgorithmIdentifier specifies the cryptographic signature algorithm with which the newly generated credential
	// will be used, and thus also the type of asymmetric key pair to be generated, e.g., RSA or Elliptic Curve.
	Alg COSEAlgorithmIdentifier `json:"alg"`
}

type PublicKeyCredentialParameters []PublicKeyCredentialParameter

// 5.10.3
// PublicKeyCredentialDescriptor contains the attributes that are specified by a caller when referring to a public key
// credential as an input parameter to the create() or get() methods. It mirrors the fields of the PublicKeyCredential
// object returned by the latter methods.
// See https://www.w3.org/TR/webauthn/#dictdef-publickeycredentialdescriptor
type PublicKeyCredentialDescriptor struct {
	// Type contains the type of the public key credential the caller is referring to.
	Type PublicKeyCredentialType `json:"type"`
	// ID contains the credential ID of the public key credential the caller is referring to.
	ID BufferSource `json:"id"`
	// Transports contains a hint as to how the client might communicate with the managing authenticator of the public
	// key credential the caller is referring to.
	// This is OPTIONAL.
	Transports AuthenticatorTransports `json:"transports,omitempty"`
}

type PublicKeyCredentialDescriptors []PublicKeyCredentialDescriptor

// 5.4.4
// AuthenticatorSelectionCriteria specifies their requirements regarding authenticator attributes.
// See https://www.w3.org/TR/webauthn/#authenticatorSelection
type AuthenticatorSelectionCriteria struct {
	// AuthenticatorAttachment makes eligible authenticators filtered to only authenticators attached with the specified
	// AuthenticatorAttachment
	AuthenticatorAttachment AuthenticatorAttachment `json:"authenticatorAttachment,omitempty"`
	// RequireResidentKey describes the Relying Parties' requirements regarding resident credentials. If the parameter
	// is set to true, the authenticator MUST create a client-side-resident public key credential source when creating a
	// public key credential.
	// Default: false
	RequireResidentKey bool `json:"requireResidentKey,omitempty"`
	// UserVerification describes the Relying Party's requirements regarding user verification for the create()
	// operation. Eligible authenticators are filtered to only those capable of satisfying this requirement.
	UserVerification UserVerificationRequirement `json:"userVerification,omitempty"`
}

// 5.4.6
// AttestationConveyancePreference specifies RP's preference regarding attestation conveyance during credential
// generation.
// See https://www.w3.org/TR/webauthn/#enumdef-attestationconveyancepreference
type AttestationConveyancePreference string

const (
	AttestationConveyancePreferenceEmpty AttestationConveyancePreference = ""
	// AttestationConveyancePreferenceNone indicates that the Relying Party is not interested in authenticator
	// attestation. For example, in order to potentially avoid having to obtain user consent to relay identifying
	// information to the Relying Party, or to save a roundtrip to an Attestation CA.
	AttestationConveyancePreferenceNone AttestationConveyancePreference = "none"
	// AttestationConveyancePreferenceIndirect indicates that the Relying Party prefers an attestation conveyance
	// yielding verifiable attestation statements, but allows the client to decide how to obtain such attestation
	// statements. The client MAY replace the authenticator-generated attestation statements with attestation statements
	// generated by an Anonymization CA, in order to protect the user’s privacy, or to assist Relying Parties with
	// attestation verification in a heterogeneous ecosystem.
	AttestationConveyancePreferenceIndirect AttestationConveyancePreference = "indirect"
	// AttestationConveyancePreferenceDirect indicates that the Relying Party wants to receive the attestation statement
	// as generated by the authenticator.
	AttestationConveyancePreferenceDirect AttestationConveyancePreference = "direct"
)

// 5.8
// AuthenticationExtensionsClientInputs is a dictionary containing the client extension output values for zero or more
// WebAuthn extensions, as defined in §9 WebAuthn Extensions.
// See https://www.w3.org/TR/webauthn/#dictdef-authenticationextensionsclientinputs
type AuthenticationExtensionsClientInputs struct{}

// PublicKeyCredentialType defines the valid credential types. It is an extension point; values can be added to it in
// the future, as more credential types are defined. The values of this enumeration are used for versioning the
// Authentication Assertion and attestation structures according to the type of the authenticator.
// Currently one credential type is defined, namely "public-key".
// See https://www.w3.org/TR/webauthn/#credentialType
type PublicKeyCredentialType string

const (
	PublicKeyCredentialTypePublicKey PublicKeyCredentialType = "public-key"
)

// 5.10.5
// COSEAlgorithmIdentifier has value which is a number identifying a cryptographic algorithm. The algorithm identifiers
// SHOULD be values registered in the IANA COSE Algorithms registry  [IANA-COSE-ALGS-REG], for instance, -7 for "ES256"
// and -257 for "RS256".
// See https://www.w3.org/TR/webauthn/#typedefdef-cosealgorithmidentifier
type COSEAlgorithmIdentifier int

const (
	COSEAlgorithmIdentifierES256 COSEAlgorithmIdentifier = -7
	COSEAlgorithmIdentifierRS256 COSEAlgorithmIdentifier = -257
)

// 5.10.4
// AuthenticatorTransport defines hints as to how clients might communicate with a particular authenticator in order to
// obtain an assertion for a specific credential. Note that these hints represent the WebAuthn Relying Party's best
// belief as to how an authenticator may be reached. A Relying Party may obtain a list of transports hints from some
// attestation statement formats or via some out-of-band mechanism; it is outside the scope of this specification to
// define that mechanism.
// Authenticators may implement various transports for communicating with clients.
// See https://www.w3.org/TR/webauthn/#credential-dictionary
type AuthenticatorTransport string

const (
	AuthenticatorTransportEmpty AuthenticatorTransport = ""
	// AuthenticatorTransportUSB indicates the respective authenticator can be contacted over removable USB.
	AuthenticatorTransportUSB AuthenticatorTransport = "usb"
	// AuthenticatorTransportNFC indicates the respective authenticator can be contacted over Near Field Communication
	// (NFC).
	AuthenticatorTransportNFC AuthenticatorTransport = "nfc"
	// AuthenticatorTransportBLE indicates the respective authenticator can be contacted over Bluetooth Smart (Bluetooth
	// Low Energy / BLE).
	AuthenticatorTransportBLE AuthenticatorTransport = "ble"
	// AuthenticatorTransportInternal indicates the respective authenticator is contacted using a client device-specific
	// transport. These authenticators are not removable from the client device.
	AuthenticatorTransportInternal AuthenticatorTransport = "internal"
)

type AuthenticatorTransports []AuthenticatorTransport

// 5.4.5
// AuthenticatorAttachment has values describe authenticators' attachment modalities. Relying Parties use this for two
// purposes:
//   - to express a preferred authenticator attachment modality when calling navigator.credentials.create() to create a
//     credential, and
//   - to inform the client of the Relying Party's best belief about how to locate the managing authenticators of the
//     credentials listed in allowCredentials when calling navigator.credentials.get().
// See https://www.w3.org/TR/webauthn/#enumdef-authenticatorattachment
type AuthenticatorAttachment string

const (
	AuthenticatorAttachmentEmpty AuthenticatorAttachment = ""
	// AuthenticatorAttachmentPlatform indicates platform attachment.
	// A platform authenticator is attached using a client device-specific transport, called platform attachment, and is
	// usually not removable from the client device. A public key credential bound to a platform authenticator is called
	// a platform credential.
	// See https://www.w3.org/TR/webauthn/#platform-attachment
	AuthenticatorAttachmentPlatform AuthenticatorAttachment = "platform"
	// AuthenticatorAttachmentCrossPlatform indicates cross-platform attachment.
	// A CrossPlatform authenticator called roaming authenticator, is attached using cross-platform transports
	// Authenticators of this class are removable from, and can "roam" among, client devices. A public key credential
	// bound to a roaming authenticator is called a roaming credential.
	// See https://www.w3.org/TR/webauthn/#cross-platform-attachment
	AuthenticatorAttachmentCrossPlatform AuthenticatorAttachment = "cross-platform"
)

// 5.10.6
// UserVerificationRequirement specifies RP's requirement of user verification for some of its operations but not for
// others, and may use this type to express its needs.
// See https://www.w3.org/TR/webauthn/#enumdef-userverificationrequirement
type UserVerificationRequirement string

const (
	UserVerificationRequirementEmpty UserVerificationRequirement = ""
	// UserVerificationRequirementRequired indicates that the Relying Party requires user verification for the operation
	// and will fail the operation if the response does not have the UV flag set.
	UserVerificationRequirementRequired UserVerificationRequirement = "required"
	// UserVerificationRequirementPreferred indicates that the Relying Party prefers user verification for the operation
	// if possible, but will not fail the operation if the response does not have the UV flag set.
	UserVerificationRequirementPreferred UserVerificationRequirement = "preferred"
	// UserVerificationRequirementDiscouraged indicates that the Relying Party prefers user verification for the
	// operation if possible, but will not fail the operation if the response does not have the UV flag set.
	UserVerificationRequirementDiscouraged UserVerificationRequirement = "discouraged"
)
