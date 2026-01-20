//go:build gui

package web

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"

	"bahmut.de/pdx-workshop-manager/config"
	"bahmut.de/pdx-workshop-manager/logging"
	"bahmut.de/pdx-workshop-manager/manager"
	"bahmut.de/pdx-workshop-manager/steam"
	"bahmut.de/pdx-workshop-manager/web/resource/static"
	"github.com/pkg/browser"
)

var window *MainWindow
var games = []*Game{
	{Identifier: 529340, Name: "Victoria 3"},
	{Identifier: 3450310, Name: "Europa Universalis V"},
}

const (
	MessageSuccess = 0
	MessageWarning = 1
	MessageError   = 2
)

type Game struct {
	Identifier uint
	Name       string
}

type Message struct {
	Shown   bool
	Level   int
	Message string
}

type MainWindow struct {
	Message       *Message
	Games         []*Game
	Game          *Game
	Configuration *config.ApplicationConfig
	Mods          []*ModFrame
}

func (w *MainWindow) RefreshMods() {
	mods := make([]*ModFrame, len(w.Configuration.Mods))
	for i, mod := range w.Configuration.Mods {
		mods[i] = NewModFrame(w, mod)
	}
	w.Mods = mods
}

func (w *MainWindow) RefreshGame() {
	for _, game := range w.Games {
		if game.Identifier == w.Configuration.Game {
			w.Game = game
		}
	}
}

func (w *MainWindow) SendMessage(message string, level int) {
	w.Message = &Message{
		Shown:   false,
		Level:   level,
		Message: message,
	}
}

//go:embed resource/template/*
var templates embed.FS

func Run() {
	configuration := loadOrSetupConfig()
	window = &MainWindow{
		Games:         games,
		Configuration: configuration,
	}
	window.RefreshMods()
	window.RefreshGame()

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Embed))))
	http.HandleFunc("GET /", main)
	http.HandleFunc("GET /guide", guide)
	http.HandleFunc("GET /game/{identifier}", changeGame)
	http.HandleFunc("GET /mod/add", addMod)
	http.HandleFunc("POST /mod/update/{index}", updateMod)
	http.HandleFunc("GET /mod/remove/{index}", removeMod)
	http.HandleFunc("GET /mod/upload/{index}", uploadMod)

	port, err := getFreePort()
	if err != nil {
		logging.Fatalf("Could not get a free port: %v", err)
	}

	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	logging.Info("Starting Webserver")
	logging.Infof("URL: %s", logging.AnsiLink(url, url))

	err = browser.OpenURL(url)
	if err != nil {
		logging.Fatalf("Could not open browser: %v", err)
	}

	logging.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func loadOrSetupConfig() *config.ApplicationConfig {
	appConfig, err := config.LoadConfig(config.DefaultFileName)
	if err != nil {
		err = nil
		appConfig, err := config.InitializeConfig(config.DefaultFileName, 529340)
		if err != nil {
			logging.Fatalf("Could not create config file: %v", err)
		}
		return appConfig
	} else {
		return appConfig
	}
}

func main(writer http.ResponseWriter, _ *http.Request) {
	if window.Message != nil {
		if window.Message.Shown {
			window.Message = nil
		} else {
			window.Message.Shown = true
		}
	}
	page := template.Must(template.ParseFS(templates, "resource/template/main.html"))
	err := page.Execute(writer, window)
	if err != nil {
		logging.Fatalf("Could not execute template: %v", err)
	}
}

func guide(writer http.ResponseWriter, _ *http.Request) {
	page := template.Must(template.ParseFS(templates, "resource/template/guide.html"))
	err := page.Execute(writer, window)
	if err != nil {
		logging.Fatalf("Could not execute template: %v", err)
	}
}

func changeGame(writer http.ResponseWriter, request *http.Request) {
	gameParameter := request.PathValue("identifier")
	game, err := strconv.ParseUint(gameParameter, 10, 32)
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could parse game identifier: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	window.Configuration.Game = uint(game)
	err = window.Configuration.Save()
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not save configuration: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	window.RefreshGame()
	window.SendMessage("Changed game successfully", MessageSuccess)
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func addMod(writer http.ResponseWriter, request *http.Request) {
	mod := &config.ModConfig{
		Identifier:          0,
		Directory:           "",
		Description:         "",
		ChangeNoteDirectory: "",
	}

	window.Configuration.Mods = append(window.Configuration.Mods, mod)
	err := window.Configuration.Save()
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not save configuration: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	window.RefreshMods()
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func removeMod(writer http.ResponseWriter, request *http.Request) {
	indexParameter := request.PathValue("index")
	index, err := strconv.Atoi(indexParameter)
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not parse mod index: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	if index < 0 || index >= len(window.Configuration.Mods) {
		window.SendMessage(fmt.Sprintf("Could not remove mod: %v", errors.New("index out of bound")), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	window.Configuration.Mods = append(window.Configuration.Mods[:index], window.Configuration.Mods[index+1:]...)
	window.RefreshMods()

	err = window.Configuration.Save()
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not save configuration: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func uploadMod(writer http.ResponseWriter, request *http.Request) {
	indexParameter := request.PathValue("index")
	index, err := strconv.Atoi(indexParameter)
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not parse mod index: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	if index < 0 || index >= len(window.Configuration.Mods) {
		window.SendMessage(fmt.Sprintf("Could not upload mod: %v", errors.New("index out of bound")), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	err = manager.Init(window.Configuration)
	defer steam.SteamAPI_Shutdown()
	if err != nil {
		window.SendMessage(fmt.Sprintf("Failed to initialize steam: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	err = manager.UploadMod(window.Configuration, window.Configuration.Mods[index])
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not upload mod: %v", err), MessageError)
	} else {
		window.SendMessage(fmt.Sprintf("Uploaded mod successfully: %d", window.Configuration.Mods[index].Identifier), MessageSuccess)
	}
	window.RefreshMods()
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func updateMod(writer http.ResponseWriter, request *http.Request) {
	indexParameter := request.PathValue("index")
	index, err := strconv.Atoi(indexParameter)
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not parse mod index: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	if index < 0 || index >= len(window.Configuration.Mods) {
		window.SendMessage(fmt.Sprintf("Could not update mod: %v", errors.New("index out of bound")), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	if err := request.ParseForm(); err != nil {
		window.SendMessage(fmt.Sprintf("Could not update mod: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	identifierValue := request.FormValue("identifier")
	directory := request.FormValue("directory")
	descriptionFile := request.FormValue("description")
	changeNoteDirectory := request.FormValue("change-note")
	identifier, err := strconv.ParseUint(identifierValue, 10, 64)
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not parse mod identifier: %v", err), MessageError)
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	window.Configuration.Mods[index].Identifier = identifier
	window.Configuration.Mods[index].Directory = directory
	window.Configuration.Mods[index].Description = descriptionFile
	window.Configuration.Mods[index].ChangeNoteDirectory = changeNoteDirectory

	err = window.Configuration.Save()
	if err != nil {
		window.SendMessage(fmt.Sprintf("Could not save configuration: %v", err), MessageError)
	} else {
		window.SendMessage(fmt.Sprintf("Updated mod successfully: %d", window.Configuration.Mods[index].Identifier), MessageSuccess)
	}
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer func(l *net.TCPListener) {
				err := l.Close()
				if err != nil {
					logging.Fatalf("Could not close listener: %v", err)
				}
			}(l)
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
