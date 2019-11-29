package doorman

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

var bitSize = 4096

// RSAKeyCertificate struct
type RSAKeyCertificate struct {
	RSAPrivateKey []byte
	Certificate   []byte
}

func newRSAKeyPair() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Error(err, "Error creating RSA key pair")
	}

	return privateKey
}

// toString() (privateKey string, publicKey string)
func generateCertificateForKey() (RSAKeyCertificate, error) {
	privateRSAKey := newRSAKeyPair()
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateRSAKey),
		},
	)

	certTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Doorman, Inc."},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	certificateDER, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateRSAKey.PublicKey, privateRSAKey)
	if err != nil {
		return RSAKeyCertificate{}, err
	}

	certificatePEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certificateDER,
		},
	)

	return RSAKeyCertificate{
		RSAPrivateKey: privateKeyPEM,
		Certificate:   certificatePEM,
	}, nil
}
