package remouse

import (
	"github.com/go-vgo/robotgo"
)

// RobotgoDriver implements Driver using the robotgo cgo library.
type RobotgoDriver struct{}

// GetSize returns the width and height of the host screen.
func (*RobotgoDriver) GetSize() (int, int, error) {
	width, height := robotgo.GetScreenSize()
	return width, height, nil
}

// Click and hold the mouse button down.
func (*RobotgoDriver) Click() error {
	robotgo.MouseToggle("down", "left")
	return nil
}

// Unclick and release the mouse button.
func (*RobotgoDriver) Unclick() error {
	robotgo.MouseToggle("up", "left")
	return nil
}

// MoveMouse sets the mouse to a specified location.
func (*RobotgoDriver) MoveMouse(x int, y int) error {
	// Reversing the x/y due to robotgo seemingly having an opposite
	// x/y concept as the typical event source of evdev, etc.
	robotgo.MoveMouse(x, y)
	return nil
}
