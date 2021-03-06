package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	hustleServer "github.com/joshk/hustle/server"
)

var (
	envPort          = os.Getenv("PORT")
	envRedisAddr     = os.Getenv("REDIS_URL")
	envRedisToGoAddr = os.Getenv("REDISTOGO_URL")
	config           = &hustleServer.Config{
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
	runAllFlag   = flag.Bool("all", true, "Run all server types (overrides 'http', 'ws', and 'stats' flags)")
	runHTTPFlag  = flag.Bool("http", false, "Run the HTTP server (overridden by 'all' flag)")
	runWSFlag    = flag.Bool("ws", false, "Run the WS server (overridden by 'all' flag)")
	runStatsFlag = flag.Bool("stats", false, "Run the Stats server (overridden by 'all' flag)")
)

func init() {
	// Send everything to stdout to play nicely with heroku
	log.SetOutput(os.Stdout)

	if len(envPort) > 0 {
		config.HTTPAddr = ":" + envPort
	}

	if hubAddr, ok := getRedisURL(envRedisAddr); ok {
		config.HubAddr = hubAddr
	}

	if hubAddr, ok := getRedisURL(envRedisToGoAddr); ok {
		config.HubAddr = hubAddr
	}

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

func getRedisURL(addr string) (string, bool) {
	// "redis://redistogo:fooauth@grideye.redistogo.com:10328"
	parsed, err := url.Parse(addr)
	if err != nil {
		log.Println("Failed to parse REDIS_URL:", err)
		return "", false
	}

	if parsed.User == nil {
		return parsed.Host, true
	}

	if passwd, ok := parsed.User.Password(); ok {
		return fmt.Sprintf("%s:%s@%s",
			parsed.User.Username(),
			passwd,
			parsed.Host,
		), true
	}

	return "", false
}

func main() {
	err := hustleServer.ProcessConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	flag.Parse()

	if *versionFlag {
		fmt.Println(hustleServer.VersionString)
		os.Exit(0)
	}

	if *revisionFlag {
		fmt.Println(hustleServer.RevisionString)
		os.Exit(0)
	}

	if *versionPlusFlag {
		fmt.Println(hustleServer.VersionPlusJSON)
		os.Exit(0)
	}

	quit := make(chan bool)

	defer func() {
		err := recover()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("hustle starting up server(s)")

	if *runAllFlag || *runHTTPFlag {
		go hustleServer.HTTPServerMain(config)
	}

	if *runAllFlag || *runWSFlag {
		go hustleServer.WSServerMain(config)
	}

	if *runAllFlag || *runStatsFlag {
		go hustleServer.StatsServerMain(config)
	}

	<-quit
}
