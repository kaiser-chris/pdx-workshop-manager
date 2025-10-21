package main

import (
	"flag"
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
	flag.Uint64Var(&modId, "mod", 0, "Configured workshop mod id or 0 for all mods (default 0)")
	flag.StringVar(&configFile, "config", "manager-config.json", "Path to the config file")
	flag.Parse()
}

func main() {
	parseArgs()

	logging.Infof("Loading configuration: %s", configFile)
	applicationConfig, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Fatalf("Failed to load config: %v", err)
	}

	logging.Info("Creating steam_appid.txt")
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

	logging.Info("Initializing Steam")
	if !steam.SteamAPI_Init() {
		logging.Fatal("Failed to initialize steam api")
	}
	defer steam.SteamAPI_Shutdown()

	// 0 = all mods
	if modId > 0 {
		logging.Infof("Start uploading mod %d", modId)
		mod := applicationConfig.GetModByIdentifier(modId)
		if mod == nil {
			logging.Fatalf("Failed to find mod %d", modId)
			os.Exit(1)
		}
		err = manager.UploadMod(applicationConfig, mod)
		if err != nil {
			logging.Fatalf("Failed to upload mod %d: %v", modId, err)
		}
		logging.Infof("Finished uploading mod: %d", modId)
	} else {
		logging.Info("Uploading all mods")
		for _, mod := range applicationConfig.Mods {
			if mod.Identifier == 0 {
				logging.Infof(" - Start uploading new mod")
			} else {
				logging.Infof(" - Start uploading mod: %d", mod.Identifier)
			}
			err = manager.UploadMod(applicationConfig, mod)
			if err != nil {
				logging.Fatalf("Failed to upload mod %d: %v", modId, err)
			}
			logging.Infof(" - Finished uploading mod: %d", mod.Identifier)
		}
		logging.Info("Finished uploading mods")
	}

}
