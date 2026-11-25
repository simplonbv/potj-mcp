package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
)

type contextKey string

// Exported so mcpserver/handler.go can read it from context
const OrganizationContextKey contextKey = "organization"

// Validate the JWT, logs the sub claim, and extract organization from
// scope for downstream use
func (v *Validator) NewTokenVerifier() sdkauth.TokenVerifier {
	return func(ctx context.Context, token string, req *http.Request) (*sdkauth.TokenInfo, error) {
		parsed, err := v.ValidateToken(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", sdkauth.ErrInvalidToken, err)
		}

		// sub claim is a traceable reference to a real user identity
		// Log the subject (not the token) on every authenticated request for auditablity
		log.Printf("authenticated request: sub=%s path=%s", parsed.Subject(), req.URL.Path)

		scopes := ExtractScopes(parsed)
		organization := DepartmentFromScopes(scopes)

		return &sdkauth.TokenInfo{
			UserID:     parsed.Subject(),
			Scopes:     scopes,
			Expiration: parsed.Expiration(),
			Extra:      map[string]any{"organization": organization},
		}, nil
	}
}

// Middleware that validates bearer tokens and enforces department scopes.
// The SDK's RequireBearerToken handles 401 and WWW-Authenticate headers.
// The inner handler enforces "at least one organization scope" requirement
// and stores the organization in context for the per-session MCP server factory.
func (v *Validator) Protect(resourceMetadataURL string, scopes []string, next http.Handler) http.Handler {
	verifier := v.NewTokenVerifier()
	sdkMiddleware := sdkauth.RequireBearerToken(verifier, &sdkauth.RequireBearerTokenOptions{
		ResourceMetadataURL: resourceMetadataURL,
	})

	return sdkMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenInfo := sdkauth.TokenInfoFromContext(r.Context())
		if tokenInfo == nil {
			http.Error(w, `{"error": "forbidden", "message": "no token info in context"}`, http.StatusForbidden)
			return
		}

		// Require at least one organization scope
		hasScope := false
		for _, required := range scopes {
			for _, granted := range tokenInfo.Scopes {
				if granted == required {
					hasScope = true
					break
				}
			}
			if hasScope {
				break
			}
		}

		if !hasScope {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(
				`Bearer resource_metadata="%s", scope="%s"`,
				resourceMetadataURL,
				strings.Join(scopes, " "),
			))
			http.Error(w, `{"error": "forbidden", "message": "insufficient scope"}`, http.StatusForbidden)
			return
		}

		// Extract department from TokenInfo and store in context for the MCP handler
		if org, ok := tokenInfo.Extra["organization"].(string); ok {
			r = r.WithContext(context.WithValue(r.Context(), OrganizationContextKey, org))
		}

		next.ServeHTTP(w, r)
	}))
}
