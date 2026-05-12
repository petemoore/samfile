// Command samfile manipulates files inside a SAM Coupé MGT floppy
// disk image: listing the directory (ls), extracting one or all
// files (cat / extract), adding a new code file (add), and
// detokenising a saved SAM BASIC program to plain text
// (basic-to-text). Run `samfile --help` for invocation details. For
// programmatic access to MGT images, import the parent package
// github.com/petemoore/samfile/v3.
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
		versionName += " [ revision: https://github.com/petemoore/samfile/commits/" + revision + " ]"
	}
	arguments, err := docopt.ParseArgs(usage(versionName), nil, versionName)
	if err != nil {
		log.Fatalf("error parsing command line arguments: %v", err)
	}
	switch {
	case arguments["cat"]:
		cat(arguments)
	case arguments["extract"]:
		extract(arguments)
	case arguments["ls"]:
		ls(arguments)
	case arguments["basic-to-text"]:
		basicToText(arguments)
	case arguments["add"]:
		add(arguments)
	case arguments["verify"]:
		format := "text"
		if v, ok := arguments["--format"].(string); ok && v != "" {
			format = v
		}
		if err := runVerify(arguments["-i"].(string), format); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("could not find a command to run")
	}
}
