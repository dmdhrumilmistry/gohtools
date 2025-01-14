package main

import (
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	// Define the peripheral device info.
	adv := adapter.Scan(func(a *bluetooth.Adapter, sr bluetooth.ScanResult) {
		log.Print(a, sr)
	})

	// adapter., bluetooth.ScanResult

	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "Go Bluetooth",
	}))

	// Start advertising
	must("start adv", adv.Start())

	println("advertising...")
	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
