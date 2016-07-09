package keymaker

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

const (
	// RSA1024 1024-bit RSA
	RSA1024 = "RSA1024"
	// RSA2048 2048-bit RSA
	RSA2048 = "RSA2048"
	// RSA4096 4096-bit RSA
	RSA4096 = "RSA4096"

	// ECDSA224 P224 elliptic curve
	ECDSA224 = "ECDSA224"
	// ECDSA256 P256 elliptic curve
	ECDSA256 = "ECDSA256"
	// ECDSA384 P384 elliptic curve
	ECDSA384 = "ECDSA384"
	// ECDSA521 P521 elliptic curve
	ECDSA521 = "ECDSA521"
)

var (
	rsaTypes = map[string]int{
		RSA1024: 1024,
		RSA2048: 2048,
		RSA4096: 4096,
	}

	ecdsaTypes = map[string]func() elliptic.Curve{
		ECDSA224: elliptic.P224,
		ECDSA256: elliptic.P256,
		ECDSA384: elliptic.P384,
		ECDSA521: elliptic.P521,
	}
)

// GeneratePrivateKey generates a private key of the given type.
// The string encodes both the algorithm type and its parameters, and
// must be one of the known types:
//  RSA1024, RSA2048, RSA4096
//  EC224, EC256, EC384, EC521
func GeneratePrivateKey(keyType string) (crypto.PrivateKey, error) {
	switch keyType {
	case RSA1024, RSA2048, RSA4096:
		return rsa.GenerateKey(rand.Reader, rsaTypes[keyType])
	case ECDSA224, ECDSA256, ECDSA384, ECDSA521:
		return ecdsa.GenerateKey(ecdsaTypes[keyType](), rand.Reader)
	default:
		return nil, fmt.Errorf("%s key type not implemented", keyType)
	}
}
