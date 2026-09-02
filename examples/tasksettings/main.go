// SPDX-License-Identifier: LGPL-3.0-or-later

// The tasksettings example combines Rosaline's everyday controls in one small
// application. Run it with: go run ./examples/tasksettings
package main

import (
	"fmt"
	"strings"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	taskName := "Write the Rosaline guide"
	category := "Documentation"
	priority := "normal"
	completion := 35.0
	status := "Ready"

	categoryBox := rosaline.ComboBox(
		&category,
		"Documentation",
		"Development",
		"Design",
		"Testing",
	).Width(28).OnChange(func(selected string) {
		status = "Category changed to " + selected
	})

	priorityGroup := rosaline.RadioGroup(
		&priority,
		rosaline.Choice("Low", "low"),
		rosaline.Choice("Normal", "normal"),
		rosaline.Choice("High", "high"),
	).Horizontal().OnChange(func(selected string) {
		status = "Priority changed to " + selected
	})

	completionSlider := rosaline.Slider(&completion, 0, 100).
		Step(5).
		OnChange(func(value float64) {
			status = fmt.Sprintf("Completion changed to %.0f%%", value)
		})
	completionBar := rosaline.ProgressBar(&completion)

	// A busy bar is useful when an operation has no known percentage. Starting
	// it in a paused state lets the buttons below demonstrate Start and Stop.
	busyBar := rosaline.ProgressBar(nil).Busy().Stop()

	usePlanningChoices := func() {
		categoryBox.SetOptions("Planning", "Research", "Documentation")
		priorityGroup.SetChoices(
			rosaline.Choice("Someday", "someday"),
			rosaline.Choice("This week", "week"),
			rosaline.Choice("Today", "today"),
		)
		status = "Loaded planning choices"
	}

	useWorkChoices := func() {
		categoryBox.SetOptions("Documentation", "Development", "Design", "Testing")
		priorityGroup.SetChoices(
			rosaline.Choice("Low", "low"),
			rosaline.Choice("Normal", "normal"),
			rosaline.Choice("High", "high"),
		)
		status = "Loaded work choices"
	}

	save := func() {
		if strings.TrimSpace(taskName) == "" {
			rosaline.Message("Missing task name", "Please enter a name for the task.")
			return
		}
		status = "Task settings saved"
		rosaline.Message(
			"Task saved",
			fmt.Sprintf(
				"%s\n\nCategory: %s\nPriority: %s\nCompletion: %.0f%%",
				taskName, category, priority, completion,
			),
		)
	}

	rosaline.RunApp(rosaline.App{
		Title:  "Rosaline Task Settings",
		Width:  700,
		Height: 650,
		Content: rosaline.Column(
			rosaline.Label("Task settings").Color(rosaline.Rose),
			rosaline.Label("Everyday controls stay connected to ordinary Go variables."),

			rosaline.Label("Task name"),
			rosaline.TextBox(&taskName).Width(42).Focus(),

			rosaline.Label("Category"),
			categoryBox,

			rosaline.Label("Priority"),
			priorityGroup,

			rosaline.LabelFunc(func() string {
				return fmt.Sprintf("Completion: %.0f%%", completion)
			}),
			completionSlider,
			completionBar,

			rosaline.Label("Unknown-duration work"),
			busyBar,
			rosaline.Row(
				rosaline.Button("Start busy bar", func() {
					busyBar.Start()
					status = "Busy animation started"
				}),
				rosaline.Button("Stop busy bar", func() {
					busyBar.Stop()
					status = "Busy animation stopped"
				}),
			),

			rosaline.Row(
				rosaline.Button("Planning choices", usePlanningChoices),
				rosaline.Button("Work choices", useWorkChoices),
				rosaline.Button("Save task", save).Primary(),
			),
			rosaline.LabelFunc(func() string { return status }),
		).Gap(9),
	})
}
