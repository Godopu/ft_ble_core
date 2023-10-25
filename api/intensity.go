package api

import (
	"encoding/json"
	"errors"
	"ft-healthcare-core/model"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GET_Intensities(c *gin.Context) {
	c.JSON(http.StatusOK, model.Intensities)
}

func POST_Intensity(c *gin.Context) {
	var payload model.Intensity
	err := c.BindJSON(&payload)
	if err != nil {
		panic(err)
	}

	log.Println(payload)
	model.Intensities = append(model.Intensities, payload)
	c.String(http.StatusOK, "Okay")
}

func PUT_Intensity(c *gin.Context) {
	log.Println("PUT")
	payload := map[string]interface{}{}
	err := c.BindJSON(&payload)
	if err != nil {
		panic(err)
	}

	idx, ok := payload["index"].(int)
	if !ok {
		panic(errors.New("invalid index error"))
	}
	log.Println(idx)

	// obj, ok := payload["intensity"].(map[string]interface{})
	// log.Println(obj)

	c.String(http.StatusOK, "OKay")
}

func DELETE_Intensity(c *gin.Context) {
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

	model.RemoveItemFromIntensities(int(idx))
	c.String(http.StatusOK, "OKay")
}

func OPTIONS_Intensity(c *gin.Context) {
	log.Println(c.Request.Method)
	c.String(http.StatusOK, "OKay")
}
