package main

import (
	"bytes"
	"encoding/json"
	"ft-healthcare-core/alg"
	"ft-healthcare-core/api"
	"ft-healthcare-core/ble"
	"ft-healthcare-core/model"
	"log"
	"net/http"
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

	ticker := time.NewTicker(time.Second * 3)
	for {
		ticker.Reset(time.Second * 3)

		measurementID, _ := uuid.NewUUID()

		trial := model.RetrieveTrialRequest()

		t := &model.TimeLog{ExpType: "LOCAL"}

		log.Println("[start ", measurementID, ") log> EMG 신호 측정", trial.EMG)
		t.Start()

		if trial.ViaCloud {
			log.Println("[OPTIMIZING] - offload EMG to Cloud")
			t.ExpType = "CLOUD"
			payload := &bytes.Buffer{}
			enc := json.NewEncoder(payload)
			err := enc.Encode(map[string]interface{}{
				"emg": trial.EMG,
			})
			if err != nil {
				return err
			}

			resp, err := http.Post(
				"http://34.122.255.36:9910/api/emg",
				"application/json",
				payload,
			)
			if err != nil {
				return err
			}

			dec := json.NewDecoder(resp.Body)
			body := map[string]interface{}{}
			err = dec.Decode(&body)
			if err != nil {
				return err
			}

			t.SetComputingTime(uint64(body["computing_delay"].(float64)))

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
