package provider

import (
	"testing"
	"time"
)

func TestRunLocalStdout(t *testing.T) {
	out, errOut, done, err := runLocal("printf 'a\\tb\\n'", 5*time.Second)
	if err != nil || !done || errOut != "" || out != "a\tb\n" {
		t.Fatalf("got %q %q %v %v", out, errOut, done, err)
	}
}

func TestRunLocalQuotingMatchesShell(t *testing.T) {
	out, _, _, _ := runLocal("echo 'tank/sl.house/users/seb'", 5*time.Second)
	if out != "tank/sl.house/users/seb\n" {
		t.Fatalf("got %q", out)
	}
}

func TestRunLocalFailureIsStderr(t *testing.T) {
	_, errOut, done, err := runLocal("echo 'dataset does not exist' >&2; exit 1", 5*time.Second)
	if err != nil || !done || errOut != "dataset does not exist\n" {
		t.Fatalf("got %q %v %v", errOut, done, err)
	}
}

func TestRunLocalSilentFailureStillReported(t *testing.T) {
	_, errOut, done, _ := runLocal("exit 3", 5*time.Second)
	if !done || errOut != "command exited with status 3" {
		t.Fatalf("got %q %v", errOut, done)
	}
}

func TestRunLocalTimeoutIsNotDone(t *testing.T) {
	_, _, done, _ := runLocal("sleep 5", 200*time.Millisecond)
	if done {
		t.Fatal("timed-out command reported as done")
	}
}

func TestCallSshCommandLocal(t *testing.T) {
	out, err := callSshCommand(&Config{local: true}, "echo %s", "'ok'")
	if err != nil || out != "ok" {
		t.Fatalf("got %q %v", out, err)
	}
	_, err = callSshCommand(&Config{local: true}, "echo 'no such pool' >&2; exit 1")
	if _, ok := err.(*PoolError); !ok {
		t.Fatalf("expected PoolError, got %T %v", err, err)
	}
}
