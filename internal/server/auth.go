package server

import (
	"net/http"

	"github.com/kudras3r/KDTog/pkg/logger"
)

func authMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	// ! TODO
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")

		if err != nil || cookie.Value != "asdasd" {
			log.Errorf("auth error : %v", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
