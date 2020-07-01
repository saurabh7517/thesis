package podfunction

import (
	"fmt"

	"github.com/saurbh7517/artifact/errorhandler"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//This is a custom type CustomPod for the newly created pod.
type CustomPod struct {
	name           string
	label          string
	nodeIP         string
	appPort        int32
	podPort        int32
	containerName  string
	containerImage string
}

//This function creates a user-defined pod
func CreateNewPod(clientset *kubernetes.Clientset) {

	var clientCon *kubernetes.Clientset = clientset

	var labelMap = map[string]string{
		"app":  "myapp",
		"type": "backend",
	}

	var newCustomPod *apiv1.Pod = &apiv1.Pod{}

	newPodConfig := &apiv1.Pod{
		TypeMeta: v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "mytestapp",
			Labels:    labelMap,
			Namespace: "default",
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "myapp",
					Image: "192.168.204.151:5000/myapp",
					Ports: []apiv1.ContainerPort{
						{
							ContainerPort: 8080,
						},
					},
				},
			},
		},
	}

	newCustomPod, err := clientCon.CoreV1().Pods("default").Create(newPodConfig)
	errorhandler.Check(err)

	fmt.Printf("Pods created %s", newCustomPod)
}

func updatePod() {

}

func deletePod() {

}
