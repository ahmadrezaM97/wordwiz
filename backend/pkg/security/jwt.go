package security

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT ...
func GenerateJWT(m map[string]interface{}, tokenExpireTime time.Duration, tokenSecretKey string) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	for key, value := range m {
		claims[key] = value
	}

	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(tokenExpireTime).Unix()

	tokenString, err = token.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractClaims extracts claims from given token
func ExtractClaims(tokenString string, tokenSecretKey string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

type TokenInfo struct {
	ID     string
	UserID string
	IP     string
}

func ParseClaims(token string, secretKey string) (result TokenInfo, err error) {
	var ok bool
	var claims jwt.MapClaims

	claims, err = ExtractClaims(token, secretKey)
	if err != nil {
		return result, err
	}

	result.ID, ok = claims["id"].(string)
	if !ok {
		err = errors.New("cannot parse 'id' field")
		return result, err
	}

	result.UserID, ok = claims["user_id"].(string)
	if !ok {
		err = errors.New("cannot parse 'user_id' field")
		return result, err
	}

	result.IP, ok = claims["ip"].(string)
	if !ok {
		err = errors.New("cannot parse 'ip' field")
		return result, err
	}

	return
}
