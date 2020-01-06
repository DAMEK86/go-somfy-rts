package main

import (
	"flag"
	"fmt"
	"github.com/damek86/go-somfy-rts"
	"github.com/damek86/go-somfy-rts/shutters"
	"github.com/warthog618/gpio"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	pin int
)

func main() {
	remoteAddrPtr := flag.Uint("remote", 0x279621, "24-bit unique somfy-rts address")
	rollingCodePtr := flag.Uint("code", 1, "rolling code / starting number for testing")
	cmdPtr := flag.String("cmd", "server", "action to run: up, down, my / stop, program or server")
	serverPortPtr := flag.Uint("port", 8001, "server port, default is 8001")
	gpioPtr := flag.Int("gpio", gpio.GPIO4, "raspberry-pi gpio to control the 422 MHz transmitter")
	flag.Parse()

	pin = *gpioPtr

	err := gpio.Open()
	defer gpio.Close()
	if err != nil {
		fmt.Println("can not open rpio", err)
	}

	service := shutters.NewService(somfy.DefaultEncryptionKey, sendCommand)

	switch strings.ToLower(*cmdPtr) {
	case "server":
		startServer(uint32(*serverPortPtr), service)
	case "up":
		service.MoveUp(uint32(*remoteAddrPtr), uint16(*rollingCodePtr))
	case "my":
	case "stop":
		service.MoveMy(uint32(*remoteAddrPtr), uint16(*rollingCodePtr))
	case "down":
		service.MoveDown(uint32(*remoteAddrPtr), uint16(*rollingCodePtr))
	case "program":
		service.Program(uint32(*remoteAddrPtr), uint16(*rollingCodePtr))
	}
}

func startServer(port uint32, service shutters.Service) {
	router := http.NewServeMux()
	controller := shutters.NewController(service)
	controller.CreateEndpoints(router)

	host := fmt.Sprintf(":%d", port)
	fmt.Printf("starting http server on %s", host)
	err := http.ListenAndServe(host, router)
	if err != nil {
		fmt.Println("can not start server", err)
	}
	os.Exit(1)
}

func sendCommand(pulseWave []somfy.Pulse) {
	tx := gpio.NewPin(pin)
	tx.Output()
	for _, pulse := range pulseWave {
		if pulse.IsHigh {
			tx.High()
		} else {
			tx.Low()
		}
		time.Sleep(pulse.Length)
	}
}
