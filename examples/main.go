package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	gorenogymodbus "github.com/michaelpeterswa/go-renogy-modbus"
)

func main() {
	renogyModbusClient, err := gorenogymodbus.NewModbusClient(log.New(os.Stdout, "test: ", log.LstdFlags), "/dev/tty.usbserial-D30F06G2")
	if err != nil {
		panic(err)
	}

	data, err := renogyModbusClient.ReadData()
	if err != nil {
		panic(err)
	}

	dci, err := gorenogymodbus.Parse(data)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(dci)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
