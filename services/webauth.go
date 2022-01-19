package services

import (
	"Dp218GO/models"
	"Dp218GO/repositories"
	"Dp218GO/utils"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// AuthService provides access to user database and sessionstore
// for user authentication in system
type AuthService struct {
	DB        repositories.UserRepo
	sessStore sessions.Store
}

const (
	sessionName = "login"
	sessionVal  = "user"
)

// NewAuthService returns new AuthService
func NewAuthService(db repositories.UserRepo, store sessions.Store) *AuthService {

	gob.Register(&models.User{})
	return &AuthService{
		DB:        db,
		sessStore: store,
	}
}

// AuthRequest contains required fields needed for authenticating user
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUp method registeres new user in system, returns error if it's failed
func (sv *AuthService) SignUp(user *models.User) error {
	pass, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = pass

	err = sv.DB.AddUser(user)

	if err != nil {
		return err
	}
	return nil
}

// SignIn method takes user from db, checks password, writes it to session and
// writes session id to cookie, returns error if it's failed
func (sv *AuthService) SignIn(w http.ResponseWriter, r *http.Request, authreq *AuthRequest) error {
	user, err := sv.DB.GetUserByEmail(authreq.Email)

	if err != nil {
		return err
	}

	if err := utils.CheckPassword(user.Password, authreq.Password); err != nil {
		return err
	}

	session, err := sv.getSessionStore().Get(r, sessionName)
	if err != nil {
		return err
	}
	user = *sanitize(&user)

	session.Values[sessionVal] = user
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

// SignOut deletes user from session, removes cookies, returns error if it's failed
func (sv *AuthService) SignOut(w http.ResponseWriter, r *http.Request) error {
	session, err := sv.getSessionStore().Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values[sessionVal] = nil
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetUserFromRequest retrieves user data from session
// returns error if is no user in session
func (sv *AuthService) GetUserFromRequest(r *http.Request) (*models.User, error) {
	sess, err := sv.sessStore.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	val, ok := sess.Values[sessionVal]
	if !ok {
		return nil, fmt.Errorf("%s", "no user in session")
	}

	var user = &models.User{}
	if user, ok = val.(*models.User); !ok {
		return nil, fmt.Errorf("%s", "no user in session")

	}

	return user, nil

}

func (sv *AuthService) getSessionStore() sessions.Store {
	return sv.sessStore
}

func sanitize(u *models.User) *models.User {
	u.Password = ""
	return u
}
