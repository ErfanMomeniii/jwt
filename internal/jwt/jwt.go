package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mehditeymorian/jwt/internal/model"
)

var errInvalidSigningMethod = errors.New("signing method dosn't match with configuration")

func Encode(encode model.Encode, key any) (string, error) {
	exp, _ := time.ParseDuration(encode.Expiration)

	claims := jwt.MapClaims{
		"iss": encode.Issuer,
		"exp": time.Now().Add(exp).Unix(),
		"sub": encode.Subject,
		"aud": encode.Audience,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(encode.Algorithm), claims)

	signedString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign claims with key: %w", err)
	}

	return signedString, nil
}

func Decode(strToken string, key any, signingMethod string) (*jwt.Token, error) {
	token, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != signingMethod {
			return nil, errInvalidSigningMethod
		}

		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil
}
