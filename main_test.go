package main

import (
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/gaffatape-io/rubberduck/header"
	rubberducktesting "github.com/gaffatape-io/rubberduck/testing"
)

// TODO(dape): add smaller tests and refactor all to use some shared functions.

func newRubberduck(t *testing.T, args []string, env []string) *exec.Cmd {
	bin := rubberducktesting.FindBinary(t)

	if args == nil {
		args = []string{}
	}
	
	cmd := exec.Command(bin, args...)
	cmd.Env = env
	return cmd
}

func startRubberduck(t *testing.T, args[] string, env []string) (io.Reader, io.Reader) {
	duck := newRubberduck(t, args, env)
	stdout, err := duck.StdoutPipe()
	if err != nil {
		t.Fatal("stdout failed", err)
	}

	stderr, err := duck.StderrPipe()
	if err != nil {
		t.Fatal("stderr failed", err)
	}

	err = duck.Start()
	if err != nil {
		t.Fatal(err)
	}
	
	return stdout, stderr
}

func TestArgs(t *testing.T) {
	args := []string{"alpha", "beta", "gamma"}
	stdout, _ := startRubberduck(t, args, nil)
	
	header, err := header.Read(stdout)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(header.Args)

	if !reflect.DeepEqual(header.Args, args) {
		t.Fatal()
	}
}

func TestEnv(t *testing.T) {
	env := []string{"a=1", "b=2"}
	stdout, _ := startRubberduck(t, nil, env)

	header, err := header.Read(stdout)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(header.Env)

	if !reflect.DeepEqual(header.Env, env) {
		t.Fatal()
	}
}

func checkOutput(t *testing.T, env []string, goldenFile string, selectPipe func(io.Reader, io.Reader) io.Reader) {
	stdout, stderr := startRubberduck(t, nil, env)
	reader := selectPipe(stdout, stderr)

	header, err := header.Read(stdout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(header)

	buf := &strings.Builder{}
	_, err = io.Copy(buf, reader)
	t.Log(buf)

	tmp, err := os.ReadFile(goldenFile)
	want := string(tmp)
	t.Log(want)
	if err != nil || buf.String() != want {
		t.Fatal(goldenFile, err, buf.String())
	}
}

func TestStdout(t *testing.T) {
	checkOutput(t, []string{"RUBBERDUCK_STDOUT=out.txt"}, "out.txt", func(stdout io.Reader, stderr io.Reader) io.Reader {
		return stdout
	})
}

func TestStderr(t *testing.T) {
	checkOutput(t, []string{"RUBBERDUCK_STDERR=err.txt"}, "err.txt", func(stdout io.Reader, stderr io.Reader) io.Reader {
		return stderr
	})
}
