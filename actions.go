package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/helpers/structHelpers"
)

func respondJSON(w http.ResponseWriter, status int, results interface{}) {

	response, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))

}

func PaisesHLP(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {
		errorToken(w, tokenError)
		return
	} else {

		db := obtenerDB(tokenAutenticacion)
		defer db.Close()

		var helper []structHelpers.Pais

		db.Find(&helper)

		respondJSON(w, http.StatusOK, helper)
	}

}

func obtenerDB(tokenAutenticacion *publico.TokenAutenticacion) *gorm.DB {

	token := *tokenAutenticacion
	tenant := token.Tenant

	return conexionBD.ConnectBD(tenant)

}

func errorToken(w http.ResponseWriter, tokenError *publico.Error) {
	errorToken := *tokenError
	respondError(w, errorToken.ErrorCodigo, errorToken.ErrorNombre)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func checkTokenValido(r *http.Request) (*publico.TokenAutenticacion, *publico.Error) {

	var tokenAutenticacion *publico.TokenAutenticacion
	var tokenError *publico.Error

	url := "http://localhost:8081/check-token"

	req, _ := http.NewRequest("GET", url, nil)

	token := r.Header.Get("Token")

	req.Header.Add("token", token)

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
