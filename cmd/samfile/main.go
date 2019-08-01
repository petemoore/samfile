package main

import (
	"log"

	docopt "github.com/docopt/docopt-go"
)

var (
	version  string // set during build with `-ldflags "-X main.version=$(git tag -l 'v*.*.*' --points-at HEAD | head -n1)"`
	revision string // set during build with `-ldflags "-X main.revision=$(git rev-parse HEAD)"`
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("samfile: ")
	versionName := "samfile"
	if version != "" {
		versionName += " " + version
	}
	if revision != "" {
		versionName += " [ revision: https://github.com/winfreddy88/samfile/commits/" + revision + " ]"
	}
	arguments, err := docopt.Parse(usage(versionName), nil, true, versionName, false, true)
	if err != nil {
		log.Fatalf("Error parsing command line arguments: %v", err)
	}
	switch {
	case arguments["extract"]:
		extract(arguments)
	case arguments["ls"]:
		ls(arguments)
	default:
		log.Fatal("Could not find a command to run")
	}
}
