package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

func Encrypt(plaintext string) (string, error) {
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		return "", err
	}

	encryptedText, err := encrypt(plaintext, secretKey)
	if err != nil {
		return "", err
	}

	encodedText := fmt.Sprintf("%s.%s", hex.EncodeToString(encryptedText), hex.EncodeToString(secretKey))
	return encodedText, nil
}

func Decrypt(encodedText string) (string, error) {
	textTokens := strings.Split(encodedText, ".")
	if len(textTokens) < 2 {
		return "", errors.New("invalid hash provided")
	}

	secretKey, err := hex.DecodeString(textTokens[1])
	if err != nil {
		return "", err
	}
	encryptedText, err := hex.DecodeString(textTokens[0])
	if err != nil {
		return "", err
	}

	plainText, err := decrypt(string(encryptedText), secretKey)
	if err != nil {
		return "", err
	}
	return plainText, nil
}

func encrypt(plaintext string, secretKey []byte) ([]byte, error) {
	aesClient, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesClient)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

func decrypt(ciphertext string, secretKey []byte) (string, error) {
	aesClient, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aesClient)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
