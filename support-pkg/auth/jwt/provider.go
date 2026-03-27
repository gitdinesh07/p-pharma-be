package jwt

import (
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"ppharma/backend/internal/domain/common"
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
	if sub == "" || role == "" {
		return nil, fmt.Errorf("missing required claims")
	}
	return &common.Principal{ID: sub, Role: role}, nil
}

func (p *Provider) GenerateToken(principal *common.Principal, expiry time.Duration) (string, error) {
	claims := jwtv5.MapClaims{
		"sub":  principal.ID,
		"role": principal.Role,
		"exp":  time.Now().Add(expiry).Unix(),
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}
