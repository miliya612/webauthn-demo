package webauthnif

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
	"math/rand"
	"time"
)

// ToBufferSource translates given string to BufferSource
func ToBufferSource(str string) (bs BufferSource) {
	//bytes := sha256.Sum256([]byte(str))
	//bs = bytes[:]
	//return
	bs = []byte(str)
	return
}

func GenChallenge() (challenge BufferSource, err error) {
	challenge = make(BufferSource, 32)
	rand.Seed(time.Now().UnixNano())
	_, err = rand.Read(challenge)
	return
}

func (bf BufferSource) Equals(abf BufferSource) bool {
	return bytes.Equal(([]byte)(bf), ([]byte)(abf))
}

// UnmarshalBinary
func (a *DecodedAttestationObject) UnmarshalBinary() error {
	if len(a.RawAuthData) < 37 {
		return errors.New("invalid authenticator data")
	}

	a.AuthData.RPIDHash = a.RawAuthData[0:32]
	a.AuthData.Flags = AuthenticatorDataFlags(a.RawAuthData[32])
	a.AuthData.SignCount = binary.BigEndian.Uint32(a.RawAuthData[33:37])

	if a.AuthData.Flags.HasAttestedCredentialData() && len(a.RawAuthData) > 37 {
		credentialIDLen := binary.BigEndian.Uint16(a.RawAuthData[53:55])
		dPubKey, err := ParseCOSE(a.RawAuthData[55+credentialIDLen:])
		if err != nil {
			return errors.New(fmt.Sprintf("unable to parse COSE key: %v", err.Error()))
		}

		attestedCredentialData := AttestedCredentialData{
			AAGUID:                     a.RawAuthData[37:53],
			CredentialIdLength:         credentialIDLen,
			CredentialID:               a.RawAuthData[55 : 55+credentialIDLen],
			CredentialPublicKey:        a.RawAuthData[55+credentialIDLen:],
			DecodedCredentialPublicKey: dPubKey,
		}

		a.AuthData.AttestedCredentialData = attestedCredentialData

	}

	return nil
}

// ParseCOSE parses a raw COSE key into a public key
func ParseCOSE(buf []byte) (interface{}, error) {
	m := make(map[int]interface{})

	cbor := codec.CborHandle{}

	if err := codec.NewDecoder(bytes.NewReader(buf), &cbor).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func ParseCOSEMap(m map[int]interface{}) (interface{}, error) {
	err := validateCOSEMap(m)
	if err != nil {
		return nil, errors.Wrap(err, "invalid COSE map")
	}
	return nil, nil
}

func validateCOSEMap(m map[int]interface{}) error {
	/*
		{
		  1: kty=2,  // EC2 key type
		  3: alg=-7, // ES256 signature algorithm
		 -1: crv=1,  // P-256 curve
		 -2: x,      // x-coordinate as byte string 32 bytes in length
		 -3: y,      // y-coordinate as byte string 32 bytes in length
		}
	 */

	// The credential public key encoded in COSE_Key format, as defined in Section 7 of [RFC8152], using the CTAP2
	// canonical CBOR encoding form. The COSE_Key-encoded credential public key MUST contain the "alg" parameter and
	// MUST NOT contain any other OPTIONAL parameters. The "alg" parameter MUST contain a COSEAlgorithmIdentifier value.
	// The encoded credential public key MUST also contain any additional REQUIRED parameters stipulated by the
	// relevant key type specification, i.e., REQUIRED for the key type "kty" and algorithm "alg" (see Section 8 of
	// [RFC8152]).

	//if kty, ok := m[1]; !ok {
	//	return errors.New("missing ec2 key type")
	//}
	//
	//if alg, ok := m[3]; !ok {
	//	return errors.New("missing es256 signature algorithm")
	//}

	return nil
}
