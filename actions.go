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
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
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

/*
func (strhelper) TableName() string {
	return codigoHelper
}*/

func getHelper(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {
		errorToken(w, tokenError)
		return
	} else {

		params := mux.Vars(r)

		db := obtenerDB(tokenAutenticacion)
		defer db.Close()

		var helper []strhelper

		//db.Raw("SELECT * FROM "+params["codigoHelper"]+" WHERE activo = 1 and deleted_at is null").Scan(&helper)

		if err := db.Table(params["codigoHelper"]).Where("activo = 1 and deleted_at is null").Select("id,nombre,codigo,descripcion").Scan(&helper).Error; err != nil {
			requestMonolitico(w, r, tokenAutenticacion, params["codigoHelper"])
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, helper)
	}

}

func requestMonolitico(w http.ResponseWriter, r *http.Request, tokenAutenticacion *publico.TokenAutenticacion, codigo string) {

	var strHlprSrv strHlprServlet
	token := *tokenAutenticacion

	strHlprSrv.Options = "HLP"
	strHlprSrv.Tenant = token.Tenant
	strHlprSrv.Token = token.Token
	strHlprSrv.Username = token.Username

	pagesJson, err := json.Marshal(token)
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

	s := string(body)
	fmt.Println("BYTES RECIBIDOS :", len(s))

	fixUtf := func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}

	var dataStruct []strhelper
	json.Unmarshal([]byte(strings.Map(fixUtf, s)), &dataStruct)

	fmt.Println("BYTES RECIBIDOS :", string(body))

	framework.RespondJSON(w, http.StatusOK, dataStruct)
}

func obtenerDB(tokenAutenticacion *publico.TokenAutenticacion) *gorm.DB {
	token := *tokenAutenticacion
	tenant := token.Tenant

	return conexionBD.ConnectBD(tenant)
}

func errorToken(w http.ResponseWriter, tokenError *publico.Error) {
	errorToken := *tokenError
	framework.RespondError(w, errorToken.ErrorCodigo, errorToken.ErrorNombre)
}

func checkTokenValido(r *http.Request) (*publico.TokenAutenticacion, *publico.Error) {

	var tokenAutenticacion *publico.TokenAutenticacion
	var tokenError *publico.Error

	url := "http://localhost:8081/check-token"

	req, _ := http.NewRequest("GET", url, nil)

	header := r.Header.Get("Authorization")

	req.Header.Add("Authorization", header)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 400 {

		// tokenAutenticacion = &(TokenAutenticacion{})
		tokenAutenticacion = new(publico.TokenAutenticacion)
		json.Unmarshal([]byte(string(body)), tokenAutenticacion)

	} else {
		tokenError = new(publico.Error)
		json.Unmarshal([]byte(string(body)), tokenError)

	}

	return tokenAutenticacion, tokenError
}
