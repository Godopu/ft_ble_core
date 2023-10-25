package main

import (
	"fmt"
	"ft-healthcare-core/alg"
	"ft-healthcare-core/api"
	"ft-healthcare-core/ble"
	"ft-healthcare-core/model"
	"log"
	"time"

	"github.com/google/uuid"
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
		log.Println("scanning...")
		characters, err := ble.Scan()
		log.Println("connected")
		if err != nil {
			log.Println(err.Error())
			continue
		}
		model.Connected()

		err = do(characters[0], characters[1])
		log.Println(err.Error())
		model.DisConnected()
	}
}

func do(r, w bluetooth.DeviceCharacteristic) error {
	ackCh := make(chan interface{}, 1)

	r.EnableNotifications(func(buf []byte) {
		for i := 0; i < len(buf) && i < len(ble.ACK); i++ {
			if buf[i] != ble.ACK[i] {
				return
			}
		}
		ackCh <- "ACK"
	})

	time.Sleep(time.Second)
	log.Println("send get device info")
	w.WriteWithoutResponse(ble.GET_DEVICE_INFO)
	time.Sleep(time.Second)

	log.Println("send set intensity 0")
	w.WriteWithoutResponse(ble.SET_INTENSITY_0)
	time.Sleep(time.Second)

	log.Println("send set frequency and dutyrate")
	w.WriteWithoutResponse(ble.SET_FREQUENCY_DUTYRATE)
	time.Sleep(time.Second)

	for {
		measurementID, _ := uuid.NewUUID()
		trial := model.RetrieveTrialRequest()
		t := &model.TimeLog{}

		log.Println("[start ", measurementID, ") log> EMG 신호 측정", trial.EMG)
		t.Start()

		if trial.ViaCloud {
			fmt.Println("via cloud!!")
		} else {
			t.SetComputingTime(alg.Detect(trial.EMG))
		}

		t.EndComputing()
		_, err := w.WriteWithoutResponse(ble.SET_INTENSITY_1)
		if err != nil {
			return err
		}

		go func(mid uuid.UUID, tl *model.TimeLog) {
			t := time.NewTimer(time.Millisecond * 100)

			select {
			case <-ackCh:
				t.Stop()
			case <-t.C:
			}

			log.Println("[end ", measurementID, ") log> EMS 작동 완료")
			tl.End()
			model.AddTimeLog(tl)
		}(measurementID, t)

		time.Sleep(time.Millisecond * 215)
		_, err = w.WriteWithoutResponse(ble.SET_INTENSITY_0)
		if err != nil {
			return err
		}
	}
}
