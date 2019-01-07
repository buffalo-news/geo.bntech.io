package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	geoip2 "github.com/oschwald/geoip2-golang"
)

var sigchan = make(chan os.Signal, 1)

var run = true
var db *geoip2.Reader

var httpCloser io.Closer

func main() {
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	//updateDB()
	var err error
	db, err = geoip2.Open("maxmind/GeoLite2-City.mmdb")
	if err != nil {
		serverLog("error parsing db")
		run = false
		return
	}
	defer db.Close()

	go httpServer()

	for run {
		select {
		case sig := <-sigchan:
			serverLog(" Caught signal %v: terminating\n", sig)
			run = false
		}
	}

	httpCloser.Close()
}

func serverLog(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}
