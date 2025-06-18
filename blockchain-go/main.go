package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	secretPhrase := "secret"
	emailBody := "Hey Bob"

	hash := getHashWithSecretPhrase(emailBody, secretPhrase)
	fmt.Printf("Input text: %s\n", emailBody)
	fmt.Printf("Secret phrase: %s\n", secretPhrase)
	fmt.Printf("Generated hash: %s\n\n", hash)

	fmt.Println("Verification examples:")

	isValid := verifyHash(emailBody, secretPhrase, hash)
	fmt.Printf("Verification with correct data: %v\n", isValid)

	tamperedEmail := emailBody + "!"
	isValid = verifyHash(tamperedEmail, secretPhrase, hash)
	fmt.Printf("Verification with tampered data: %v\n", isValid)

	wrongSecret := "wrong"
	isValid = verifyHash(emailBody, wrongSecret, hash)
	fmt.Printf("Verification with wrong secret: %v\n", isValid)
}

func getHashWithSecretPhrase(inputData, secretPhrase string) string {
	combined := inputData + secretPhrase

	hasher := sha256.New()

	hasher.Write([]byte(combined))

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func verifyHash(inputData, secretPhrase, hash string) bool {
	calculatedHash := getHashWithSecretPhrase(inputData, secretPhrase)

	return calculatedHash == hash
}
