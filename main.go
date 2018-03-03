package main

import "flag"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"text/template"
)

const Version = "0.3.0"

func usageAndExit(s string, code int) {
	if s != "" {
		fmt.Fprintf(os.Stderr, "!! %s\n", s)
	}
	fmt.Fprint(os.Stderr, "usage: envtpl [options] [template]\n")
	fmt.Fprint(os.Stderr, "options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "%s version: %s (%s on %s/%s; %s)\n", os.Args[0], Version, runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler)
	os.Exit(code)
}

func init() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func main() {
	keepTemplate := flag.Bool("keep-template", false, "Keep the template file instead of deleting it")
	showHelp := flag.Bool("help", false, "Show this help")
	flag.Parse()

	if *showHelp {
		usageAndExit("", 0)
	}

	files := flag.Args()

	if len(files) == 0 {
		usageAndExit("missing [template]", 1)
	}

	checkFilenames(files)

	for _, file := range files {
		renderTemplate(file, *keepTemplate)
	}
}

func checkFilenames(files []string) {
	for _, file := range files {
		if len(file) < 4 || file[len(file)-4:] != ".tpl" {
			usageAndExit(fmt.Sprintf("%q does not end with .tpl", file), 1)
		}
	}
}

func renderTemplate(input string, keepTemplate bool) {
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
	if !keepTemplate {
		os.Remove(input)
	}
}

func getEnvironMap() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		x := strings.SplitN(e, "=", 2)
		env[x[0]] = x[1]
	}
	return env
}
