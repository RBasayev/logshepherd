package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var timeUp = make(chan bool)

func shaltThouPass(line string, show []string, hide []string) bool {
	// by default the line goes into the log
	pass := true
	matchesShow := false
	matchesHide := false

	// but if there are filters in "show", the default is to not log
	if len(show) > 0 {
		pass = false
		matchesShow = func() bool {
			for _, filter := range show {
				if strings.Contains(line, filter) {
					return true
				}
			}
			return false
		}()
	}

	// "show" has priority over "hide"
	if !matchesShow {
		matchesHide = func() bool {
			for _, filter := range hide {
				if strings.Contains(line, filter) {
					return true
				}
			}
			return false
		}()
	}

	pass = (pass || matchesShow) && !matchesHide
	return pass
}

func dumpBefore(outFile *os.File, buffer *[]string, bufferSize int, dumpTimeout int, dumpCountdown *int) {
	dumpBlock := "** v BEGIN DUMP BEFORE TRIGGER v *********************************\n"
	l := len(*buffer) - 1
	for i := 0; i < bufferSize; i++ {
		if i == l {
			// last line of the buffer is the trigger
			break
		}
		if (*buffer)[i] != "" {
			dumpBlock += "**  " + (*buffer)[i] + "\n"
		}
	}
	dumpBlock += "** ^  END DUMP BEFORE TRIGGER  ^ *********************************\n"

	outFile.WriteString(dumpBlock)

	// empty buffer, that's why we need to pass this slice by reference
	*buffer = make([]string, bufferSize)
	*dumpCountdown = bufferSize + 1
	if dumpTimeout > 0 {
		timer := time.NewTimer(time.Duration(dumpTimeout) * time.Second)
		go func() {
			<-timer.C
			timeUp <- true
		}()
	}
}

func dumpAfter(outFile *os.File, line string, dumpCountdown *int, bufferSize int) bool {
	if *dumpCountdown == (bufferSize + 1) {
		outFile.WriteString("** v BEGIN DUMP AFTER TRIGGER v **********************************\n")
	}

	select {
	case <-timeUp:
		// quickly checking if time's up ( timer started in dumpBefore() )
		if *dumpCountdown > 0 {
			outFile.WriteString("**  -- not waiting for more lines (timeout) --")
			outFile.WriteString("** ^  END DUMP AFTER TRIGGER  ^ **********************************\n")
			*dumpCountdown = 0
		}
	default:
	}

	if *dumpCountdown > 0 {
		outFile.WriteString("**  " + line + "\n")

		if *dumpCountdown == 1 {
			outFile.WriteString("** ^  END DUMP AFTER TRIGGER  ^ **********************************\n")
		}
		*dumpCountdown--

		return true
	}

	// The time's up already (we don't know when it's up - 10 miliseconds or 20 minutes ago)
	// or no AFTER dump is due. We didn't output the line, so the line still needs
	// to be matched to filters and considered for logging, that's why we
	// return FALSE here - the calling function will handle this.
	return false

}

func processStream(stream inputDef) {
	// open input and output and prepare for work
	inFile, err := pipeOpenOrCreate(stream.Input)
	if !isOK(err) {
		fmt.Println("Could not open SOURCE PIPE. Ending thread for " + stream.ID)
		return
	}
	outFile, err := fileOpenOrCreate(stream.Output)
	if !isOK(err) {
		fmt.Println("Could not open FILTERED LOG file. Ending thread for " + stream.ID)
		return
	}
	defer inFile.Close()
	defer outFile.Close()

	lineBuffer := make([]string, stream.DumpBuffer)

	var scanner = bufio.NewReader(inFile)

	var line string
	var dump_countdown int = 0
	if stream.DumpTimeout < 1 {
		stream.DumpTimeout = 60
	}

	for {
		// ReadLine() will block the loop until a line comes in
		byteline, _, err := scanner.ReadLine()
		line = string(byteline)
		if err == io.EOF {
			// we shouldn't end up here
			fmt.Print(".")
			continue
		}

		if line != "" {

			pass := shaltThouPass(line, stream.Filters["show"], stream.Filters["hide"])
			// dump_countdown > 0 means that we're dumping the AFTER trigger part,
			// everything is being written into the filtered log.
			if dump_countdown < 1 {
				// not in AFTER trigger mode
				// whatever comes in, goes into buffer
				lineBuffer = lineBuffer[1:]
				lineBuffer = append(lineBuffer, line)

				// if the line matches one of the dump trigger patterns,
				// we're dumping the buffer into the log
				for _, dumpString := range stream.DumpUpon {
					if strings.Contains(line, dumpString) {
						dumpBefore(outFile, &lineBuffer, stream.DumpBuffer, stream.DumpTimeout, &dump_countdown)
						// if the dump has been triggered, we want to print the line anyway,
						// i.e., the dump trigger is independent of the "show" and "hide" filters
						pass = true
						line = "** TRIGGER: ** " + line
						break
					}
				}
			} else {
				// in AFTER trigger mode
				line_written := dumpAfter(outFile, line, &dump_countdown, stream.DumpBuffer)
				if line_written {
					// the line has been written in dumpAfter(), no need to write again
					pass = false
				}
			}

			if pass {
				outFile.WriteString(line + "\n")
			}

			if stream.FullOutput {
				ts := strconv.FormatInt(time.Now().UnixNano(), 10)
				msg := line
				channels[stream.ID] <- []string{ts, msg}
			}
		}
		outFile, err = considerRotating(outFile, stream.RotateAt)
		if !isOK(err) {
			fmt.Println("Could not open FILTERED LOG after ROTATION. Ending thread for " + stream.ID)
			return
		}
	}
}
