package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joshk/hustle"
	"github.com/kelseyhightower/envconfig"
)

type configSpec struct {
	HTTPAddr string
	// HTTPSAddr string
	WSAddr string
	// WSSAddr   string
	StatsAddr string
}

var (
	config = &configSpec{
		HTTPAddr: ":8661",
		// HTTPSAddr: ":8662",
		WSAddr: ":8663",
		// WSSAddr:   ":8664",
		StatsAddr: ":8665",
	}

	versionFlag     = flag.Bool("version", false, "Print version and exit")
	revisionFlag    = flag.Bool("revision", false, "Print revision and exit")
	versionPlusFlag = flag.Bool("version+", false, "Print version and revision and exit")
)

func init() {
	flag.StringVar(&config.HTTPAddr, "http-addr", config.HTTPAddr, "HTTP Server address")
	flag.StringVar(&config.WSAddr, "ws-addr", config.WSAddr, "WS Server address")
	flag.StringVar(&config.StatsAddr, "stats-addr", config.StatsAddr, "Stats Server address")
}

func main() {
	err := envconfig.Process("hustle", config)
	if err != nil {
		log.Fatal(err.Error())
	}

	flag.Parse()

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

	quit := make(chan bool)

	go hustle.HTTPServerMain(config.HTTPAddr)
	go hustle.WSServerMain(config.WSAddr)
	go hustle.StatsServerMain(config.StatsAddr)

	<-quit
}
