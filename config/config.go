package config

import (
	"encoding/json"
	"fmt"
	"os"

	"bahmut.de/pdx-workshop-manager/logging"
)

const (
	DefaultFileName = "manager-config.json"
)

type ApplicationConfig struct {
	configFilePath string
	Game           uint         `json:"game"`
	Mods           []*ModConfig `json:"mods"`
}

type ModConfig struct {
	Identifier          uint64 `json:"id"`
	Directory           string `json:"directory"`
	Description         string `json:"description"`
	ChangeNoteDirectory string `json:"change-note-directory"`
}

func LoadConfig(path string) (*ApplicationConfig, error) {
	// Read config
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logging.Fatal(err)
		}
	}(file)

	// Decode json
	decoder := json.NewDecoder(file)
	var config ApplicationConfig
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config.configFilePath = path

	return &config, nil
}

func InitializeConfig(configFilePath string, game uint) (*ApplicationConfig, error) {
	config := &ApplicationConfig{
		configFilePath: configFilePath,
		Game:           game,
		Mods:           make([]*ModConfig, 0),
	}
	err := config.Save()
	if err != nil {
		return nil, fmt.Errorf("failed to save created config file: %w", err)
	}
	return config, nil
}

func (config *ApplicationConfig) GetModByIdentifier(identifier uint64) *ModConfig {
	for _, mod := range config.Mods {
		if mod.Identifier == identifier {
			return mod
		}
	}
	return nil
}

func (config *ApplicationConfig) Save() error {
	content, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	err = os.WriteFile(config.configFilePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
