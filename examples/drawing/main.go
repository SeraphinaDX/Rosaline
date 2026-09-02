// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	heart := rosaline.NewPath().
		MoveTo(0, 34).
		CubicTo(-58, -4, -64, 74, 0, 112).
		CubicTo(64, 74, 58, -4, 0, 34).
		Close()

	canvas := rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(rosaline.Hex("#fff6fb"))
		c.ClipRect(24, 24, 632, 352)

		// A clipped background pattern makes the drawing area visible.
		for x := -300.0; x < 900; x += 28 {
			c.Line(x, 20, x+250, 380, 2, rosaline.Hex("#f8d7e7"))
		}

		// Push and Pop keep each petal's transform independent.
		for rotation := 0.0; rotation < 360; rotation += 45 {
			c.Push()
			c.Translate(340, 190)
			c.Rotate(rotation)
			c.Translate(0, -126)
			c.Scale(0.45, 0.45)
			c.FillPath(heart, rosaline.SoftRose)
			c.Pop()
		}

		c.Push()
		c.Translate(340, 122)
		c.FillPath(heart, rosaline.Rose)
		c.StrokePath(heart, 4, rosaline.Hex("#8e2856"))
		c.Pop()

		c.ResetClip()
		c.Rect(24, 24, 632, 352, 3, rosaline.Hex("#b73e74"))
		c.Text("Paths · Bézier curves · transforms · clipping", 112, 390,
			rosaline.TextStyle{Color: rosaline.Hex("#7d244f"), Size: 17})
	}).Size(680, 430).Expand()

	rosaline.RunApp(rosaline.App{
		Title:   "Rosaline Advanced Drawing",
		Width:   760,
		Height:  520,
		Content: canvas,
	})
}
