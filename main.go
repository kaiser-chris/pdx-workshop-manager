package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"bahmut.de/pdx-workshop-manager/config"
	"bahmut.de/pdx-workshop-manager/logging"
	"bahmut.de/pdx-workshop-manager/manager"
	"bahmut.de/pdx-workshop-manager/steam"
)

var modId uint64 = 0
var configFile string

func parseArgs() {
	flag.CommandLine.Init("", flag.ExitOnError)
	flag.Uint64Var(&modId, "mod", 0, "Configured workshop mod id or 0 for all mods")
	flag.StringVar(&configFile, "config", "manager-config.json", "path to config file")
	flag.Parse()
}

func main() {
	parseArgs()

	applicationConfig, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Fatalf("Failed to load config: %v", err)
	}

	executablePath, err := os.Executable()
	if err != nil {
		logging.Fatalf("Could not get executable path: %v", err)
	}
	if os.WriteFile(
		filepath.Join(filepath.Dir(executablePath), "steam_appid.txt"),
		[]byte(strconv.FormatUint(uint64(applicationConfig.Game), 10)),
		0644,
	) != nil {
		logging.Fatalf("Failed to write to steam_appid.txt: %v", err)
	}

	if !steam.SteamAPI_Init() {
		fmt.Println("Failed to initialize steam api")
		os.Exit(1)
	}
	defer steam.SteamAPI_Shutdown()

	// 0 = all mods
	if modId > 0 {
		mod := applicationConfig.GetModByIdentifier(modId)
		if mod == nil {
			logging.Fatalf("Failed to find mod %d", modId)
			os.Exit(1)
		}
		err = manager.UploadMod(applicationConfig, mod)
		if err != nil {
			logging.Fatalf("Failed to upload mod %d: %v", modId, err)
		}
	} else {
		for _, mod := range applicationConfig.Mods {
			err = manager.UploadMod(applicationConfig, mod)
			if err != nil {
				logging.Fatalf("Failed to upload mod %d: %v", modId, err)
			}
		}
	}
}
