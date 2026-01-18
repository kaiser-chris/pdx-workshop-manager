package cmd

import (
	"fmt"

	"bahmut.de/pdx-workshop-manager/config"
	"bahmut.de/pdx-workshop-manager/logging"
	"bahmut.de/pdx-workshop-manager/manager"
	"bahmut.de/pdx-workshop-manager/steam"
)

const AllMods uint64 = 0

func Run(configFile string, modId uint64) error {
	logging.Infof("Loading configuration: %s", configFile)
	applicationConfig, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Errorf("Failed to load config: %v", err)
		return err
	}

	logging.Info("Initializing Steam")
	err = manager.Init(applicationConfig)
	defer steam.SteamAPI_Shutdown()
	if err != nil {
		logging.Errorf("Failed to initialize steam: %v", err)
		return err
	}

	if modId > AllMods {
		logging.Infof("Start uploading mod %d", modId)
		mod := applicationConfig.GetModByIdentifier(modId)
		if mod == nil {
			logging.Errorf("Failed to find mod %d", modId)
			return fmt.Errorf("failed to find mod %d", modId)
		}
		err = manager.UploadMod(applicationConfig, mod)
		if err != nil {
			logging.Errorf("Failed to upload mod %d: %v", modId, err)
			return err
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
				logging.Errorf("Failed to upload mod %d: %v", modId, err)
				return err
			}
			logging.Infof(" - Finished uploading mod: %d", mod.Identifier)
		}
		logging.Info("Finished uploading mods")
	}

	return nil
}
