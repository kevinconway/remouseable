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
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestFileEvdevIterator_Next(t *testing.T) {
	sentinelErr := fmt.Errorf("test")
	tests := []struct {
		name        string
		want        bool
		wantErr     bool
		expectedErr error
		wantRead    bool
		readBytes   []byte
		readErr     error
		closeErr    error
	}{
		{
			name:        "read error",
			want:        false,
			wantErr:     true,
			expectedErr: sentinelErr,
			wantRead:    true,
			readBytes:   nil,
			readErr:     sentinelErr,
			closeErr:    nil,
		},
		{
			name:        "read data",
			want:        true,
			wantErr:     false,
			expectedErr: nil,
			wantRead:    true,
			readBytes:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			readErr:     nil,
			closeErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			src := NewMockReadCloser(ctrl)
			it := &FileEvdevIterator{
				Source: src,
			}
			if tt.wantRead {
				src.EXPECT().Read(gomock.Any()).Do(func(b []byte) {
					copy(b, tt.readBytes)
				}).Return(0, tt.readErr)
			}
			src.EXPECT().Close().Return(tt.closeErr).AnyTimes()
			require.Equal(t, tt.want, it.Next())
			if tt.wantErr {
				require.Equal(t, tt.expectedErr, it.Close())
				require.Equal(t, false, it.Next())
			}
		})
	}
}

func TestSelectingEvdevIterator_Next(t *testing.T) {
	type fields struct {
		Selection []uint16
	}
	tests := []struct {
		name     string
		fields   fields
		source   []EvdevEvent
		expected []EvdevEvent
	}{
		{
			name:     "empty source",
			fields:   fields{Selection: []uint16{0}},
			source:   []EvdevEvent{},
			expected: []EvdevEvent{},
		},
		{
			name:     "full set",
			fields:   fields{Selection: []uint16{0}},
			source:   []EvdevEvent{{}, {}, {}},
			expected: []EvdevEvent{{}, {}, {}},
		},
		{
			name:     "partial set",
			fields:   fields{Selection: []uint16{0}},
			source:   []EvdevEvent{{Type: 1}, {Type: 1}, {Type: 0}},
			expected: []EvdevEvent{{Type: 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wrapped := NewMockEvdevIterator(ctrl)
			it := &SelectingEvdevIterator{
				Wrapped:   wrapped,
				Selection: tt.fields.Selection,
			}
			for _, s := range tt.source {
				wrapped.EXPECT().Next().Return(true)
				wrapped.EXPECT().Current().Return(s)
			}
			wrapped.EXPECT().Next().Return(false)
			wrapped.EXPECT().Close().Return(nil)
			results := make([]EvdevEvent, 0, len(tt.expected))
			for it.Next() {
				results = append(results, it.Current())
			}
			require.Equal(t, nil, it.Close())
			require.ElementsMatch(t, tt.expected, results)
		})
	}
}

func TestFilteringEvdevIterator_Next(t *testing.T) {
	type fields struct {
		Filter []uint16
	}
	tests := []struct {
		name     string
		fields   fields
		source   []EvdevEvent
		expected []EvdevEvent
	}{
		{
			name:     "empty source",
			fields:   fields{Filter: []uint16{0}},
			source:   []EvdevEvent{},
			expected: []EvdevEvent{},
		},
		{
			name:     "all filtered",
			fields:   fields{Filter: []uint16{0}},
			source:   []EvdevEvent{{}, {}, {}},
			expected: []EvdevEvent{},
		},
		{
			name:     "partial filter",
			fields:   fields{Filter: []uint16{0}},
			source:   []EvdevEvent{{Type: 1}, {Type: 1}, {Type: 0}},
			expected: []EvdevEvent{{Type: 1}, {Type: 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wrapped := NewMockEvdevIterator(ctrl)
			it := &FilteringEvdevIterator{
				Wrapped: wrapped,
				Filter:  tt.fields.Filter,
			}
			for _, s := range tt.source {
				wrapped.EXPECT().Next().Return(true)
				wrapped.EXPECT().Current().Return(s)
			}
			wrapped.EXPECT().Next().Return(false)
			wrapped.EXPECT().Close().Return(nil)
			results := make([]EvdevEvent, 0, len(tt.expected))
			for it.Next() {
				results = append(results, it.Current())
			}
			require.Equal(t, nil, it.Close())
			require.ElementsMatch(t, tt.expected, results)
		})
	}
}
