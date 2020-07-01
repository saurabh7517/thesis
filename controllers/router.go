package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/saurbh7517/artifact/errorhandler"
)

func RegisterController() {
	pc := newPodController()

	http.Handle("/pod", *pc)
	http.Handle("/watch", *pc)
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	err := enc.Encode(data)
	errorhandler.Check(err)
}
