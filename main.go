package main

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

// Env holds references to useful objects in router funcs
type Env struct {
	applyPath    string
	useKustomize bool
}

func main() {
	// Verify env vars are set
	gitRepo := os.Getenv("GIT_SYNC_CHECKOUT")
	if len(gitRepo) == 0 {
		panic("GIT_SYNC_CHECKOUT env var is not set")
	}
	subDir := os.Getenv("RECURRANT_SUBDIR")
	if len(subDir) == 0 {
		panic("RECURRANT_SUBDIR env var is not set")
	}

	applyPath := path.Join(gitRepo, subDir)

	useKustomize := false
	useKustomizeEnv := os.Getenv("RECURRANT_USE_KUSTOMIZE")
	if useKustomizeEnv == "true" {
		useKustomize = true
	}

	env := &Env{applyPath: applyPath, useKustomize: useKustomize}

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
