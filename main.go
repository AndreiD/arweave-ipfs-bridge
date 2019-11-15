package main

import (
	"aif/configs"
	"aif/utils/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"
)

const version = "1.0 Alpha"

var router *gin.Engine

// Configuration .
var configuration *configs.ViperConfiguration

func init() {

	configuration = configs.NewConfiguration()
	configuration.Init()

	debug := configuration.GetBool("debug")
	log.Init(debug)

	log.Println("==================================================")
	log.Println("Starting IPFS - Arweave Bridge version: " + version)
	log.Println("==================================================")
	log.Println()

	// check if IPFS daemon is running
	out, err := exec.Command("ipfs", "swarm", "peers").CombinedOutput()
	if err != nil {
		log.Fatalf("please check if IPFS daemon is running. Error %s", string(out))
	}
}

func main() {

	router = gin.New()

	if configuration.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(corsMiddleware())
	router.Use(configurationMiddleware(configuration))

	InitializeRouter()

	server := &http.Server{
		Addr:           configuration.Get("server.host") + ":" + strconv.Itoa(configuration.GetInt("server.port")),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 10, // 1Mb
	}
	server.SetKeepAlivesEnabled(true)

	// Serve'em
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	log.Printf("Running on %s:%s", configuration.Get("server.host"), strconv.Itoa(configuration.GetInt("server.port")))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("initiated server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}
	log.Println("server exiting. bye!")
}

// requestIDMiddleware adds x-request-id
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// configurationMiddleware will add the configuration to the context
func configurationMiddleware(config *configs.ViperConfiguration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("configuration", config)
		c.Next()
	}
}
