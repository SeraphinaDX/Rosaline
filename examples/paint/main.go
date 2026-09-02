// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/SeraphinaDX/Rosaline"
)

type point struct {
	x float64
	y float64
}

type stroke struct {
	points []point
}

func main() {
	var strokes []stroke
	painting := false
	dirty := false
	status := "Hold the left mouse button and drag to draw. Save from the File menu."

	canvas := rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(rosaline.White)
		for _, line := range strokes {
			if len(line.points) == 0 {
				continue
			}
			if len(line.points) == 1 {
				p := line.points[0]
				c.FillCircle(p.x, p.y, 3, rosaline.Rose)
				continue
			}

			path := rosaline.NewPath().MoveTo(line.points[0].x, line.points[0].y)
			for _, point := range line.points[1:] {
				path.LineTo(point.x, point.y)
			}
			c.StrokePath(path, 6, rosaline.Rose)
		}
	}).Size(680, 400).Expand()

	savePNG := func() {
		path, ok := rosaline.SaveFileDialog(rosaline.FileDialogOptions{
			Title:            "Save Painting as PNG",
			InitialFile:      "rosaline-painting.png",
			DefaultExtension: ".png",
			Filters: []rosaline.FileFilter{
				{Name: "PNG images", Extensions: []string{"png"}},
			},
		})
		if !ok {
			return
		}
		if err := canvas.Picture().SavePNG(path); err != nil {
			rosaline.Error("Could not save painting", err.Error())
			return
		}
		dirty = false
		status = "Saved PNG: " + path
		rosaline.Message("Painting saved", "Saved a PNG image to:\n"+path)
	}

	saveAVIF := func() {
		path, ok := rosaline.SaveFileDialog(rosaline.FileDialogOptions{
			Title:            "Save Painting as AVIF",
			InitialFile:      "rosaline-painting.avif",
			DefaultExtension: ".avif",
			Filters: []rosaline.FileFilter{
				{Name: "AVIF images", Extensions: []string{"avif"}},
			},
		})
		if !ok {
			return
		}
		if err := canvas.Picture().SaveAVIF(path); err != nil {
			rosaline.Error("Could not save painting", err.Error())
			return
		}
		dirty = false
		status = "Saved AVIF: " + path
		rosaline.Message("Painting saved", "Saved an AVIF image to:\n"+path)
	}

	clearPainting := func() {
		if dirty && !rosaline.Confirm("Clear painting?", "Your unsaved drawing will be removed.") {
			return
		}
		strokes = nil
		painting = false
		dirty = false
		status = "The canvas is clear. Draw something new!"
		canvas.Redraw()
	}

	canvas.OnMouseDown(func(event rosaline.MouseEvent) {
		if event.Button != rosaline.MouseLeft {
			return
		}
		painting = true
		dirty = true
		strokes = append(strokes, stroke{
			points: []point{{x: event.X, y: event.Y}},
		})
		status = fmt.Sprintf("Drawing at %.0f, %.0f", event.X, event.Y)
	})

	canvas.OnMouseMove(func(event rosaline.MouseEvent) {
		if !event.Dragging || event.Button != rosaline.MouseLeft {
			painting = false
			return
		}
		if !painting || len(strokes) == 0 {
			return
		}

		current := &strokes[len(strokes)-1]
		current.points = append(current.points, point{x: event.X, y: event.Y})
		status = fmt.Sprintf("Drawing at %.0f, %.0f", event.X, event.Y)
	})

	canvas.OnMouseUp(func(event rosaline.MouseEvent) {
		if event.Button == rosaline.MouseLeft {
			painting = false
			status = fmt.Sprintf("Finished at %.0f, %.0f", event.X, event.Y)
		}
	})

	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Paint",
		Width:  760,
		Height: 560,
		Menu: rosaline.MenuBar(
			rosaline.Menu("File",
				rosaline.MenuItem("Save as PNG…", savePNG).Shortcut("Ctrl+S"),
				rosaline.MenuItem("Save as AVIF…", saveAVIF).Shortcut("Ctrl+Shift+S"),
				rosaline.MenuSeparator(),
				rosaline.MenuItem("Quit", rosaline.Quit).Shortcut("Ctrl+Q"),
			),
			rosaline.Menu("Edit",
				rosaline.MenuItem("Clear painting", clearPainting).Shortcut("Ctrl+L"),
			),
		),
		Content: rosaline.Column(
			rosaline.Label("Rosaline Paint").Color(rosaline.Rose),
			rosaline.LabelFunc(func() string { return status }),
			canvas,
			rosaline.Row(
				rosaline.Button("Save PNG", savePNG),
				rosaline.Button("Save AVIF", saveAVIF).Primary(),
				rosaline.Button("Clear", clearPainting),
			),
		).Gap(10).Expand(),
	})
}
