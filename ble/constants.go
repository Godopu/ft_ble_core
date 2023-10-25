package ble

type BleDevStatus int

const (
	BLE_STATUS_CONNECTED BleDevStatus = iota
	BLE_STATUS_DISCONNECTED
)

var CONST_SVC_UUID = "0000fff0-0000-1000-8000-00805f9b34fb"
