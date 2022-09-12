package main

import (
	"fmt"
	"os"
	"time"
)

func writeFullToFile(id string) {
	fullFilePath := fullOutputURL.Path + fmt.Sprintf("%c", os.PathSeparator) + id + "-full.output"
	fullFile, err := fileOpenOrCreate(fullFilePath)
	if !isOK(err) {
		fmt.Println("Could not open FULL LOG file. Ending thread for " + id)
		return
	}
	// in essence - polling for updates in the channel every 50ms
	dur, _ := time.ParseDuration("50ms")
	var shuttle []string
	var contents string
	var countDown int = fullOutputBufferSize
	var startTime int64

	for {
		select {
		case shuttle = <-channels[id]:
			countDown--
			if startTime == int64(0) {
				startTime = time.Now().Unix()
			}
			contents += shuttle[0] + " | " + shuttle[1] + "\n"
			if countDown < 1 {
				fullFile.WriteString(contents)
				fullFile.Sync()
				countDown = fullOutputBufferSize
				startTime = int64(0)
				contents = ""
			}
		default:
			if (startTime != int64(0)) && ((time.Now().Unix() + fullOutputBufferTimeout) > startTime) {
				fullFile.WriteString(contents)
				fullFile.Sync()
				countDown = fullOutputBufferSize
				startTime = int64(0)
				contents = ""
			}
			time.Sleep(dur)
		}
		// rotating regardless of current activity
		// may be moved to 'default:' above - to rotate when idle
		// TODO: makes sense to compress the rotated file
		// go zipRotatedFile(fullFilePath+".rotated."+timestamp)
		fullFile, err = considerRotating(fullFile, fullOutputRotateAt)
		if !isOK(err) {
			fmt.Println("Could not open FILTERED LOG after ROTATION. Ending thread for " + id)
			return
		}

	}

}
