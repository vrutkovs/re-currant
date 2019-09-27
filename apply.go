package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func (e *Env) apply(c *gin.Context) {
	// Find the resulting dir name based on GIT_SYNC_REPO env var
	gitRepo := os.Getenv("GIT_SYNC_REPO")
	gitDirSlice := strings.Split(gitRepo, "/")
	gitDir := gitDirSlice[len(gitDirSlice)-1]

	// Make sure we run in a subdir
	subDir := os.Getenv("RECURRANT_SUBDIR")
	applyPath := path.Join(gitDir, subDir)

	// Run `oc apply -k <dir-name>`
	cmd := exec.Command("oc", "apply", "-k", applyPath)
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(stdout))
}
