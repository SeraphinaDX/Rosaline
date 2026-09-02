// SPDX-License-Identifier: LGPL-3.0-or-later

// The keyboard example combines canvas key events with window shortcuts.
// Run it with: go run ./examples/keyboard
package main

import (
	"fmt"

	"github.com/SeraphinaDX/Rosaline"
)

type flower struct {
	x float64
	y float64
}

func main() {
	const (
		canvasWidth  = 700
		canvasHeight = 400
	)

	x, y := 350.0, 210.0
	flowers := make([]flower, 0)
	status := "The garden has keyboard focus"
	lastKey := "None"

	var canvas *rosaline.CanvasWidget
	canvas = rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(rosaline.Hex("#fff7fb"))

		// A quiet grid makes movement easy to see.
		for lineX := 20.0; lineX < canvasWidth; lineX += 40 {
			c.Line(lineX, 45, lineX, canvasHeight, 1, rosaline.Hex("#f9dce9"))
		}
		for lineY := 45.0; lineY < canvasHeight; lineY += 40 {
			c.Line(0, lineY, canvasWidth, lineY, 1, rosaline.Hex("#f9dce9"))
		}

		for _, planted := range flowers {
			c.Line(planted.x, planted.y+7, planted.x, planted.y+24, 3, rosaline.Hex("#60966d"))
			c.FillCircle(planted.x-7, planted.y, 7, rosaline.SoftRose)
			c.FillCircle(planted.x+7, planted.y, 7, rosaline.SoftRose)
			c.FillCircle(planted.x, planted.y-7, 7, rosaline.SoftRose)
			c.FillCircle(planted.x, planted.y+7, 7, rosaline.SoftRose)
			c.FillCircle(planted.x, planted.y, 5, rosaline.Hex("#ffd166"))
		}

		// The dark ring marks the current planting position.
		c.FillCircle(x, y, 13, rosaline.Rose)
		c.Circle(x, y, 18, 3, rosaline.Hex("#7d244f"))
		c.Text("Keyboard Garden", 18, 13, rosaline.TextStyle{
			Color: rosaline.Hex("#7d244f"),
			Size:  20,
		})
	}).Size(canvasWidth, canvasHeight).Expand().Focus()

	canvas.OnKeyDown(func(event rosaline.KeyEvent) {
		// Let window shortcuts such as Primary+S pass through without treating
		// their letter as a movement key.
		if event.Primary || event.Control || event.Alt {
			lastKey = event.Key.String()
			return
		}

		step := 12.0
		if event.Shift {
			step = 30
		}

		switch event.Key {
		case rosaline.KeyLeft, rosaline.Key("a"):
			x -= step
			status = "Moved left"
		case rosaline.KeyRight, rosaline.Key("d"):
			x += step
			status = "Moved right"
		case rosaline.KeyUp, rosaline.Key("w"):
			y -= step
			status = "Moved up"
		case rosaline.KeyDown, rosaline.Key("s"):
			y += step
			status = "Moved down"
		case rosaline.KeySpace:
			flowers = append(flowers, flower{x: x, y: y})
			status = "Planted a flower"
		case rosaline.KeyBackspace, rosaline.KeyDelete:
			if len(flowers) > 0 {
				flowers = flowers[:len(flowers)-1]
				status = "Removed the newest flower"
			}
		default:
			status = "Pressed " + event.Key.String()
		}

		x = min(max(x, 20), canvasWidth-20)
		y = min(max(y, 65), canvasHeight-28)
		lastKey = event.Key.String()
	})

	canvas.OnKeyUp(func(event rosaline.KeyEvent) {
		lastKey = event.Key.String() + " released"
	})

	reset := func() {
		x, y = 350, 210
		flowers = nil
		status = "Garden reset"
		canvas.Redraw()
	}

	save := func() {
		path, ok := rosaline.SaveFileDialog(rosaline.FileDialogOptions{
			Title:            "Save Keyboard Garden",
			InitialFile:      "keyboard-garden.png",
			DefaultExtension: ".png",
			Filters: []rosaline.FileFilter{
				{Name: "PNG image", Extensions: []string{".png"}},
			},
		})
		if !ok {
			status = "Save cancelled"
			return
		}
		if err := canvas.Picture().SavePNG(path); err != nil {
			rosaline.Error("Could not save garden", err.Error())
			status = "Save failed"
			return
		}
		status = "Saved " + path
	}

	help := func() {
		rosaline.Message(
			"Keyboard Garden",
			"Move with the arrow keys or WASD.\nHold Shift to move farther.\nSpace plants a flower.\nBackspace removes the newest flower.",
		)
		status = "Help opened with F1"
	}

	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Keyboard Garden",
		Width:  780,
		Height: 580,
		Shortcuts: rosaline.Shortcuts(
			rosaline.Shortcut("Primary+R", reset),
			rosaline.Shortcut("Primary+S", save),
			rosaline.Shortcut("F1", help),
		),
		Content: rosaline.Column(
			rosaline.Label("Keyboard input").Color(rosaline.Rose),
			rosaline.Label("Move with arrows or WASD · Shift moves farther · Space plants"),
			canvas,
			rosaline.Row(
				rosaline.Button("Reset", reset),
				rosaline.Button("Save PNG", save).Primary(),
				rosaline.Button("Help", help),
			),
			rosaline.LabelFunc(func() string {
				return fmt.Sprintf("%s · Last key: %s · Flowers: %d", status, lastKey, len(flowers))
			}),
			rosaline.Label("Shortcuts: Primary+R reset · Primary+S save · F1 help"),
		).Gap(9).Expand(),
	})
}
