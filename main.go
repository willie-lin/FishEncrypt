package main

import (
	"FishEncrypt/pkg/fish"
	"log"
)
func main() {

	key := "secret key"
	message := "汉字"

	enc, err := fish.Encrypt(key, message)
	if err != nil {
		log.Fatalf("Failed to encrypt: %s", err)
	}
	log.Printf("Encrypted: %s => %s", message, enc)

	dec, err := fish.Decrypt(key, enc)
	if err != nil {
		log.Fatalf("Failed to decrypt: %s", err)
	}
	log.Printf("Decrypted: %s => %s", enc, dec)
	
}
