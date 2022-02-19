package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io"
)

type Credential struct {
	Key string
}

// shaKey
func (cred *Credential) shaKey() (res []byte) {
	h := sha1.New()
	h.Write([]byte(cred.Key))
	res = h.Sum(nil)
	res = res[:16]

	return res
}

// Encrypt ...
func (cred *Credential) Encrypt(textString string) (string, error) {
	key := cred.shaKey()
	text := []byte(textString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	b := base64.StdEncoding.EncodeToString(text)
	cipherText := make([]byte, aes.BlockSize+len(b))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(b))

	res := hex.EncodeToString(cipherText)
	return res, nil
}

// EncryptString ...
func (cred *Credential) EncryptString(textString string) string {
	res, _ := cred.Encrypt(textString)

	return res
}
