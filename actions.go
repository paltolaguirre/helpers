package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
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
	case "estadocivil":
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
	case "tipoconcepto":
		return "PURAPUBLICA"

	case "tipodecalculo":
		return "PURAPUBLICA"

	case "siradigtipoimpuesto":
		return "PURAPUBLICA"

	case "siradigtipooperacion":
		return "PURAPUBLICA"

	case "function":
		return "PURAPRIVADA"

	case "tipocalculoautomatico":
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

func getImporteEnLetras(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		fmt.Println("La URL accedida: " + r.URL.String())

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		params := mux.Vars(r)

		numero := params["numero"]

		var getimporteenletras string

		row := db.Raw( "select public.getImporteEnLetras(?)", numero).Row()

		row.Scan(&getimporteenletras)

		framework.RespondJSON(w, http.StatusOK, getimporteenletras)
	}

}

func getHelperConcepto(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		fmt.Println("La URL accedida: " + r.URL.String())

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		p_tipoconcepto := r.URL.Query()["tipoconcepto"]
		p_solonovedades := r.URL.Query()["solonovedades"]

		var conceptos []structConcepto.Concepto

		var sql string
		var arrayCondiciones []string

		if p_tipoconcepto != nil {
			condicionconcepto := "(tipoconcepto.codigo = '" + p_tipoconcepto[0] + "')"
			arrayCondiciones = append(arrayCondiciones, condicionconcepto)
		}

		if p_solonovedades != nil && p_solonovedades[0] != "false" {
			condicionsolonovedades := "(esnovedad = " + p_solonovedades[0] + ")"
			arrayCondiciones = append(arrayCondiciones, condicionsolonovedades)
		}
		condicion := ""
		if len(arrayCondiciones) > 0 {
			condicion = " WHERE "
			for i := 0; i < len(arrayCondiciones); i++ {
				condicion = condicion + arrayCondiciones[i]
				if i+1 != len(arrayCondiciones) {
					condicion = condicion + " AND "
				}
			}
		}
		sql = "SELECT * FROM CONCEPTO INNER JOIN tipoconcepto ON concepto.tipoconceptoid = tipoconcepto.id" + condicion

		db.Set("gorm:auto_preload", true).Raw(sql).Scan(&conceptos)

		framework.RespondJSON(w, http.StatusOK, conceptos)
	}
}

func getHelperTipoimpuestoganancias(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		fmt.Println("La URL accedida: " + r.URL.String())

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		p_tipoconcepto := r.URL.Query()["tipoconcepto"]

		var tipoimpuestoganancias []structConcepto.Tipoimpuestoganancias
		condicion := ""

		if p_tipoconcepto != nil {
			condicion = " WHERE APLICA" + strings.ReplaceAll(p_tipoconcepto[0], "_", "") + " = true"

		}
		sql := "SELECT * FROM TIPOIMPUESTOGANANCIAS" + condicion

		db.Set("gorm:auto_preload", true).Raw(sql).Scan(&tipoimpuestoganancias)

		framework.RespondJSON(w, http.StatusOK, tipoimpuestoganancias)
	}
}

func getHelperFunction(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		fmt.Println("La URL accedida: " + r.URL.String())

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		p_tipo := r.URL.Query()["tipoformulas"][0]

		var helpers []structHelper.HelperFunction

		var sql string

		condicion := ""

		if p_tipo != "sistema" {
			condicion = "and tabla2.type = 'generic'"
		}

		sql = "select tabla2.name as ID, tabla2.name as nombre, tabla2.name as codigo, tabla2.description as descripcion from " + tenant + ".function as tabla2 left join " + tenant + ".param as p on p.functionname = tabla2.name and p.deleted_at is null where p.id is null and tabla2.deleted_at is null and tabla2.result = 'number' " + condicion

		db.Set("gorm:auto_preload", true).Raw(sql).Scan(&helpers)

		framework.RespondJSON(w, http.StatusOK, helpers)
	}
}

func getHelperConceptoafip(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		fmt.Println("La URL accedida: " + r.URL.String())

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		p_tipo := r.URL.Query()["tipoconcepto"]

		var conceptofipHelpers []structHelper.Helper
		var consulta string = ""

		if p_tipo != nil {
			p_tipoconcepto := p_tipo[0]
			if p_tipoconcepto == "DESCUENTO" {
				p_tipoconcepto = "IMPORTE_REMUNERATIVO"
			}
			consulta = "INNER JOIN TIPOCONCEPTO TC ON TC.ID = CA.TIPOCONCEPTOID WHERE TC.CODIGO = '" + p_tipoconcepto + "'"
		}

		sql := "SELECT ca.id as id, ca.codigo || ' - ' || ca.nombre as nombre, ca.codigo as codigo, ca.descripcion as descripcion FROM CONCEPTOAFIP CA " + consulta
		if err := db.Set("gorm:auto_preload", true).Raw(sql).Scan(&conceptofipHelpers).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		} else {
			framework.RespondJSON(w, http.StatusOK, conceptofipHelpers)
		}
	}
}

func getHelperLegajo(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		fmt.Println("La URL accedida: " + r.URL.String())

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		var legajosHelpers []structHelper.Helper

		sql := "select id, legajo || ' - ' || nombre || ', ' || apellido as nombre, codigo, descripcion from legajo where deleted_at is null and activo = 1"

		if err := db.Set("gorm:auto_preload", true).Raw(sql).Scan(&legajosHelpers).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		} else {
			framework.RespondJSON(w, http.StatusOK, legajosHelpers)
		}
	}
}

func getHelperObrasocial(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		fmt.Println("La URL accedida: " + r.URL.String())

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		var obrasocialHelpers []structHelper.Helper

		sql := "select id, codigo || ' - ' || nombre as nombre, codigo, descripcion from obrasocial where deleted_at is null and activo = 1"

		if err := db.Set("gorm:auto_preload", true).Raw(sql).Scan(&obrasocialHelpers).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		} else {
			framework.RespondJSON(w, http.StatusOK, obrasocialHelpers)
		}
	}
}
