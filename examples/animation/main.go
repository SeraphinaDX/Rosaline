// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"fmt"
	"math"
	"time"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	const (
		canvasWidth  = 680
		canvasHeight = 380
		ballRadius   = 24
	)

	x, y := 110.0, 100.0
	dx, dy := 3.2, 2.4
	angle := 0.0
	seconds := 0
	status := "Playing"
	tip := "The canvas is redrawn by a 60 FPS animation timer."

	var canvas *rosaline.CanvasWidget
	canvas = rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(rosaline.Hex("#fff7fb"))

		// The soft horizontal lines give the scene a little depth.
		for lineY := 60.0; lineY < canvasHeight; lineY += 64 {
			c.Line(0, lineY, canvasWidth, lineY, 1, rosaline.Hex("#f8d9e7"))
		}

		// Three small petals orbit the bouncing rose.
		for i := 0; i < 3; i++ {
			petalAngle := angle + float64(i)*2*math.Pi/3
			petalX := x + math.Cos(petalAngle)*42
			petalY := y + math.Sin(petalAngle)*42
			c.FillCircle(petalX, petalY, 7, rosaline.SoftRose)
		}

		c.FillCircle(x, y, ballRadius, rosaline.Rose)
		c.Circle(x, y, ballRadius+4, 2, rosaline.Hex("#9d2f62"))
		c.Text("Rosaline in motion", 18, 16, rosaline.TextStyle{
			Color: rosaline.Hex("#7d244f"),
			Size:  18,
		})
	}).Size(canvasWidth, canvasHeight).Expand()

	animation := rosaline.Animate(60, func() {
		x += dx
		y += dy
		angle += 0.05

		if x-ballRadius <= 0 || x+ballRadius >= canvasWidth {
			dx = -dx
			x = min(max(x, ballRadius), canvasWidth-ballRadius)
		}
		if y-ballRadius <= 52 || y+ballRadius >= canvasHeight {
			dy = -dy
			y = min(max(y, 52+ballRadius), canvasHeight-ballRadius)
		}

		canvas.Redraw()
	})

	clock := rosaline.Every(time.Second, func() {
		seconds++
	})

	newTip := rosaline.After(5*time.Second, func() {
		tip = "Try Pause, then Resume—the same timer continues safely."
	})

	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Animation",
		Width:  760,
		Height: 560,
		Timers: []*rosaline.Timer{animation, clock, newTip},
		Content: rosaline.Column(
			rosaline.Label("Timers and animation").Color(rosaline.Rose),
			rosaline.LabelFunc(func() string {
				return fmt.Sprintf("%s · running for %d seconds", status, seconds)
			}),
			canvas,
			rosaline.Row(
				rosaline.Button("Pause", func() {
					animation.Stop()
					status = "Paused"
				}),
				rosaline.Button("Resume", func() {
					animation.Start()
					status = "Playing"
				}).Primary(),
				rosaline.Button("Restart tip timer", func() {
					tip = "A fresh five-second timer has started."
					newTip.Restart()
				}),
			),
			rosaline.LabelFunc(func() string { return tip }),
		).Gap(10).Expand(),
	})
}
