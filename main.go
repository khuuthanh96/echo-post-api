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

const (
	maxMsgPerPath = 30
)

func main() {
	router := gin.New()
	router.LoadHTMLFiles("./index.html")

	rootPath := struct {
		mu   sync.Mutex
		Msgs []Message
	}{}

	customPaths := struct {
		mu sync.Mutex
		m  map[string][]Message
	}{
		m: make(map[string][]Message),
	}

	router.Use(
		gin.Logger(),
	)

	router.GET("/:custom_path", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"path":     c.Param("custom_path"),
			"messages": customPaths.m[c.Param("custom_path")],
		})
	})

	router.POST("/:custom_path", func(c *gin.Context) {
		key := c.Param("custom_path")

		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("err", err.Error())
			c.JSON(http.StatusOK, gin.H{})

			return
		}

		msg := createMsg(data)

		customPaths.mu.Lock()

		customPathMsgs := customPaths.m[key]
		customPaths.m[key] = addMsg(customPathMsgs, msg)

		customPaths.mu.Unlock()

		c.JSON(http.StatusOK, gin.H{})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"path":     "",
			"messages": rootPath.Msgs,
		})
	})

	router.POST("/", func(c *gin.Context) {
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("err", err.Error())
			c.JSON(http.StatusOK, gin.H{})

			return
		}

		msg := createMsg(data)

		rootPath.mu.Lock()

		rootPath.Msgs = addMsg(rootPath.Msgs, msg)

		rootPath.mu.Unlock()

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

func addMsg(msgs []Message, newMsg Message) []Message {
	msgs = append(msgs, newMsg)

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].CreatedTimestamp > msgs[j].CreatedTimestamp
	})

	// remove old msgs
	if len(msgs) > maxMsgPerPath {
		msgs = msgs[:maxMsgPerPath]
	}

	return msgs
}

func createMsg(data []byte) Message {
	now := time.Now()

	return Message{
		Payload:          string(data),
		CreatedAt:        now.Format(time.ANSIC),
		CreatedTimestamp: now.Unix(),
	}
}
