package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zond/tesla"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

func main() {
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					onStart()
					a.Send(paint.Event{})
				}
			case paint.Event:
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}

func logf(format string, args ...interface{}) {
	format = fmt.Sprintf("tehomelink: %s", format)
	log.Printf(format, args...)
}

func onStart() {
	logf("starting")

	client, err := tesla.NewClient(
		&tesla.Auth{
			ClientID:     "e4a9949fcfa04068f59abb5a658f2bac0a3428e4652315490b659d5ab3f35a9e",
			ClientSecret: "c75f14bbadc8bee3a7594412c31416f8300256d7668ea7e6e7f06727bfb9d220",
			Email:        "your email",
			Password:     "your password",
		})
	if err != nil {
		logf("NewClient(...): %v, %v", client, err)
		panic(err)
	}

	vehicles, err := client.Vehicles()
	if err != nil {
		logf("Vehicles(): %v, %v", vehicles, err)
		panic(err)
	}

	if len(vehicles) != 1 {
		logf("not exactly 1 vehicle")
		panic("not exactly 1 vehicle")
	}

	vehicle := vehicles[0]
	status, err := vehicle.MobileEnabled()
	if err != nil {
		logf("MobileEnabled(): %v, %v", status, err)
		panic(err)
	}

	if !status {
		logf("not mobile enabled!")
		panic("not mobile enabled")
	}

	sock, err := vehicle.WebSocket()
	if err != nil {
		logf("WebSocket(): %v, %v", sock, err)
		panic(err)
	}
	defer sock.Close()

	go func() {
		for m := range sock.Output {
			logf("sock.Output => %+v", m)
		}
	}()

	if !sock.HomelinkNearby() {
		logf("HomelinkNearby() => false")
		panic(fmt.Errorf("HomelinkNearby() => false"))
	}

	sock.ActivateHomelink()

	os.Exit(0)
}
