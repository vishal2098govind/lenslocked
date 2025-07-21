package middlewares

import (
	"fmt"
	"net/http"

	"github.com/vishal2098govind/lenslocked/context"
	"github.com/vishal2098govind/lenslocked/cookies"
	sessionM "github.com/vishal2098govind/lenslocked/models/session"
)

type UsersMiddleware struct {
	SessionService *sessionM.SessionService
}

func (umw *UsersMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionC, err := cookies.ReadCookie(r, cookies.CookieSession)
		if err != nil || sessionC == "" {
			next.ServeHTTP(w, r)
			return
		}

		res, err := umw.SessionService.User(sessionM.GetUserIdRequest{
			Token: sessionC,
		})
		if err != nil || res.User == nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, res.User)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// RequireUser assumes that UserMiddleware.SetUser is already being called
// before reaching here
func (umw *UsersMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
