//go:build gui

package gui

import (
	"fmt"
	"os"

	"bahmut.de/pdx-workshop-manager/config"
	"bahmut.de/pdx-workshop-manager/manager"
	"bahmut.de/pdx-workshop-manager/steam"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/tree"
)

type Topbar struct {
	window  *MainWindow
	toolbar *core.Toolbar
}

func NewTopbar(window *MainWindow) *Topbar {
	return &Topbar{
		window: window,
	}
}

func (m *Topbar) Render() {
	if m.toolbar == nil {
		m.toolbar = core.NewToolbar(m.window.body)
	} else {
		m.toolbar.Delete()
		m.window.body.Update()
		m.toolbar = core.NewToolbar(m.window.body)
	}

	m.toolbar.Maker(func(p *tree.Plan) {
		tree.Add(p, func(w *core.Chooser) {
			w.SetType(core.ChooserOutlined)
			gameItems := make([]core.ChooserItem, len(games))
			i := 0
			for id, name := range games {
				gameItems[i] = core.ChooserItem{
					Value: id,
					Text:  name,
				}
				i++
			}
			w.SetItems(gameItems...)
			w.OnChange(func(e events.Event) {
				oldGame := m.window.configuration.Game
				m.window.configuration.Game = w.CurrentItem.Value.(uint)
				if !m.window.saveConfig() {
					w.SetCurrentValue(oldGame)
				}
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Upload Selected Mods")
			w.SetIcon(icons.PlayArrow)
			w.OnClick(func(e events.Event) {
				err := manager.Init(m.window.configuration)
				defer steam.SteamAPI_Shutdown()
				if err != nil {
					core.ErrorSnackbar(m.window.body, err, "Failed to initialize steam")
					return
				}

				for _, mod := range m.window.mods {
					if !mod.upload {
						continue
					}
					core.MessageSnackbar(m.window.body, fmt.Sprintf("Start Mod Upload: %d", mod.configuration.Identifier))
					err := manager.UploadMod(m.window.configuration, mod.configuration)
					if err != nil {
						core.ErrorSnackbar(m.window.body, err, "Mod could not be uploaded")
						continue
					}
					mod.upload = false
					m.window.RenderMods()
					core.MessageSnackbar(m.window.body, fmt.Sprintf("Finished Mod Upload: %d", mod.configuration.Identifier))
				}

			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Add new Mod")
			w.SetIcon(icons.Add)
			w.OnClick(func(e events.Event) {
				addedMod := &config.ModConfig{}
				m.window.configuration.Mods = append(m.window.configuration.Mods, addedMod)
				m.window.mods = append(m.window.mods, NewModFrame(m.window, addedMod))
				if m.window.saveConfig() {
					m.window.RenderMods()
				}
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Close")
			w.SetIcon(icons.Close)
			w.OnClick(func(e events.Event) {
				os.Exit(0)
			})
		})
	})
}
