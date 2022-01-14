package routing

import (
	"Dp-218_GO_micro/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var keyOrderRoutes = []Route{
	{
		Uri:     `/orders`,
		Method:  http.MethodGet,
		Handler: getAllOrders,
	},
}

//AddOrderHandler adds routes to the router from the list of routes.
func AddOrderHandler(router *mux.Router, order *services.OrderService) {
	orderService = order
	for _, rt := range keyOrderRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := orderService.GetAllOrders()
	if err != nil {
		ServerErrorRender(FormatJSON, w)
		fmt.Println(err)
		return
	}

	EncodeAnswer(FormatJSON, w, orders)
}
