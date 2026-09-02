// SPDX-License-Identifier: LGPL-3.0-or-later

// The background example renders a picture without freezing the window.
// Run it with: go run ./examples/background
package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/SeraphinaDX/Rosaline"
)

const (
	bloomWidth  = 720
	bloomHeight = 420
)

func main() {
	progress := 0.0
	status := "Preparing the first bloom..."
	var finished *rosaline.Picture

	preview := rosaline.Image(nil).
		Placeholder("Your generated bloom will appear here.").
		Expand()

	var task *rosaline.Task
	task = rosaline.Background(func(ctx context.Context, report *rosaline.TaskReporter) error {
		pixels, err := paintBloom(ctx, report)
		if err != nil {
			return err
		}

		picture := rosaline.NewPicture(pixels)
		report.Post(func() {
			finished = picture
			preview.SetImage(picture)
			status = "Bloom complete — ready to save"
		})
		return nil
	}).OnProgress(func(update rosaline.TaskProgress) {
		progress = update.Percent
		status = update.Message
	}).OnDone(func(err error) {
		switch {
		case errors.Is(err, context.Canceled):
			status = "Generation cancelled"
		case err != nil:
			status = "Generation failed"
			rosaline.Error("Could not create bloom", err.Error())
		default:
			progress = 100
		}
	}).AutoStart()

	start := func() {
		if task.Running() {
			return
		}
		progress = 0
		status = "Preparing a fresh bloom..."
		task.Start()
	}

	save := func() {
		if finished == nil {
			rosaline.Message("Background Bloom", "Create a bloom before saving it.")
			return
		}
		path, ok := rosaline.SaveFileDialog(rosaline.FileDialogOptions{
			Title:            "Save Background Bloom",
			InitialFile:      "background-bloom.png",
			DefaultExtension: ".png",
			Filters: []rosaline.FileFilter{
				{Name: "PNG image", Extensions: []string{".png"}},
			},
		})
		if !ok {
			return
		}
		if err := finished.SavePNG(path); err != nil {
			rosaline.Error("Could not save bloom", err.Error())
			return
		}
		status = "Saved " + path
	}

	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Background Bloom",
		Width:  800,
		Height: 610,
		Tasks:  []*rosaline.Task{task},
		Shortcuts: rosaline.Shortcuts(
			rosaline.Shortcut("Primary+R", start),
			rosaline.Shortcut("Primary+S", save),
			rosaline.Shortcut("Escape", task.Cancel),
		),
		Content: rosaline.Column(
			rosaline.Label("Background Bloom").Color(rosaline.Rose),
			rosaline.Label("A responsive window while pure Go paints every pixel"),
			preview,
			rosaline.ProgressBar(&progress).Length(bloomWidth),
			rosaline.Row(
				rosaline.Button("Create another", start).Primary(),
				rosaline.Button("Cancel", task.Cancel),
				rosaline.Button("Save PNG", save),
			),
			rosaline.LabelFunc(func() string {
				return fmt.Sprintf("%s · %.0f%%", status, progress)
			}),
			rosaline.Label("Shortcuts: Primary+R create · Escape cancel · Primary+S save"),
		).Gap(9).Expand(),
	})
}

func paintBloom(ctx context.Context, report *rosaline.TaskReporter) (*image.RGBA, error) {
	pixels := image.NewRGBA(image.Rect(0, 0, bloomWidth, bloomHeight))
	for y := range bloomHeight {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		for x := range bloomWidth {
			nx := (float64(x) - bloomWidth/2) / (bloomWidth / 2)
			ny := (float64(y) - bloomHeight/2) / (bloomHeight / 2)
			radius := math.Hypot(nx, ny)
			angle := math.Atan2(ny, nx)
			petals := 0.5 + 0.5*math.Cos(8*angle-13*radius)
			rose := math.Exp(-2.6*radius*radius) * (0.45 + 0.55*petals)
			glow := math.Exp(-12 * math.Pow(radius-0.42, 2))
			stars := 0.5 + 0.5*math.Sin(float64(x)*0.08+math.Sin(float64(y)*0.05))

			red := 24 + 210*rose + 26*glow
			green := 14 + 65*rose + 52*glow + 8*stars
			blue := 34 + 122*rose + 76*glow + 10*stars
			pixels.SetRGBA(x, y, color.RGBA{
				R: uint8(min(255, red)),
				G: uint8(min(255, green)),
				B: uint8(min(255, blue)),
				A: 255,
			})
		}

		percent := 100 * float64(y+1) / bloomHeight
		if !report.Report(percent, fmt.Sprintf("Painting light row %d of %d", y+1, bloomHeight)) {
			return nil, ctx.Err()
		}
		time.Sleep(2 * time.Millisecond) // Makes progress easy to see in the demo.
	}
	return pixels, nil
}
