//go:build gui

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

type ModFrame struct {
	window        *MainWindow
	configuration *config.ModConfig
	frame         *core.Frame
	upload        bool
}

func NewModFrame(window *MainWindow, configuration *config.ModConfig) *ModFrame {
	return &ModFrame{
		window:        window,
		configuration: configuration,
		upload:        false,
	}
}

func (m *ModFrame) Delete() {
	old := m.window.configuration.Mods
	index := slices.Index(old, m.configuration)
	if index < 0 || index > len(old) {
		core.ErrorSnackbar(m.window.body, errors.New("mod index out of bound"), "Error removing mod")
	}
	newModList := append(old[:index], old[index+1:]...)
	m.window.configuration.Mods = newModList
	if m.window.saveConfig() {
		if m.frame != nil {
			m.frame.Delete()
		}
		index := slices.Index(m.window.mods, m)
		if index < 0 || index > len(m.window.mods) {
			core.ErrorSnackbar(m.window.body, errors.New("mod index out of bound"), "Error removing mod")
		}
		m.window.mods = append(m.window.mods[:index], m.window.mods[index+1:]...)
		m.window.body.Update()
	} else {
		m.window.configuration.Mods = old
	}
}

func (m *ModFrame) Destroy() {
	if m.frame != nil {
		m.frame.Delete()
	}
}

func (m *ModFrame) Render() {
	if m.frame == nil {
		m.frame = core.NewFrame(m.window.body)
	} else {
		m.frame.Delete()
		m.window.body.Update()
		m.frame = core.NewFrame(m.window.body)
	}

	m.frame.Styler(func(s *styles.Style) {
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

	modToolbar := core.NewToolbar(m.frame)
	modConfigFrame := core.NewFrame(m.frame)
	modConfigFrame.Styler(func(s *styles.Style) {
		s.Display = styles.Grid
		s.Columns = 2
		s.Justify.Content = styles.Center
		s.Align.Items = styles.Center
		s.Min.X.Set(90, units.UnitPw)
		s.Gap.Set(units.Em(0.5))
	})

	m.createModLabel("Identifier", modConfigFrame)
	identifierInput := core.NewTextField(modConfigFrame)
	identifierInput.SetText(strconv.FormatUint(m.configuration.Identifier, 10))
	identifierInput.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(80, units.UnitPw)
	})
	identifierInput.OnChange(func(e events.Event) {
		identifier, err := strconv.ParseUint(identifierInput.Text(), 10, 64)
		if err != nil {
			core.ErrorSnackbar(m.window.body, err, "Invalid mod identifier: Mod identifiers are a number")
			identifierInput.SetText(strconv.FormatUint(m.configuration.Identifier, 10))
		} else {
			old := m.configuration.Identifier
			m.configuration.Identifier = identifier
			if !m.window.saveConfig() {
				m.configuration.Identifier = old
				identifierInput.SetText(strconv.FormatUint(old, 10))
			}
		}
	})

	m.createModLabel("Mod Directory", modConfigFrame)
	directoryInput := core.NewTextField(modConfigFrame)
	directoryInput.SetText(m.configuration.Directory)
	directoryInput.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(80, units.UnitPw)
	})
	directoryInput.OnChange(func(e events.Event) {
		old := m.configuration.Directory
		m.configuration.Directory = directoryInput.Text()
		if !m.window.saveConfig() {
			m.configuration.Directory = old
			directoryInput.SetText(old)
		}
	})

	m.createModLabel("Description File", modConfigFrame)
	descriptionInput := core.NewTextField(modConfigFrame)
	descriptionInput.SetText(m.configuration.Description)
	descriptionInput.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(80, units.UnitPw)
	})
	descriptionInput.OnChange(func(e events.Event) {
		old := m.configuration.Description
		m.configuration.Directory = descriptionInput.Text()
		if !m.window.saveConfig() {
			m.configuration.Description = old
			descriptionInput.SetText(old)
		}
	})

	m.createModLabel("Change Note Directory", modConfigFrame)
	changeNoteDirectoryInput := core.NewTextField(modConfigFrame)
	changeNoteDirectoryInput.SetText(m.configuration.ChangeNoteDirectory)
	changeNoteDirectoryInput.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(80, units.UnitPw)
	})
	changeNoteDirectoryInput.OnChange(func(e events.Event) {
		old := m.configuration.ChangeNoteDirectory
		m.configuration.ChangeNoteDirectory = changeNoteDirectoryInput.Text()
		if !m.window.saveConfig() {
			m.configuration.ChangeNoteDirectory = old
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
			w.SetChecked(m.upload)
			w.OnChange(func(e events.Event) {
				if w.IsChecked() {
					if strings.TrimSpace(directoryInput.Text()) == "" {
						core.ErrorSnackbar(m.window.body, errors.New("no directory set"), "Cannot upload mod")
						w.SetChecked(false)
						return
					}
					if _, err := os.Stat(directoryInput.Text()); err != nil {
						if os.IsNotExist(err) {
							core.ErrorSnackbar(m.window.body, errors.New("directory does not exist"), "Cannot upload mod")
							w.SetChecked(false)
							return
						}
					}
				}
				m.upload = w.IsChecked()
			})
		})
		tree.Add(p, func(w *core.Button) {
			w.SetText("Remove Mod")
			w.SetIcon(icons.Delete)
			w.OnClick(func(e events.Event) {
				m.Delete()
			})
		})
	})

	m.window.body.Update()
}

func (m *ModFrame) createModLabel(text string, frame *core.Frame) {
	descriptionLabel := core.NewText(frame)
	descriptionLabel.SetText(text)
	descriptionLabel.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
		s.Min.X.Set(20, units.UnitPw)
	})
}
