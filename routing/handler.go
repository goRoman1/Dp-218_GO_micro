package routing

import (
	"Dp-218_GO_micro/configs"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Uri     string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

// available formats of http response representation
const (
	FormatJSON = iota
	FormatHTML
)

var (
	// HTMLPath - path where html templates are stored
	HTMLPath = configs.TEMPLATES_PATH + "html/"
	// MainPageHTML - path to main page template
	MainPageHTML = HTMLPath + "main-page.html"
	// ErrorPageHTML - path to error page template
	ErrorPageHTML = HTMLPath + "error.html"
	//APIprefix - endpoint prefix to receive result in API (json) format
	APIprefix = "/api/v1"
)

// NewRouter - initialize main http router of the application
func NewRouter() *mux.Router {

	router := mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/",
		http.FileServer(http.Dir(configs.TEMPLATES_PATH))))

	router.HandleFunc("/", showHomePage)
	router.HandleFunc("/login", showLoginPage)
	return router
}

func showHomePage(w http.ResponseWriter, r *http.Request) {
	EncodeAnswer(FormatHTML, w, nil, MainPageHTML)
}

func showLoginPage(w http.ResponseWriter, r *http.Request) {
	// not needed if filter is applied
	user := GetUserFromContext(r)
	if user != nil {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	EncodeAnswer(FormatHTML, w, nil, HTMLPath+"login-registration.html")
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	EncodeError(FormatHTML, w, ErrorRenderer(fmt.Errorf("method not allowed"), "Not allowed", http.StatusMethodNotAllowed))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	EncodeError(FormatHTML, w, ErrorRenderer(fmt.Errorf("resource not found"), "Not found", http.StatusNotFound))
}

// ServerErrorRender - renders error page with server error
func ServerErrorRender(format int, w http.ResponseWriter) {
	EncodeError(format, w, ErrorRenderer(fmt.Errorf("server error"), "Internal server error", http.StatusInternalServerError))
}

// EncodeError - renders general error page with passed error info
func EncodeError(format int, w http.ResponseWriter, respErr *ResponseStatus) {
	var err error
	switch format {
	case FormatJSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(respErr)
	case FormatHTML:
		w.Header().Set("Content-Type", "text/html")
		var tmpl *template.Template
		if tmpl, err = template.ParseFiles(ErrorPageHTML); err == nil {
			err = tmpl.Execute(w, respErr)
		}
	default:
		err = fmt.Errorf("format error")
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(respErr.StatusCode)
}

// EncodeAnswer - renders given answer structure into given htmlTemplate using given format
func EncodeAnswer(format int, w http.ResponseWriter, answer interface{}, htmlTemplates ...string) {
	var err error

	switch format {

	case FormatJSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(answer)

	case FormatHTML:
		if len(htmlTemplates) == 0 {
			EncodeError(format, w, &ResponseStatus{StatusText: "Encode error", Message: "no html temlates", StatusCode: http.StatusInternalServerError})
			return
		}
		w.Header().Set("Content-Type", "text/html")
		var tmpl *template.Template
		if tmpl, err = template.ParseFiles(htmlTemplates[0]); err == nil {
			err = tmpl.Execute(w, answer)
		}

	default:
		err = fmt.Errorf("format error")
	}

	if err != nil {
		EncodeError(format, w, &ResponseStatus{Err: err, StatusText: "Encode error", Message: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}

}

// DecodeRequest - receives information from given request and decodes it into given format
func DecodeRequest(format int, w http.ResponseWriter, r *http.Request, requestData interface{}, htmlDecoder func(r *http.Request, dataToDecode interface{}) error) {
	var err error

	switch format {

	case FormatJSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewDecoder(r.Body).Decode(requestData)

	case FormatHTML:
		w.Header().Set("Content-Type", "text/html")
		err = htmlDecoder(r, requestData)

	default:
		err = fmt.Errorf("format error")
	}

	if err != nil {
		EncodeError(FormatJSON, w, ErrorRenderer(err, "Bad request", http.StatusBadRequest))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetFormatFromRequest - get format html/json based on request URI
func GetFormatFromRequest(r *http.Request) int {
	if strings.Contains(r.RequestURI, APIprefix) {
		return FormatJSON
	}
	return FormatHTML
}

// GetParameterFromRequest - get parameter from the request by name in needed type with given convertToType function
func GetParameterFromRequest(r *http.Request, paramName string, convertToType func(strData string) (interface{}, error)) (interface{}, error) {
	if r.Method == http.MethodGet && r.Form == nil || r.Method != http.MethodGet && r.PostForm == nil {
		r.ParseForm()
	}

	result := ""
	switch r.Method {
	case http.MethodGet:
		if _, ok := r.Form[paramName]; !ok {
			return "", fmt.Errorf("no such field <%s> in request", paramName)
		}
		result = r.FormValue(paramName)
	default:
		if _, ok := r.PostForm[paramName]; !ok {
			return "", fmt.Errorf("no such field <%s> in request", paramName)
		}
		result = r.PostFormValue(paramName)
	}

	return convertToType(result)
}
