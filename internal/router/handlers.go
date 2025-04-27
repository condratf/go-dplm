package router

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/condratf/go-musthave-diploma-tpl/internal/repository"
	"github.com/condratf/go-musthave-diploma-tpl/internal/utils"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func (app *App) setSession(w http.ResponseWriter, login string) {
	session, _ := app.sessionStore.Get(nil, "session")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 10,
		HttpOnly: true,
	}
	session.Values["user"] = login
	session.Save(nil, w)
}

func (app *App) registerUserHandler(w http.ResponseWriter, r *http.Request) {
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

	err = app.userRepo.CreateUser(req.Login, string(hashedPassword))
	if err == repository.ErrUserExists {
		http.Error(w, "Login already exists", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	app.setSession(w, req.Login)
	w.WriteHeader(http.StatusOK)
}
func (app *App) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := app.userRepo.GetUserPassword(req.Login)
	if err == repository.ErrUserNotFound || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	app.setSession(w, req.Login)
	w.WriteHeader(http.StatusOK)
}

func (app *App) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "session")
	delete(session.Values, "user")
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
}

func (app *App) uploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "session")
	login, ok := session.Values["user"].(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	order := strings.TrimSpace(string(bodyBytes))

	if !regexp.MustCompile(`^\d+$`).MatchString(order) {
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	if !utils.IsValidLuhn(order) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	err = app.ordersRepo.UploadOrder(login, order)
	if err != nil {
		if errors.Is(err, repository.ErrOrderAlreadyUploadedBySameUser) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, repository.ErrOrderAlreadyUploadedByAnotherUser) {
			http.Error(w, "Order already uploaded by another user", http.StatusConflict)
			return
		}
		http.Error(w, "Error uploading order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (app *App) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "session")
	login, ok := session.Values["user"].(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := app.ordersRepo.GetOrders(login)
	if err != nil {
		http.Error(w, "Error getting orders", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (app *App) getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "session")
	login := session.Values["user"].(string)

	balance, err := app.ordersRepo.GetBalance(login)
	if err != nil {
		http.Error(w, "Error getting balance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(balance)
}

func (app *App) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "session")
	login := session.Values["user"].(string)

	var req struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := app.ordersRepo.Withdraw(login, req.Amount)
	if err == repository.ErrInsufficientFunds {
		http.Error(w, "Not enough money", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Error withdrawing money", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) getWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.sessionStore.Get(r, "session")
	login := session.Values["user"].(string)

	withdrawals, err := app.ordersRepo.GetWithdrawals(login)
	if err != nil {
		http.Error(w, "Error getting withdrawals", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(withdrawals)
}
