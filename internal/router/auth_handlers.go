package router

import (
	"encoding/json"
	"net/http"

	"github.com/condratf/go-musthave-diploma-tpl/internal/repository"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func (appRouter *AppRouter) setSession(w http.ResponseWriter, r *http.Request, login string) {
	session, _ := appRouter.sessionStore.Get(r, "session")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 10,
		HttpOnly: true,
	}
	session.Values["user"] = login
	session.Save(nil, w)
}

func (appRouter *AppRouter) checkSession(r *http.Request) (string, bool) {
	session, err := appRouter.sessionStore.Get(r, "session")
	if err != nil {
		return "", false
	}
	login, ok := session.Values["user"].(string)
	return login, ok
}

func (appRouter *AppRouter) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	err = appRouter.userRepo.CreateUser(req.Login, string(hashedPassword))
	if err == repository.ErrUserExists {
		http.Error(w, "Login already exists", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	appRouter.setSession(w, r, req.Login)
	w.WriteHeader(http.StatusOK)
}

func (appRouter *AppRouter) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := appRouter.userRepo.GetUserPassword(req.Login)
	if err == repository.ErrUserNotFound || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	appRouter.setSession(w, r, req.Login)
	w.WriteHeader(http.StatusOK)
}

func (appRouter *AppRouter) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := appRouter.sessionStore.Get(r, "session")
	delete(session.Values, "user")
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
}
