package server

import (
	"encoding/base64"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func basicAuthHandleWrap(next HandlerFunc, username, password string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, log *zap.Logger) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("No Authorization Basic"))
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[len("Basic "):])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != username || pair[1] != password {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Authorization Basic is not match"))
			return
		}
		next(w, r, log)
	}
}

func tokenAuthHandleWrap(next HandlerFunc, token string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, log *zap.Logger) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.Header().Set("WWW-Authenticate", `Bearer realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("No Authorization Bearer"))
			return
		}
		if auth[len("Bearer "):] != token {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Authorization Bearer is not match"))
			return
		}
		next(w, r, log)
	}
}
