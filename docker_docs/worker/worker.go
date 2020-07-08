package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	nats "github.com/nats-io/nats.go"
)

var wg sync.WaitGroup

func main() {

	subjectToRecv := os.Args[1]
	nc, err := nats.Connect("192.168.166.159:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to a topic to receive data for processing
	wg.Add(1)

	_, err = nc.Subscribe(subjectToRecv, handleIncomingMsg)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()

}

func handleIncomingMsg(msg *nats.Msg) {
	var resultByte *[]byte = processMsg(msg.Data)
	// Use the response
	log.Printf("Reply: %s", msg.Data)
	msg.Respond(*resultByte)
	wg.Done()
}

func processMsg(msg []byte) *[]byte {
	stringData := string(msg)
	var lines []string = strings.Split(stringData, "\n")
	words := make([]string, 0, 50)
	wordmap := make(map[string]int)
	var isPresent bool
	for _, k := range lines {
		words = strings.Split(k, " ")
		for _, k := range words {
			_, isPresent = wordmap[k]
			if isPresent == false {
				wordmap[k] = 1
			} else {
				wordmap[k] = wordmap[k] + 1
			}
		}
	}
	var result string = ""
	for k, v := range wordmap {
		result = result + k + " " + strconv.Itoa(v) + "\n"
	}

	var resultBytes []byte = []byte(result)
	return &resultBytes

}
