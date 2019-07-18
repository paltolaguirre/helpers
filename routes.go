package main

import "github.com/gorilla/mux"
import "net/http"

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)

	}

	return router
}

var routes = Routes{
	Route{
		"Healthy",
		"GET",
		"/api/helper/healthy",
		Healthy,
	},
	Route{
		"getHelper",
		"GET",
		"/api/helper/helpers/{codigoHelper}",
		getHelper,
	},
	Route{
		"getHelperId",
		"GET",
		"/api/helper/helpers/{codigoHelper}/{id}",
		getHelperId,
	},
	Route{
		"HealthyEmpresa",
		"GET",
		"/api/empresa/healthy",
		Healthy,
	},
	Route{
		"getEmpresaId",
		"GET",
		"/api/empresa/empresas/{id}",
		getEmpresaId,
	},
}
