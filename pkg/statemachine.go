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

// EvdevStateMachine converts and EvdevIterator into significant state events.
type EvdevStateMachine struct {
	Iterator          EvdevIterator
	PressureThreshold int
	x                 int
	xChanged          bool
	y                 int
	yChanged          bool
	clicked           bool
	current           StateChange
}

// next pushes the state machine one step. The return value is whether or not
// a new state was achieved in the step.
func (it *EvdevStateMachine) next(raw EvdevEvent) bool {
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

type RateLimitStateMachine struct {
	Wrapped StateMachine
	Rate    time.Duration
	current StateChange
	last    time.Time
	now     func() time.Time
}

func NewRateLimitStateMachine(rate time.Duration, sm StateMachine) *RateLimitStateMachine {
	return &RateLimitStateMachine{
		Wrapped: sm,
		Rate:    rate,
		now:     time.Now,
	}
}

// Next progresses the iterator. It returns false when there are no more
// elements to iterate or when the iterator encountered an error.
func (m *RateLimitStateMachine) Next() bool {
	for {
		if ok := m.Wrapped.Next(); !ok {
			return false
		}
		now := time.Now()
		next := m.Wrapped.Current()
		if m.current == nil {
			m.last = now
			m.current = next
			return true
		}
		if next.Type() != m.current.Type() {
			time.Sleep(now.Sub(m.last))
			m.last = time.Now()
			m.current = next
			return true
		}
		if now.Sub(m.last) > m.Rate {
			m.last = now
			m.current = next
			return true
		}
	}
}

// Current returns the active element of the iterator. This should only be
// called if Next() returned a true.
func (m *RateLimitStateMachine) Current() StateChange {
	return m.current
}

// Close must be called before discarding the iterator. If the iterator
// exited cleanly then the error is nil. The error is non-nil if either the
// iterator encountered an internal error and stopped early or if it failed
// to close.
func (m *RateLimitStateMachine) Close() error {
	return m.Wrapped.Close()
}
