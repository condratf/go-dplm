package router

import "net/http"

func (appRouter *AppRouter) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := appRouter.sessionStore.Get(r, "session")
		if session.Values["user"] == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
