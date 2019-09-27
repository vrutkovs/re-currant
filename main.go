package main

import (
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// Env holds references to useful objects in router funcs
type Env struct {
	applyPath    string
	useKustomize bool
}

func main() {
	// setup webhook listener
	r := gin.Default()

	gitRepo := os.Getenv("GIT_SYNC_REPO")
	if len(gitRepo) == 0 {
		panic("GIT_SYNC_REPO env var is not set")
	}
	subDir := os.Getenv("RECURRANT_SUBDIR")
	if len(subDir) == 0 {
		panic("RECURRANT_SUBDIR env var is not set")
	}

	gitDirSlice := strings.Split(gitRepo, "/")
	applyPath := path.Join(gitDirSlice[len(gitDirSlice)-1], subDir)

	useKustomize := false
	useKustomizeEnv := os.Getenv("RECURRANT_USE_KUSTOMIZE")
	if useKustomizeEnv == "true" {
		useKustomize = true
	}

	env := &Env{applyPath: applyPath, useKustomize: useKustomize}

	// liveness healthcheck
	r.GET("/healthz", env.healthz)

	// git repo got updated
	r.POST("/reload", env.reload)

	// git repo got updated
	r.POST("/apply", env.apply)

	r.Run(":8080")
}
