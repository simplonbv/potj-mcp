package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	config      *oauth2.Config
	audience    string
	stateCookie string
	pkceCookie  string
}

func NewOAuthHandler(domain, clientID, clientSecret, callbackURL, audience string) *OAuthHandler {
	endpoint := oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("https://%s/authorize", domain),
		TokenURL: fmt.Sprintf("https://%s/oauth/token", domain),
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint:     endpoint,
	}

	return &OAuthHandler{
		config:      cfg,
		audience:    audience,
		stateCookie: "oauth_state",
		pkceCookie:  "pkce_cookie",
	}
}

// Redirect the user to Auth0's hosted logi page with PKCE
func (o *OAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestedScope := r.URL.Query().Get("scope")
	scopes := []string{"openid", "profile"}
	if requestedScope != "" {
		scopes = append(scopes, strings.Split(requestedScope, " ")...)
	}

	state, err := generateRandom()
	if err != nil {
		http.Error(w, "failed to generate state", http.StatusInternalServerError)
		return
	}

	verifier, err := generateRandom()
	if err != nil {
		http.Error(w, "failed to generate PKCE verifier", http.StatusInternalServerError)
		return
	}
	challenge := generateS256Challenge(verifier)

	// Store state cookie to validate on callback.
	// Prevents CSRF
	http.SetCookie(w, &http.Cookie{
		Name:     o.stateCookie,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Store PKCE verifier in cookie for use during token exchange
	http.SetCookie(w, &http.Cookie{
		Name:     o.pkceCookie,
		Value:    verifier,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	authURL := o.config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("audience", o.audience),
		oauth2.SetAuthURLParam("scope", strings.Join(scopes, " ")),
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Exchange the authorization code for a token
func (o *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie((o.stateCookie))
	if err != nil {
		http.Error(w, "missing state cookie", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "invalid state parameter", http.StatusBadRequest)
		return
	}

	pkceVerifier, err := r.Cookie(o.pkceCookie)
	if err != nil {
		http.Error(w, "missing PKCE verifier cookie", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	token, err := o.config.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", pkceVerifier.Value),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to exchange code: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token.AccessToken,
		"token_type":   token.TokenType,
	})
}

func generateRandom() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateS256Challenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
