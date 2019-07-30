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

	server := http.ListenAndServe(":"+configuracion.Puertomicroserviciohelpers, router)

	fmt.Println("Microservicio Helpers escuchando en el puerto: " + configuracion.Puertomicroserviciohelpers)

	log.Fatal(server)

}
