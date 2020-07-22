package connection

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

var (
	Nc *nats.Conn
)

func GetNatsConnection() {
	var err error
	Nc, err = nats.Connect("192.168.204.151:30008")
	if err != nil {
		fmt.Println("Cannot establish connection to NATS server")
	}
}
