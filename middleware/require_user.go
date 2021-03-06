package middleware

import (
	"net/http"
	"strings"

	"github.com/imattf/go-courses/gallery/context"
	"github.com/imattf/go-courses/gallery/models"
)

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// If the user is requesting a static assets or image
		// we will not need to lookup the current user, so we skip
		// doing that.
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets") ||
			strings.HasPrefix(path, "/images") {
			next(w, r)
			return
		}

		//if the user is logged in...
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}

		//user is found...
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

// RequireUser assumes the User middleware has already been run
// otherwise it will not work correctly.
type RequireUser struct {
	User
}

// Apply assumes the User middleware has already been run
// otherwise it will not work correctly.
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// Stuff in ApplyFn (which is a function as a parameter), gets ran when we call
// the ApplyFn function itself with a call to one of the next() calls

// AppleFn assumes the User middleware has already been run
// otherwise it will not work correctly.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
