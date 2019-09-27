package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// setup webhook listener
	r := gin.Default()

	// liveness healthcheck
	// r.GET("/healthz", env.incoming)

	// git repo got updated
	// r.POST("/reload", env.travisMessage)

	r.Run(":8080")
}
