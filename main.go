package main

import (
	"fmt"

	"github.com/saurbh7517/artifact/connection"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {

	var clientset *kubernetes.Clientset = connection.CreateConnection()

	// access the API to list pods
	pods, _ := clientset.CoreV1().Pods("").List(v1.ListOptions{})

	// podfunction.CreateNewPod(clientset)

	// deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	// storageClient := clientset.StorageV1().RESTClient()
	// request := storageClient.Post()

	// request.Body()

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

}
