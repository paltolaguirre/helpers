package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/framework/configuracion"
)

type strhelper struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
}

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

type strempresa struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	Domicilio	string `json:"domicilio"`
	Localidad	string `json:"localidad"`
	Cuit		string `json:"cuit"`
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

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, "helper", 0, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		var helper []strhelper

		//db.Raw("SELECT * FROM "+params["codigoHelper"]+" WHERE activo = 1 and deleted_at is null").Scan(&helper)

		var requestMono requestMono

		if obtenerTablaPrivada(params["codigoHelper"]) == "MIXTA" {
			if err := db.Raw(crearQueryMixta(params["codigoHelper"], tokenAutenticacion.Tenant, r)).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "PURAPUBLICA" {
			if err := db.Raw(crearQueryPublica(params["codigoHelper"], r)).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "PURAPRIVADA" {
			if err := db.Raw(crearQueryPrivada(params["codigoHelper"], tokenAutenticacion.Tenant, r)).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "MONOLITICO" {
			if err := requestMono.requestMonolitico(w, r, tokenAutenticacion, params["codigoHelper"], "").Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
			}
		}
	}

}

func getHelperId(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)

		helper_id := params["id"]

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, "helper", 0, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		var helper strhelper

		var requestMono requestMono

		if err := db.Raw(" select * from (" + crearQueryMixta(params["codigoHelper"], tokenAutenticacion.Tenant, r) + ") as tabla where tabla.id = " + helper_id).Scan(&helper).Error; err != nil {
			if err := requestMono.requestMonolitico(w, r, tokenAutenticacion, params["codigoHelper"], helper_id).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			return
		}

		framework.RespondJSON(w, http.StatusOK, helper)
	}

}

func obtenerFiltroTabla(codigo string, r *http.Request) string {
	var queryFilter string = ""
	switch os := codigo; os {
	case "provincia":
		{
			if r.URL.Query()["paisid"] != nil {
				queryFilter += " and paisid =" + r.URL.Query()["paisid"][0]
			}
		}
	case "localidad":
		{
			if r.URL.Query()["provinciaid"] != nil {
				queryFilter += " and provinciaid =" + r.URL.Query()["provinciaid"][0]
			}
		}
	default:
		return ""
	}
	return queryFilter
}

//TODO MIGRAR TODO ESTO AL ARCHIVO DE CONFIGURACION
func obtenerTablaPrivada(concepto string) string {
	switch os := concepto; os {
	case "legajo":
		return "PURAPRIVADA"
	case "concepto":
		return "MIXTA"
	case "pais":
		return "PURAPUBLICA"
	case "provincia":
		return "PURAPUBLICA"
	case "localidad":
		return "PURAPUBLICA"
	case "cuenta":
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
		return "PURAPUBLICA"
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
	default:
		return "MIXTA"
	}
}

//id,nombre,codigo,descripcion"
func crearQueryMixta(codigo string, tenant string, r *http.Request) string {
	return crearQueryPublica(codigo, r) + " union all " + crearQueryPrivada(codigo, tenant, r)
}

func crearQueryPublica(codigo string, r *http.Request) string {
	return "select * from public." + codigo + " as tabla1 where tabla1.deleted_at is null and activo = 1 " + obtenerFiltroTabla(codigo, r)
}

func crearQueryPrivada(codigo string, tenant string, r *http.Request) string {
	return "select * from " + tenant + "." + codigo + " as tabla2 where tabla2.deleted_at is null and activo = 1 " + obtenerFiltroTabla(codigo, r)
}

func (s *requestMono) requestMonolitico(w http.ResponseWriter, r *http.Request, tokenAutenticacion *publico.Security, codigo string, id string) *requestMono {

	var strHlprSrv strHlprServlet
	token := *tokenAutenticacion

	strHlprSrv.Options = "HLP"
	strHlprSrv.Tenant = token.Tenant
	strHlprSrv.Token = token.Token
	strHlprSrv.Username = token.Username
	strHlprSrv.Id = id

	pagesJson, err := json.Marshal(strHlprSrv)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := configuracion.GetUrlMonolitico() + codigo + "GoServlet"

	fmt.Println("URL:>", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pagesJson))

	if err != nil {
		fmt.Println("Error: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	str := string(body)
	fmt.Println("BYTES RECIBIDOS :", len(str))

	fixUtf := func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}

	var dataStruct []strhelper
	json.Unmarshal([]byte(strings.Map(fixUtf, str)), &dataStruct)

	fmt.Println("BYTES RECIBIDOS :", string(body))

	//Para que el json que devuelva quede acorde al que devuelve go
	if len(dataStruct) == 1 {
		framework.RespondJSON(w, http.StatusOK, dataStruct[0])
	} else {
		framework.RespondJSON(w, http.StatusOK, dataStruct)
	}
	return s
}

func AutomigrateTablasPrivadas(db *gorm.DB) {

}

func getEmpresaId(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)

		helper_id := params["id"]

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, "helper", 0, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		var empresa strempresa

		/*var requestMono requestMono

		if err := db.Raw(" select * from (" + crearQueryMixta(params["codigoHelper"], tokenAutenticacion.Tenant) + ") as tabla where tabla.id = " + helper_id).Scan(&helper).Error; err != nil {
			if err := requestMono.requestMonolitico(w, r, tokenAutenticacion, params["codigoHelper"], helper_id).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			return
		}*/
		id, err := strconv.Atoi(helper_id)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		empresa.ID = id
		empresa.Nombre = "Mi Empresa"
		empresa.Codigo = "TNT_914"
		empresa.Descripcion = "Empresa Online confiable de venta de garantias"
		empresa.Domicilio = "Av. Siempre Viva 1234"
		empresa.Localidad = "C.A.B.A"
		empresa.Cuit = "12-12123123-1"

		framework.RespondJSON(w, http.StatusOK, empresa)
	}

}