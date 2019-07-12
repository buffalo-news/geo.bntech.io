package main

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func httpServer() {
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/ip", ipRoute)
	http.HandleFunc("/stop", stopRoute)

	//var err error
	httpCloser, _ = listenAndServeWithClose(":3000", nil)
	serverLog("HTTP Server started\n")
}

type jsonResponse struct {
	Status  int
	Message string
	Body    interface{}
}

func jsonRespond(w http.ResponseWriter, status int, msg string, body interface{}) {
	response := jsonResponse{status, msg, body}
	settingsJSON, _ := json.Marshal(response)
	io.WriteString(w, string(settingsJSON))
}

func listenAndServeWithClose(addr string, handler http.Handler) (io.Closer, error) {

	var (
		listener  net.Listener
		srvCloser io.Closer
		err       error
	)

	srv := &http.Server{Addr: addr, Handler: handler}

	if addr == "" {
		addr = ":http"
	}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			serverLog("HTTP Server Error - %v\n", err)
		} else {
			serverLog("HTTP call\n")
		}
	}()

	srvCloser = listener
	return srvCloser, nil
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
