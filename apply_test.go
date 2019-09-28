package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...) // #nosec
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

var testCases = []struct {
	Env      *Env
	Cmd      string
	Args     string
	ExitCode string
}{
	{setupEnv(), "kubectl", "apply -f checkout", ""},
	{&Env{
		applyPath:     "foo/bar",
		useKustomize:  false,
		useOC:         false,
		customCommand: "",
	}, "kubectl", "apply -f foo/bar", ""},
	{&Env{
		applyPath:     "checkout",
		useKustomize:  true,
		useOC:         false,
		customCommand: "",
	}, "kubectl", "apply -k checkout", ""},
	{&Env{
		applyPath:     "checkout",
		useKustomize:  false,
		useOC:         true,
		customCommand: "",
	}, "oc", "apply -f checkout", ""},
	{&Env{
		applyPath:     "checkout",
		useKustomize:  true,
		useOC:         true,
		customCommand: "",
	}, "oc", "apply -k checkout", ""},
	{&Env{
		applyPath:     "foo/bar",
		useKustomize:  true,
		useOC:         true,
		customCommand: "",
	}, "oc", "apply -k foo/bar", ""},
	{&Env{
		applyPath:     "foo/bar",
		useKustomize:  true,
		useOC:         true,
		customCommand: "helm release",
	}, "sh", "-c helm release", ""},
}

func TestApply(t *testing.T) {
	for i, tc := range testCases {
		tc := tc // make scopelint happy
		t.Run(string(i), func(_ *testing.T) {
			execCommand = fakeExecCommand
			defer func() { execCommand = exec.Command }()

			body := map[string]string(
				map[string]string{
					"args":      tc.Args,
					"cmd":       tc.Cmd,
					"exit_code": tc.ExitCode,
				},
			)

			router := setupRouter(tc.Env)
			w := performRequest(router, "POST", "/apply")

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.Nil(t, err)
			assert.Equal(t, body, response)
		})
	}
}
