//go:build gui

package web

import (
	"bahmut.de/pdx-workshop-manager/config"
)

type ModFrame struct {
	Window        *MainWindow
	Configuration *config.ModConfig
	Upload        bool
}

func NewModFrame(window *MainWindow, configuration *config.ModConfig) *ModFrame {
	return &ModFrame{
		Window:        window,
		Configuration: configuration,
		Upload:        false,
	}
}

func (m *ModFrame) Delete() {

}

func (m *ModFrame) Destroy() {
}

func (m *ModFrame) Render() {

}
