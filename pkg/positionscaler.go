package remouse

const (
	// DefaultTabletHeight is the standard max height value that can be
	// measured on a remarkable tablet.
	DefaultTabletHeight = 20951
	// DefaultTabletWidth is the standard max width value that can be measured
	// on a remarkable tablet.
	DefaultTabletWidth = 15725
)

// RightPositionScaler converts points from a right-horizontally positioned
// tablet to a differently sized screen.
type RightPositionScaler struct {
	TabletWidth  int
	TabletHeight int
	ScreenWidth  int
	ScreenHeight int
}

// ScalePosition resolves based on a hoizontal position of the tablet.
func (s *RightPositionScaler) ScalePosition(x int, y int) (int, int) {
	// Reverse the coordinates to account for the horizontal position.
	x, y = y, x
	tabHeight, tabWidth := s.TabletWidth, s.TabletHeight
	// Determine the scaling factor by calculating a proportion of the
	// screen height and width that will be applied to the x/y values.
	scaleHeight := float64(s.ScreenHeight) / float64(tabHeight)
	scaleWidth := float64(s.ScreenWidth) / float64(tabWidth)
	// Apply the scaling factor to the points.
	return int(scaleHeight * float64(x)), int(scaleWidth * float64(y))
}

// LeftPositionScaler converts points from a left-horizontally positioned
// tablet to a differently sized screen.
type LeftPositionScaler struct {
	TabletWidth  int
	TabletHeight int
	ScreenWidth  int
	ScreenHeight int
}

// ScalePosition resolves based on a hoizontal position of the tablet.
func (s *LeftPositionScaler) ScalePosition(x int, y int) (int, int) {
	// Reverse the coordinates and adjust for left to account for the
	// horizontal position.
	x, y = s.TabletWidth-y, s.TabletHeight-x
	tabHeight, tabWidth := s.TabletWidth, s.TabletHeight
	// Determine the scaling factor by calculating a proportion of the
	// screen height and width that will be applied to the x/y values.
	scaleHeight := float64(s.ScreenHeight) / float64(tabHeight)
	scaleWidth := float64(s.ScreenWidth) / float64(tabWidth)
	// Apply the scaling factor to the points.
	return int(scaleHeight * float64(x)), int(scaleWidth * float64(y))
}

// VerticalPositionScaler converts points from a vertically positioned
// tablet to a differently sized screen.
type VerticalPositionScaler struct {
	TabletWidth  int
	TabletHeight int
	ScreenWidth  int
	ScreenHeight int
}

// ScalePosition resolves based on a vertical position of the tablet.
func (s *VerticalPositionScaler) ScalePosition(x int, y int) (int, int) {
	x = s.TabletHeight - x
	scaleHeight := float64(s.ScreenHeight) / float64(s.TabletHeight)
	scaleWidth := float64(s.ScreenWidth) / float64(s.TabletWidth)
	return int(scaleHeight * float64(x)), int(scaleWidth * float64(y))
}
