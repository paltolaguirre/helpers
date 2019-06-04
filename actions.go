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

/*
func (strhelper) TableName() string {
	return codigoHelper
}*/

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func getHelper(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)

		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, "helper", 0, AutomigrateTablasPrivadas)
		defer db.Close()

		var helper []strhelper

		//db.Raw("SELECT * FROM "+params["codigoHelper"]+" WHERE activo = 1 and deleted_at is null").Scan(&helper)

		var requestMono requestMono

		if obtenerTablaPrivada(params["codigoHelper"]) == "MIXTA" {
			if err := db.Raw(crearQueryMixta(params["codigoHelper"], tokenAutenticacion.Tenant)).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "PURAPUBLICA" {
			if err := db.Raw(crearQueryPublica(params["codigoHelper"])).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "PURAPRIVADA" {
			if err := db.Raw(crearQueryPrivada(params["codigoHelper"], tokenAutenticacion.Tenant)).Scan(&helper).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			} else {
				framework.RespondJSON(w, http.StatusOK, helper)
			}
		}

		if obtenerTablaPrivada(params["codigoHelper"]) == "MONOLITICO" {
			if err := requestMono.requestMonolitico(w, r, tokenAutenticacion, params["codigoHelper"], "").Error; err != nil {
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
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, "helper", 0, AutomigrateTablasPrivadas)
		defer db.Close()

		var helper strhelper

		var requestMono requestMono

		if err := db.Raw(" select * from (" + crearQueryMixta(params["codigoHelper"], tokenAutenticacion.Tenant) + ") as tabla where tabla.id = " + helper_id).Scan(&helper).Error; err != nil {
			if err := requestMono.requestMonolitico(w, r, tokenAutenticacion, params["codigoHelper"], helper_id).Error; err != nil {
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
		return "MIXTA"
	case "pais":
		return "PURAPUBLICA"
	case "provincia":
		return "PURAPUBLICA"
	case "cuenta":
		return "MONOLITICO"
	default:
		return "MIXTA"
	}
}

//id,nombre,codigo,descripcion"
func crearQueryMixta(concepto string, tenant string) string {
	return crearQueryPublica(concepto) + " union all " + crearQueryPrivada(concepto, tenant)
}

func crearQueryPublica(concepto string) string {
	return "select * from public." + concepto + " as tabla1 where tabla1.deleted_at is null and activo = 1"
}

func crearQueryPrivada(concepto string, tenant string) string {
	return "select * from " + tenant + "." + concepto + " as tabla2 where tabla2.deleted_at is null and activo = 1"
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
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

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
