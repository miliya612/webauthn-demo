package attestation

import "github.com/miliya612/webauthn-demo/webauthnif"

type AttVerifyFunc func(webauthnif.DecodedAttestationObject, [32]byte) error

var AttVerifiers = make(map[webauthnif.AttestationStatementFormatIdentifier]AttVerifyFunc)

func RegisterAttVerifier(fmt webauthnif.AttestationStatementFormatIdentifier, f AttVerifyFunc) {
	AttVerifiers[fmt] = f
}
