package helper

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestExecuteCommandKeepsExitCode(t *testing.T) {
	for _, code := range []string{"0", "1", "2", "42"} {
		er := ExecuteCommand("sh -c 'echo out; exit "+code+"'", 5, true)
		if strconv.Itoa(er.ReturnCode) != code {
			t.Errorf("exit %s: got ReturnCode %d, output %q", code, er.ReturnCode, er.Output)
		}
		if !strings.Contains(er.Output, "out") {
			t.Errorf("exit %s: output lost: %q", code, er.Output)
		}
	}
}

func TestExecuteCommandTimeout(t *testing.T) {
	start := time.Now()
	er := ExecuteCommand("sleep 30", 1, true)
	if er.ReturnCode != 124 {
		t.Errorf("timeout: got ReturnCode %d, want 124, output %q", er.ReturnCode, er.Output)
	}
	if d := time.Since(start); d > 10*time.Second {
		t.Errorf("command was not killed after the timeout, took %s", d)
	}
	if !strings.Contains(er.Output, "timed out after 1s") {
		t.Errorf("output does not mention the timeout: %q", er.Output)
	}
}

func TestExecuteCommandNoTimeoutAndMissingCommand(t *testing.T) {
	if er := ExecuteCommand("true", 0, true); er.ReturnCode != 0 {
		t.Errorf("timeout 0 must mean no timeout: %+v", er)
	}
	if er := ExecuteCommand("/nonexistent/command", 5, true); er.ReturnCode != 1 {
		t.Errorf("command that can not be started: got %d, want 1", er.ReturnCode)
	}
}
