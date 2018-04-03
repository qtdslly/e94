package rsa

import (
	"common/logger"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

/*
	A Signer is can create signatures with private key.
*/
type Signer interface {
	Sign(data []byte) ([]byte, error)
}

/*
	A Verifier is verify a signature with the public key.
*/
type Verifier interface {
	Verify(data []byte, sig []byte) error
}

type RsaPublicKey struct {
	*rsa.PublicKey
}

type RsaPrivateKey struct {
	*rsa.PrivateKey
}

/*
	Sign signs data with rsa-sha1
*/
func (r *RsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha1.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA1, d)
}

/*
	Unsign verifies the message using a rsa-sha1 signature
*/
func (r *RsaPublicKey) Verify(message []byte, sig []byte) error {
	h := sha1.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA1, d, sig)
}

/*
	loads and parses a PEM encoded private key file to a Verifier.
*/
func LoadPublicKey(path string) (Verifier, error) {
	pubkey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// parse the key file to a Verifier
	block, _ := pem.Decode(pubkey)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	// TODO, pay attention, thsi could be changed
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newUnsignerFromKey(rawkey)
}

func LoadPublicKeyFromStr(path string) (Verifier, error) {
	// parse the key file to a Verifier
	block, _ := pem.Decode([]byte(path))
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	// TODO, pay attention, thsi could be changed
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newUnsignerFromKey(rawkey)
}

/*
	loads an parses a PEM encoded private key file to a signer
*/
func LoadPrivateKey(path string) (Signer, error) {
	pubkey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubkey)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	// TODO, pay attention, this cloud be chang
	case "PRIVATE KEY":
		rsa, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

func LoadPrivateKeyFromStr(str string) (Signer, error) {
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	// TODO, pay attention, this cloud be chang
	case "PRIVATE KEY":
		rsa, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &RsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

func newUnsignerFromKey(k interface{}) (Verifier, error) {
	var sshKey Verifier
	switch t := k.(type) {
	case *rsa.PublicKey:
		sshKey = &RsaPublicKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}
