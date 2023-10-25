package main

import (
	"fmt"
	"ft-healthcare-core/api"
	"ft-healthcare-core/ble"
	"ft-healthcare-core/model"
	"time"

	"tinygo.org/x/bluetooth"
)

func main() {
	go func() {
		srv := api.NewWebServer(":9910")
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	for {
		fmt.Println("scanning...")
		characters, err := ble.Scan()
		fmt.Println("connected")
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		model.Connected()

		fmt.Println("Start communication")
		characters[0].EnableNotifications(func(buf []byte) {
			fmt.Print("read> len(buf):", len(buf), "] - ")
			for i := 0; i < len(buf); i++ {
				fmt.Printf("%x ", buf[i])
			}
			fmt.Println()
		})

		err = do(characters[1])
		fmt.Println(err.Error())
		model.DisConnected()
	}
}

func do(dc bluetooth.DeviceCharacteristic) error {

	time.Sleep(time.Second)
	fmt.Println("send get device info")
	dc.WriteWithoutResponse(ble.GET_DEVICE_INFO)
	time.Sleep(time.Second)

	fmt.Println("send set intensity 0")
	dc.WriteWithoutResponse(ble.SET_INTENSITY_0)
	time.Sleep(time.Second)

	fmt.Println("send set frequency and dutyrate")
	dc.WriteWithoutResponse(ble.SET_FREQUENCY_DUTYRATE)
	time.Sleep(time.Second)

	var idx = 0
	var cmd string
	for {
		fmt.Print("Enter>")
		fmt.Scanln(&cmd)
		if cmd == "exit" {
			break
		}
		// _, err := dc.WriteWithoutResponse(ble.SET_INTENSITIES[idx])
		_, err := dc.WriteWithoutResponse(ble.SET_INTENSITY_1)
		time.Sleep(time.Millisecond * 220)
		_, err = dc.WriteWithoutResponse(ble.SET_INTENSITY_0)
		idx = (idx + 1) % 9
		if err != nil {
			return err
		}
	}

	return nil
}
