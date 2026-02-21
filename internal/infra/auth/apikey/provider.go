package apikey

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"ppharma/backend/internal/domain/common"
)

type StaticSecretProvider struct {
	keys []common.InternalAPIKeySecret
}

func NewStaticSecretProvider(keys []common.InternalAPIKeySecret) *StaticSecretProvider {
	return &StaticSecretProvider{keys: keys}
}

func (s *StaticSecretProvider) InternalAPIKeys(_ context.Context) ([]common.InternalAPIKeySecret, error) {
	return s.keys, nil
}

type hashedSecret struct {
	keyID  string
	hash   [32]byte
	scopes map[string]struct{}
}

type Authenticator struct {
	secrets []hashedSecret
}

func NewAuthenticator(ctx context.Context, sp common.SecretProvider) (*Authenticator, error) {
	keys, err := sp.InternalAPIKeys(ctx)
	if err != nil {
		return nil, err
	}
	secrets := make([]hashedSecret, 0, len(keys))
	for _, key := range keys {
		if key.KeyID == "" || key.RawKey == "" {
			return nil, fmt.Errorf("invalid key config")
		}
		h := sha256.Sum256([]byte(key.RawKey))
		scopeMap := make(map[string]struct{}, len(key.Scopes))
		for _, scope := range key.Scopes {
			scopeMap[scope] = struct{}{}
		}
		secrets = append(secrets, hashedSecret{keyID: key.KeyID, hash: h, scopes: scopeMap})
	}
	return &Authenticator{secrets: secrets}, nil
}

func (a *Authenticator) Authenticate(_ context.Context, rawKey string) (*common.APIKeyPrincipal, error) {
	if rawKey == "" {
		return nil, fmt.Errorf("missing api key")
	}
	input := sha256.Sum256([]byte(rawKey))
	for _, s := range a.secrets {
		if subtle.ConstantTimeCompare(input[:], s.hash[:]) == 1 {
			return &common.APIKeyPrincipal{KeyID: s.keyID, Scopes: s.scopes}, nil
		}
	}
	return nil, fmt.Errorf("invalid api key: %s", hex.EncodeToString(input[:4]))
}
