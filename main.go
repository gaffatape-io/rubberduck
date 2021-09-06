package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func log(tag string, data []string) {
	fmt.Println(tag + ":")
	for _, d := range data {
		fmt.Println("  ", d)
	}
	fmt.Println()
}

func logfile(w io.Writer, f string) {
	content, err := os.ReadFile(f)
	if err != nil {
		os.Exit(125)
	}
	fmt.Fprint(w, string(content))
}

func logfileenv(w io.Writer, envVar string) {
	f := os.Getenv(envVar)
	if f != "" {
		logfile(w, f)
	}
}

func exitenv(envVar string) {
	status, err := strconv.Atoi(os.Getenv(envVar))
	if err != nil {
		status = 0
	}
	os.Exit(status)
}

// rubberduck is a small command line utility intended for testing
// os/exec based code, it will log all arguments, environment vars
// and the content of any file found in the RUBBERDUCK_STDERR,
// RUBBERDUCK_STDOUT variables. The exit code is normally zero but
// can be controlled through the RUBBERDUCK_STATUS environment variable.
func main() {
	log("args", os.Args[1:])
	log("env", os.Environ())
	logfileenv(os.Stdout, "RUBBERDUCK_STDOUT")
	os.Stdout.Sync()
	logfileenv(os.Stderr, "RUBBERDUCK_STDERR")
	os.Stderr.Sync()

	exitenv("RUBBERDUCK_STATUS")
}
