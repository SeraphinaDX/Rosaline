// SPDX-License-Identifier: LGPL-3.0-or-later

// Package rosaline makes small graphical Go applications easy to build.
//
// The package deliberately keeps window setup, the event loop, layout, and the
// platform backend out of beginner programs. A complete application can be as
// small as:
//
//	rosaline.Run(rosaline.Label("Hello, world!"))
//
// Use RunApp when you want to set the title, initial window size, or theme.
// TextBox, TextArea, and CheckBox bind directly to ordinary Go variables.
package rosaline
