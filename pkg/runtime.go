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

import "fmt"

// Runtime binds the various domain elements into an application.
type Runtime struct {
	StateMachine   StateMachine
	PositionScaler PositionScaler
	Driver         Driver
	err            error
}

// Next executes one step of the runtime loop.
func (r *Runtime) Next() bool {
	if r.err != nil {
		// Attempt to prevent re-entry after an error is encountered.
		return false
	}
	if !r.StateMachine.Next() {
		// Stop iteration if the state machine has completed all iterations.
		return false
	}
	change := r.StateMachine.Current()
	switch change.Type() {
	case ChangeTypeMove:
		evt := change.(*StateChangeMove)
		if err := r.Driver.MoveMouse(r.PositionScaler.ScalePosition(evt.X, evt.Y)); err != nil {
			r.err = err
			return false
		}
		return true
	case ChangeTypeDrag:
		evt := change.(*StateChangeDrag)
		if err := r.Driver.DragMouse(r.PositionScaler.ScalePosition(evt.X, evt.Y)); err != nil {
			r.err = err
			return false
		}
		return true
	case ChangeTypeClick:
		if err := r.Driver.Click(r.StateMachine.Eraser()); err != nil {
			r.err = err
			return false
		}
		return true
	case ChangeTypeUnclick:
		if err := r.Driver.Unclick(r.StateMachine.Eraser()); err != nil {
			r.err = err
			return false
		}
		return true
	default:
		r.err = fmt.Errorf("encountered unhandled state machine event %s", change.Type())
		return false
	}
}

// Close the runtime and any internal resources.
func (r *Runtime) Close() error {
	err := r.StateMachine.Close()
	if r.err != nil {
		return r.err
	}
	return err
}
