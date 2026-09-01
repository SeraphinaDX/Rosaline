// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

// Theme contains semantic colors used by Rosaline widgets.
type Theme struct {
	Background Color
	Surface    Color
	Primary    Color
	Text       Color
	Muted      Color
	Border     Color
	Danger     Color
	Success    Color
}

// DefaultTheme is Rosaline's light rose theme.
var DefaultTheme = Theme{
	Background: Hex("#fff8fc"),
	Surface:    Hex("#ffffff"),
	Primary:    Hex("#c43f7a"),
	Text:       Hex("#2a1722"),
	Muted:      Hex("#7d6874"),
	Border:     Hex("#d9b8ca"),
	Danger:     Hex("#b4234d"),
	Success:    Hex("#267a50"),
}
