package webauthnif

// 5.5
// PublicKeyCredentialRequest supplies get() with the data it needs to generate an assertion. Its challenge member MUST
// be present, while its other members are OPTIONAL.
// See
type PublicKeyCredentialRequest struct {
	// Challenge represents a challenge that the selected authenticator signs, along with other data, when producing an
	// authentication assertion. See the §13.1 Cryptographic Challenges security consideration.
	Challenge        BufferSource                         `json:"challenge"`
	// Timeout is an OPTIONAL member specifies a time, in milliseconds, that the caller is willing to wait for the call
	// to complete. The value is treated as a hint, and MAY be overridden by the client.
	Timeout          int                                  `json:"timeout,omitempty"`
	// RPID is an OPTIONAL member specifies the relying party identifier claimed by the caller. If omitted, its value
	// will be the CredentialsContainer object’s relevant settings object's origin's effective domain.
	RPID             string                               `json:"rpId"`
	// AllowCredentials is an OPTIONAL member contains a list of PublicKeyCredentialDescriptor objects representing
	// public key credentials acceptable to the caller, in descending order of the caller’s preference (the first item
	// in the list is the most preferred credential, and so on down the list).
	AllowCredentials PublicKeyCredentialDescriptors       `json:"allowCredentials"`
	// UserVerification describes the Relying Party's requirements regarding user verification for the get() operation.
	// Eligible authenticators are filtered to only those capable of satisfying this requirement.
	UserVerification UserVerificationRequirement          `json:"userVerification"`
	// Extensions OPTIONAL member contains additional parameters requesting additional processing by the client and
	// authenticator. For example, if transaction confirmation is sought from the user, then the prompt string might be
	// included as an extension.
	Extensions       AuthenticationExtensionsClientInputs `json:"extensions"`
}
