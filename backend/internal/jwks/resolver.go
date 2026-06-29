package jwks

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
)

type jwk struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type Resolver struct {
	mu   sync.RWMutex
	url  string
	keys map[string]*rsa.PublicKey
}

func NewResolver(jwksURL string) *Resolver {
	return &Resolver{
		url:  jwksURL,
		keys: make(map[string]*rsa.PublicKey),
	}
}

// GetKey returns the RSA public key for the given kid.
// It uses the in-memory cache and only fetches from Microsoft if the kid is missing.
func (r *Resolver) GetKey(kid string) (*rsa.PublicKey, error) {
	r.mu.RLock()
	key, ok := r.keys[kid]
	r.mu.RUnlock()
	if ok {
		return key, nil
	}

	if err := r.fetch(); err != nil {
		return nil, fmt.Errorf("jwks: failed to refresh keys: %w", err)
	}

	r.mu.RLock()
	key, ok = r.keys[kid]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("jwks: key with kid %q not found in Microsoft JWKS", kid)
	}
	return key, nil
}

func (r *Resolver) fetch() error {
	resp, err := http.Get(r.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var payload jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}

	fresh := make(map[string]*rsa.PublicKey, len(payload.Keys))
	for _, k := range payload.Keys {
		pub, err := buildRSAPublicKey(k.N, k.E)
		if err != nil {
			continue
		}
		fresh[k.Kid] = pub
	}

	r.mu.Lock()
	r.keys = fresh
	r.mu.Unlock()
	return nil
}

func buildRSAPublicKey(nB64, eB64 string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, fmt.Errorf("jwks: failed to decode modulus: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, fmt.Errorf("jwks: failed to decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	var eInt int
	for _, b := range eBytes {
		eInt = eInt<<8 + int(b)
	}

	return &rsa.PublicKey{N: n, E: eInt}, nil
}
