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
	status := "Hold the left mouse button and drag to draw."

	canvas := rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(rosaline.White)
		for _, line := range strokes {
			if len(line.points) == 1 {
				p := line.points[0]
				c.FillCircle(p.x, p.y, 3, rosaline.Rose)
			}
			for i := 1; i < len(line.points); i++ {
				from := line.points[i-1]
				to := line.points[i]
				c.Line(from.x, from.y, to.x, to.y, 6, rosaline.Rose)
			}
		}
	}).Size(680, 400).Expand()

	canvas.OnMouseDown(func(event rosaline.MouseEvent) {
		if event.Button != rosaline.MouseLeft {
			return
		}
		painting = true
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
		Content: rosaline.Column(
			rosaline.Label("Interactive canvas").Color(rosaline.Rose),
			rosaline.LabelFunc(func() string { return status }),
			canvas,
			rosaline.Button("Clear drawing", func() {
				strokes = nil
				painting = false
				status = "The canvas is clear. Draw something new!"
				canvas.Redraw()
			}),
		).Gap(10).Expand(),
	})
}
