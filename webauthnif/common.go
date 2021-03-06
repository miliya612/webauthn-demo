package webauthnif

// BufferSource is used to represent byte array which replaces BufferSource.
// https://heycam.github.io/webidl/#BufferSource
// As a cryptographic protocol, Web Authentication is dependent upon randomized challenges to avoid replay attacks.
// Therefore, the values of both PublicKeyCredentialCreationOptions.challenge and
// PublicKeyCredentialRequestOptions.challenge MUST be randomly generated by Relying Parties in an environment
// they trust (e.g., on the server-side), and the returned challenge value in the client’s response MUST match what was
// generated. This SHOULD be done in a fashion that does not rely upon a client’s behavior, e.g., the Relying Party
// SHOULD store the challenge temporarily until the operation is complete. Tolerating a mismatch will compromise the
// security of the protocol. In order to prevent replay attacks, the challenges MUST contain enough entropy to make
// guessing them infeasible. Challenges SHOULD therefore be at least 16 bytes long.
// See https://www.w3.org/TR/webauthn/#cryptographic-challenges
type BufferSource []byte

// PublicKeyCredentialEntity describes a user account, or a WebAuthn Relying Party, with which a public key credential
// is associated.
// See https://www.w3.org/TR/webauthn/#dictionary-pkcredentialentity
type PublicKeyCredentialEntity struct {
	// Name is A human-palatable name for the entity. Its function depends on what the PublicKeyCredentialEntity
	// represents:
	//   - When inherited by PublicKeyCredentialRpEntity it is a human-palatable identifier for the Relying Party,
	//     intended only for display. For example, "ACME Corporation", "Wonderful Widgets, Inc." or "ОАО Примертех".
	//     - Relying Parties SHOULD perform enforcement, as prescribed in Section 2.3 of  [RFC8266] for the Nickname
	//       Profile of the PRECIS FreeformClass [RFC8264], when setting name's value, or displaying the value to the
	//       user.
	//     - Clients SHOULD perform enforcement, as prescribed in Section 2.3 of [RFC8266] for the Nickname Profile of
	//       the PRECIS FreeformClass [RFC8264], on name's value prior to displaying the value to the user or including
	//       the value as a parameter of the authenticatorMakeCredential operation.
	//   - When inherited by PublicKeyCredentialUserEntity, it is a human-palatable identifier for a user account. It is
	//     intended only for display, i.e., aiding the user in determining the difference between user accounts with
	//     similar displayNames. For example, "alexm", "alex.p.mueller@example.com" or "+14255551234".
	//     - The Relying Party MAY let the user choose this value. The Relying Party SHOULD perform enforcement, as
	//       prescribed in Section 3.4.3 of [RFC8265] for the UsernameCasePreserved Profile of the PRECIS
	//       IdentifierClass [RFC8264], when setting name's value, or displaying the value to the user.
	//     - Clients SHOULD perform enforcement, as prescribed in Section 3.4.3 of [RFC8265] for the
	//       UsernameCasePreserved Profile of the PRECIS IdentifierClass [RFC8264], on name's value prior to displaying
	//       the value to the user or including the value as a parameter of the authenticatorMakeCredential operation.
	// When clients, client platforms, or authenticators display a name's value, they should always use UI elements to
	// provide a clear boundary around the displayed value, and not allow overflow into other elements [css-overflow-3].
	// Authenticators MUST accept and store a 64-byte minimum length for a name member’s value. Authenticators MAY
	// truncate a name member’s value to a length equal to or greater than 64 bytes.
	Name string `json:"name"`
	// Icon is a serialized URL which resolves to an image associated with the entity. For example, this could be a
	// user’s avatar or a Relying Party's logo. This URL MUST be an a priori authenticated URL. Authenticators MUST
	// accept and store a 128-byte minimum length for an icon member’s value. Authenticators MAY ignore an icon member’s
	// value if its length is greater than 128 bytes. The URL’s scheme MAY be "data" to avoid fetches of the URL, at the
	// cost of needing more storage.
	Icon string `json:"icon,omitempty"`
}
