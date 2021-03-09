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

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestStateMachineEmptyIterator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	it := NewMockEvdevIterator(ctrl)
	sm := &EvdevStateMachine{
		Iterator: it,
	}

	it.EXPECT().Next().Return(false).AnyTimes()
	it.EXPECT().Close().Return(nil)

	require.False(t, sm.Next())
	require.False(t, sm.Next())
	require.Nil(t, sm.Close())
}

func TestEvdevStateMachine_next(t *testing.T) {
	type fields struct {
		Iterator          EvdevIterator
		PressureThreshold int
		x                 int
		xChanged          bool
		y                 int
		yChanged          bool
		clicked           bool
		current           StateChange
	}
	type args struct {
		raw EvdevEvent
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        bool
		wantMachine *EvdevStateMachine
	}{
		{
			name: "skips non-ABS event",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{Type: EV_LED},
			},
			want: false,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
		},
		{
			name: "x event without y",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_X,
					Value: 1,
				},
			},
			want: false,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          true,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
		},
		{
			name: "y event without x",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_Y,
					Value: 1,
				},
			},
			want: false,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 1,
				yChanged:          true,
				clicked:           false,
				current:           nil,
			},
		},
		{
			name: "x to y",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          true,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_Y,
					Value: 1,
				},
			},
			want: true,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          false,
				y:                 1,
				yChanged:          false,
				clicked:           false,
				current:           &StateChangeMove{X: 1, Y: 1},
			},
		},
		{
			name: "y to x",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 1,
				yChanged:          true,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_X,
					Value: 1,
				},
			},
			want: true,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          false,
				y:                 1,
				yChanged:          false,
				clicked:           false,
				current:           &StateChangeMove{X: 1, Y: 1},
			},
		},
		{
			name: "click",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_PRESSURE,
					Value: 2000,
				},
			},
			want: true,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           true,
				current:           &StateChangeClick{},
			},
		},
		{
			name: "click while clicked",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           true,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_PRESSURE,
					Value: 2000,
				},
			},
			want: false,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           true,
				current:           nil,
			},
		},
		{
			name: "unclick",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           true,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_PRESSURE,
					Value: 500,
				},
			},
			want: true,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           &StateChangeUnclick{},
			},
		},
		{
			name: "unclick while unclicked",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_PRESSURE,
					Value: 500,
				},
			},
			want: false,
			wantMachine: &EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &EvdevStateMachine{
				Iterator:          tt.fields.Iterator,
				PressureThreshold: tt.fields.PressureThreshold,
				x:                 tt.fields.x,
				xChanged:          tt.fields.xChanged,
				y:                 tt.fields.y,
				yChanged:          tt.fields.yChanged,
				clicked:           tt.fields.clicked,
				current:           tt.fields.current,
			}
			require.Equal(t, tt.want, it.next(tt.args.raw))
			require.Equal(t, *tt.wantMachine, *it)
		})
	}
}

func TestDraggingEvdevStateMachine_next(t *testing.T) {
	type fields struct {
		Iterator          EvdevIterator
		PressureThreshold int
		x                 int
		xChanged          bool
		y                 int
		yChanged          bool
		clicked           bool
		current           StateChange
	}
	type args struct {
		raw EvdevEvent
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        bool
		wantMachine *DraggingEvdevStateMachine
	}{
		{
			name: "skips non-ABS event",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{Type: EV_LED},
			},
			want: false,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			}},
		},
		{
			name: "x event without y",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_X,
					Value: 1,
				},
			},
			want: false,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          true,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			}},
		},
		{
			name: "y event without x",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_Y,
					Value: 1,
				},
			},
			want: false,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 1,
				yChanged:          true,
				clicked:           false,
				current:           nil,
			}},
		},
		{
			name: "x to y",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          true,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_Y,
					Value: 1,
				},
			},
			want: true,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          false,
				y:                 1,
				yChanged:          false,
				clicked:           false,
				current:           &StateChangeMove{X: 1, Y: 1},
			}},
		},
		{
			name: "y to x",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 1,
				yChanged:          true,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_X,
					Value: 1,
				},
			},
			want: true,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          false,
				y:                 1,
				yChanged:          false,
				clicked:           false,
				current:           &StateChangeMove{X: 1, Y: 1},
			}},
		},
		{
			name: "move while clicked",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 1,
				yChanged:          true,
				clicked:           true,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_X,
					Value: 1,
				},
			},
			want: true,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 1,
				xChanged:          false,
				y:                 1,
				yChanged:          false,
				clicked:           true,
				current:           &StateChangeDrag{X: 1, Y: 1},
			}},
		},
		{
			name: "click",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_PRESSURE,
					Value: 2000,
				},
			},
			want: true,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           true,
				current:           &StateChangeClick{},
			}},
		},
		{
			name: "click while clicked",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           true,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_PRESSURE,
					Value: 2000,
				},
			},
			want: false,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           true,
				current:           nil,
			}},
		},
		{
			name: "unclick",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           true,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_PRESSURE,
					Value: 500,
				},
			},
			want: true,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           &StateChangeUnclick{},
			}},
		},
		{
			name: "unclick while unclicked",
			fields: fields{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			},
			args: args{
				raw: EvdevEvent{
					Type:  EV_ABS,
					Code:  ABS_PRESSURE,
					Value: 500,
				},
			},
			want: false,
			wantMachine: &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          nil,
				PressureThreshold: 1000,
				x:                 0,
				xChanged:          false,
				y:                 0,
				yChanged:          false,
				clicked:           false,
				current:           nil,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &DraggingEvdevStateMachine{&EvdevStateMachine{
				Iterator:          tt.fields.Iterator,
				PressureThreshold: tt.fields.PressureThreshold,
				x:                 tt.fields.x,
				xChanged:          tt.fields.xChanged,
				y:                 tt.fields.y,
				yChanged:          tt.fields.yChanged,
				clicked:           tt.fields.clicked,
				current:           tt.fields.current,
			}}
			require.Equal(t, tt.want, it.next(tt.args.raw))
			require.Equal(t, *tt.wantMachine, *it)
		})
	}
}
