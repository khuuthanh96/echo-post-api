package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Payload          interface{} `json:"payload"`
	CreatedAt        string      `json:"created_at`
	CreatedTimestamp int64       `json:"-"`
}

type Messages struct {
	mu   sync.Mutex
	Msgs []Message
}

func main() {
	router := gin.New()
	router.LoadHTMLFiles("./index.html")

	msgs := Messages{}
	maxMsg := 30

	router.Use(
		gin.Logger(),
	)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"messages": msgs.Msgs,
		})
	})

	router.POST("/", func(c *gin.Context) {
		now := time.Now()
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("err", err.Error())
			c.JSON(http.StatusOK, gin.H{})

			return
		}

		msg := Message{
			Payload:          string(data),
			CreatedAt:        now.Format(time.ANSIC),
			CreatedTimestamp: now.Unix(),
		}

		msgs.mu.Lock()
		msgs.Msgs = append(msgs.Msgs, msg)

		sort.Slice(msgs.Msgs, func(i, j int) bool {
			return msgs.Msgs[i].CreatedTimestamp > msgs.Msgs[j].CreatedTimestamp
		})

		// remove old msgs
		if len(msgs.Msgs) > maxMsg {
			msgs.Msgs = msgs.Msgs[:maxMsg]
		}

		msgs.mu.Unlock()

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
