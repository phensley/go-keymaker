package keymaker

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

// BuildClientTLSConfig constructs a tls.Config from the given parts
func BuildClientTLSConfig(certPEM, keyPEM, bundlePEM []byte) (*tls.Config, error) {
	certificate, err := decodeCertificate(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	caBundle, err := decodeCABundle(bundlePEM)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:      caBundle,
		Certificates: []tls.Certificate{*certificate},
	}, nil
}

// BuildDroneTLSConfig constructs a tls.Config from the given parts
func BuildDroneTLSConfig(certPEM, keyPEM, bundlePEM []byte, clientAuth string) (*tls.Config, error) {
	certificate, err := decodeCertificate(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	caBundle, err := decodeCABundle(bundlePEM)
	if err != nil {
		return nil, err
	}

	var authType tls.ClientAuthType

	switch clientAuth {
	case "none":
		authType = tls.NoClientCert
	case "any":
		authType = tls.RequireAnyClientCert
	case "request":
		authType = tls.RequestClientCert
	case "require-and-verify":
		authType = tls.RequireAndVerifyClientCert
	default:
		return nil, fmt.Errorf("Unknown client auth type %v", clientAuth)
	}

	return &tls.Config{
		RootCAs:                  caBundle,
		ClientCAs:                caBundle,
		Certificates:             []tls.Certificate{*certificate},
		PreferServerCipherSuites: true,
		ClientAuth:               authType,
		MinVersion:               tls.VersionTLS12,

		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}, nil
}

func decodeCertificate(certPEM, keyPEM []byte) (*tls.Certificate, error) {
	certAndKey, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	certAndKey.Leaf, err = x509.ParseCertificate(certAndKey.Certificate[0])
	if err != nil {
		return nil, err
	}

	return &certAndKey, nil
}

func decodeCABundle(bundlePEM []byte) (*x509.CertPool, error) {
	caBundle := x509.NewCertPool()
	if ok := caBundle.AppendCertsFromPEM(bundlePEM); !ok {
		return nil, fmt.Errorf("failed to append certs from PEM")
	}
	return caBundle, nil
}
