package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUsername contextKey = "username"
	ContextKeyRoles    contextKey = "roles"
)

type jwks struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
}

type KeycloakClaims struct {
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	PreferredUsername string `json:"preferred_username"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	keycloakURL   string
	keycloakRealm string
	httpClient    *http.Client
	keyCache      map[string]interface{}
}

func NewAuthMiddleware(keycloakURL, keycloakRealm string) *AuthMiddleware {
	return &AuthMiddleware{
		keycloakURL:   keycloakURL,
		keycloakRealm: keycloakRealm,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		keyCache:      make(map[string]interface{}),
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeUnauthorized(w, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			writeUnauthorized(w, "invalid authorization header format")
			return
		}

		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(tokenStr, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("missing kid in token header")
			}

			return m.getPublicKey(kid)
		})

		if err != nil || !token.Valid {
			writeUnauthorized(w, "invalid token")
			return
		}

		claims, ok := token.Claims.(*KeycloakClaims)
		if !ok {
			writeUnauthorized(w, "invalid token claims")
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.Subject)
		ctx = context.WithValue(ctx, ContextKeyUsername, claims.PreferredUsername)
		ctx = context.WithValue(ctx, ContextKeyRoles, claims.RealmAccess.Roles)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value(ContextKeyRoles).([]string)
			if !ok {
				writeForbidden(w)
				return
			}

			roleMap := make(map[string]bool)
			for _, role := range userRoles {
				roleMap[role] = true
			}

			for _, required := range roles {
				if roleMap[required] {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeForbidden(w)
		})
	}
}

func (m *AuthMiddleware) getPublicKey(kid string) (interface{}, error) {
	if key, ok := m.keyCache[kid]; ok {
		return key, nil
	}

	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", m.keycloakURL, m.keycloakRealm)

	resp, err := m.httpClient.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var keys jwks
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	for _, k := range keys.Keys {
		if k.Kid == kid {
			pubKey, err := parseRSAPublicKey(k.N, k.E)
			if err != nil {
				return nil, err
			}
			m.keyCache[kid] = pubKey
			return pubKey, nil
		}
	}

	return nil, fmt.Errorf("key not found for kid: %s", kid)
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
}
