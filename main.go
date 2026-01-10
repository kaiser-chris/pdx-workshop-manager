package main

import (
	"flag"

	"bahmut.de/pdx-workshop-manager/cmd"
	"bahmut.de/pdx-workshop-manager/config"
	"bahmut.de/pdx-workshop-manager/gui"
	"bahmut.de/pdx-workshop-manager/logging"
)

var modId uint64 = 0
var configFile string

func parseArgs() int {
	flag.CommandLine.Init("", flag.ExitOnError)
	flag.Uint64Var(&modId, "mod", cmd.AllMods, "Configured workshop mod id or 0 for all mods (default 0)")
	flag.StringVar(&configFile, "config", config.DefaultFileName, "Path to the config file")
	flag.Parse()
	return len(flag.Args())
}

func main() {
	if parseArgs() == 0 {
		err := gui.Run()
		if err != nil {
			logging.Errorf("Critical Error: %v", err)
		}
	} else {
		err := cmd.Run(configFile, modId)
		if err != nil {
			logging.Errorf("Error: %v", err)
		} else {
			logging.Infof("Upload successful")
		}
	}
}
