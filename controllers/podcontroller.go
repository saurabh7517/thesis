package controllers

import (
	"fmt"
	"net/http"
)

type podController struct {
	urlpath string
}

func (pc podController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/pod" {
		switch r.Method {
		case http.MethodGet:
			//get all running pods
		case http.MethodPost:
			// start a pod depending on congiguration
		case http.MethodPut:
			//update pod
		case http.MethodDelete:
			//delete pod
		default:
			fmt.Println("Wrong Verb Request")
		}
	}

}

func newPodController() *podController {
	return &podController{urlpath: "/pod"}
}
