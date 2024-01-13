package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

type UserClaim struct {
	jwt.RegisteredClaims
	ID string
}

func SetJWTSecretKey(secret string) {
	secretKey = []byte(secret)
}

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		ID: userID,
	})
	return token.SignedString(secretKey)
}

func GetClaimsFromToken(tokenString string) (UserClaim, error) {
	var userClaim UserClaim
	token, err := jwt.ParseWithClaims(tokenString, &userClaim, func(token *jwt.Token) (interface{}, error) {

		return secretKey, nil
	})

	if err != nil {
		return UserClaim{}, err
	}

	if token.Valid {
		return userClaim, nil
	}
	return UserClaim{}, err
}
