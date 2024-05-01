package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/models"
	"github.com/gleb-korostelev/gophermart.git/internal/service/utils"
	"github.com/gleb-korostelev/gophermart.git/tools/logger"
)

func (svc *APIService) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		logger.Infof("Invalid Content-Type")
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		logger.Infof("Bad request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if newUser.Login == "" || newUser.Password == "" {
		logger.Infof("Login and password required: %v", err)
		http.Error(w, "Login and password required", http.StatusBadRequest)
		return
	}
	err = svc.store.Auth(context.Background(), newUser)
	if err != nil {
		if err == config.ErrGone {
			logger.Infof("This user was deleted: %v", err)
			http.Error(w, config.ErrGone.Error(), http.StatusGone)
			return
		}
		logger.Infof("Internal server error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	utils.SetJWTInCookie(w, newUser.Login)
	w.WriteHeader(http.StatusOK)
	logger.Infof("User authenticated")
	w.Write([]byte("User authenticated"))
}
