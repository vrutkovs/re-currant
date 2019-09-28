package main

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

// Env holds references to useful objects in router funcs
type Env struct {
	applyPath     string
	useKustomize  bool
	customCommand string
	useOC         bool
}

const gitCheckoutPath = "checkout"

func main() {
	// Verify env vars are set
	subDir := os.Getenv("RECURRANT_SUBDIR")
	if len(subDir) == 0 {
		subDir = "."
	}

	applyPath := path.Join(gitCheckoutPath, subDir)

	useKustomize := false
	useKustomizeEnv := os.Getenv("RECURRANT_USE_KUSTOMIZE")
	if useKustomizeEnv == "true" {
		useKustomize = true
	}

	useOC := false
	useOCEnv := os.Getenv("RECURRANT_USE_OC")
	if useOCEnv == "true" {
		useOC = true
	}

	env := &Env{
		applyPath:     applyPath,
		useKustomize:  useKustomize,
		useOC:         useOC,
		customCommand: os.Getenv("RECURRANT_COMMAND"),
	}

	// setup Gin
	r := gin.New()
	r.Use(gin.Recovery())
	// Private route - won't be logged
	// liveness healthcheck
	private := r.Group("/")
	private.GET("/healthz", env.healthz)

	// Public routes
	// git repo got updated
	public := r.Group("/")
	public.Use(gin.Logger())
	public.POST("/reload", env.reload)

	// git repo got updated
	public.POST("/apply", env.apply)

	r.Run(":8080")
}
