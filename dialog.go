// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// Message displays a simple informational dialog.
func Message(title, text string) {
	tk.MessageBox(tk.Title(title), tk.Msg(text), tk.Type("ok"), tk.Icon("info"))
}
