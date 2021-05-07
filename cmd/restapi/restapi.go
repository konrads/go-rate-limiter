package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konrads/go-rate-limiter/pkg/decorator"
	"github.com/konrads/go-rate-limiter/pkg/leakybucket"
	"github.com/konrads/go-rate-limiter/pkg/model"
)

// main cmd entrypoint to start up the server.
func main() {
	restUri := flag.String("rest-uri", "0.0.0.0:8080", "rest uri")
	limitConf := flag.String("limit-conf", "limits.json", "config file for limit")
	flag.Parse()

	// fetch the conf file
	limitConfFile, err := os.Open(*limitConf)
	if err != nil {
		log.Fatalf("Failed to load config file: %s", *limitConf)
	}
	byteValue, _ := ioutil.ReadAll(limitConfFile)
	var limitRules []model.LimitRule
	err = json.Unmarshal([]byte(byteValue), &limitRules)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", *limitConf)
	}

	log.Printf(`Starting restapi service with params:
	- restUri:    %s
	- limitConf:  %s
	- limitRules: %v
	`, *restUri, *limitConf, limitRules)

	// create limiter, note, requests are limited per IP not request, hence can reuse
	var lb leakybucket.LeakyBucket = leakybucket.NewSafeLeakyBucket(limitRules)

	var pingHandler gin.HandlerFunc = func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
	var healthHandler gin.HandlerFunc = func(c *gin.Context) {
		c.JSON(200, gin.H{
			"health": "good",
		})
	}

	// setup cleanup ticker
	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)
	defer func() {
		done <- true
		defer ticker.Stop()
	}()
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				lb.Cleanup(t)
				log.Printf("...cleanup, storage stats: %d", lb.Stats())
			}
		}
	}()

	// init gin server and routes
	r := gin.Default()
	r.GET("/ping", pingHandler)
	r.GET("/health", healthHandler)
	r.GET("/pingLimited", decorator.Decorate(&lb, &pingHandler))
	r.GET("/healthLimited", decorator.Decorate(&lb, &healthHandler))
	r.Run(*restUri)
}
