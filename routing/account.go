package routing

import (
	"Dp218GO/internal/validation"
	"Dp218GO/models"
	"Dp218GO/services"
	"Dp218GO/utils"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var accountService *services.AccountService
var accountIDKey = "accID"

var keyAccountRoutes = []Route{
	{
		Uri:     `/accounts`,
		Method:  http.MethodGet,
		Handler: getAllAccounts,
	},
	{
		Uri:     `/account/{` + accountIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getAccountInfo,
	},
	{
		Uri:     `/account/{` + accountIDKey + `}`,
		Method:  http.MethodPost,
		Handler: updateAccountInfo,
	},
	{
		Uri:     `/account`,
		Method:  http.MethodGet,
		Handler: createAccountPage,
	},
	{
		Uri:     `/account`,
		Method:  http.MethodPost,
		Handler: createAccount,
	},
}

// AddAccountHandler - add endpoints for money accounts to http router
func AddAccountHandler(router *mux.Router, service *services.AccountService) {
	accountService = service
	accountRouter := router.NewRoute().Subrouter()
	accountRouter.Use(FilterAuth(authenticationService))

	for _, rt := range keyAccountRoutes {
		accountRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		accountRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllAccounts(w http.ResponseWriter, r *http.Request) {

	var accounts *models.AccountList
	var err error
	format := GetFormatFromRequest(r)

	user := GetUserFromContext(r)
	if user == nil {
		EncodeError(format, w, ErrorRendererDefault(errors.New("not authorized")))
		return
	}

	accounts, err = accountService.GetAccountsByOwner(*user)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, accounts, HTMLPath+"accounts.html")
}

func getAccountInfo(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	accID, err := strconv.Atoi(mux.Vars(r)[accountIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	accData, err := accountService.GetAccountOutputStructByID(accID)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, accData, HTMLPath+"account.html")
}

func updateAccountInfo(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	accID, err := strconv.Atoi(mux.Vars(r)[accountIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	account, err := accountService.GetAccountByID(accID)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	actionType, err := GetParameterFromRequest(r, "ActionType", utils.ConvertStringToString())
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	switch actionType {
	case "AddMoneyToAccount":
		moneyAmount, err := GetParameterFromRequest(r, "MoneyAmount", utils.ConvertStringToFloat())
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
		err = accountService.AddMoneyToAccount(account, int(moneyAmount.(float64)*100))
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
	case "TakeMoneyFromAccount":
		moneyAmount, err := GetParameterFromRequest(r, "MoneyAmount", utils.ConvertStringToFloat())
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
		err = accountService.TakeMoneyFromAccount(account, int(moneyAmount.(float64)*100))
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
	default:
		return
	}

	getAccountInfo(w, r)
}

func createAccountPage(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	user := GetUserFromContext(r)

	tmpl, err := template.ParseFiles("templates/html/account-add.html")
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	tmpl.Execute(w, user)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	user := GetUserFromContext(r)

	accReq := validation.CreateAccountRequest{
		Name:   r.FormValue("name"),
		Number: r.FormValue("number"),
	}
	if err := accReq.Validate(); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return

	}

	account := models.Account{
		Name:   accReq.Name,
		Number: accReq.Number,
		User:   *user,
	}

	if err := accountService.AddAccount(&account); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	http.Redirect(w, r, "/accounts", http.StatusFound)
}
