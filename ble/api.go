package ble

import (
	"fmt"
	"log"

	"tinygo.org/x/bluetooth"
)

func handleError() {
	if r := recover(); r != nil {
		log.Println(r)
	}
}
func Scan() ([]bluetooth.DeviceCharacteristic, error) {
	defer handleError()
	SERVICE_UUID, err := bluetooth.ParseUUID(CONST_SVC_UUID)
	_ = SERVICE_UUID

	if err != nil {
		return nil, err
	}

	var characteristic []bluetooth.DeviceCharacteristic = nil

	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		// println("found device:", device.Address.String(), device.RSSI, device.LocalName(), string(device.AdvertisementPayload.Bytes()))
		if len(device.LocalName()) == 0 {
			return
		}

		if device.AdvertisementPayload.HasServiceUUID(SERVICE_UUID) && device.LocalName() == "spap2110" {
			var err error
			dev, err = adapter.Connect(device.Address, bluetooth.ConnectionParams{})
			if err != nil {
				panic(err)
			}

			srvcs, err := dev.DiscoverServices([]bluetooth.UUID{})
			if err != nil {
				panic(err)
			}

			fmt.Println("svcs:", srvcs)

			characteristic, err = srvcs[1].DiscoverCharacteristics(nil)
			if err != nil {
				panic(err)
			}
			fmt.Println("Connected!!")
			adapter.StopScan()
		}
	})

	// character[0].EnableNotifications(func(buf []byte) {
	// 	fmt.Println("read:", len(buf), buf)
	// })

	// fmt.Println("Receiving the message")

	// do(character[1])
	// for {
	// 	var cmd string
	// 	fmt.Scanln(&cmd)
	// 	msg := makeMessage(cmd)
	// 	n, err := character[1].WriteWithoutResponse(msg)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("sent length:", n)
	// }
	return characteristic, err
}
