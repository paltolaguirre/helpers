package main

import (
	"fmt"
	"net/http"

	"github.com/xubiosueldos/conexionBD/Helper/structHelper"
	"github.com/xubiosueldos/monoliticComunication"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD"

	"github.com/xubiosueldos/framework"
)

type strHlprServlet struct {
	//	gorm.Model
	Username string `json:"username"`
	Tenant   string `json:"tenant"`
	Token    string `json:"token"`
	Options  string `json:"options"`
	Id       string `json:"id"`
}

type requestMono struct {
	Value interface{}
	Error error
}

/*
func (strhelper) TableName() string {
	return codigoHelper
}*/

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy."))
}

func getHelper(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		fmt.Println("La URL accedida: " + r.URL.String() + "/" + params["codigoHelper"])

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		var helper []structHelper.Helper

		//db.Raw("SELECT * FROM "+params["codigoHelper"]+" WHERE activo = 1 and deleted_at is null").Scan(&helper)

		if obtenerTablaPrivada(params["codigoHelper"]) == "MIXTA" {
			if err := db.Raw(crearQueryMixta(params["codigoHelper"], tokenAutenticacion.Tenant)).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "PURAPUBLICA" {
			if err := db.Raw(crearQueryPublica(params["codigoHelper"])).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "PURAPRIVADA" {
			if err := db.Raw(crearQueryPrivada(params["codigoHelper"], tokenAutenticacion.Tenant)).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "MONOLITICO" {
			if err := monoliticComunication.Gethelpers(w, r, tokenAutenticacion, params["codigoHelper"], "").Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

}

func getHelperId(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)

		helper_id := params["id"]
		helper_codigo := params["codigoHelper"]

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db)

		var helper []structHelper.Helper

		if err := db.Raw(" select * from (" + crearQueryMixta(helper_codigo, tokenAutenticacion.Tenant) + ") as tabla where tabla.id = " + helper_id).Scan(&helper).Error; err != nil {

			if err := monoliticComunication.Gethelpers(w, r, tokenAutenticacion, helper_codigo, helper_id).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			return
		}

		framework.RespondJSON(w, http.StatusOK, helper)
	}

}

//TODO MIGRAR TODO ESTO AL ARCHIVO DE CONFIGURACION
func obtenerTablaPrivada(concepto string) string {
	switch os := concepto; os {
	case "legajo":
		return "PURAPRIVADA"
	case "concepto":
		return "PURAPRIVADA"
	case "pais":
		return "PURAPUBLICA"
	case "provincia":
		return "PURAPUBLICA"
	case "localidad":
		return "PURAPUBLICA"
	case "cuenta":
		return "MONOLITICO"
	case "banco":
		return "MONOLITICO"
	case "liquidacion":
		return "PURAPRIVADA"
	case "hijo":
		return "PURAPRIVADA"
	case "conyuge":
		return "PURAPRIVADA"
	case "obrasocial":
		return "PURAPUBLICA"
	case "condicion":
		return "PURAPUBLICA"
	case "centrodecosto":
		return "MONOLITICO"
	case "condicionsiniestrado":
		return "PURAPUBLICA"
	case "conveniocolectivo":
		return "PURAPUBLICA"
	case "modalidadcontratacion":
		return "PURAPUBLICA"
	case "situacion":
		return "PURAPUBLICA"
	case "zona":
		return "PURAPUBLICA"

	case "liquidacioncondicionpago":
		return "PURAPUBLICA"

	case "liquidaciontipo":
		return "PURAPUBLICA"
	default:
		return "MIXTA"
	}
}

//id,nombre,codigo,descripcion"
func crearQueryMixta(codigo string, tenant string) string {
	return crearQueryPublica(codigo) + " union all " + crearQueryPrivada(codigo, tenant)
}

func crearQueryPublica(codigo string) string {
	return "select * from public." + codigo + " as tabla1 where tabla1.deleted_at is null and activo = 1"
}

func crearQueryPrivada(codigo string, tenant string) string {
	return "select * from " + tenant + "." + codigo + " as tabla2 where tabla2.deleted_at is null and activo = 1"
}

func getEmpresaId(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		fmt.Println("La URL accedida: " + r.URL.String())

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		dataempresa := monoliticComunication.Obtenerdatosempresa(w, r, tokenAutenticacion)
		framework.RespondJSON(w, http.StatusOK, dataempresa)
	}

}
