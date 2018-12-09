package attestation

import "github.com/miliya612/webauthn-demo/webauthnif"

func init() {
	RegisterAttVerifier(webauthnif.AttestationStatementFormatPacked, verifyPacked)
}

func verifyPacked(attObj webauthnif.DecodedAttestationObject, clientDataHash [32]byte) error {
	return nil
}
