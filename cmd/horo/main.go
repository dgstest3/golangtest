package main

import (
	"flag"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	"github.com/Sirupsen/logrus"
	horo "testore.me/horo"
)

var (
	log    = logrus.New()
	debug  bool
	addr   string
	dbpath string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&addr, "addr", "localhost:8888", "Listen Sever Address")
	flag.StringVar(&dbpath, "db", "horo.db", "Path to DB")
}

func main() {
	flag.Parse()

	// Init logger
	if debug {
		log.Level = logrus.DebugLevel
	}

	// Setting app environment
	app := &horo.App{
		Log:    log,
		DBPath: dbpath,
	}

	// Creating and running server
	server := horo.NewServer(app)
	graceful.Run(addr, 5*time.Second, server)
}
