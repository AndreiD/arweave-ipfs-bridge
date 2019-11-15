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
	"os/signal"
	"strconv"
	"time"
)

const version = "1.0 Alpha"

var router *gin.Engine

// Configuration .
var Configuration *configs.ViperConfiguration

func init() {

	Configuration = configs.NewConfiguration()
	Configuration.Init()

	debug := Configuration.GetBool("debug")
	log.Init(debug)

	log.Println("==================================================")
	log.Println("Starting IPFS - Arweave Bridge version: " + version)
	log.Println("==================================================")
	log.Println()

}

func main() {

	router = gin.New()

	if Configuration.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(corsMiddleware())
	router.Use(configurationMiddleware(Configuration))

	InitializeRouter()

	server := &http.Server{
		Addr:           Configuration.Get("server.host") + ":" + strconv.Itoa(Configuration.GetInt("server.port")),
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

	log.Printf("Running on %s:%s", Configuration.Get("server.host"), strconv.Itoa(Configuration.GetInt("server.port")))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
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
