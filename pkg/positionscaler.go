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

const (
	// DefaultTabletHeight is the standard max height value that can be measured
	// on a remarkable tablet. Height is the measure of the maximum x coordinate
	// value of the tablet screen. The tablet screen is actually oriented
	// horizontally with the origin in the upper left corner when the top of the
	// device (power button) is on the right. Note that this magic number value
	// is not documented anywhere but was discovered by printing out the X value
	// events from evdev and using the stylus to draw a line to the edge of the
	// device screen.
	DefaultTabletHeight = 15725
	// DefaultTabletWidth is the standard max width value that can be measured
	// on a remarkable tablet. Width is the measure of the maximum y coordinate
	// value of the tablet screen. The tablet screen is actually oriented
	// horizontally with the origin in the upper left corner when the top of the
	// device (power button) is on the right. Note that this magic number value
	// is not documented anywhere but was discovered by printing out the X value
	// events from evdev and using the stylus to draw a line to the edge of the
	// device screen.
	DefaultTabletWidth = 20967
)

// RightPositionScaler converts points from a right-horizontally positioned
// tablet to a differently sized screen.
type RightPositionScaler struct {
	TabletWidth  int
	TabletHeight int
	ScreenWidth  int
	ScreenHeight int
	ScreenOffsetX int
	ScreenOffsetY int
}

// ScalePosition resolves based on a hoizontal position of the tablet.
func (s *RightPositionScaler) ScalePosition(x int, y int) (int, int) {
	// A horizontal orientation of the tablet with the top on the right is
	// actually the natural orientation of the screen on the device. This
	// orientation places the origin at the upper left corner which matches the
	// origin of the host screen. Because this orientation is the most "natural"
	// it has the simplest scaling policy of directly translating x and y values
	// using the proportional screen size as a scaling factor.
	scaleX := float64(s.ScreenWidth) / float64(s.TabletWidth)
	scaleY := float64(s.ScreenHeight) / float64(s.TabletHeight)
	// Apply the scaling factor to the points.
	return s.ScreenOffsetX + int(scaleX * float64(x)), s.ScreenOffsetY + int(scaleY * float64(y))
}

// LeftPositionScaler converts points from a left-horizontally positioned
// tablet to a differently sized screen.
type LeftPositionScaler struct {
	TabletWidth  int
	TabletHeight int
	ScreenWidth  int
	ScreenHeight int
	ScreenOffsetX int
	ScreenOffsetY int
}

// ScalePosition resolves based on a hoizontal position of the tablet.
func (s *LeftPositionScaler) ScalePosition(x int, y int) (int, int) {
	// Because the tablet is oriented opposite a typical screen we need
	// to adjust the x and y values by translating them. For example, the tablet
	// coordinate (0,0) is the bottom right corner of the tablet when oriented
	// left. However, the equivalent screen coordinate would be
	// (ScreenHeight, ScreenWidth). To resolve this conflice we subtract the
	// x and y values of the tablet from the maximum values so that (0,0)
	// becomes (max, max) and (max, max) becomes (0, 0).
	x, y = s.TabletWidth-x, s.TabletHeight-y
	scaleX := float64(s.ScreenWidth) / float64(s.TabletWidth)
	scaleY := float64(s.ScreenHeight) / float64(s.TabletHeight)
	// Apply the scaling factor to the points.
	return s.ScreenOffsetX + int(scaleX * float64(x)), s.ScreenOffsetY + int(scaleY * float64(y))
}

// VerticalPositionScaler converts points from a vertically positioned
// tablet to a differently sized screen.
type VerticalPositionScaler struct {
	TabletWidth  int
	TabletHeight int
	ScreenWidth  int
	ScreenHeight int
	ScreenOffsetX int
	ScreenOffsetY int
}

// ScalePosition resolves based on a vertical position of the tablet.
func (s *VerticalPositionScaler) ScalePosition(x int, y int) (int, int) {
	// A vertical position of the tablet is the "natural" position of the device
	// with the buttons on bottom and power button on top. However, this is not
	// the natural orientation of the tablet screen which is actually oriented
	// horizontally with buttons on the left and power button on the right.
	//
	// Because of this, x traversal on the tablet becomes y traversal on the
	// screen and vice versa. Additionally, the 90 degree rotation requires that
	// we translate coordinate values to account for the different origin
	// positions. To translate we subtract the x value from the tablet width
	// so that the tablet (0,0) which is the bottom right corner becomes the
	// screen (max, 0) which is also the bottom right corner. Likewise the tablet
	// (max, 0) value becomes the screen (0,0) value which represents the upper
	// left corner. The tablet y values are not adjusted as they are directly equal to
	// the corresponding screen x values without additional translation.
	x, y = y, s.TabletWidth-x
	scaleX := float64(s.ScreenWidth) / float64(s.TabletHeight)
	scaleY := float64(s.ScreenHeight) / float64(s.TabletWidth)
	return s.ScreenOffsetX + int(scaleX * float64(x)), s.ScreenOffsetY + int(scaleY * float64(y))
}
