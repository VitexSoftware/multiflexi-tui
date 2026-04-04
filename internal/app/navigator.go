package app

import tea "github.com/charmbracelet/bubbletea"

// ViewState captures a snapshot of a view for the navigation stack.
type ViewState struct {
	View     tea.Model
	MenuIdx  int // which menu item was active
}

// Navigator manages a stack of views for back-navigation.
type Navigator struct {
	stack []ViewState
}

// Push adds a new view to the stack.
func (n *Navigator) Push(s ViewState) {
	n.stack = append(n.stack, s)
}

// Pop removes and returns the top view. Returns zero value if empty.
func (n *Navigator) Pop() (ViewState, bool) {
	if len(n.stack) == 0 {
		return ViewState{}, false
	}
	top := n.stack[len(n.stack)-1]
	n.stack = n.stack[:len(n.stack)-1]
	return top, true
}

// Current returns the top of the stack without removing it.
func (n *Navigator) Current() (ViewState, bool) {
	if len(n.stack) == 0 {
		return ViewState{}, false
	}
	return n.stack[len(n.stack)-1], true
}

// Depth returns the number of items on the stack.
func (n *Navigator) Depth() int {
	return len(n.stack)
}

// Clear empties the stack.
func (n *Navigator) Clear() {
	n.stack = nil
}
