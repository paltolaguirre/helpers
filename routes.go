package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

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
		"getHelperConcepto",
		"GET",
		"/api/helper/helpers/concepto",
		getHelperConcepto,
	},
	Route{
		"getHelperFunction",
		"GET",
		"/api/helper/helpers/function",
		getHelperFunction,
	},
	Route{
		"getHelperTipoimpuestoganancias",
		"GET",
		"/api/helper/helpers/tipoimpuestosganancias",
		getHelperTipoimpuestoganancias,
	},
	Route{
		"getHelperConceptoafip",
		"GET",
		"/api/helper/helpers/conceptoafip",
		getHelperConceptoafip,
	},
	Route{
		"getHelperLegajo",
		"GET",
		"/api/helper/helpers/legajo",
		getHelperLegajo,
	},
	Route{
		"getHelperObrasocial",
		"GET",
		"/api/helper/helpers/obrasocial",
		getHelperObrasocial,
	},
	Route{
		"getHelperZona",
		"GET",
		"/api/helper/helpers/zona",
		getHelperZona,
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
		"/api/empresa/empresas",
		getEmpresaId,
	},
	Route{
		"getImporteEnLetras",
		"GET",
		"/api/helper/importeenletras/{numero}",
		getImporteEnLetras,
	},
}
