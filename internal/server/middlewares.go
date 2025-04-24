package server

import (
	"net/http"

	"github.com/kudras3r/KDTog/pkg/logger"
)

var tokens map[string]string

func authMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc := GLOC + "authMiddleware()"
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value != "123" {
			log.Warnf("auth error at %s: %v", loc, err)
			http.Redirect(w, r, "/auth", http.StatusPermanentRedirect)
			return
		}
		log.Infof("%s checked auth_token %s for %s", loc, cookie.Value, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
