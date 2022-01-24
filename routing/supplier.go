package routing

import (
	"Dp218GO/models"
	"Dp218GO/services"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
)

var supplierService *services.SupplierService

var supplierKeyRoutes = []Route{
	{
		Uri:     `/models`,
		Method:  http.MethodGet,
		Handler: getModels,
	},
	{
		Uri:     `/models`,
		Method:  http.MethodPost,
		Handler: createModel,
	},
	{
		Uri:     `/price/{id}`,
		Method:  http.MethodPost,
		Handler: editPrice,
	},
	{
		Uri:     `/upload/{id}`,
		Method:  http.MethodPost,
		Handler: uploadFile,
	},
	{
		Uri:     `/model/{id}`,
		Method:  http.MethodPost,
		Handler: addSuppliersScooter,
	},
	{
		Uri:     `/delete/{id}`,
		Method:  http.MethodPost,
		Handler: deleteSuppliersScooter,
	},
}

// FileHeader - the multipart.FileHeader uses the following struct
type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
}

// AddSupplierHandler - add endpoints for working with supplier scooters and models to http router
func AddSupplierHandler(router *mux.Router, service *services.SupplierService) {
	supplierService = service
	supplierRouter := router.NewRoute().Subrouter()
	supplierRouter.Use(FilterAuth(authenticationService))

	for _, rt := range supplierKeyRoutes {
		supplierRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		supplierRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getModels(w http.ResponseWriter, r *http.Request) {
	var modelList = &models.ScooterModelDTOList{}
	var err error
	format := GetFormatFromRequest(r)

	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	modelList, err = supplierService.GetModels()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, modelList, HTMLPath+"supplier.html")
}

func createModel(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	model := &models.ScooterModelDTO{}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	modelName := r.FormValue("modelName")
	maxWeight := r.FormValue("maxWeight")
	speed := r.FormValue("speed")
	price := r.FormValue("price")

	intMaxWeight, err := strconv.Atoi(maxWeight)
	if err != nil {
		log.Println(err)
	}

	intSpeed, err := strconv.Atoi(speed)
	if err != nil {
		log.Println(err)
	}
	intPrice, err := strconv.Atoi(price)
	if err != nil {
		log.Println(err)
	}
	model.ModelName = modelName
	model.MaxWeight = intMaxWeight
	model.Speed = intSpeed
	model.Price = intPrice

	if err := supplierService.AddModel(model); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	http.Redirect(w, r, "http://localhost:8080/models", http.StatusFound)
}

func editPrice(w http.ResponseWriter, r *http.Request) {
	model := &models.ScooterModelDTO{}
	format := GetFormatFromRequest(r)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	editedPrice := r.FormValue("priceInput")
	intPrice, err := strconv.Atoi(editedPrice)
	if err != nil {
		log.Println(err)
	}

	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	modelData, err := supplierService.SelectModel(modelId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	model = modelData
	model.Price = intPrice

	if err := supplierService.ChangePrice(model); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	http.Redirect(w, r, "http://localhost:8080/models", http.StatusFound)
}

func addSuppliersScooter(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)
	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	scooterSerial := r.FormValue("newScooter")

	if err := supplierService.AddSuppliersScooter(modelId, scooterSerial); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	http.Redirect(w, r, "http://localhost:8080/models", http.StatusFound)
}

func deleteSuppliersScooter(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	scooterId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	err = supplierService.DeleteSuppliersScooter(scooterId)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	http.Redirect(w, r, "http://localhost:8080/models", http.StatusFound)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	filepath := "./internal/" + handler.Filename
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	modelId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	supplierService.InsertScootersToDb(modelId, filepath)

	http.Redirect(w, r, "http://localhost:8080/models", http.StatusFound)
}
