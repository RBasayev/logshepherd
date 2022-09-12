package main

import (
	"fmt"
	"log"
)

func isOK(err error) (e bool) {
	if err != nil {
		fmt.Printf("Non-fatal Error: %s\n", err)
		return false
	}
	return true
}

func bail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
