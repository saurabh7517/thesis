package main

import (
	"net/http"

	"github.com/saurbh7517/artifact/connection"
	"github.com/saurbh7517/artifact/controllers"
)

func main() {

	connection.GetNatsConnection()

	controllers.RegisterController()

	http.ListenAndServe(":8080", nil)
}
