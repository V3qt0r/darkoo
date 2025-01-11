package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"log"

	"darkoo/apperrors"

	"golang.org/x/crypto/scrypt"
)


func hashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)

	if err != nil {
		log.Print("Error creating byte slice")
		return "", apperrors.NewInternal()
	}

	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)

	if err != nil {
		log.Print("Error hashing password")
		return "", apperrors.NewInternal()
	}

	hashedPassword := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPassword, nil
}


func comparePassword(storedPassword, providedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")

	if len(pwsalt) < 2 {
		log.Print("Did not provide a valid hash")
		return false, fmt.Errorf("Did not provide a valid hash")
	}

	salt, err := hex.DecodeString(pwsalt[1])

	if err != nil {
		log.Print("Unable to verify password")
		return false, fmt.Errorf("Unable to verify password")
	}

	shash, err := scrypt.Key([]byte(providedPassword), salt, 32768, 8, 1, 32)

	if err != nil {
		log.Print("Error hashing supplied password")
		return false, fmt.Errorf("Error hashing supplied password")
	}

	return hex.EncodeToString(shash) == pwsalt[0], nil
}