package middlewares

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"log"
	"net/http"
	"strings"
)

type Decrypter struct {
	privateKey *rsa.PrivateKey
}

func NewDecrypter(key *rsa.PrivateKey) *Decrypter {
	return &Decrypter{
		privateKey: key,
	}
}

func (d *Decrypter) Decrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			var body []byte
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("failed to read request body: %v", err)
				return
			}
			decryptedBody := make([]byte, 0, len(body))
			var nextBlockLength int
			for i := 0; i < len(body); i += d.privateKey.PublicKey.Size() {
				nextBlockLength = i + d.privateKey.PublicKey.Size()
				if nextBlockLength > len(body) {
					nextBlockLength = len(body)
				}
				block, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, d.privateKey, body[i:nextBlockLength], []byte("yandex"))
				if err != nil {
					log.Printf("failed to decrypt request body: %v", err)
					return
				}
				decryptedBody = append(decryptedBody, block...)
			}
			r.Body = io.NopCloser(strings.NewReader(string(decryptedBody)))
		}
		next.ServeHTTP(w, r)
	})
}
