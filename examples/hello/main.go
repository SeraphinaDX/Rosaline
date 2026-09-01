// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	rosaline.Run(
		rosaline.Column(
			rosaline.Label("Hello from Rosaline!"),
			rosaline.Label("You just made a graphical Go application."),
			rosaline.Button("Say hello", func() {
				rosaline.Message("Rosaline", "Hello, Britney!")
			}).Primary(),
		).Gap(12),
	)
}
