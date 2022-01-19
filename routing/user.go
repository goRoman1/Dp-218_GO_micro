package routing

import (
	"Dp218GO/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"Dp218GO/services"

	"github.com/gorilla/mux"
)

var userService *services.UserService
var userIDKey = "userID"

var keyUserRoutes = []Route{
	{
		Uri:     `/home`,
		Method:  http.MethodGet,
		Handler: getUserPage,
	},
	{
		Uri:     `/users`,
		Method:  http.MethodGet,
		Handler: getAllUsers,
	},
	{
		Uri:     `/users`,
		Method:  http.MethodPost,
		Handler: allUsersOperation,
	},
	{
		Uri:     `/user/{` + userIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getUser,
	},
	{
		Uri:     `/user`,
		Method:  http.MethodPost,
		Handler: createUser,
	},
	{
		Uri:     `/user/{` + userIDKey + `}`,
		Method:  http.MethodPost,
		Handler: updateUser,
	},
	{
		Uri:     `/user/{` + userIDKey + `}`,
		Method:  http.MethodDelete,
		Handler: deleteUser,
	},
}

type userWithRoleList struct {
	models.User
}

// ListOfRoles - returns slice of all roles to render them on template
func (ur *userWithRoleList) ListOfRoles() []models.Role {
	roles, _ := userService.GetAllRoles()
	return roles.Roles
}

// AddUserHandler - add endpoints for working with users to http router
func AddUserHandler(router *mux.Router, service *services.UserService) {
	userService = service
	userRouter := router.NewRoute().Subrouter()
	userRouter.Use(FilterAuth(authenticationService))

	for _, rt := range keyUserRoutes {
		userRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		userRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	user := &models.User{}
	DecodeRequest(FormatJSON, w, r, user, nil)

	if err := userService.AddUser(user); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, user)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {

	var users = &models.UserList{}
	var err error
	format := GetFormatFromRequest(r)

	r.ParseForm()
	searchData := r.FormValue("SearchData")

	if len(searchData) == 0 {
		users, err = userService.GetAllUsers()
	} else {
		users, err = userService.FindUsersByLoginNameSurname(searchData)
	}
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, users, HTMLPath+"user-list.html")
}

func getUserPage(w http.ResponseWriter, r *http.Request) {
	// check can be omited if filter is applied to route
	user := GetUserFromContext(r)
	if user == nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(errors.New("not authorized")))
		return
	}

	EncodeAnswer(FormatHTML, w, user, HTMLPath+"home.html")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userID, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	user, err := userService.GetUserByID(userID)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, &userWithRoleList{user}, HTMLPath+"user-edit.html")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userID, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	err = userService.DeleteUser(userID)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}
	EncodeAnswer(format, w, ErrorRenderer(fmt.Errorf(""), "success", http.StatusOK))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userID, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	userData, err := userService.GetUserByID(userID)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	DecodeRequest(format, w, r, &userData, decodeUserUpdateRequest)
	userData, err = userService.UpdateUser(userID, userData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, &userWithRoleList{userData}, HTMLPath+"user-edit.html")
}

func allUsersOperation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	r.ParseForm()
	if _, ok := r.Form["ActionType"]; !ok {

		return
	}
	actionType := r.FormValue("ActionType")
	switch actionType {
	case "BlockUser":
		userID, err := strconv.Atoi(r.FormValue("UserID"))
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
		err = userService.ChangeUsersBlockStatus(userID)
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
	default:
		EncodeError(format, w, ErrorRendererDefault(fmt.Errorf("unknown users operation")))
	}
	getAllUsers(w, r)
}

func decodeUserUpdateRequest(r *http.Request, data interface{}) error {

	var err error
	r.ParseForm()
	userData := data.(*models.User)

	if _, ok := r.Form["LoginEmail"]; ok {
		userData.LoginEmail = r.FormValue("LoginEmail")

	}
	if _, ok := r.Form["UserName"]; ok {
		userData.UserName = r.FormValue("UserName")
	}
	if _, ok := r.Form["UserSurname"]; ok {
		userData.UserSurname = r.FormValue("UserSurname")
	}
	if _, ok := r.Form["RoleID"]; ok {

		var roleID int
		roleID, err = strconv.Atoi(r.FormValue("RoleID"))
		if err != nil {
			return err
		}
		userData.Role, err = userService.GetRoleByID(roleID)
		if err != nil {
			return err
		}
	}
	if _, ok := r.Form["IsBlocked"]; ok {
		userData.IsBlocked, _ = strconv.ParseBool(r.FormValue("IsBlocked"))
	}

	return nil
}
