package gui

import (
	"errors"
	"os"
	"slices"
	"strconv"
	"strings"

	"bahmut.de/pdx-workshop-manager/config"
	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

var applicationConfig = &config.ApplicationConfig{}
var games = map[uint]string{
	529340:  "Victoria 3",
	3450310: "Europa Universalis V",
}
var uploadMods = map[uint64]bool{}

func Run() error {
	body := core.NewBody("PDX Workshop Manager")
	loadOrSetupConfig(body)
	setupTopbarFrame(body)
	for _, mod := range applicationConfig.Mods {
		setupModConfigurationFrame(
			body,
			mod,
		)
	}
	body.RunMainWindow()
	return nil
}

func loadOrSetupConfig(body *core.Body) {
	appConfig, err := config.LoadConfig(config.DefaultFileName)
	if err != nil {
		err = nil
		appConfig, err := config.InitializeConfig(config.DefaultFileName, 529340)
		if err != nil {
			core.ErrorSnackbar(body, err, "Error initializing configuration file")
		}
		applicationConfig = appConfig
	} else {
		applicationConfig = appConfig
	}
	for _, mod := range applicationConfig.Mods {
		uploadMods[mod.Identifier] = false
	}
}

func setupModConfigurationFrame(
	body *core.Body,
	mod *config.ModConfig,
) {
	modFrame := core.NewFrame(body)
	modFrame.Styler(func(s *styles.Style) {
		s.Justify.Content = styles.Center
		s.Align.Items = styles.Center
		s.Display = styles.Grid
		s.Columns = 1
		s.Grow.Set(1, 0)
		s.Gap.Set(units.Em(1))
		s.Border.Radius = styles.BorderRadiusLarge
		s.Border.Style.Set(styles.BorderSolid)
		s.Border.Width.Set(units.Dp(1))
		s.Border.Color.Set(colors.Scheme.Outline)
		s.Padding.Set(units.Em(1))
		s.Margin.Set(units.Em(0.5))
	})

	modToolbar := core.NewToolbar(modFrame)

	modConfigFrame := core.NewFrame(modFrame)
	modConfigFrame.Styler(func(s *styles.Style) {
		s.Display = styles.Grid
		s.Columns = 2
		s.Justify.Content = styles.Center
		s.Align.Items = styles.Center
		s.Min.X.Set(90, units.UnitPw)
		s.Gap.Set(units.Em(0.5))
	})

	createModLabel("Identifier", modConfigFrame)
	identifierInput := core.NewTextField(modConfigFrame)
	identifierInput.SetText(strconv.FormatUint(mod.Identifier, 10))
	identifierInput.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(80, units.UnitPw)
	})
	identifierInput.OnChange(func(e events.Event) {
		identifier, err := strconv.ParseUint(identifierInput.Text(), 10, 64)
		if err != nil {
			core.ErrorSnackbar(body, err, "Invalid mod identifier: Mod identifiers are a number")
			identifierInput.SetText(strconv.FormatUint(mod.Identifier, 10))
		} else {
			old := mod.Identifier
			mod.Identifier = identifier
			if !saveConfig(body) {
				mod.Identifier = old
				identifierInput.SetText(strconv.FormatUint(old, 10))
			}
		}
	})

	createModLabel("Mod Directory", modConfigFrame)
	directoryInput := core.NewTextField(modConfigFrame)
	directoryInput.SetText(mod.Directory)
	directoryInput.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(80, units.UnitPw)
	})
	directoryInput.OnChange(func(e events.Event) {
		old := mod.Directory
		mod.Directory = directoryInput.Text()
		if !saveConfig(body) {
			mod.Directory = old
			directoryInput.SetText(old)
		}
	})

	createModLabel("Description File", modConfigFrame)
	descriptionInput := core.NewTextField(modConfigFrame)
	descriptionInput.SetText(mod.Description)
	descriptionInput.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(80, units.UnitPw)
	})
	descriptionInput.OnChange(func(e events.Event) {
		old := mod.Description
		mod.Directory = descriptionInput.Text()
		if !saveConfig(body) {
			mod.Description = old
			descriptionInput.SetText(old)
		}
	})

	createModLabel("Change Note Directory", modConfigFrame)
	changeNoteDirectoryInput := core.NewTextField(modConfigFrame)
	changeNoteDirectoryInput.SetText(mod.ChangeNoteDirectory)
	changeNoteDirectoryInput.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(80, units.UnitPw)
	})
	changeNoteDirectoryInput.OnChange(func(e events.Event) {
		old := mod.ChangeNoteDirectory
		mod.ChangeNoteDirectory = changeNoteDirectoryInput.Text()
		if !saveConfig(body) {
			mod.ChangeNoteDirectory = old
			changeNoteDirectoryInput.SetText(old)
		}
	})

	modToolbar.Maker(func(p *tree.Plan) {
		tree.Add(p, func(w *core.Switch) {
			w.SetText("Upload Mod")
			w.Styler(func(s *styles.Style) {
				s.Padding.Right.Set(1, units.UnitEm)
				s.Padding.Left.Set(1, units.UnitEm)
			})
			w.OnChange(func(e events.Event) {
				identifier, err := strconv.ParseUint(identifierInput.Text(), 10, 64)
				if err != nil {
					core.ErrorSnackbar(body, err, "Invalid mod identifier: Mod identifiers are a number")
					identifierInput.SetText("0")
				}
				if w.IsChecked() {
					if strings.TrimSpace(directoryInput.Text()) == "" {
						core.ErrorSnackbar(body, errors.New("no directory set"), "Cannot upload mod")
						w.SetChecked(false)
						return
					}
					if _, err := os.Stat(directoryInput.Text()); err != nil {
						if os.IsNotExist(err) {
							core.ErrorSnackbar(body, errors.New("directory does not exist"), "Cannot upload mod")
							w.SetChecked(false)
							return
						}
					}
				}
				uploadMods[identifier] = w.IsChecked()
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Remove Mod")
			w.SetIcon(icons.Delete)
			w.OnClick(func(e events.Event) {
				old := applicationConfig.Mods
				index := slices.Index(old, mod)
				if index < 0 || index > len(old) {
					core.ErrorSnackbar(body, errors.New("mod index out of bound"), "Error removing mod")
				}
				newModList := append(old[:index], old[index+1:]...)
				applicationConfig.Mods = newModList
				if saveConfig(body) {
					modFrame.Delete()
					body.Update()
				} else {
					applicationConfig.Mods = old
				}
			})
		})
	})
}

func setupTopbarFrame(body *core.Body) {
	topbar := core.NewToolbar(body)
	topbar.Maker(func(p *tree.Plan) {
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
				oldGame := applicationConfig.Game
				applicationConfig.Game = w.CurrentItem.Value.(uint)
				if !saveConfig(body) {
					w.SetCurrentValue(oldGame)
				}
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Upload Selected Mods")
			w.SetIcon(icons.PlayArrow)
			w.OnClick(func(e events.Event) {
				core.MessageSnackbar(body, "Button clicked")
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Add new Mod")
			w.SetIcon(icons.Add)
			w.OnClick(func(e events.Event) {
				addedMod := &config.ModConfig{}
				applicationConfig.Mods = append(applicationConfig.Mods, addedMod)
				if saveConfig(body) {
					setupModConfigurationFrame(body, addedMod)
					body.Update()
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

func createModLabel(text string, frame *core.Frame) {
	descriptionLabel := core.NewText(frame)
	descriptionLabel.SetText(text)
	descriptionLabel.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(20, units.UnitPw)
	})
}

func saveConfig(body *core.Body) bool {
	err := applicationConfig.Save()
	if err != nil {
		core.ErrorSnackbar(body, err, "Error saving configuration file")
		return false
	} else {
		core.MessageSnackbar(body, "Configuration file saved")
		return true
	}
}
