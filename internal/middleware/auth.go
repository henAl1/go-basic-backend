package middleware

import (
    "net/http"
    "myapp/config"
)

const ApiKeyHeader = "X-API-Key"

type AuthMiddleware struct {
    config *config.Config
}

func NewAuthMiddleware(config *config.Config) *AuthMiddleware {
    return &AuthMiddleware{
        config: config,
    }
}

func (m *AuthMiddleware) RequireAPIKey(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get(ApiKeyHeader)
        if apiKey == "" {
            http.Error(w, "API key is missing", http.StatusUnauthorized)
            return
        }

        if apiKey != m.config.APIKey {
            http.Error(w, "Invalid API key", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    }
}
