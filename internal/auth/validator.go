package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Validator struct {
	jwksUrl   string
	jwksCache *jwk.Cache
	audience  string
	issuer    string
}

func NewValidator(ctx context.Context, issuer, jwksUrl, audience string) (*Validator, error) {
	cache := jwk.NewCache(ctx)
	if err := cache.Register(jwksUrl); err != nil {
		return nil, fmt.Errorf("failed to register JWKS URL: %w", err)
	}

	// Eagerly fetch JWKs on startup.
	// Fails fast if unreachable
	if _, err := cache.Refresh(ctx, jwksUrl); err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	return &Validator{
		jwksUrl:   jwksUrl,
		jwksCache: cache,
		audience:  audience,
		issuer:    issuer,
	}, nil
}

func (v *Validator) ValidateToken(ctx context.Context, tokenString string) (jwt.Token, error) {

	keySet, err := v.jwksCache.Get(ctx, v.jwksUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}

	token, err := jwt.Parse([]byte(tokenString),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}

func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil
}

// Returns individual scopes from token's scope claim.
// Scope claim represents what was explicitly requested
// and approved for the session
func ExtractScopes(token jwt.Token) []string {
	scopeVal, ok := token.Get("scope")
	if !ok {
		return nil
	}
	return strings.Fields(fmt.Sprintf("%v", scopeVal))
}

// Map a department scope to a department name
func DepartmentFromScopes(scopes []string) string {
	for _, s := range scopes {
		switch s {
		case "read:engineering":
			return "engineering"
		case "read:hr":
			return "hr"
		}
	}
	return ""
}
