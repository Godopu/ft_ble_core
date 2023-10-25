package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"ft-healthcare-core/model"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewWebServer(addr string) *http.Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: createRouter(),
	}

	return srv
}

func createRouter() *gin.Engine {
	// Create a new gin router for api
	// What is difference between gin.Default() and gin.New()?
	// https://stackoverflow.com/questions/44318441/what-is-difference-between-gin-default-and-gin-new

	apiEngine := gin.New()
	apiGroup := apiEngine.Group("/api")
	{
		apiGroup.GET("/meta", GET_MetaInformation)
		apiGroup.GET("/intensity", GET_Intensities)
		apiGroup.POST("/intensity", POST_Intensity)
		apiGroup.DELETE("/intensity", DELETE_Intensity)
		apiGroup.OPTIONS("/intensity", OPTIONS_Intensity)
	}

	// create a new gin router for static files
	staticEngine := gin.New()
	staticEngine.Static("/", "./web")

	// Create a new gin router
	r := gin.Default()
	// r can accept all messages from apiEngine and staticEngine
	r.Any("/*any", func(c *gin.Context) {
		defer handleError(c)
		w := c.Writer
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		path := c.Param("any")

		if strings.HasPrefix(path, "/api") {
			apiEngine.ServeHTTP(c.Writer, c.Request)
		} else {
			staticEngine.HandleContext(c)
		}

	})

	// Return the router
	return r
}

func GET_MetaInformation(c *gin.Context) {
	c.JSON(http.StatusOK, model.MetaInfo)
}

func GET_Intensities(c *gin.Context) {
	c.JSON(http.StatusOK, model.Intensities)
}

func POST_Intensity(c *gin.Context) {
	var payload model.Intensity
	err := c.BindJSON(&payload)
	if err != nil {
		panic(err)
	}

	fmt.Println(payload)
	model.Intensities = append(model.Intensities, payload)
	c.String(http.StatusOK, "Okay")
}

func PUT_Intensity(c *gin.Context) {
	fmt.Println("PUT")
	payload := map[string]interface{}{}
	err := c.BindJSON(&payload)
	if err != nil {
		panic(err)
	}

	idx, ok := payload["index"].(int)
	if !ok {
		panic(errors.New("invalid index error"))
	}
	fmt.Println(idx)

	// obj, ok := payload["intensity"].(map[string]interface{})
	// fmt.Println(obj)

	c.String(http.StatusOK, "OKay")
}

func DELETE_Intensity(c *gin.Context) {
	fmt.Println("DELETE")
	payload := map[string]interface{}{}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	fmt.Println(payload)

	idx, ok := payload["index"].(float64)
	if !ok {
		panic(errors.New("invalid index error"))
	}

	model.RemoveItemFromIntensities(int(idx))
	c.String(http.StatusOK, "OKay")
}

func OPTIONS_Intensity(c *gin.Context) {
	fmt.Println(c.Request.Method)
	c.String(http.StatusOK, "OKay")
}
