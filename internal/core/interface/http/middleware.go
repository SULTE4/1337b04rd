package http

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
)

// generateToken creates a random token string
func generateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "fallback-token"
	}
	return hex.EncodeToString(b)
}

func (h *Handler) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			// ip     = r.RemoteAddr
			// proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		// h.logger.Info("received request", slog.String("ip", ip), slog.String("proto", proto), slog.String("method", method), slog.String("uri", uri))
		h.logger.Info("received request", slog.String("method", method), slog.String("uri", uri))

		next.ServeHTTP(w, r)
	})
}

// TokenMiddleware ensures a user has a token cookie when creating post or comment
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("token")
		if err != nil {
			token := generateToken()
			// make new user etc.
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
			})
		}

		next.ServeHTTP(w, r)
	})
}
