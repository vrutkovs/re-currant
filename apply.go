package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
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

	var command string
	var commandArgs []string

	if len(e.customCommand) > 0 {
		// Run customCommand in a shell
		command = "sh"
		commandArgs = []string{"-c", e.customCommand}
	} else {
		// Run `kubectl apply -k <dir-name>`
		command = "kubectl"
		commandArgs = make([]string, 3)
		commandArgs[0] = "apply"
		if e.useKustomize {
			commandArgs[1] = "-k"
		} else {
			commandArgs[1] = "-f"
		}
		commandArgs[2] = e.applyPath
	}

	log.Printf("Running %s %v", command, strings.Join(commandArgs, " "))

	cmd := exec.Command(command, commandArgs...)
	output, err := cmd.CombinedOutput()
	log.Printf(string(output))

	if err != nil {
		log.Printf(err.Error())
		return
	}

	c.String(200, "done")
}
