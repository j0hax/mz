package statusbar

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Blur implements tview.Primitive.
func (s *StatusBar) Blur() {
	s.field.Blur()
}

// Draw implements tview.Primitive.
func (s *StatusBar) Draw(screen tcell.Screen) {
	s.field.Draw(screen)
}

// Focus implements tview.Primitive.
func (s *StatusBar) Focus(delegate func(p tview.Primitive)) {
	s.field.Focus(delegate)
}

// GetRect implements tview.Primitive.
func (s *StatusBar) GetRect() (int, int, int, int) {
	return s.field.GetRect()
}

// HasFocus implements tview.Primitive.
func (s *StatusBar) HasFocus() bool {
	return s.field.HasFocus()
}

// InputHandler implements tview.Primitive.
func (s *StatusBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return s.field.InputHandler()
}

// MouseHandler implements tview.Primitive.
func (s *StatusBar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return s.field.MouseHandler()
}

// SetRect implements tview.Primitive.
func (s *StatusBar) SetRect(x int, y int, width int, height int) {
	s.field.SetRect(x, y, width, height)
}
