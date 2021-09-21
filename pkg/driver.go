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

import "github.com/kevinconway/remouseable/pkg/internal/robotgo"

// RobotgoDriver implements Driver using the robotgo cgo library.
type RobotgoDriver struct{}

// GetSize returns the width and height of the host screen.
func (*RobotgoDriver) GetSize() (int, int, error) {
	width, height := robotgo.GetScreenSize()
	return width, height, nil
}

// Click and hold the mouse button down.
func (*RobotgoDriver) Click(eraser bool) error {
	if eraser{
		robotgo.MouseToggle("down", "right")
	}else{
		robotgo.MouseToggle("down", "left")
	}
	return nil
}

// Unclick and release the mouse button.
func (*RobotgoDriver) Unclick(eraser bool) error {
	if eraser{
		robotgo.MouseToggle("up", "right")
	}else{
		robotgo.MouseToggle("up", "left")
	}
	return nil
}

// MoveMouse sets the mouse to a specified location.
func (*RobotgoDriver) MoveMouse(x int, y int) error {
	// Reversing the x/y due to robotgo seemingly having an opposite
	// x/y concept as the typical event source of evdev, etc.
	robotgo.MoveMouse(x, y)
	return nil
}

// DragMouse sets the mouse to a specified location while dragging a screen element.
func (*RobotgoDriver) DragMouse(x int, y int) error {
	// Reversing the x/y due to robotgo seemingly having an opposite
	// x/y concept as the typical event source of evdev, etc.
	robotgo.DragMouse(x, y)
	return nil
}
