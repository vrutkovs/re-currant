package main

import (
	"io"
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

func (e *Env) reload(c *gin.Context) {
	// Iterate over /proc/<id>/comm and find `git-sync`
	err := filepath.Walk("/proc", func(path string, info os.FileInfo, err error) error {

		// TODO: Use regexp here
		if strings.Count(path, "/") == 3 {
			if strings.Contains(path, "/comm") {
				pid, err := strconv.Atoi(path[6:strings.LastIndex(path, "/")])
				if err != nil {
					log.Println(err)
					return nil
				}

				// Read comm
				f, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				if string(f) == processName {
					log.Printf("Restarting git-sync process via sending HUP to pid %d", pid)
					proc, err := os.FindProcess(pid)
					if err != nil {
						log.Println(err)
					}
					// Kill the process
					proc.Signal(syscall.SIGHUP)

					// return fake error to stop the walk
					return io.EOF
				}
			}
		}
		return nil
	})
	if err != nil {
		if err == io.EOF {
			// Not an error, just a signal when we are done
			err = nil
		} else {
			log.Fatal(err)
		}
	}
	c.JSON(http.StatusOK, "{'message': 'ok'}")
}
