package userrouter

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type userService interface {
	RegisterUser(login, password, email string) error
	LoginUser(login, password string) error
}

type userRouter struct {
	sessionStore *sessions.CookieStore
	userService  userService
}

type UserRouter interface {
	RegisterUserHandler(w http.ResponseWriter, r *http.Request)
	LoginUserHandler(w http.ResponseWriter, r *http.Request)
	LogoutUserHandler(w http.ResponseWriter, r *http.Request)
	CheckSession(r *http.Request) (string, bool)
}
