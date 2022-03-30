package notifier

import (
	"crypto/rand"
	"encoding/hex"
)

// func GenerateToken(email string) string {
// 	bcrypt.Gene
// 	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Hash to store:", string(hash))

// 	return string(hash)

// }

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
