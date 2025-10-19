package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"bahmut.de/pdx-workshop-manager/steam"
	"bahmut.de/pdx-workshop-manager/util"
)

var fAppId uint32 = 480
var fPublishedFileID uint64 = 0
var fItemTitle string = ""
var fItemDescription string = ""
var fItemContent string = ""
var fItemPreview string = ""
var fChangeNote string = ""

func parseArgs() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s: -c dir [options]\n", os.Args[0])
		util.PrintDefaults()
	}

	flag.CommandLine.Init("", flag.ExitOnError)

	flag.Var((*util.Uint32Flag)(&fAppId), "a", "0:`AppID` of the game")
	flag.Uint64Var(&fPublishedFileID, "i", 0, "1:Existing workshop item `id` or 0 for new item")
	flag.StringVar(&fItemTitle, "t", "", "2:Item `title`")
	flag.StringVar(&fItemDescription, "d", "", "3:Item `description`")
	flag.StringVar(&fItemContent, "c", "", "4:Path to content `dir` for upload")
	flag.StringVar(&fItemPreview, "p", "", "5:Path to preview `image` for upload")
	flag.StringVar(&fChangeNote, "n", "", "6:Change `note`")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func createItem() {
	var steamError bool = false

	fmt.Println("Creating new item...")

	hSteamAPICall := steam.SteamUGC().CreateItem(
		uint(fAppId),
		steam.K_EWorkshopFileTypeCommunity,
	)

	CreateItemResult := steam.NewCreateItemResult_t()

	for true {
		if steam.SteamUtils().IsAPICallCompleted(hSteamAPICall, &steamError) {
			steam.SteamUtils().GetAPICallResult(
				hSteamAPICall,
				CreateItemResult.Swigcptr(),
				steam.Sizeof_CreateItemResult_t,
				steam.CreateItemResult_tK_iCallback,
				&steamError,
			)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if steamError {
		fmt.Println("Steam api call CreateItem() failed:", steam.SteamUtils().GetAPICallFailureReason(hSteamAPICall))
		os.Exit(1)
	}

	if CreateItemResult.GetM_eResult() != 1 {
		fmt.Println("Item creation failed:", CreateItemResult.GetM_eResult())
		os.Exit(1)
	}

	if CreateItemResult.GetM_bUserNeedsToAcceptWorkshopLegalAgreement() {
		fmt.Println("To make your item public you need to agree to the workshop terms of service <http://steamcommunity.com/sharedfiles/workshoplegalagreement>")
	}

	fmt.Println("Acquired new workshop item id:", CreateItemResult.GetM_nPublishedFileId())
	fPublishedFileID = CreateItemResult.GetM_nPublishedFileId()
	updateItem()
}

func updateItem() {
	var steamError bool = false

	fmt.Println("Updating item...")

	UGCUpdateHandle := steam.SteamUGC().StartItemUpdate(uint(fAppId), fPublishedFileID)

	if len(fItemTitle) > 0 {
		steam.SteamUGC().SetItemTitle(UGCUpdateHandle, fItemTitle)
	}

	if len(fItemDescription) > 0 {
		steam.SteamUGC().SetItemDescription(UGCUpdateHandle, fItemDescription)
	}

	if len(fItemContent) > 0 {
		steam.SteamUGC().SetItemContent(UGCUpdateHandle, util.ArgSlice(filepath.Abs(fItemContent))[0].(string))
	}

	if len(fItemPreview) > 0 {
		steam.SteamUGC().SetItemPreview(UGCUpdateHandle, util.ArgSlice(filepath.Abs(fItemPreview))[0].(string))
	}

	hSteamAPICall := steam.SteamUGC().SubmitItemUpdate(UGCUpdateHandle, fChangeNote)

	SubmitItemUpdateResult := steam.NewSubmitItemUpdateResult_t()

	for true {
		if steam.SteamUtils().IsAPICallCompleted(hSteamAPICall, &steamError) {
			steam.SteamUtils().GetAPICallResult(
				hSteamAPICall,
				SubmitItemUpdateResult.Swigcptr(),
				steam.Sizeof_SubmitItemUpdateResult_t,
				steam.SubmitItemUpdateResult_tK_iCallback,
				&steamError,
			)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if steamError {
		fmt.Println("Steam api call SubmitItemUpdate() failed:", steam.SteamUtils().GetAPICallFailureReason(hSteamAPICall))
		os.Exit(1)
	}

	if SubmitItemUpdateResult.GetM_eResult() != 1 {
		fmt.Println("Item update failed:", SubmitItemUpdateResult.GetM_eResult())
		os.Exit(1)
	}

	fmt.Println("Update complete")
}

func main() {
	parseArgs()

	ep, err := os.Executable()
	if err != nil {
		fmt.Println("Could not get executable path")
		os.Exit(1)
	}
	if ioutil.WriteFile(
		filepath.Join(filepath.Dir(ep), "steam_appid.txt"),
		[]byte(strconv.FormatUint(uint64(fAppId), 10)),
		0644,
	) != nil {
		fmt.Println("Failed to write to steam_appid.txt")
	}

	if !steam.SteamAPI_Init() {
		fmt.Println("Failed to initialize steam api")
		os.Exit(1)
	}
	defer steam.SteamAPI_Shutdown()

	if fPublishedFileID > 0 {
		updateItem()
	} else {
		createItem()
	}
}
