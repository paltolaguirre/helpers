package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/framework"
)

type strhelper struct {
	//	gorm.Model
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	//	Activo      int    `json:"activo"`
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

		db.Table(params["codigoHelper"]).Where("activo = 1 and deleted_at is null").Select("id,nombre,codigo,descripcion").Scan(&helper)

		framework.RespondJSON(w, http.StatusOK, helper)
	}

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
