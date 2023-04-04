package rest

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
)

const blockLength = 128

type Encryptor struct {
	publicKey *rsa.PublicKey
}

func NewEncryptor(publicKey *rsa.PublicKey) *Encryptor {
	return &Encryptor{
		publicKey: publicKey,
	}
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	if e.publicKey == nil {
		return data, nil
	}
	encrypted := make([]byte, 0, len(data))
	var nextBlockLength int
	for i := 0; i < len(data); i += blockLength {
		nextBlockLength = i + blockLength
		if nextBlockLength > len(data) {
			nextBlockLength = len(data)
		}
		block, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, e.publicKey, data[i:nextBlockLength], []byte("yandex"))
		if err != nil {
			return nil, err
		}
		encrypted = append(encrypted, block...)
	}
	log.Printf("Encrypted data '%s': %s", string(data), encrypted)
	return encrypted, nil
}
