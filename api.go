package main

import (
	"net"
	"net/http"
)

func ipRoute(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-IP") == "" {
		jsonRespond(w, 500, "X-IP is blank", nil)
		return
	}

	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(r.Header.Get("X-IP"))
	record, err := db.City(ip)
	if err != nil {
		serverLog("error parsing ip: %v", err)
		jsonRespond(w, 500, "error parsing ip", r.Header.Get("X-IP"))
		return
	}

	jsonRespond(w, 200, "success", record)
}

func stopRoute(w http.ResponseWriter, r *http.Request) {
	jsonRespond(w, 200, "success", nil)
	run = false
}
