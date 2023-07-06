package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(
		gin.Logger(),
	)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.POST("/", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("err", err.Error())
		} else {
			fmt.Printf("\n\nrequest payload: %v\n", string(data))
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	port := "8080"
	args := os.Args

	if len(args) > 1 {
		if _, err := strconv.ParseInt(args[1], 10, 64); err == nil {
			port = args[1]
		}
	}

	router.Run(":" + port)
}
