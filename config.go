package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func readConfig(configPath string) (config logshepherdConf) {
	var in []byte

	in, err := os.ReadFile(configPath)
	bail(err)
	yaml.Unmarshal(in, &config)
	bail(err)
	return
}

// func reloadConfig(configPath string) (reloaded bool) {
// 	return true
// }
