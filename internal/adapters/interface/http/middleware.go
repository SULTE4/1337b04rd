package http

import (
	"1337b04rd/internal/core/util"
	"context"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		// h.logger.Info("received request", slog.String("ip", ip), slog.String("proto", proto), slog.String("method", method), slog.String("uri", uri))
		h.logger.Info("received request", slog.String("method", method), slog.String("uri", uri))

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		var claims *util.Claims

		if err != nil { // no token → new user
			token, newClaims, err := h.userService.NewUser(r)
			if err != nil {
				h.serverError(w, r, err)
				return
			}
			claims = newClaims
			setTokenCookie(w, token)
		} else {
			claims, err = util.VerifyJWT(cookie.Value, "supersecrethahaha")
			if err != nil {
				// invalid/expired → new token
				token, newClaims, err := h.userService.NewUser(r)
				if err != nil {
					h.serverError(w, r, err)
					return
				}
				claims = newClaims
				setTokenCookie(w, token)
			}
		}

		// attach claims to request context so handlers can use them
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()), // 1 week
	})
}
