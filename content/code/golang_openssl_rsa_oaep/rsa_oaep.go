package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	// ErrPEMDecode is used to be propagated up the call stack.
	// It is returned when something goes wrong in the crypt module.
	ErrPEMDecode = fmt.Errorf("failed to decode PEM block containing private key")

	reader = rand.Reader //nolint:golint,gochecknoglobals
	empty  = []byte("")  //nolint:golint,gochecknoglobals
)

// EncryptRSAOAEP encrypt some data using the rsa public key argument.
// It uses the rsa.EncryptOAEP function under the hood
func EncryptRSAOAEP(key *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), reader, key, data, empty)
}

// DecodeRSAKey decode a slice of bytes into an rsa.PrivateKey.
// It uses the x509.ParsePKCS1PrivateKey function under the hood
func DecodeRSAKey(bytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, ErrPEMDecode
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// ReadRSAKeyFromFile read a file and decode the content into a rsa.PrivateKey
func ReadRSAKeyFromFile(path string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	return DecodeRSAKey(bytes)
}

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read content from stdin")
	}
	key, err := ReadRSAKeyFromFile("rsa.priv")
	if err != nil {
		log.Fatalf("failed to read RSA key %q", err)
	}
	cipher, err := EncryptRSAOAEP(&key.PublicKey, bytes)
	if err != nil {
		log.Fatalf("failed to encrypt data %q", err)
	}
	fmt.Printf("%s", cipher)
}
