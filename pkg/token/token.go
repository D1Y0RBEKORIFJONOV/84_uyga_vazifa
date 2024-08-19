package tokens

import (
	"cors/internal/config"
	userentity "cors/internal/entity/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"net/http"
	"strings"
	"time"
)

func NewAccessToken(user *userentity.User) (string, error) {
	cfg := config.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.UserID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(cfg.Token.AccessTTL).Unix()
	claims["role"] = "user"

	tokenString, err := token.SignedString([]byte(config.Token()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func NewRefreshToken(user *userentity.User) (string, error) {
	cfg := config.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.UserID
	claims["exp"] = time.Now().Add(cfg.Token.RefreshTTL).Unix()
	claims["role"] = "user"
	claims["email"] = user.Email

	tokenString, err := token.SignedString([]byte(config.Token()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateTokens(user *userentity.User) (string, string, error) {
	accessToken, err := NewAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := NewRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ExtractClaim(tokenStr string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Token()), nil
	}
	token, err = jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, err
	}

	return claims, nil
}

func GetIdFromToken(r *http.Request) (string, int) {
	var softToken string
	token := r.Header.Get("Authorization")

	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := ExtractClaim(softToken)
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["uid"]), 0
}

func GetEmailFromToken(r *http.Request) (string, int) {
	var softToken string
	token := r.Header.Get("Authorization")

	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := ExtractClaim(softToken)
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["email"]), 0
}
