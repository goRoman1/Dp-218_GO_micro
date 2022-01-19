package routing

import (
	"Dp218GO/internal/validation"
	"Dp218GO/models"
	"Dp218GO/services"
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type userKey string

var (
	ukey                  userKey = "user"
	authenticationService *services.AuthService
	// ErrSignUp error returned to client if registering failed
	ErrSignUp = errors.New("signup error")
	// ErrSignIn error returned to client if authentication failed
	ErrSignIn = errors.New("signin error")
)

//AddAuthHandler registeres endpoints for authentication
func AddAuthHandler(router *mux.Router, service *services.AuthService) {
	authenticationService = service
	router.Path("/signup").HandlerFunc(SignUp(authenticationService)).Methods(http.MethodPost)
	router.Path("/signin").HandlerFunc(SignIn(authenticationService)).Methods(http.MethodPost)
	router.Path("/signout").HandlerFunc(SignOut(authenticationService)).Methods(http.MethodGet)
}

//SignUp is handler for signup authentication service method
func SignUp(sv *services.AuthService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		valReq := validation.SignUpUserRequest{
			LoginEmail:  r.FormValue("email"),
			UserName:    r.FormValue("name"),
			UserSurname: r.FormValue("surname"),
			Password:    r.FormValue("password"),
		}
		if err := valReq.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := &models.User{
			LoginEmail:  valReq.LoginEmail,
			IsBlocked:   true,
			UserName:    valReq.UserName,
			UserSurname: valReq.UserSurname,
			Role:        models.Role{ID: 2},
			Password:    valReq.Password,
		}

		if err := sv.SignUp(user); err != nil {

			http.Error(w, ErrSignUp.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

//SignIn is handler for signin authentication service method
func SignIn(sv *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		valReq := validation.SignInUserRequest{
			LoginEmail: r.FormValue("email"),
			Password:   r.FormValue("password"),
		}
		if err := valReq.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &services.AuthRequest{
			Email:    valReq.LoginEmail,
			Password: valReq.Password,
		}

		if err := sv.SignIn(w, r, req); err != nil {

			EncodeError(FormatHTML, w, &ResponseStatus{
				Err:        ErrSignIn,
				StatusCode: http.StatusForbidden,
				StatusText: ErrSignIn.Error(),
				Message:    "cant get user" + err.Error(),
			})
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

//SignOut is handler for signout authentication service method
func SignOut(sv *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := sv.SignOut(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// FilterAuth  is middleware checks if user is authenticated
// writes user to context for retrieving if chaining middleware is present
// shows error if user is not authenticated
func FilterAuth(sv *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := sv.GetUserFromRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			newReq := r.WithContext(context.WithValue(r.Context(), ukey, user))

			next.ServeHTTP(w, newReq)
		})
	}
}

// GetUserFromContext retrieves user from context
func GetUserFromContext(r *http.Request) *models.User {
	val := r.Context().Value(ukey)
	user, ok := val.(*models.User)

	if ok {
		return user
	}
	return nil
}
