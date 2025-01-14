package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zechao158/ecomm/config"
)

func CreateJWT(secret []byte, userID uuid.UUID) (string, error) {
	expiration := time.Second * time.Duration(config.ENVs.JWTExpirationSecoond)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID.String(),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	return token.SignedString(secret)
}
