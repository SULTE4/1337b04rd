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
			token, newClaims, err := h.userService.NewUser(r.FormValue("name"))
			if err != nil {
				h.logger.Error(err.Error())
				h.serverError(w, r, err)
				return
			}
			claims = newClaims
			setTokenCookie(w, token)
		} else {
			claims, err = util.VerifyJWT(cookie.Value, "supersecrethahaha")
			if err != nil {
				// invalid/expired → new token
				token, newClaims, err := h.userService.NewUser(r.FormValue("name"))
				if err != nil {
					h.logger.Error(err.Error())
					h.serverError(w, r, err)
					return
				}
				claims = newClaims
				setTokenCookie(w, token)
			} else {
				exists, err := h.userService.Exists(claims.UserID)
				if err != nil {
					h.logger.Error(err.Error())
					h.serverError(w, r, err)
					return
				}
				if !exists {
					token, newClaims, err := h.userService.NewUser(r.FormValue("name"))
					if err != nil {
						h.logger.Error(err.Error())
						h.serverError(w, r, err)
						return
					}
					claims = newClaims
					setTokenCookie(w, token)
				}
			}
		}
		// claims.UserID, err = h.userService.GetUserIDByToken(cookie.Value)
		// if err != nil {
		// 	h.logger.Error(err.Error())
		// 	h.serverError(w, r, err)
		// 	return
		// }

		if r.Method == http.MethodPost {
			if r.URL.Path == "/submit-post" {
				if err := r.ParseMultipartForm(10 << 20); err == nil {
					newName := r.FormValue("name")
					if newName != "" && newName != claims.Username {
						err := h.userService.UpdateUsername(claims.UserID, newName)
						if err != nil {
							h.logger.Error(err.Error())
							h.serverError(w, r, err)
							return
						}
						newToken, err := util.CreateJWT(
							newName,
							claims.Avatar,
							"supersecrethahaha",
							claims.UserID,
							24*time.Hour,
						)
						if err != nil {
							h.logger.Error(err.Error())
							h.serverError(w, r, err)
							return
						}
						// parse back into claims (so context has updated username too)
						updatedClaims, err := util.VerifyJWT(newToken, "supersecrethahaha")
						if err != nil {
							h.logger.Error(err.Error())
							h.serverError(w, r, err)
							return
						}
						claims = updatedClaims

						setTokenCookie(w, newToken)
					}
				}
			}
		}
		// attach claims to request context so handlers can use them
		ctx := context.WithValue(r.Context(), "user", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()), // 1 week
	})
}
