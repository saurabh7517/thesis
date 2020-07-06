package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/saurbh7517/artifact/connection"
	"github.com/saurbh7517/artifact/errorhandler"
)

//RegisterController is creating controllers for authentic requests
func RegisterController() {

	conn := connection.CreateConnection()

	pc := newPodController(conn)

	http.Handle("/fileupload", *pc)
	http.Handle("/watch", *pc)
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	err := enc.Encode(data)
	errorhandler.Check(err)
}
