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
// Canvas mouse callbacks make drawing programs interactive without exposing
// platform event types.
// Images, scroll areas, menus, and file dialogs provide the groundwork for
// complete desktop applications while ordinary file I/O remains normal Go.
// App-owned timers support delayed work, repeating updates, and canvas
// animation without exposing the private event loop.
// Paths, transforms, clipping, off-screen rendering, and PNG or AVIF export
// support complete graphical applications with one consistent drawing API.
// Tabs organize larger interfaces into composable pages. Lists provide native
// scrolling, keyboard selection, activation, and programmatic item updates.
// Tables display ordinary slices of strings under named columns with native
// scrolling, keyboard behavior, selection, and activation.
package rosaline
