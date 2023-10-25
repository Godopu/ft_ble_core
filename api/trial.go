package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"ft-healthcare-core/model"
	"io"
	"log"
	"net/http"

	"math/rand"

	"github.com/gin-gonic/gin"
)

func GET_Trials(c *gin.Context) {
	c.JSON(http.StatusOK, model.Trials)
}

func POST_Trial(c *gin.Context) {
	var payload model.Trial
	dec := json.NewDecoder(c.Request.Body)
	err := dec.Decode(&payload)
	if err != nil {
		panic(err)
	}

	log.Println(payload)
	model.Trials = append(model.Trials, &payload)
	c.String(http.StatusOK, "Okay")
}

func POST_Try(c *gin.Context) {
	var payload model.Trial

	dec := json.NewDecoder(c.Request.Body)
	err := dec.Decode(&payload)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 4; i++ {
		if payload.EMG[i] == 0 {
			payload.EMG[i] = rand.Intn(1000)
		}
	}

	model.RequestTrial(&payload)

	c.String(http.StatusOK, "Okay")
}

func POST_START(c *gin.Context) {
	model.ClearLogStack()
	c.String(http.StatusOK, "Okay")
}

func POST_END(c *gin.Context) {
	for _, e := range model.LogStack {
		fmt.Println(e.String())
	}
	c.JSON(http.StatusOK, model.LogStack)
}

func PUT_Trial(c *gin.Context) {
	log.Println("PUT")
	payload := map[string]interface{}{}

	dec := json.NewDecoder(c.Request.Body)
	err := dec.Decode(&payload)
	if err != nil {
		panic(err)
	}

	idx, ok := payload["index"].(int)
	if !ok {
		panic(errors.New("invalid index error"))
	}
	log.Println(idx)

	obj, ok := payload["trial"].(map[string]interface{})
	log.Println(obj)

	c.String(http.StatusOK, "OKay")
}

func DELETE_Trial(c *gin.Context) {
	log.Println("DELETE")
	payload := map[string]interface{}{}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	log.Println(payload)

	idx, ok := payload["index"].(float64)
	if !ok {
		panic(errors.New("invalid index error"))
	}

	model.RemoveItemFromTrials(int(idx))
	c.String(http.StatusOK, "OKay")
}

func OPTIONS_Trial(c *gin.Context) {
	log.Println(c.Request.Method)
	c.String(http.StatusOK, "OKay")
}
