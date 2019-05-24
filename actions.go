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
)

type strhelper struct {
	ID          string `json:"id"`
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
}

type requestMono struct {
	Value interface{}
	Error error
}

/*
func (strhelper) TableName() string {
	return codigoHelper
}*/

func getHelper(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)

		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, "helper", 0, AutomigrateTablasPrivadas)
		defer db.Close()

		var helper []strhelper

		//db.Raw("SELECT * FROM "+params["codigoHelper"]+" WHERE activo = 1 and deleted_at is null").Scan(&helper)

		var requestMono requestMono

		if err := db.Table(params["codigoHelper"]).Where("activo = 1 and deleted_at is null").Select("id,nombre,codigo,descripcion").Scan(&helper).Error; err != nil {
			if err := requestMono.requestMonolitico(w, r, tokenAutenticacion, params["codigoHelper"]).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			return
		}

		framework.RespondJSON(w, http.StatusOK, helper)
	}

}

func (s *requestMono) requestMonolitico(w http.ResponseWriter, r *http.Request, tokenAutenticacion *publico.Security, codigo string) *requestMono {

	var strHlprSrv strHlprServlet
	token := *tokenAutenticacion

	strHlprSrv.Options = "HLP"
	strHlprSrv.Tenant = token.Tenant
	strHlprSrv.Token = token.Token
	strHlprSrv.Username = token.Username

	pagesJson, err := json.Marshal(strHlprSrv)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := "https://localhost:8443/NXV/" + codigo + "GoServlet"

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

	framework.RespondJSON(w, http.StatusOK, dataStruct)
	return s
}

func AutomigrateTablasPrivadas(db *gorm.DB) {

}
