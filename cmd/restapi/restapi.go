package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/konrads/go-rate-limiter/pkg/db"
	"github.com/konrads/go-rate-limiter/pkg/decorator"
	"github.com/konrads/go-rate-limiter/pkg/limiter"
	"github.com/konrads/go-rate-limiter/pkg/model"
)

// main cmd entrypoint to start up the server.
func main() {
	restUri := flag.String("rest-uri", "0.0.0.0:8080", "rest uri")
	limitConf := flag.String("limit-conf", "limits.json", "config file for limit")
	flag.Parse()

	log.Printf(`Starting restapi service with params:
	- restUri:   %s
	- limitConf: %s
	`, *restUri, *limitConf)

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

	// init memory db
	var db db.DB = db.NewMemDb()

	// create limiter, note, requests are limited per IP not request, hence can reuse
	l := limiter.NewLimiter(limitRules, &db)

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

	// init server and routes
	r := gin.Default()
	r.GET("/ping", pingHandler)
	r.GET("/health", healthHandler)
	r.GET("/pingLimited", decorator.Decorate(&l, &pingHandler))
	r.GET("/healthLimited", decorator.Decorate(&l, &healthHandler))
	r.Run(*restUri)
}
