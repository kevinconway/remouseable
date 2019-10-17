package remouse

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
	case ChangeTypeClick:
		if err := r.Driver.Click(); err != nil {
			r.err = err
			return false
		}
		return true
	case ChangeTypeUnclick:
		if err := r.Driver.Unclick(); err != nil {
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
