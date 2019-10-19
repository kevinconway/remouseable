package remouse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRightPositionScaler_ScalePosition(t *testing.T) {
	type fields struct {
		TabletWidth  int
		TabletHeight int
		ScreenWidth  int
		ScreenHeight int
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantX  int
		wantY  int
	}{
		{
			name: "square scale factor 1",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 100,
				ScreenWidth:  100,
				ScreenHeight: 100,
			},
			args: args{
				x: 50,
				y: 50,
			},
			wantX: 50,
			wantY: 50,
		},
		{
			name: "square scale factor 2",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 100,
				ScreenWidth:  200,
				ScreenHeight: 200,
			},
			args: args{
				x: 50,
				y: 50,
			},
			wantX: 100,
			wantY: 100,
		},
		{
			name: "non-square",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 200,
				ScreenWidth:  400,
				ScreenHeight: 200,
			},
			args: args{
				x: 50,
				y: 100,
			},
			wantX: 200,
			wantY: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RightPositionScaler{
				TabletWidth:  tt.fields.TabletWidth,
				TabletHeight: tt.fields.TabletHeight,
				ScreenWidth:  tt.fields.ScreenWidth,
				ScreenHeight: tt.fields.ScreenHeight,
			}
			gotX, gotY := s.ScalePosition(tt.args.x, tt.args.y)
			require.Equal(t, tt.wantX, gotX)
			require.Equal(t, tt.wantY, gotY)
		})
	}
}

func TestLeftPositionScaler_ScalePosition(t *testing.T) {
	type fields struct {
		TabletWidth  int
		TabletHeight int
		ScreenWidth  int
		ScreenHeight int
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantX  int
		wantY  int
	}{
		{
			name: "square scale factor 1",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 100,
				ScreenWidth:  100,
				ScreenHeight: 100,
			},
			args: args{
				x: 50,
				y: 50,
			},
			wantX: 50,
			wantY: 50,
		},
		{
			name: "square scale factor 2",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 100,
				ScreenWidth:  200,
				ScreenHeight: 200,
			},
			args: args{
				x: 50,
				y: 50,
			},
			wantX: 100,
			wantY: 100,
		},
		{
			name: "non-square",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 200,
				ScreenWidth:  400,
				ScreenHeight: 200,
			},
			args: args{
				x: 50,
				y: 100,
			},
			wantX: 200,
			wantY: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LeftPositionScaler{
				TabletWidth:  tt.fields.TabletWidth,
				TabletHeight: tt.fields.TabletHeight,
				ScreenWidth:  tt.fields.ScreenWidth,
				ScreenHeight: tt.fields.ScreenHeight,
			}
			gotX, gotY := s.ScalePosition(tt.args.x, tt.args.y)
			require.Equal(t, tt.wantX, gotX)
			require.Equal(t, tt.wantY, gotY)
		})
	}
}

func TestVerticalPositionScaler_ScalePosition(t *testing.T) {
	type fields struct {
		TabletWidth  int
		TabletHeight int
		ScreenWidth  int
		ScreenHeight int
	}
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantX  int
		wantY  int
	}{
		{
			name: "square scale factor 1",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 100,
				ScreenWidth:  100,
				ScreenHeight: 100,
			},
			args: args{
				x: 50,
				y: 50,
			},
			wantX: 50,
			wantY: 50,
		},
		{
			name: "square scale factor 2",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 100,
				ScreenWidth:  200,
				ScreenHeight: 200,
			},
			args: args{
				x: 50,
				y: 50,
			},
			wantX: 100,
			wantY: 100,
		},
		{
			name: "non-square",
			fields: fields{
				TabletWidth:  100,
				TabletHeight: 200,
				ScreenWidth:  400,
				ScreenHeight: 200,
			},
			args: args{
				x: 50,
				y: 100,
			},
			wantX: 200,
			wantY: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &VerticalPositionScaler{
				TabletWidth:  tt.fields.TabletWidth,
				TabletHeight: tt.fields.TabletHeight,
				ScreenWidth:  tt.fields.ScreenWidth,
				ScreenHeight: tt.fields.ScreenHeight,
			}
			gotX, gotY := s.ScalePosition(tt.args.x, tt.args.y)
			require.Equal(t, tt.wantX, gotX)
			require.Equal(t, tt.wantY, gotY)
		})
	}
}
