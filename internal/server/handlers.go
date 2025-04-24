package server

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kudras3r/KDTog/internal/chat/ws"
	"github.com/kudras3r/KDTog/internal/storage"
	"github.com/kudras3r/KDTog/pkg/logger"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func indexHandler(log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("index handler called")
		http.ServeFile(w, r, "static/index.html")
	}
}

func authHandler(log *logger.Logger, db storage.Storage) func(http.ResponseWriter, *http.Request) {
	loc := GLOC + "getAuthHandler()"
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			log.Infof("GET at /auth: %s", r.RemoteAddr)
			http.ServeFile(w, r, "static/auth.html")

		case "POST":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Errorf("cannot read body at: %s: %v", loc, err)
				http.Error(w, "cannot read body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var creds Credentials
			if err = json.Unmarshal(body, &creds); err != nil {
				log.Errorf("cannot unmarshal json at: %s: %v", loc, err)
				http.Error(w, "cannot parse json", http.StatusBadRequest)
				return
			}

			phash := sha256.Sum256([]byte(creds.Password))
			uid, err := db.GetIDByName(creds.Username)
			if err != nil {
				log.Warnf("user with username %s not found at: %s", creds.Username, loc)
				http.Error(w, "invalid auth for user"+creds.Username, http.StatusBadRequest)
				return
			}
			if phash == db.GetPHashByID(uid) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"status":"success","message":"Успешный вход"}`)
			}

		}
	}
}

func WSHandler(hub *ws.Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	}
}
