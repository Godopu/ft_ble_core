package api

import (
	"ft-healthcare-core/model"
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

		apiGroup.GET("/trials", GET_Trials)
		apiGroup.POST("/trial", POST_Trial)
		apiGroup.POST("/try", POST_Try)
		apiGroup.PUT("/trial", PUT_Trial)
		apiGroup.DELETE("/trial", DELETE_Trial)
		apiGroup.OPTIONS("/trial", OPTIONS_Trial)

		apiGroup.POST("/start", POST_START)
		apiGroup.POST("/end", POST_END)
	}

	// create a new gin router for static files
	staticEngine := gin.New()
	staticEngine.Static("/", "./web")

	// Create a new gin router
	r := gin.New()
	// gin.
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
