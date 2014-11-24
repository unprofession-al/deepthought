package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Log struct {
	Timestamp string        `json:"timestamp"`
	Status    int           `json:"status"`
	Method    string        `json:"method"`
	Request   string        `json:"request"`
	Latency   time.Duration `json:"latency"`
}

func LogJSON() gin.HandlerFunc {
	out := log.New(os.Stdout, "", 0)

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()

		l := &Log{
			Timestamp: end.Format("2006/01/02-15:04:05.000"),
			Status:    c.Writer.Status(),
			Latency:   end.Sub(start),
			Method:    c.Request.Method,
			Request:   c.Request.URL.Path,
		}

		b, _ := json.Marshal(l)

		out.Println(string(b))
	}
}

func SetCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if origin := c.Request.Header.Get("Origin"); origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}
