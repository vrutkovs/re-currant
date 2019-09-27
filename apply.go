package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func (e *Env) apply(c *gin.Context) {
	// Find the resulting dir name based on GIT_SYNC_REPO env var
	gitRepo := os.Getenv("GIT_SYNC_REPO")
	gitDirSlice := strings.Split(gitRepo, "/")
	gitDir := gitDirSlice[len(gitDirSlice)-1]

	// Run `oc apply -k <dir-name>`
	cmd := exec.Command("oc", "apply", "-k", gitDir)
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(stdout))
}
