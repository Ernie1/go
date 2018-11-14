package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default() // an Engine instance with the Logger and Recovery middleware already attached

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /?a=1&b=2
	router.GET("/", func(c *gin.Context) {
		a := c.DefaultQuery("a", "Default") // the keyed url query value if it exists, otherwise the specified defaultValue
		b := c.Query("b")                   //  the keyed url query value if it exists, otherwise an empty string

		c.String(http.StatusOK, "Hello %s %s\n", a, b) // writes the given string into the response body
	})
	router.Run(":8080") // attaches the router to a http.Server and starts listening and serving HTTP requests
}
