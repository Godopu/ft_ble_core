package ble

import "tinygo.org/x/bluetooth"

func init() {
	// init
	adapter = bluetooth.DefaultAdapter
	must("enable BLE stack", adapter.Enable())
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
