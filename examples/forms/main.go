// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"strings"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	var name string
	var email string
	var password string
	var about string
	var newsletter bool

	submit := func() {
		if strings.TrimSpace(name) == "" {
			rosaline.Message("Missing name", "Please enter your name.")
			return
		}
		if strings.TrimSpace(email) == "" {
			rosaline.Message("Missing email", "Please enter your email address.")
			return
		}

		message := "Welcome, " + name + "!"
		if newsletter {
			message += "\n\nYou chose to receive project news."
		}
		rosaline.Message("Form submitted", message)
	}

	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Form",
		Width:  520,
		Height: 610,
		Content: rosaline.Column(
			rosaline.Label("Create your profile").Color(rosaline.Rose),
			rosaline.Label("Name"),
			rosaline.TextBox(&name).
				Placeholder("Your name").
				Focus(),
			rosaline.Label("Email"),
			rosaline.TextBox(&email).
				Placeholder("you@example.com"),
			rosaline.Label("Password"),
			rosaline.TextBox(&password).
				Placeholder("Choose a password").
				Password().
				OnSubmit(func(string) { submit() }),
			rosaline.Label("About you"),
			rosaline.TextArea(&about).Size(44, 6),
			rosaline.CheckBox("Send me Rosaline project news", &newsletter),
			rosaline.Button("Create profile", submit).Primary(),
		).Gap(10),
	})
}
