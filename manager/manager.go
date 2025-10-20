package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"bahmut.de/pdx-workshop-manager/config"
	"bahmut.de/pdx-workshop-manager/logging"
	"bahmut.de/pdx-workshop-manager/steam"
)

type ModUploadData struct {
	Game        uint
	Description string
	ChangeNote  string
	Thumbnail   string
	Metadata    *ModMetadata
	Config      *config.ModConfig
}

type ModMetadata struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Tags    []string `json:"tags"`
}

func UploadMod(appConfig *config.ApplicationConfig, modConfig *config.ModConfig) error {
	data, err := createModUploadData(modConfig, appConfig.Game)
	if err != nil {
		return err
	}

	if modConfig.Identifier == 0 {
		identifier, err := createMod(appConfig.Game)
		if err != nil {
			return err
		}
		data.Config.Identifier = identifier
	}

	err = appConfig.Save()
	if err != nil {
		return err
	}

	err = uploadModData(data)
	if err != nil {
		return err
	}
	return nil
}

func createModUploadData(config *config.ModConfig, game uint) (*ModUploadData, error) {
	uploadData := &ModUploadData{}
	uploadData.Config = config
	uploadData.Game = game

	// Read metadata file
	metadataPath := filepath.Join(config.Directory, ".metadata", "metadata.json")
	metadataFile, err := os.Open(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer func(metadataFile *os.File) {
		err := metadataFile.Close()
		if err != nil {
			logging.Fatal(err)
		}
	}(metadataFile)

	// Decode metadata json
	decoder := json.NewDecoder(metadataFile)
	var metadata ModMetadata
	err = decoder.Decode(&metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata file: %w", err)
	}

	uploadData.Metadata = &metadata
	uploadData.Thumbnail = filepath.Join(config.Directory, "thumbnail.png")

	if config.Description != "" {
		content, err := os.ReadFile(config.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to read description file %s: %w", config.Description, err)
		}
		uploadData.Description = string(content)
	}

	if config.ChangeNoteDirectory != "" {
		changeNotePath := filepath.Join(config.ChangeNoteDirectory, metadata.Version+".bbcode")
		content, err := os.ReadFile(changeNotePath)
		if err != nil {
			logging.Warnf("failed to read changeNote file %s: %v", changeNotePath, err)
		}
		uploadData.ChangeNote = string(content)
	}

	return uploadData, nil
}

func createMod(game uint) (uint64, error) {
	var steamError = false

	apiCall := steam.SteamUGC().CreateItem(
		game,
		steam.K_EWorkshopFileTypeCommunity,
	)

	result := steam.NewCreateItemResult_t()

	for true {
		if steam.SteamUtils().IsAPICallCompleted(apiCall, &steamError) {
			steam.SteamUtils().GetAPICallResult(
				apiCall,
				result.Swigcptr(),
				steam.Sizeof_CreateItemResult_t,
				steam.CreateItemResult_tK_iCallback,
				&steamError,
			)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if steamError {
		return 0, fmt.Errorf("steam API call failed: %s", steam.SteamUtils().GetAPICallFailureReason(apiCall))
	}

	if result.GetM_eResult() != 1 {
		return 0, fmt.Errorf("item creation failed: %v", result.GetM_eResult())
	}

	if result.GetM_bUserNeedsToAcceptWorkshopLegalAgreement() {
		return 0, errors.New("to make your item public you need to agree to the workshop terms of service <http://steamcommunity.com/sharedfiles/workshoplegalagreement>")
	}

	return result.GetM_nPublishedFileId(), nil
}

func uploadModData(data *ModUploadData) error {
	var steamError = false

	handle := steam.SteamUGC().StartItemUpdate(data.Game, data.Config.Identifier)

	steam.SteamUGC().SetItemTitle(handle, data.Metadata.Name)

	contentPath, err := filepath.Abs(data.Config.Directory)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute content path: %w", err)
	}
	steam.SteamUGC().SetItemContent(handle, contentPath)

	thumbnailPath, err := filepath.Abs(data.Thumbnail)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute thumbnail path: %w", err)
	}
	steam.SteamUGC().SetItemPreview(handle, thumbnailPath)

	if len(data.Metadata.Tags) > 0 {
		tagArray := steam.NewSteamParamStringArray(data.Metadata.Tags)
		steam.SteamUGC().SetItemTagsExtension(handle, tagArray)
	}

	if data.Description != "" {
		steam.SteamUGC().SetItemDescription(handle, data.Description)
	}

	updateResult := steam.NewSubmitItemUpdateResult_t()

	apiCall := steam.SteamUGC().SubmitItemUpdate(handle, data.ChangeNote)
	for {
		if steam.SteamUtils().IsAPICallCompleted(apiCall, &steamError) {
			steam.SteamUtils().GetAPICallResult(
				apiCall,
				updateResult.Swigcptr(),
				steam.Sizeof_SubmitItemUpdateResult_t,
				steam.SubmitItemUpdateResult_tK_iCallback,
				&steamError,
			)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if updateResult.GetM_eResult() != steam.K_EResultOK {
		return fmt.Errorf("steam API call failed: %v", updateResult.GetM_eResult())
	}

	if steamError {
		return fmt.Errorf("steam API call failed: %s", steam.SteamUtils().GetAPICallFailureReason(apiCall))
	}

	return nil
}
