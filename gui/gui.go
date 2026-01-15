//go:build gui

package gui

import (
	"bahmut.de/pdx-workshop-manager/config"
	"cogentcore.org/core/core"
)

var games = map[uint]string{
	529340:  "Victoria 3",
	3450310: "Europa Universalis V",
}

var window *MainWindow

type MainWindow struct {
	body          *core.Body
	configuration *config.ApplicationConfig
	mods          []*ModFrame
	topbar        *Topbar
}

func (mw *MainWindow) Render() {
	mw.topbar.Render()
	mw.RenderMods()
}

func (mw *MainWindow) RenderMods() {
	for _, mod := range mw.mods {
		mod.Render()
	}
}

func Run() error {
	body := core.NewBody("PDX Workshop Manager")
	configuration := loadOrSetupConfig(body)
	window = &MainWindow{
		body:          body,
		configuration: configuration,
	}
	window.topbar = NewTopbar(window)
	window.mods = make([]*ModFrame, len(window.configuration.Mods))
	for index, mod := range window.configuration.Mods {
		window.mods[index] = NewModFrame(window, mod)
	}
	window.Render()
	body.RunMainWindow()
	return nil
}

func loadOrSetupConfig(body *core.Body) *config.ApplicationConfig {
	appConfig, err := config.LoadConfig(config.DefaultFileName)
	if err != nil {
		err = nil
		appConfig, err := config.InitializeConfig(config.DefaultFileName, 529340)
		if err != nil {
			core.ErrorSnackbar(body, err, "Error initializing configuration file")
		}
		return appConfig
	} else {
		return appConfig
	}
}

func (mw *MainWindow) saveConfig() bool {
	err := mw.configuration.Save()
	if err != nil {
		core.ErrorSnackbar(mw.body, err, "Error saving configuration file")
		return false
	} else {
		core.MessageSnackbar(mw.body, "Configuration file saved")
		return true
	}
}
