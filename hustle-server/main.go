package main

import (
	"flag"
	"fmt"
	"github.com/joshk/hustle"
	"os"
)

var (
	versionFlag     = flag.Bool("version", false, "Print version and exit")
	revisionFlag    = flag.Bool("revision", false, "Print revision and exit")
	versionPlusFlag = flag.Bool("version+", false, "Print version and revision and exit")

	addrFlag = flag.String("a", os.Getenv("HUSTLE_ADDR"), "Server address [default ':8966']")
)

func main() {
	flag.Parse()

	if *addrFlag == "" {
		*addrFlag = ":8966"
	}

	if *versionFlag {
		fmt.Println(hustle.VersionString)
		os.Exit(0)
	}

	if *revisionFlag {
		fmt.Println(hustle.RevisionString)
		os.Exit(0)
	}

	if *versionPlusFlag {
		fmt.Println(hustle.VersionPlusJSON)
		os.Exit(0)
	}

	hustle.Serve(*addrFlag)
}
