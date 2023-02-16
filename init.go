package main

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strconv"
)

var config logshepherdConf
var channels map[string]chan []string
var settingsFile string
var routes []inputDef
var fullOutputURL *url.URL
var fullOutputBufferSize int = 10
var fullOutputBufferTimeout int64 = int64(60)
var fullOutputRotateAt int = 50
var Version string = "0.3.today"

func init() {
	// setting the "globals"
	settingsFile = "logshepherd.yaml"
	if len(os.Args) > 1 {
		settingsFile = os.Args[1]
	}

	ver := []string{"version", "-v", "-ver", "--ver", "-version", "--version"}
	for _, v := range ver {
		if settingsFile == v {
			fmt.Printf("logshepherd version:\n%s\n", Version)
			os.Exit(0)
		}
	}

	config = readConfig(settingsFile)
	// TODO: verify and/or set to default mandatory settings

	runtime.GOMAXPROCS(config.Threads)

	routes = config.Routes

	channels = make(map[string]chan []string, len(routes))

	var err error
	fullOutputURL, err = url.Parse(config.OutputFull["path"])
	bail(err)

	buffer, err := strconv.Atoi(config.OutputFull["write_buffer"])
	if isOK(err) {
		fullOutputBufferSize = buffer
	}

	timeout, err := strconv.Atoi(config.OutputFull["write_timeout"])
	if isOK(err) {
		fullOutputBufferTimeout = int64(timeout)
	}

	rotateCap, err := strconv.Atoi(config.OutputFull["cap"])
	if isOK(err) {
		fullOutputRotateAt = rotateCap
	}
}
