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

}

// DecodedAuthenticatorResponse represents the result of running UTF-8 decode on the value of
// AuthenticatorResponse.clientDataJSON
type DecodedAuthenticatorResponse struct {
	ClientData ClientData `json:"clientData"`
}

// ClientData represents
type ClientData struct {
	Type string `json:"type"`
	Challenge BufferSource `json:"challenge"`
	Origin string `json:"origin"`
}