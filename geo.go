package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	geoip2 "github.com/oschwald/geoip2-golang"
)

// Defining the settings structure
type settings struct {
	License string `json:license,omitempty`
}

// Settings for the application are global so we can access them from anywhere
var Settings settings

var sigchan = make(chan os.Signal, 1)

var run = true
var db *geoip2.Reader

var httpCloser io.Closer

func main() {
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize anything the API needs
	initialize()

	// Update the db to the newest copy
	updateDB()
	var err error
	db, err = geoip2.Open("maxmind/GeoLite2-City.mmdb")
	if err != nil {
		serverLog("error parsing db\n")
		os.Exit(1)
	}
	defer db.Close()

	// Start the API
	go httpServer()

	// Run until shutdown is signaled, then close
	for run {
		select {
		case sig := <-sigchan:
			serverLog("Caught signal %v: terminating\n", sig)
			run = false
		}
	}

	httpCloser.Close()
}

func initialize() {
	// Open our settingFile
	settingsFile, err := os.Open("./settings.json")

	// If we os.Open returns an error then handle it
	if err != nil {
		serverLog(err.Error())
		os.Exit(1)
	}

	// Read the settings into the global settings
	settings, _ := ioutil.ReadAll(settingsFile)
	err = json.Unmarshal(settings, &Settings)

	// If there is an error in the json structure handle it
	if err != nil {
		serverLog(err.Error())
		os.Exit(1)
	}

	// Defer the closing of our settingFile so that we can parse it later on
	defer settingsFile.Close()
}

func serverLog(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}
