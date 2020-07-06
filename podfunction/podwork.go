package podfunction

import (
	"fmt"

	"github.com/saurbh7517/artifact/errorhandler"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	c         chan *CustomPod
	numOfPods int
	// wg        sync.WaitGroup
)

//CustomPod struct contains the details of the newly created Ppod.
type CustomPod struct {
	tempPod   *apiv1.Pod
	clientset *kubernetes.Clientset
}

//CreateNewPod creates a new pod for the client
func CreateNewPod(clientset *kubernetes.Clientset, c chan *CustomPod) {

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
	c <- &CustomPod{tempPod: newCustomPod, clientset: clientset}

	errorhandler.Check(err)
	fmt.Printf("Pods created %s", newCustomPod.Name)
	// wg.Done()
}

//WatchPod will watch the status of the newly created pod
func WatchPod(clientset *kubernetes.Clientset) {

	if numOfPods > 0 {
		customPodArray := make([]*CustomPod, numOfPods)

		// watchList := make([]*v1.ListOptions, numOfPods)

		for i := 0; i < numOfPods; i++ {
			customPodArray[i] = <-c
			fmt.Printf("Pod %s is being created... Hold on!!", customPodArray[i].tempPod.Name)
			// watchList = append(watchList, &v1.ListOptions{
			// 	Watch:           true,
			// 	ResourceVersion: customPodArray[i].tempPod.ResourceVersion,
			// 	FieldSelector:   fields.Set{"metadata.name": "test-pod"}.AsSelector().String(),
			// 	LabelSelector:   labels.Everything().String(),
			// })
		}

		//Setting the watchers for all pods
		// for i := 0; i < numOfPods; i++ {
		// 	clientset.CoreV1().Pods("default").Watch(*watchList[i])
		// }

		factory := informers.NewSharedInformerFactory(clientset, 2)
		PodInformer := factory.Core().V1().Pods().Informer()
		PodInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: sendData,
		})

	} else {
		fmt.Print("The number of pods is 0 therefore nothing to watch......")
	}

}

func sendData(obj interface{}) {
	fmt.Println(obj)

}

func updatePod() {

}

func deletePod() {

}

//ProcessFileMap will consume each block generated
func ProcessFileMap(temp *map[int][]byte, clientset *kubernetes.Clientset) {
	c := make(chan *CustomPod)
	numOfPods = len(*temp)
	for k := range *temp {
		fmt.Print("Creating a pod for block : ", k)
		// wg.Add(1)
		go CreateNewPod(clientset, c)
	}

	WatchPod(clientset)
	// wg.Wait()

}
