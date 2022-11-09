package centrifuge_api

import (
	"github.com/golang-jwt/jwt/v4"
)

const TokenHmacSecret = "46b38493-147e-4e3f-86e0-dc5ec54f5133"

type authTokenClaims struct {
	ID                   string   `json:"id"`   // 유저 ID
	Name                 string   `json:"name"` // 유저 이름
	Role                 []string `json:"role"` // 유저 역할
	Base                 string   `json:"base"` // base url
	Sub                  string   `json:"sub"`
	Channels             []string `json:"channels"`
	jwt.RegisteredClaims          // 표준 토큰 Claims
}

func IssueJWT(id, name, role, base string, channels []string) (string, error) {
	claims := authTokenClaims{
		ID:       id,
		Name:     name,
		Role:     []string{"user", "admin"},
		Base:     base,
		Sub:      id,
		Channels: channels, // client side subscription시 nil 전달
		// RegisteredClaims: jwt.RegisteredClaims{
		// 	// A usual scenario is to set the expiration time relative to the current time
		// 	// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 20)),
		// 	ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		// 	IssuedAt:  jwt.NewNumericDate(time.Now()),
		// 	NotBefore: jwt.NewNumericDate(time.Now()),
		// },
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(TokenHmacSecret))

	return ss, err
}
