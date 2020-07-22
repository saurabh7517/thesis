package podfunction

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/saurbh7517/artifact/connection"
	"github.com/saurbh7517/artifact/errorhandler"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	c         chan *CustomPod
	numOfPods int
	// podNameList []*CustomPod
	podMap  map[string]*CustomPod
	wg      sync.WaitGroup
	fileMap map[int][]byte
)

//CustomPod struct contains the details of the newly created Ppod.
type CustomPod struct {
	podName      string
	isPodRunning bool
	podId        int
}

func (cp *CustomPod) changeRunningStatus() {
	if cp.isPodRunning == true {
		cp.isPodRunning = false
	} else {
		cp.isPodRunning = true
	}
}

//CreateNewPod creates a new pod for the client
func CreateNewPod(clientset *kubernetes.Clientset, podName string) {

	var clientCon *kubernetes.Clientset = clientset

	var labelMap = map[string]string{
		"app": "myapp",
	}
	argsList := make([]string, 0, 1)
	argsList = append(argsList, podName)

	// var newCustomPod *apiv1.Pod = &apiv1.Pod{}

	// testConfig := &apiv1.Pod{
	// 	TypeMeta: v1.TypeMeta{
	// 		Kind:       "Pod",
	// 		APIVersion: "v1",
	// 	},
	// 	ObjectMeta: v1.ObjectMeta{
	// 		Name:      podName,
	// 		Labels:    labelMap,
	// 		Namespace: "default",
	// 	},
	// 	Spec: apiv1.PodSpec{
	// 		Containers: []apiv1.Container{
	// 			{
	// 				Name:  "myapp",
	// 				Image: "192.168.204.151:5000/myapp",
	// 				// Args:  argsList,
	// 				Ports: []apiv1.ContainerPort{
	// 					{
	// 						ContainerPort: 8080,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	workerconfig := &apiv1.Pod{
		TypeMeta: v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      podName,
			Labels:    labelMap,
			Namespace: "default",
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "workersmall",
					Image: "192.168.204.151:5000/workersmall",
					Args:  argsList,
				},
			},
		},
	}

	_, err := clientCon.CoreV1().Pods("default").Create(workerconfig)
	// c <- &CustomPod{tempPod: newCustomPod, clientset: clientset}

	errorhandler.Check(err)
	// fmt.Println("Pods created ", newCustomPod.Name)
	wg.Done()

}

//WatchPod will watch the status of the newly created pod
func WatchPod(clientset *kubernetes.Clientset) {

	//another way of watching
	watchList := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", apiv1.NamespaceDefault, fields.Everything())
	_, controller := cache.NewInformer(
		watchList,
		&apiv1.Pod{},
		time.Second*30,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    postPodCreation,
			UpdateFunc: updateData,
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)

}

func postPodCreation(obj interface{}) {
	// fmt.Println(obj)
	pod := obj.(*apiv1.Pod)
	var podstatus *apiv1.PodStatus = &pod.Status
	var podphase *apiv1.PodPhase = &podstatus.Phase
	// pod.Status ==
	fmt.Printf("Newly created Pod name %s in phase %s \n", pod.Name, *podphase)
	fmt.Println("Checking for more resources.....")

}

func sendData(podName string) {
	var conn *nats.Conn = connection.Nc

	// Create a unique subject name for replies.
	uniqueReplyTo := nats.NewInbox()

	cp := podMap[podName]
	data, _ := fileMap[cp.podId]

	// Listen for a single response
	sub, err := conn.SubscribeSync(uniqueReplyTo)
	if err != nil {
		log.Fatal(err)
	}

	// Send the request.
	// If processing is synchronous, use Request() which returns the response message.
	if err := conn.PublishRequest(podName, uniqueReplyTo, data); err != nil {
		log.Fatal(err)
	}

	// Read the reply
	msg, err := sub.NextMsg(time.Second * 10)
	if err != nil {
		log.Fatal(err)
	}

	// Use the response
	log.Printf("Reply from %s : %s", podName, msg.Data)

}

func updateData(oldobj interface{}, obj interface{}) {
	pod := obj.(*apiv1.Pod)
	// pod.Status ==

	var podstatus *apiv1.PodStatus = &pod.Status
	var podphase *apiv1.PodPhase = &podstatus.Phase
	if pod.Name == podMap[pod.Name].podName && *podphase == "Running" {
		cp := podMap[pod.Name]
		if cp.isPodRunning == false {

			//start the data sending process
			sendData(cp.podName)
			fmt.Println("Changing isPodRunning status to true")

			cp.isPodRunning = true

		}
	}

	if *podphase == "Running" {
		fmt.Printf("Updated pod status with pod name %s and with status %s", pod.Name, *podphase)
		fmt.Println("Looking for more updates.....")
	}

}

func updatePod() {

}

//DeletePods will delete all the pods
func DeletePods(clientset *kubernetes.Clientset) {
	for x, _ := range podMap {
		fmt.Println("Deleting pod", x)
		clientset.CoreV1().Pods("default").Delete(x, &v1.DeleteOptions{})
	}

}

//ProcessFileMap will consume each block generated
func ProcessFileMap(temp map[int][]byte, clientset *kubernetes.Clientset) {

	fileMap = temp

	WatchPod(clientset)
	// c = make(chan *CustomPod)
	numOfPods = len(temp)
	// var podName string
	podMap = make(map[string]*CustomPod, 200)
	for k := range temp {
		fmt.Println("Creating a pod for block : ", k)
		wg.Add(1)

		podName := "workerpod" + strconv.Itoa(k)
		go CreateNewPod(clientset, podName)
		podMap[podName] = &CustomPod{podName: podName, isPodRunning: false, podId: k}
		// podNameList = append(podNameList, &CustomPod{podName: podName, isPodRunning: false})
	}

	wg.Wait()

}
