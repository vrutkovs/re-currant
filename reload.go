package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

const processName = "git-sync"

func findProcessPid() (int, error) {
	// Find all files which match /proc/<id>/comm
	matches, err := filepath.Glob("/proc/*/comm")
	if err != nil {
		return 0, err
	}

	var pid int
	for _, path := range matches {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			// Failed to read comm contents
			continue
		}
		if string(f) == processName {
			pid, err = strconv.Atoi(path[len("/proc")+1 : strings.LastIndex(path, "/")])
			if err == nil {
				return pid, nil
			}
		}
	}
	return 0, errors.New("failed to find pid of the process")
}

func (e *Env) reload(c *gin.Context) {
	pid, err := findProcessPid()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}
	// Kill the process
	proc.Signal(syscall.SIGHUP)

	c.JSON(http.StatusOK, gin.H{
		"message": "sidecar reloaded",
	})
}
