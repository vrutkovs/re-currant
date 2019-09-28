package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func fakefilePathGlob(pattern string) (matches []string, err error) {
	// Add two fake paths. Last one is real and should return syscall.Getpid()
	result := make([]string, 3)
	result[0] = "/proc/123/comm"
	result[1] = "/proc/345/comm"
	result[2] = fmt.Sprintf("/proc/%d/comm", syscall.Getpid())
	return result, nil
}

func fakeioutilReadFile(filename string) ([]byte, error) {
	switch filename {
	case "/proc/123/comm":
		return []byte("foo"), errors.New("bar")
	case "/proc/345/comm":
		return []byte("nope"), nil
	default:
		return []byte("git-sync"), nil
	}
}

func waitSig(t *testing.T, c <-chan os.Signal, sig os.Signal) {
	select {
	case s := <-c:
		if s != sig {
			t.Fatalf("signal was %v, want %v", s, sig)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout waiting for %v", sig)
	}
}

func TestReload(t *testing.T) {
	body := gin.H{
		"message": "sidecar reloaded",
	}

	filePathGlob = fakefilePathGlob
	defer func() { filePathGlob = filepath.Glob }()
	ioutilReadFile = fakeioutilReadFile
	defer func() { ioutilReadFile = ioutil.ReadFile }()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	defer signal.Stop(c)

	env := setupEnv()
	router := setupRouter(env)

	//lint:ignore SA2002 this is fine
	go waitSig(t, c, syscall.SIGHUP)
	w := performRequest(router, "POST", "/reload")

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	message, exists := response["message"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["message"], message)
}
