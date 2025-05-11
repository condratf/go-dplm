package userrouter

import (
	"encoding/json"
	"net/http"

	"github.com/condratf/go-musthave-diploma-tpl/internal/errors_custom"
	"github.com/gorilla/sessions"
)

func NewUserRouter(sessionStore *sessions.CookieStore, userService userService) UserRouter {
	return &userRouter{
		sessionStore: sessionStore,
		userService:  userService,
	}
}

func (u *userRouter) setSession(w http.ResponseWriter, r *http.Request, login string) {
	session, _ := u.sessionStore.Get(r, "session")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 10,
		HttpOnly: true,
	}
	session.Values["user"] = login
	session.Save(nil, w)
}

func (u *userRouter) CheckSession(r *http.Request) (string, bool) {
	session, err := u.sessionStore.Get(r, "session")
	if err != nil {
		return "", false
	}
	login, ok := session.Values["user"].(string)
	return login, ok
}

func (u *userRouter) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := u.userService.RegisterUser(req.Login, req.Password)
	if err == errors_custom.ErrUserExists {
		http.Error(w, "Login already exists", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	u.setSession(w, r, req.Login)
	w.WriteHeader(http.StatusOK)
}

func (u *userRouter) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := u.userService.LoginUser(req.Login, req.Password)
	if err == errors_custom.ErrInvalidAuth {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	u.setSession(w, r, req.Login)
	w.WriteHeader(http.StatusOK)
}

func (u *userRouter) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := u.sessionStore.Get(r, "session")
	delete(session.Values, "user")
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
}
