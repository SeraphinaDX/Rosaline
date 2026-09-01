// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	count := rosaline.NewState(0)

	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Counter",
		Width:  420,
		Height: 240,
		Content: rosaline.Column(
			rosaline.Label("A tiny counter"),
			rosaline.LabelFunc(func() string {
				return fmt.Sprintf("Count: %d", count.Get())
			}).Color(rosaline.Rose),
			rosaline.Row(
				rosaline.Button("Subtract", func() {
					count.Update(func(n int) int { return n - 1 })
				}),
				rosaline.Button("Add", func() {
					count.Update(func(n int) int { return n + 1 })
				}).Primary(),
			),
		).Gap(14),
	})
}
