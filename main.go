package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xubiosueldos/framework/configuracion"
)

func main() {

	configuracion := configuracion.GetInstance()

	router := newRouter()

	fmt.Println("Microservicio Helpers escuchando en el puerto: " + configuracion.Puertomicroserviciohelpers)

	server := http.ListenAndServe(":"+configuracion.Puertomicroserviciohelpers, router)

	log.Fatal(server)

}
