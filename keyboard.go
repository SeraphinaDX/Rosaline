// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"runtime"
	"strings"
	"unicode/utf8"

	tk "modernc.org/tk9.0"
)

// Key identifies a keyboard key without exposing backend-specific names.
// Printable keys use a lowercase value; Text preserves the text actually
// produced by the event.
type Key string

const (
	KeyUnknown   Key = ""
	KeyBackspace Key = "Backspace"
	KeyTab       Key = "Tab"
	KeyEnter     Key = "Enter"
	KeyEscape    Key = "Escape"
	KeySpace     Key = "Space"
	KeyDelete    Key = "Delete"
	KeyInsert    Key = "Insert"
	KeyHome      Key = "Home"
	KeyEnd       Key = "End"
	KeyPageUp    Key = "PageUp"
	KeyPageDown  Key = "PageDown"
	KeyLeft      Key = "Left"
	KeyUp        Key = "Up"
	KeyRight     Key = "Right"
	KeyDown      Key = "Down"
	KeyShift     Key = "Shift"
	KeyControl   Key = "Control"
	KeyAlt       Key = "Alt"
	KeySuper     Key = "Super"
	KeyCapsLock  Key = "CapsLock"
	KeyF1        Key = "F1"
	KeyF2        Key = "F2"
	KeyF3        Key = "F3"
	KeyF4        Key = "F4"
	KeyF5        Key = "F5"
	KeyF6        Key = "F6"
	KeyF7        Key = "F7"
	KeyF8        Key = "F8"
	KeyF9        Key = "F9"
	KeyF10       Key = "F10"
	KeyF11       Key = "F11"
	KeyF12       Key = "F12"
)

// String returns the friendly name of a key.
func (k Key) String() string { return string(k) }

// KeyEvent describes a key press or release. Text contains produced text such
// as "a" or "A" and is empty for keys such as Left and Escape.
type KeyEvent struct {
	Key     Key
	Text    string
	Shift   bool
	Control bool
	Alt     bool
	Primary bool
}

// Is reports whether this event belongs to key.
func (e KeyEvent) Is(key Key) bool { return e.Key == key }

// KeyShortcut connects a key combination to an ordinary func(). Create one
// with Shortcut and place it in App.Shortcuts or WindowOptions.Shortcuts.
type KeyShortcut struct {
	keys    string
	onPress func()
}

// Shortcut creates a window shortcut such as "Primary+S", "Ctrl+Shift+S",
// or "Escape". Primary means Control on Linux and Windows and Command on
// macOS.
func Shortcut(keys string, onPress func()) KeyShortcut {
	if onPress == nil {
		onPress = func() {}
	}
	return KeyShortcut{keys: strings.TrimSpace(keys), onPress: onPress}
}

// Keys returns the human-readable key combination supplied to Shortcut.
func (s KeyShortcut) Keys() string { return s.keys }

// Shortcuts collects shortcut values without requiring a slice literal.
func Shortcuts(shortcuts ...KeyShortcut) []KeyShortcut {
	return append([]KeyShortcut(nil), shortcuts...)
}

func keyEvent(event *tk.Event) KeyEvent {
	return keyEventForOS(event, runtime.GOOS)
}

func keyEventForOS(event *tk.Event, goos string) KeyEvent {
	if event == nil {
		return KeyEvent{}
	}
	shift := event.State&tk.ModifierShift != 0
	control := event.State&tk.ModifierControl != 0
	altModifier := tk.ModifierAlt
	primary := control
	if goos == "darwin" {
		// Tk maps Command to Mod1 and Option to Mod2 on macOS.
		primary = event.State&tk.ModifierMod1 != 0
		altModifier = tk.ModifierMod2
	}
	return KeyEvent{
		Key:     keyFromKeysym(event.Keysym),
		Text:    event.Char,
		Shift:   shift,
		Control: control,
		Alt:     event.State&altModifier != 0,
		Primary: primary,
	}
}

func keyFromKeysym(keysym string) Key {
	switch keysym {
	case "BackSpace":
		return KeyBackspace
	case "Tab", "ISO_Left_Tab":
		return KeyTab
	case "Return", "KP_Enter":
		return KeyEnter
	case "Escape":
		return KeyEscape
	case "space":
		return KeySpace
	case "Delete", "KP_Delete":
		return KeyDelete
	case "Insert", "KP_Insert":
		return KeyInsert
	case "Home", "KP_Home":
		return KeyHome
	case "End", "KP_End":
		return KeyEnd
	case "Prior", "KP_Prior":
		return KeyPageUp
	case "Next", "KP_Next":
		return KeyPageDown
	case "Left", "KP_Left":
		return KeyLeft
	case "Up", "KP_Up":
		return KeyUp
	case "Right", "KP_Right":
		return KeyRight
	case "Down", "KP_Down":
		return KeyDown
	case "Shift_L", "Shift_R":
		return KeyShift
	case "Control_L", "Control_R":
		return KeyControl
	case "Alt_L", "Alt_R", "Option_L", "Option_R":
		return KeyAlt
	case "Meta_L", "Meta_R", "Super_L", "Super_R", "Win_L", "Win_R", "Command":
		return KeySuper
	case "Caps_Lock":
		return KeyCapsLock
	case "":
		return KeyUnknown
	}
	if utf8.RuneCountInString(keysym) == 1 {
		return Key(strings.ToLower(keysym))
	}
	return Key(keysym)
}

