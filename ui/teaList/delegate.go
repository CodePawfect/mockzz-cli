package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = true

	// Optionally customize key bindings
	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{}
	}
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{}
	}

	return d
}
