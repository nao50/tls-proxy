package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

// Generate self-signed certificate Root CA cert (ca.pem) and its key (ca-key.pem)
func main() {

	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}

	ca := x509.Certificate{
		IsCA:         true,
		SerialNumber: big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"Japan"},
			Organization: []string{"nananao.com"},
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
		EmailAddresses:        []string{"KryptoKings@random.com"},
	}

	// Create Certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &ca, &ca, &priv.PublicKey, priv)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}

	certOut, err := os.Create("certs/ca.pem")
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}); err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}

	keyOut, err := os.OpenFile("certs/ca-key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}
	defer keyOut.Close()

	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}

	if err := pem.Encode(keyOut, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}); err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}
}
