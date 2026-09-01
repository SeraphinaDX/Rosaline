// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Canvas",
		Width:  640,
		Height: 440,
		Content: rosaline.Column(
			rosaline.Label("Canvas drawing"),
			rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
				c.Clear(rosaline.Hex("#fff8fc"))
				c.FillRect(24, 24, 130, 80, rosaline.SoftRose)
				c.Rect(24, 24, 130, 80, 3, rosaline.Rose)
				c.FillCircle(245, 65, 42, rosaline.Rose)
				c.Line(315, 25, 430, 105, 5, rosaline.Hex("#6d3c58"))
				c.Text("Drawing should feel like Go.", 24, 145, rosaline.TextStyle{
					Color: rosaline.Hex("#2a1722"),
					Size:  18,
				})
			}).Size(560, 260).Expand(),
		).Gap(10).Expand(),
	})
}
