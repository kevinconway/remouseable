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

// EvdevStateMachine converts and EvdevIterator into significant state events.
type EvdevStateMachine struct {
	Iterator          EvdevIterator
	PressureThreshold int
	x                 int
	xChanged          bool
	y                 int
	yChanged          bool
	clicked           bool
	eraser            bool
	current           StateChange
	PenFunc           string
	EraserFunc        string
}

// next pushes the state machine one step. The return value is whether or not
// a new state was achieved in the step.
func (it *EvdevStateMachine) next(raw EvdevEvent) bool {
	if raw.Type == EV_KEY {
		if raw.Code == 321 {
			it.eraser = true
		} else if raw.Code == 320 {
			it.eraser = false
		}
	}
	if raw.Type != EV_ABS {
		return false
	}
	switch raw.Code {
	case ABS_X:
		it.x = int(raw.Value)
		it.xChanged = true
	case ABS_Y:
		it.y = int(raw.Value)
		it.yChanged = true
	case ABS_PRESSURE:
		if int(raw.Value) > it.PressureThreshold && !it.clicked {
			it.clicked = true
			it.current = &StateChangeClick{}
			return true
		}
		if int(raw.Value) < it.PressureThreshold && it.clicked {
			it.clicked = false
			it.current = &StateChangeUnclick{}
			return true
		}
	default:
	}
	if it.xChanged && it.yChanged {
		it.xChanged = false
		it.yChanged = false
		it.current = &StateChangeMove{X: it.x, Y: it.y}
		return true
	}
	return false
}

// Next consumes from the raw event iterator until a new state is achieved.
func (it *EvdevStateMachine) Next() bool {
	for it.Iterator.Next() {
		raw := it.Iterator.Current()
		if it.next(raw) {
			return true
		}
	}
	return false
}

// Current returns the iterator value.
func (it *EvdevStateMachine) Current() StateChange {
	return it.current
}

// Returns if eraser
func (it *EvdevStateMachine) Eraser() bool {
	return it.eraser
}

// Returns Pen function
func (it *EvdevStateMachine) GetPenFunc() string {
	return it.PenFunc
}

// Returns Eraser function
func (it *EvdevStateMachine) GetEraserFunc() string {
	return it.EraserFunc
}

// Close the underlying source and return any errors.
func (it *EvdevStateMachine) Close() error {
	return it.Iterator.Close()
}

type DraggingEvdevStateMachine struct {
	*EvdevStateMachine
}

// next pushes the state machine one step. The return value is whether or not
// a new state was achieved in the step.
func (it *DraggingEvdevStateMachine) next(raw EvdevEvent) bool {
	if ok := it.EvdevStateMachine.next(raw); !ok {
		return false
	}
	switch ev := it.current.(type) {
	case *StateChangeMove:
		if it.clicked {
			it.current = &StateChangeDrag{X: ev.X, Y: ev.Y}
		}
	default:
		break
	}
	return true
}

// Next consumes from the raw event iterator until a new state is achieved.
func (it *DraggingEvdevStateMachine) Next() bool {
	for it.Iterator.Next() {
		raw := it.Iterator.Current()
		if it.next(raw) {
			return true
		}
	}
	return false
}
