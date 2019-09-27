package main

import (
	"github.com/gin-gonic/gin"
)

// Env holds references to useful objects in router funcs
type Env struct{}

func main() {
	// setup webhook listener
	r := gin.Default()

	env := &Env{}

	// liveness healthcheck
	r.GET("/healthz", env.healthz)

	// git repo got updated
	r.POST("/reload", env.reload)

	// git repo got updated
	r.POST("/apply", env.apply)

	r.Run(":8080")
}
