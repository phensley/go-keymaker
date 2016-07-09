package keymaker

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
)

// PKCS8 ASN1 encoding
type pkcs8 struct {
	version    int
	Algorithm  pkix.AlgorithmIdentifier
	PrivateKey []byte
}

var (
	asn1Null = asn1.RawValue{Tag: 5}

	oidRSA       = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	algorithmRSA = pkix.AlgorithmIdentifier{
		Algorithm:  oidRSA,
		Parameters: asn1Null,
	}

	oidECDSA       = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	algorithmECDSA = pkix.AlgorithmIdentifier{
		Algorithm:  oidECDSA,
		Parameters: asn1Null,
	}
)

// MarshalPKCS8PrivateKey encodes a key in PKCS#8 binary
func MarshalPKCS8PrivateKey(key crypto.PrivateKey) ([]byte, error) {
	switch keyType := key.(type) {
	case *rsa.PrivateKey:
		bytes := x509.MarshalPKCS1PrivateKey(key.(*rsa.PrivateKey))
		return asn1.Marshal(pkcs8{0, algorithmRSA, bytes})

	case *ecdsa.PrivateKey:
		bytes, err := x509.MarshalECPrivateKey(key.(*ecdsa.PrivateKey))
		if err != nil {
			return nil, err
		}
		return asn1.Marshal(pkcs8{0, algorithmECDSA, bytes})

	default:
		return nil, fmt.Errorf("unsupported private key %v", keyType)
	}
}

// UnmarshalPEMPrivateKey decodes a private key from the PEM bytes.
// It returns the decoded private key along with any remaining bytes.
func UnmarshalPEMPrivateKey(raw []byte) (crypto.PrivateKey, []byte, error) {
	decoded, rest := pem.Decode(raw)
	if decoded == nil {
		return nil, nil, fmt.Errorf("failed to decode PEM")
	}

	var key crypto.PrivateKey
	var err error

	// Choose parser based on PEM header
	switch decoded.Type {
	case "PRIVATE KEY":
		key, err = x509.ParsePKCS8PrivateKey(decoded.Bytes)
	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(decoded.Bytes)
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(decoded.Bytes)
	default:
		err = fmt.Errorf("%s is not a supported PEM private key type", decoded.Type)
	}

	return key, rest, err
}
