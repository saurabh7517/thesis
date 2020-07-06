package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/saurbh7517/artifact/errorhandler"
	"github.com/saurbh7517/artifact/models"
	"github.com/saurbh7517/artifact/podfunction"
	"k8s.io/client-go/kubernetes"
)

type podController struct {
	urlpath   string
	clientset *kubernetes.Clientset
}

func (pc podController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/fileupload" {
		switch r.Method {
		case http.MethodGet:
			//get all running pods
		case http.MethodPost:
			// upload file from the client
			pc.uploadFile(w, r)
		case http.MethodPut:
			//update pod
		case http.MethodDelete:
			//delete pods
			pc.removeAllPods()
		default:
			fmt.Println("Wrong Verb Request")
		}
	}

}

func (pc *podController) uploadFile(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) //max memory of 10MB is allowed for this file Parser
	errorhandler.Check(err)
	//reading file from client
	file, handler, err := r.FormFile("myfile")
	if err != nil {
		fmt.Println("Error retrieving file from form data")
		fmt.Println(err)
		return
	}
	defer file.Close()

	//creating a temporary file
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("", "uploadfile.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	// service

	myNewFile := models.NewFile(handler.Filename, int(handler.Size), 300, tempFile)

	var mapPtr *map[int][]byte = myNewFile.ProcessFileByNewLine()

	podfunction.ProcessFileMap(mapPtr, pc.clientset)
}

func (pc *podController) removeAllPods() {
	podfunction.DeletePods(pc.clientset)
}

func newPodController(conn *kubernetes.Clientset) *podController {
	return &podController{urlpath: "/fileupload", clientset: conn}
}