func mountKeyboard(ctx *mountContext, window *tk.Window, onKeyDown, onKeyUp func(KeyEvent), shortcuts []KeyShortcut) {
	if ctx == nil || window == nil {
		return
	}
	if onKeyDown != nil {
		tk.Bind(window, "<KeyPress>", tk.Command(func(event *tk.Event) {
			ctx.flush()
			onKeyDown(keyEvent(event))
			ctx.refresh()
		}))
	}
	if onKeyUp != nil {
		tk.Bind(window, "<KeyRelease>", tk.Command(func(event *tk.Event) {
			ctx.flush()
			onKeyUp(keyEvent(event))
			ctx.refresh()
		}))
	}
	for _, shortcut := range shortcuts {
		bindShortcut(ctx, window, shortcut.keys, shortcut.onPress)
	}
}

func bindShortcut(ctx *mountContext, window *tk.Window, keys string, onPress func()) bool {
	sequence, ok := shortcutSequence(keys)
	if !ok || ctx == nil || window == nil || onPress == nil {
		return false
	}
	tk.Bind(window, sequence, tk.Command(func(event *tk.Event) {
		ctx.flush()
		onPress()
		ctx.refresh()
		event.SetReturnCodeBreak()
	}))
	return true
}

func shortcutSequence(shortcut string) (string, bool) {
	return shortcutSequenceForOS(shortcut, runtime.GOOS)
}

func shortcutSequenceForOS(shortcut, goos string) (string, bool) {
	parts := strings.Split(shortcut, "+")
	if len(parts) == 0 {
		return "", false
	}
	for index := range parts {
		parts[index] = strings.TrimSpace(parts[index])
		if parts[index] == "" {
			return "", false
		}
	}

	modifiers := make([]string, 0, len(parts)-1)
	hasShift := false
	hasCommandModifier := false
	for _, part := range parts[:len(parts)-1] {
		switch strings.ToLower(part) {
		case "primary":
			hasCommandModifier = true
			if goos == "darwin" {
				modifiers = append(modifiers, "Command")
			} else {
				modifiers = append(modifiers, "Control")
			}
		case "ctrl", "control":
			hasCommandModifier = true
			modifiers = append(modifiers, "Control")
		case "shift":
			modifiers = append(modifiers, "Shift")
			hasShift = true
		case "alt":
			hasCommandModifier = true
			if goos == "darwin" {
				modifiers = append(modifiers, "Option")
			} else {
				modifiers = append(modifiers, "Alt")
			}
		case "option":
			if goos != "darwin" {
				return "", false
			}
			hasCommandModifier = true
			modifiers = append(modifiers, "Option")
		case "cmd", "command":
			if goos != "darwin" {
				return "", false
			}
			hasCommandModifier = true
			modifiers = append(modifiers, "Command")
		case "super":
			hasCommandModifier = true
			if goos == "darwin" {
				modifiers = append(modifiers, "Command")
			} else {
				modifiers = append(modifiers, "Mod4")
			}
		default:
			return "", false
		}
	}

	key, named := shortcutKeysym(parts[len(parts)-1])
	if key == "" {
		return "", false
	}
	if !named && !hasCommandModifier {
		// Bare and Shift-only printable keys would interfere with text entry.
		return "", false
	}
	if key == "Tab" && !hasCommandModifier {
		// Tab and Shift+Tab belong to Rosaline's focus traversal.
		return "", false
	}
	if utf8.RuneCountInString(key) == 1 {
		if hasShift {
			key = strings.ToUpper(key)
		} else {
			key = strings.ToLower(key)
		}
	}
	return "<" + strings.Join(append(modifiers, key), "-") + ">", true
}

func shortcutKeysym(key string) (keysym string, named bool) {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "enter", "return":
		return "Return", true
	case "esc", "escape":
		return "Escape", true
	case "space":
		return "space", true
	case "backspace":
		return "BackSpace", true
	case "delete", "del":
		return "Delete", true
	case "insert", "ins":
		return "Insert", true
	case "home":
		return "Home", true
	case "end":
		return "End", true
	case "pageup", "pgup":
		return "Prior", true
	case "pagedown", "pgdown", "pgdn":
		return "Next", true
	case "left", "arrowleft":
		return "Left", true
	case "up", "arrowup":
		return "Up", true
	case "right", "arrowright":
		return "Right", true
	case "down", "arrowdown":
		return "Down", true
	case "tab":
		return "Tab", true
	case "plus":
		return "plus", true
	case "minus":
		return "minus", true
	}
	trimmed := strings.TrimSpace(key)
	upper := strings.ToUpper(trimmed)
	if len(upper) >= 2 && upper[0] == 'F' {
		var number int
		for _, digit := range upper[1:] {
			if digit < '0' || digit > '9' {
				return "", false
			}
			number = number*10 + int(digit-'0')
		}
		if number >= 1 && number <= 24 {
			return upper, true
		}
	}
	if utf8.RuneCountInString(trimmed) == 1 {
		return trimmed, false
	}
	return "", false
}

func shortcutDisplay(shortcut string) string {
	return shortcutDisplayForOS(shortcut, runtime.GOOS)
}

func shortcutDisplayForOS(shortcut, goos string) string {
	parts := strings.Split(shortcut, "+")
	for index, part := range parts {
		switch strings.ToLower(strings.TrimSpace(part)) {
		case "primary":
			if goos == "darwin" {
				parts[index] = "Command"
			} else {
				parts[index] = "Ctrl"
			}
		case "alt":
			if goos == "darwin" {
				parts[index] = "Option"
			}
		case "super":
			if goos == "darwin" {
				parts[index] = "Command"
			}
		}
	}
	return strings.Join(parts, "+")
}
