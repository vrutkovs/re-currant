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
	dat, err := ioutil.ReadFile(".git/HEAD")
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

	// Run `oc apply -k <dir-name>`
	command := "oc"
	commandArgs := []string{"apply", "-k", applyPath}
	log.Printf("Running %s%v", command, commandArgs)

	cmd := exec.Command(command, commandArgs...)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf(err.Error())
		return
	}

	log.Printf(string(stdout))
}
