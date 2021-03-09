// This file is part of remouseable.
//
// remouseable is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as published
// by the Free Software Foundation.
//
// remouseable is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with remouseable.  If not, see <https://www.gnu.org/licenses/>.

package remouseable

import "time"

// EvdevEvent is a container type for raw evdev events. It is structured
// such that it can be used with the encoding/binary package to unmarshal
// raw events from the binary format.
type EvdevEvent struct {
	// Time of the event
	Time time.Time
	// Type is one of the EV_* named constants in evdevcodes.go
	Type uint16
	// Code is the relevant event constant from evdevcodes.go
	Code uint16
	// Numeric value of the event. Dependent on the event type.
	Value int32
}

// EvdevIterator represents a source of EvdevEvent instances. This is generally
// sourced from a /dev/input/event* file but may have alternative
// implementations for special use cases. For example, alternatives may be
// streaming events from a network source or replaying of static data for
// testing.
type EvdevIterator interface {
	// Next progresses the iterator. It returns false when there are no more
	// elements to iterate or when the iterator encountered an error.
	Next() bool
	// Current returns the active element of the iterator. This should only be
	// called if Next() returned a true.
	Current() EvdevEvent
	// Close must be called before discarding the iterator. If the iterator
	// exited cleanly then the error is nil. The error is non-nil if either the
	// iterator encountered an internal error and stopped early or if it failed
	// to close.
	Close() error
}

const (
	// StateChangeMove represents a move of the x and y for the mouse.
	ChangeTypeMove = "MOVE"
	// StateChangeDrag represents a move of the x and y for the mouse when clicked.
	ChangeTypeDrag = "DRAG"
	// ChangeTypeClick indicates that the stylus is touching the tablet.
	ChangeTypeClick = "CLICK"
	// ChangeTypeUnclick indicates the stylus is no longer touching the tablet.
	ChangeTypeUnclick = "UNCLICK"
)

// StateChangeMove contains mouse movement data.
type StateChangeMove struct {
	X int
	Y int
}

// Type returns the specific change type.
func (*StateChangeMove) Type() string {
	return ChangeTypeMove
}

// StateChangeDrag contains mouse movement data when clicked.
type StateChangeDrag struct {
	X int
	Y int
}

// Type returns the specific change type.
func (*StateChangeDrag) Type() string {
	return ChangeTypeDrag
}

// StateChangeClick contains mouse click data.
type StateChangeClick struct{}

// Type returns the specific change type.
func (*StateChangeClick) Type() string {
	return ChangeTypeClick
}

// StateChangeUnclick contains mouse click data.
type StateChangeUnclick struct{}

// Type returns the specific change type.
func (*StateChangeUnclick) Type() string {
	return ChangeTypeUnclick
}

// StateChange is a type for switching on the kind of change in order to convert
// the generic change type into a specific change type.
type StateChange interface {
	Type() string
}

// StateMachine is a specialized version of the EvdevIterator that only emits
// events on significant changes of the machine.
type StateMachine interface {
	// Next progresses the iterator. It returns false when there are no more
	// elements to iterate or when the iterator encountered an error.
	Next() bool
	// Current returns the active element of the iterator. This should only be
	// called if Next() returned a true.
	Current() StateChange
	// Close must be called before discarding the iterator. If the iterator
	// exited cleanly then the error is nil. The error is non-nil if either the
	// iterator encountered an internal error and stopped early or if it failed
	// to close.
	Close() error
}

// PositionScaler implements scaling rules for converting x/y coordinates
// between differently sized screens.
type PositionScaler interface {
	ScalePosition(x int, y int) (int, int)
}

// Driver is used to control a host system.
type Driver interface {
	MoveMouse(x int, y int) error
	DragMouse(x int, y int) error
	Click() error
	Unclick() error
	GetSize() (width int, height int, err error)
}
