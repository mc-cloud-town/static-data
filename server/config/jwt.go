package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var (
	jwtPublicKey  *ecdsa.PublicKey
	jwtPrivateKey *ecdsa.PrivateKey
)

func GetJwtPrivateKey() *ecdsa.PrivateKey {
	return jwtPrivateKey
}

func GetJwtPublicKey() *ecdsa.PublicKey {
	return jwtPublicKey
}

func LoadPrivateKey(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return fmt.Errorf("failed to decode private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	jwtPrivateKey = key
	return nil
}

func LoadPublicKey(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return fmt.Errorf("failed to decode public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	jwtPublicKey = key.(*ecdsa.PublicKey)
	return nil
}
