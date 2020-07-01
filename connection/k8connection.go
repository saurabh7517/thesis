package connection

import (
	"flag"
	"path/filepath"

	"github.com/saurbh7517/artifact/errorhandler"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// CreateConnction creates a connection to the kubernetes api-server
func CreateConnection() *kubernetes.Clientset {
	var kubeconfig *string
	home := homedir.HomeDir()
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "/home/saurabh/.kube/config")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	errorhandler.Check(err)
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)

	return clientset
}
