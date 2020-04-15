package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

const (
	keySize        = 2048
	typePrivateKey = "RSA PRIVATE KEY"
)

func generateKey() (*rsa.PrivateKey, ssh.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}

func marshalRSAPrivate(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type: typePrivateKey, Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
}

func marshalRSAPublic(pub ssh.PublicKey) []byte {
	return bytes.TrimSuffix(ssh.MarshalAuthorizedKey(pub), []byte{'\n'})
}

func unmarshalRSAPrivate(bytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func unmarshalRSAPublic(bytes []byte) (ssh.PublicKey, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	return pub, err
}

func generateCert(pub ssh.PublicKey) *ssh.Certificate {
	permissions := ssh.Permissions{
		CriticalOptions: map[string]string{},
		Extensions:      map[string]string{"permit-agent-forwarding": ""},
	}
	return &ssh.Certificate{
		CertType: ssh.UserCert, Permissions: permissions, Key: pub,
	}
}

func generateSignerFromKey(priv *rsa.PrivateKey) (ssh.Signer, error) {
	return ssh.NewSignerFromKey(priv)
}

func generateSignerFromBytes(bytes []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKey(bytes)
}

func generateAndSign() (*rsa.PrivateKey, *ssh.Certificate, error) {
	priv, pub, err := generateKey()
	if err != nil {
		return nil, nil, err
	}
	signer, err := generateSignerFromKey(priv)
	if err != nil {
		return nil, nil, err
	}
	cert := generateCert(pub)
	return priv, cert, cert.SignCert(rand.Reader, signer)
}

func marshalCert(cert *ssh.Certificate) []byte {
	return ssh.MarshalAuthorizedKey(cert)
}

func unmarshalCert(bytes []byte) (*ssh.Certificate, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		return nil, err
	}
	cert, ok := pub.(*ssh.Certificate)
	if !ok {
		return nil, fmt.Errorf("failed to cast to certificate")
	}
	return cert, nil
}

func main() {
	priv, pub, err := generateKey()
	if err != nil {
		log.Fatalf("failed to generate RSA keys %q", err)
	}
	if _, err := unmarshalRSAPrivate(marshalRSAPrivate(priv)); err != nil {
		log.Fatalf("failed to marshal unmarshal private key: %q", err)
	}
	if _, err := unmarshalRSAPublic(marshalRSAPublic(pub)); err != nil {
		log.Fatalf("failed to marshal unmarshal public key: %q", err)
	}

	_, cert, err := generateAndSign()
	if err != nil {
		log.Fatalf("failed to generate and sign certificate: %q", err)
	}
	if _, err := unmarshalCert(marshalCert(cert)); err != nil {
		log.Fatalf("failed to marshal unmarshal certificate: %q", err)
	}
}
