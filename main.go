package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"text/template"
)

const Version = "0.2.0"

func usageAndExit(s string, code int) {
	if s != "" {
		fmt.Fprintf(os.Stderr, "!! %s\n", s)
	}
	fmt.Fprint(os.Stderr, "usage: envtpl [template]\n")
	fmt.Fprintf(os.Stderr, "%s version: %s (%s on %s/%s; %s)\n", os.Args[0], Version, runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler)
	os.Exit(code)
}

func init() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usageAndExit("missing [template]", 1)
	}
	input := args[0]
	if input == "--help" {
		usageAndExit("", 0)
	}
	if len(input) < 4 || input[len(input)-4:] != ".tpl" {
		usageAndExit("[template] does not end with .tpl", 1)
	}
	t, err := template.ParseFiles(input)
	if err != nil {
		usageAndExit(err.Error(), 1)
	}
	var b bytes.Buffer
	err = t.Option("missingkey=zero").Execute(&b, getEnvironMap())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	ioutil.WriteFile(input[:len(input)-4], b.Bytes(), 0644)
	os.Remove(input)
}

func getEnvironMap() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		x := strings.SplitN(e, "=", 2)
		env[x[0]] = x[1]
	}
	return env
}
