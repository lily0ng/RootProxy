package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type SelfSignedOptions struct {
	CommonName string
	ValidFor   time.Duration
}

type Generated struct {
	CertPEM []byte
	KeyPEM  []byte
}

func GenerateSelfSigned(opts SelfSignedOptions) (Generated, error) {
	if opts.CommonName == "" {
		opts.CommonName = "RootProxy"
	}
	if opts.ValidFor == 0 {
		opts.ValidFor = 365 * 24 * time.Hour
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return Generated{}, err
	}

	now := time.Now()
	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			CommonName: opts.CommonName,
		},
		NotBefore:             now.Add(-1 * time.Minute),
		NotAfter:              now.Add(opts.ValidFor),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		return Generated{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return Generated{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}
