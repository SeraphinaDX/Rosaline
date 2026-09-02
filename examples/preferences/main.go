// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/SeraphinaDX/Rosaline"
)

type palette struct {
	name        string
	description string
	background  rosaline.Color
	accent      rosaline.Color
	text        rosaline.Color
}

func main() {
	palettes := []palette{
		{"Rosaline", "Warm rose with a soft paper background.", rosaline.Hex("#fff8fc"), rosaline.Hex("#c43f7a"), rosaline.Hex("#2a1722")},
		{"Lavender", "Calm violet tones for an evening workspace.", rosaline.Hex("#f8f5ff"), rosaline.Hex("#7756b3"), rosaline.Hex("#261d36")},
		{"Ocean", "Clear blue with cool, quiet surfaces.", rosaline.Hex("#f3f9fc"), rosaline.Hex("#267aa5"), rosaline.Hex("#17303c")},
		{"Forest", "Natural greens with a gentle cream background.", rosaline.Hex("#f7faf4"), rosaline.Hex("#477b51"), rosaline.Hex("#1e3222")},
	}

	selectedPalette := 0
	fontSizes := []string{"Small — 12 px", "Comfortable — 14 px", "Large — 16 px", "Extra large — 18 px"}
	selectedFont := 1
	displayName := "Seraphina"
	autosave := true
	showTips := true
	autoUpdates := true
	status := "Ready"

	var preview *rosaline.CanvasWidget
	preview = rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		choice := palettes[selectedPalette]
		c.Clear(choice.background)
		c.FillRect(18, 18, 430, 132, rosaline.White)
		c.FillRect(18, 18, 12, 132, choice.accent)
		c.Text("Rosaline Preferences", 52, 45, rosaline.TextStyle{Color: choice.text, Size: 21})
		c.Text(choice.description, 52, 80, rosaline.TextStyle{Color: choice.text, Size: 13})
		c.FillRect(52, 108, 120, 27, choice.accent)
		c.Text("Preview", 83, 114, rosaline.TextStyle{Color: rosaline.White, Size: 13})
	}).Size(470, 168)

	themeList := rosaline.List("Rosaline", "Lavender", "Ocean", "Forest").
		Size(18, 6).
		OnSelect(func(index int, value string) {
			selectedPalette = index
			status = "Selected the " + value + " palette"
			preview.Redraw()
		}).
		OnActivate(func(index int, value string) {
			rosaline.Message("Palette selected", value+" will be used after you save.")
		})

	fontList := rosaline.List(fontSizes...).
		Size(28, 6).
		OnSelect(func(index int, value string) {
			selectedFont = index
			status = "Editor text size: " + value
		})
	fontList.Select(selectedFont)

	tabs := rosaline.Tabs(
		rosaline.Tab("Appearance", rosaline.Column(
			rosaline.Label("Choose a color palette").Color(rosaline.Rose),
			rosaline.Label("Use the arrow keys to move, or press Enter to activate a choice."),
			rosaline.Row(
				themeList,
				rosaline.Column(
					rosaline.LabelFunc(func() string { return palettes[selectedPalette].name + " preview" }),
					preview,
				).Gap(8).Expand(),
			).Gap(14).Expand(),
		).Gap(10).Expand()),
		rosaline.Tab("Editor", rosaline.Column(
			rosaline.Label("Editor preferences").Color(rosaline.Rose),
			rosaline.Label("Display name"),
			rosaline.TextBox(&displayName).Placeholder("Your name").Width(34),
			rosaline.CheckBox("Save documents automatically", &autosave),
			rosaline.CheckBox("Show beginner tips", &showTips),
			rosaline.Label("Text size"),
			fontList,
		).Gap(10).Expand()),
		rosaline.Tab("Updates", rosaline.Column(
			rosaline.Label("Keep Rosaline current").Color(rosaline.Rose),
			rosaline.CheckBox("Check for updates when the app starts", &autoUpdates),
			rosaline.Label("This demonstration does not connect to the internet."),
			rosaline.Button("Check now", func() {
				status = "Rosaline is up to date"
				rosaline.Message("Updates", "You are using the newest demonstration version.")
			}),
		).Gap(10).Expand()),
		rosaline.Tab("About", rosaline.Column(
			rosaline.Label("Rosaline v0.6.0").Color(rosaline.Rose),
			rosaline.Label("A small, beginner-friendly GUI and graphics library for Go."),
			rosaline.Label("Pure Go · Linux first-class · LGPL-3.0-or-later"),
		).Gap(12).Expand()),
	).Expand().OnChange(func(index int, title string) {
		status = fmt.Sprintf("Opened %s (tab %d)", title, index+1)
	})

	restoreDefaults := func() {
		displayName = "Seraphina"
		autosave = true
		showTips = true
		autoUpdates = true
		selectedPalette = 0
		selectedFont = 1
		themeList.Select(selectedPalette)
		fontList.Select(selectedFont)
		tabs.Select(0)
		preview.Redraw()
		status = "Restored the default preferences"
	}

	save := func() {
		_, fontSize, _ := fontList.Selected()
		message := fmt.Sprintf(
			"Saved preferences for %s.\n\nPalette: %s\nText size: %s",
			displayName,
			palettes[selectedPalette].name,
			fontSize,
		)
		status = "Preferences saved"
		rosaline.Message("Preferences saved", message)
	}

	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Preferences",
		Width:  820,
		Height: 600,
		Content: rosaline.Column(
			rosaline.Label("Preferences").Color(rosaline.Rose),
			rosaline.Label("Tabs organize related settings; lists make choices easy to scan."),
			tabs,
			rosaline.Row(
				rosaline.Button("Restore defaults", restoreDefaults),
				rosaline.Button("Save preferences", save).Primary(),
			),
			rosaline.LabelFunc(func() string { return status }),
		).Gap(10).Expand(),
	})
}
