package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joshk/hustle"
)

var (
	config = &hustle.Config{
		HTTPAddr:    ":8661",
		HTTPPubAddr: "localhost:8661",
		// HTTPSAddr: ":8662",
		HubAddr:   ":6379",
		WSAddr:    ":8663",
		WSPubAddr: "localhost:8663",
		// WSSAddr:   ":8664",
		StatsAddr:    ":8665",
		StatsPubAddr: "localhost:8665",
	}

	versionFlag     = flag.Bool("version", false, "Print version and exit")
	revisionFlag    = flag.Bool("revision", false, "Print revision and exit")
	versionPlusFlag = flag.Bool("version+", false,
		"Print version and revision and exit")
)

func init() {
	flag.StringVar(&config.HTTPAddr, "http-addr", config.HTTPAddr,
		"HTTP Server address")
	flag.StringVar(&config.HTTPPubAddr, "http-public-addr", config.HTTPAddr,
		"HTTP Public server address (reachable from distant lands)")
	flag.StringVar(&config.HubAddr, "hub-addr", config.HubAddr,
		"Redis Hub address")
	flag.StringVar(&config.WSAddr, "ws-addr", config.WSAddr,
		"WS Server address")
	flag.StringVar(&config.WSPubAddr, "ws-public-addr", config.WSPubAddr,
		"WS Public server address (reachable from distant lands)")
	flag.StringVar(&config.StatsAddr, "stats-addr", config.StatsAddr,
		"Stats Server address")
	flag.StringVar(&config.StatsPubAddr, "stats-public-addr",
		config.StatsPubAddr,
		"Stats Public Server address (reachable from distant lands)")
}

func main() {
	err := hustle.ProcessConfig(config)
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

	defer func() {
		err := recover()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go hustle.HTTPServerMain(config)
	go hustle.WSServerMain(config)
	go hustle.StatsServerMain(config)

	<-quit
}
