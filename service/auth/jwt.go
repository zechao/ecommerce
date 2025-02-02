package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zechao158/ecomm/config"
	httputil "github.com/zechao158/ecomm/http"
	"github.com/zechao158/ecomm/types"
)

const (
	UserIDKey    = "userID"
	expiredAtKey = "expiredAt"
)

func CreateJWT(secret []byte, userID uuid.UUID) (string, error) {
	expiration := time.Second * time.Duration(config.ENVs.JWTExpirationSecoond)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		UserIDKey:    userID.String(),
		expiredAtKey: time.Now().Add(expiration).Unix(),
	})

	return token.SignedString(secret)
}

func AuthMiddleware(store types.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.WriteError(w, http.StatusUnauthorized, fmt.Errorf("missing Authorization header"))
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				httputil.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid Authorization header"))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := getJWTMapClaims(tokenString)
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			userID, ok := claims[UserIDKey].(string)
			if !ok {
				httputil.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
				return
			}

			user, err := store.GetByID(r.Context(), uuid.MustParse(userID), false)
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, fmt.Errorf("token invalid"))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, user)
			r = r.WithContext(ctx)

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func UserFromContext(ctx context.Context) (*types.User, bool) {
	user, ok := ctx.Value(UserIDKey).(*types.User)
	return user, ok
}

func getJWTMapClaims(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ENVs.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	expiredAt := int64(claims[expiredAtKey].(float64))
	if time.Now().Unix() > expiredAt {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
