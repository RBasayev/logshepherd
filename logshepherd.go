package main

import (
	"fmt"
	"time"
)

func main() {
	for _, input := range routes {

		fmt.Printf("Starting %s. Filters: %v\n", input.ID, input.Filters)

		if input.FullOutput {

			channels[input.ID] = make(chan []string, fullOutputBufferSize+30)

			// potential schemes could be "bolt" or "rrd" (round-robin-db)
			if fullOutputURL.Scheme == "file" {
				// launch full output listener goroutine
				go writeFullToFile(input.ID)
			}
		}
		go processStream(input)

		// It basically goes like this:
		// for every route, if full output needs to be saved,
		// a channel is created and a writeFullToFile goroutine
		// is launched and listens to the channel. Then the
		// processStream goroutine is being launched which
		// sends a slice of timestamp and the incoming log line
		// into the channel (only of full output is true).
	}

	for {
		d, _ := time.ParseDuration("300ms")
		time.Sleep(d)
	}

}
