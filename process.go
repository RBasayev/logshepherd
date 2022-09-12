package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

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

	//var byteline []byte
	var line string

	for {
		byteline, _, err := scanner.ReadLine()
		line = string(byteline)
		if err == io.EOF {
			// we shouldn't end up here
			fmt.Print(".")
			continue
		}

		if line != "" {
			// whatever comes in, goes into buffer
			lineBuffer = lineBuffer[1:]
			lineBuffer = append(lineBuffer, line)

			// by default the line goes into the log
			pass := true
			matchesShow := false
			matchesHide := false

			// but if there are filters in "show", the default is to not log
			if len(stream.Filters["show"]) > 0 {
				pass = false
				matchesShow = func() bool {
					for _, filter := range stream.Filters["show"] {
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
					for _, filter := range stream.Filters["hide"] {
						if strings.Contains(line, filter) {
							return true
						}
					}
					return false
				}()
			}

			pass = (pass || matchesShow) && !matchesHide

			// if the line matches one of the dump trigger patterns,
			// we're dumping the buffer into the log
			// (the last line of the dump is then the trigger)
			for _, dumpString := range stream.DumpUpon {
				if strings.Contains(line, dumpString) {

					dumpBlock := "** v BEGIN DUMP BEFORE TRIGGER v *********************************\n"
					for i := 0; i < stream.DumpBuffer; i++ {
						if lineBuffer[i] != "" {
							dumpBlock += "**  " + lineBuffer[i] + "\n"
						}
					}
					dumpBlock += "** ^ END DUMP BEFORE TRIGGER ^ ***********************************\n"

					outFile.WriteString(dumpBlock)
					// empty buffer
					lineBuffer = make([]string, stream.DumpBuffer)
					// TODO: dump the same amount of lines AFTER TRIGGER (important for Java)

					pass = false
					break
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
