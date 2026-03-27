package jwt

import (
	"fmt"
	"time"

	"ppharma/backend/internal/domain/common"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Provider struct {
	secret []byte
}

func NewProvider(secret string) *Provider {
	return &Provider{secret: []byte(secret)}
}

func (p *Provider) ParseAccessToken(token string) (*common.Principal, error) {
	parsed, err := jwtv5.Parse(token, func(t *jwtv5.Token) (any, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return p.secret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := parsed.Claims.(jwtv5.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	sub, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	email, _ := claims["email"].(string)
	mobile, _ := claims["mobile"].(string)
	if sub == "" || role == "" {
		return nil, fmt.Errorf("missing required claims")
	}
	return &common.Principal{ID: sub, Role: role, Email: email, Mobile: mobile}, nil
}

func (p *Provider) GenerateToken(principal *common.Principal, expiry time.Duration) (string, error) {
	claims := jwtv5.MapClaims{
		"sub":  principal.ID,
		"role": principal.Role,
		"exp":  time.Now().Add(expiry).Unix(),
	}
	if principal.Email != "" {
		claims["email"] = principal.Email
	}
	if principal.Mobile != "" {
		claims["mobile"] = principal.Mobile
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}
