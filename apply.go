package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func (e *Env) apply(c *gin.Context) {
	// Read current commit from .git/FETCH_HEAD
	dat, err := ioutil.ReadFile(".git/FETCH_HEAD")
	if err != nil {
		panic(fmt.Sprintf("error reading git head: %v", err))
	}
	log.Printf("Checked out commit %s", string(dat))

	// Find the resulting dir name based on GIT_SYNC_REPO env var
	gitRepo := os.Getenv("GIT_SYNC_REPO")
	gitDirSlice := strings.Split(gitRepo, "/")
	gitDir := gitDirSlice[len(gitDirSlice)-1]

	// Make sure we run in a subdir
	subDir := os.Getenv("RECURRANT_SUBDIR")
	applyPath := path.Join(gitDir, subDir)

	useKustomize := false
	useKustomizeEnv := os.Getenv("RECURRANT_USE_KUSTOMIZE")
	if useKustomizeEnv == "true" {
		useKustomize = true
	}

	// Run `oc apply -k <dir-name>`
	command := "oc"
	commandArgs := make([]string, 3)
	commandArgs[0] = "apply"
	if useKustomize {
		commandArgs[1] = "-k"
	} else {
		commandArgs[1] = "-f"
	}
	commandArgs[2] = applyPath
	log.Printf("Running %s %v", command, strings.Join(commandArgs, " "))

	cmd := exec.Command(command, commandArgs...)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf(err.Error())
		return
	}

	log.Printf(string(stdout))
	c.String(200, "done")
}
