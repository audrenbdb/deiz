package crypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
)

const errPassPhraseNotProvided Error = "passphrase not found"

type Error string

func (e Error) Error() string {
	return string(e)
}

type service struct{}

func getHashedPassphrase() (string, error) {
	passphrase := os.Getenv("CIPHER_PASSPHRASE")
	if passphrase == "" {
		return "", errPassPhraseNotProvided
	}
	return createHash(passphrase), nil
}

//createHash encodes string with md5 hash.
func createHash(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func (d *service) EncryptString(ctx context.Context, str string) ([]byte, error) {
	var b []byte

	hashedPhrase, err := getHashedPassphrase()
	if err != nil {
		return b, err
	}

	block, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		return b, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return b, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return b, err
	}

	strBytes := gcm.Seal(nonce, nonce, []byte(str), nil)
	return strBytes, err
}

func (d *service) DecryptBytes(ctx context.Context, strBytes []byte) (string, error) {
	hashedPhrase, err := getHashedPassphrase()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := strBytes[:nonceSize], strBytes[nonceSize:]
	b, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NewService() *service {
	return &service{}
}
