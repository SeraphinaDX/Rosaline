// SPDX-License-Identifier: LGPL-3.0-or-later

// The calculator example combines Rosaline's layout and presentation tools.
// Run it with: go run ./examples/calculator
package main

import (
	"math"
	"strconv"
	"strings"

	"github.com/SeraphinaDX/Rosaline"
)

type calculator struct {
	display    string
	expression string
	status     string
	left       float64
	operator   string
	replace    bool
	error      bool
}

func newCalculator() *calculator {
	return &calculator{display: "0", status: "Ready", replace: true}
}

func (c *calculator) digit(digit string) {
	if c.error || c.replace {
		c.display = digit
		c.replace = false
		c.error = false
	} else if len(c.display) < 16 {
		if c.display == "0" {
			c.display = digit
		} else {
			c.display += digit
		}
	}
	c.status = "Entering a number"
}

func (c *calculator) decimal() {
	if c.error || c.replace {
		c.display = "0."
		c.replace = false
		c.error = false
	} else if !strings.Contains(c.display, ".") {
		c.display += "."
	}
	c.status = "Decimal number"
}

func (c *calculator) choose(operator string) {
	if c.error {
		c.clear()
	}
	if c.operator != "" && !c.replace {
		if !c.solve() {
			return
		}
	}
	c.left = c.number()
	c.operator = operator
	c.expression = c.display + " " + operator
	c.replace = true
	c.status = "Choose the next number"
}

func (c *calculator) equals() {
	if c.operator == "" || c.error {
		return
	}
	right := c.number()
	left := c.left
	operator := c.operator
	if !c.solve() {
		return
	}
	c.expression = formatNumber(left) + " " + operator + " " + formatNumber(right) + " ="
	c.status = "Result"
}

func (c *calculator) solve() bool {
	right := c.number()
	result := c.left
	switch c.operator {
	case "+":
		result += right
	case "−":
		result -= right
	case "×":
		result *= right
	case "÷":
		if right == 0 {
			c.display = "Cannot divide by zero"
			c.expression = ""
			c.status = "Please clear and try again"
			c.operator = ""
			c.replace = true
			c.error = true
			return false
		}
		result /= right
	}
	c.display = formatNumber(result)
	c.left = result
	c.operator = ""
	c.replace = true
	return true
}

func (c *calculator) percent() {
	if c.error {
		return
	}
	c.display = formatNumber(c.number() / 100)
	c.replace = true
	c.status = "Converted to a percentage"
}

func (c *calculator) sign() {
	if c.error || c.display == "0" {
		return
	}
	if strings.HasPrefix(c.display, "-") {
		c.display = strings.TrimPrefix(c.display, "-")
	} else {
		c.display = "-" + c.display
	}
	c.status = "Changed the sign"
}

func (c *calculator) backspace() {
	if c.error || c.replace {
		return
	}
	c.display = strings.TrimSuffix(c.display, c.display[len(c.display)-1:])
	if c.display == "" || c.display == "-" {
		c.display = "0"
		c.replace = true
	}
	c.status = "Removed one digit"
}

func (c *calculator) clear() {
	c.display = "0"
	c.expression = ""
	c.status = "Ready"
	c.left = 0
	c.operator = ""
	c.replace = true
	c.error = false
}

func (c *calculator) number() float64 {
	value, _ := strconv.ParseFloat(c.display, 64)
	return value
}

func formatNumber(value float64) string {
	if math.Abs(value) < 1e-12 {
		value = 0
	}
	return strconv.FormatFloat(value, 'g', 12, 64)
}

func main() {
	model := newCalculator()
	theme := rosaline.DefaultTheme
	theme.Background = rosaline.Hex("#170e1c")
	theme.Surface = rosaline.Hex("#2a1830")
	theme.Primary = rosaline.Hex("#ef5da8")
	theme.Text = rosaline.Hex("#fff2fa")
	theme.Muted = rosaline.Hex("#c8a6bd")
	theme.Border = rosaline.Hex("#74415f")
	theme.Danger = rosaline.Hex("#ff6f91")
	theme.Success = rosaline.Hex("#79d4aa")

	button := func(text string, action func()) *rosaline.ButtonWidget {
		return rosaline.Button(text, action)
	}
	digit := func(value string) *rosaline.ButtonWidget {
		return button(value, func() { model.digit(value) })
	}
	operation := func(label, value string) *rosaline.ButtonWidget {
		return button(label, func() { model.choose(value) }).Primary()
	}

	keypad := rosaline.Grid(4,
		button("AC", model.clear),
		button("±", model.sign),
		button("%", model.percent),
		operation("÷", "÷"),
		digit("7"), digit("8"), digit("9"), operation("×", "×"),
		digit("4"), digit("5"), digit("6"), operation("−", "−"),
		digit("1"), digit("2"), digit("3"), operation("+", "+"),
		digit("0"), button(".", model.decimal), button("⌫", model.backspace),
		button("=", model.equals).Primary(),
	).Gap(9).Expand()

	panel := rosaline.Card(
		rosaline.Column(
			rosaline.Row(
				rosaline.Label("ROSALINE").Bold().Color(theme.Primary),
				rosaline.Spring(),
				rosaline.Label("CALCULATOR").Color(theme.Muted),
			),
			rosaline.LabelFunc(func() string { return model.expression }).
				TextAlign(rosaline.AlignEnd).
				Color(theme.Muted),
			rosaline.MinSize(
				rosaline.LabelFunc(func() string { return model.display }).
					FontSize(34).
					Bold().
					TextAlign(rosaline.AlignEnd),
				320, 60,
			),
			rosaline.Separator(),
			keypad,
			rosaline.Row(
				rosaline.LabelFunc(func() string { return model.status }).Color(theme.Muted),
				rosaline.Spring(),
				rosaline.Label("Enter = · Esc clear").Color(theme.Muted),
			),
		).Gap(12).Expand(),
	).Padding(20)

	background := rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		c.Clear(theme.Background)
		c.FillCircle(75, 90, 120, rosaline.Hex("#3d1735"))
		c.FillCircle(470, 570, 180, rosaline.Hex("#30183d"))
		for y := 30.0; y < 680; y += 55 {
			for x := 25.0; x < 540; x += 55 {
				c.FillCircle(x, y, 2, rosaline.Hex("#8f5275"))
			}
		}
	}).Size(540, 680).Background(theme.Background).Expand()

	onKey := func(event rosaline.KeyEvent) {
		if event.Primary || event.Control || event.Alt {
			return
		}
		switch event.Text {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			model.digit(event.Text)
		case ".", ",":
			model.decimal()
		case "+":
			model.choose("+")
		case "-":
			model.choose("−")
		case "*":
			model.choose("×")
		case "/":
			model.choose("÷")
		case "%":
			model.percent()
		}
	}

	rosaline.RunApp(rosaline.App{
		Title:     "Rosaline Calculator",
		Width:     560,
		Height:    700,
		Padding:   8,
		Theme:     theme,
		OnKeyDown: onKey,
		Shortcuts: rosaline.Shortcuts(
			rosaline.Shortcut("Enter", model.equals),
			rosaline.Shortcut("Escape", model.clear),
			rosaline.Shortcut("Backspace", model.backspace),
		),
		Content: rosaline.Stack(
			background,
			rosaline.Center(rosaline.Size(panel, 390, 570)),
		).Expand(),
	})
}
