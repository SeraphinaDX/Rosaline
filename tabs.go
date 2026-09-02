// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import tk "modernc.org/tk9.0"

// TabPage is one named page inside Tabs.
type TabPage struct {
	title   string
	content Widget
}

// Tab creates one named tab page.
func Tab(title string, content Widget) *TabPage {
	if title == "" {
		title = "Tab"
	}
	return &TabPage{title: title, content: content}
}

// Title returns the text displayed on the tab.
func (t *TabPage) Title() string {
	if t == nil {
		return ""
	}
	return t.title
}

// TabsWidget displays one of several named pages at a time.
type TabsWidget struct {
	pages       []*TabPage
	selected    int
	expand      bool
	onChange    func(index int, title string)
	notebook    *tk.TNotebookWidget
	pageWindows []*tk.Window
	ctx         *mountContext
}

// Tabs creates a tabbed interface. The first non-nil page is selected by
// default.
func Tabs(pages ...*TabPage) *TabsWidget {
	clean := make([]*TabPage, 0, len(pages))
	for _, page := range pages {
		if page != nil {
			clean = append(clean, page)
		}
	}
	selected := -1
	if len(clean) != 0 {
		selected = 0
	}
	return &TabsWidget{pages: clean, selected: selected}
}

// Expand asks the tabs to use available layout space.
func (t *TabsWidget) Expand() *TabsWidget {
	t.expand = true
	return t
}

// OnChange runs after the user or application changes the selected tab.
func (t *TabsWidget) OnChange(handler func(index int, title string)) *TabsWidget {
	t.onChange = handler
	return t
}

// Select displays the page at index. Invalid indices have no effect. When the
// tabs are mounted, a changed selection runs OnChange.
func (t *TabsWidget) Select(index int) {
	if t == nil || index < 0 || index >= len(t.pages) || index == t.selected {
		return
	}
	t.selected = index
	if t.notebook != nil && index < len(t.pageWindows) {
		t.notebook.Select(t.pageWindows[index])
		t.notifyChange()
	}
}

// Selected returns the selected page index and title. ok is false when there
// are no pages.
func (t *TabsWidget) Selected() (index int, title string, ok bool) {
	if t == nil || t.selected < 0 || t.selected >= len(t.pages) {
		return -1, "", false
	}
	return t.selected, t.pages[t.selected].title, true
}

// Pages returns a copy of the configured tab pages.
func (t *TabsWidget) Pages() []*TabPage {
	if t == nil {
		return nil
	}
	return append([]*TabPage(nil), t.pages...)
}

func (t *TabsWidget) readSelection() bool {
	if t.notebook == nil {
		return false
	}
	selectedWindow := t.notebook.Select(nil)
	for index, window := range t.pageWindows {
		if window.String() == selectedWindow {
			changed := index != t.selected
			t.selected = index
			return changed
		}
	}
	return false
}

func (t *TabsWidget) notifyChange() {
	if t.ctx != nil {
		t.ctx.flush()
	}
	if t.onChange != nil {
		if index, title, ok := t.Selected(); ok {
			t.onChange(index, title)
		}
	}
	if t.ctx != nil {
		t.ctx.refresh()
	}
}

func (t *TabsWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	t.ctx = ctx
	t.notebook = parent.TNotebook(takeFocusOption(true))
	ctx.addFocusable(t.notebook.Window, false)

	for pageIndex, page := range t.pages {
		frame := t.notebook.Frame(
			tk.Background(ctx.theme.Background.String()),
			tk.Borderwidth(0),
		)
		content := page.content
		if content == nil {
			content = Label("This tab has no content.")
		}
		var mounted mountedWidget
		index := pageIndex
		ctx.withFocusCondition(func() bool {
			return t.selected == index
		}, func() {
			mounted = content.mount(ctx, frame.Window)
		})
		options := []tk.Opt{
			mounted.window,
			tk.Anchor("nw"),
			tk.Fill("both"),
			tk.Expand(true),
			tk.Padx(14),
			tk.Pady(14),
		}
		tk.Pack(options...)
		t.notebook.Add(frame, tk.Txt(page.title), tk.Padding("8 5"))
		t.pageWindows = append(t.pageWindows, frame.Window)
	}

	if t.selected >= 0 && t.selected < len(t.pageWindows) {
		t.notebook.Select(t.pageWindows[t.selected])
	}
	tk.Bind(t.notebook, "<<NotebookTabChanged>>", tk.Command(func() {
		if t.readSelection() {
			t.notifyChange()
		}
	}))

	return mountedWidget{window: t.notebook.Window, expandX: t.expand, expandY: t.expand}
}
