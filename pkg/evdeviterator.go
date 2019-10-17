package remouse

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

type rawEvent struct {
	// The time values must be uint32 to work with the tablet build. The default
	// syscall.Timeval uses uint64 on a 64bit platform so we must adapt here.
	Sec   uint32
	Usec  uint32
	Type  uint16
	Code  uint16
	Value int32
}

// FileEvdevIterator implements the EvdevIterator interface by consuming from
// an io.ReadCloser.
type FileEvdevIterator struct {
	Source  io.ReadCloser
	err     error
	current EvdevEvent
}

// Next reads an event from the file source.
func (it *FileEvdevIterator) Next() bool {
	if it.err != nil {
		// Prevent re-entry after an error.
		return false
	}

	evt := rawEvent{}
	size := binary.Size(evt)
	buf := make([]byte, size)

	if _, err := it.Source.Read(buf); err != nil {
		it.err = err
		return false
	}

	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &evt); err != nil {
		it.err = err
		return false
	}

	it.current = EvdevEvent{
		Time:  time.Unix(int64(evt.Sec), int64(evt.Usec)),
		Type:  evt.Type,
		Code:  evt.Code,
		Value: evt.Value,
	}
	return true
}

// Current returns the iterator value.
func (it *FileEvdevIterator) Current() EvdevEvent {
	return it.current
}

// Close the underlying source and return any errors.
func (it *FileEvdevIterator) Close() error {
	err := it.Source.Close()
	if it.err == nil {
		return err
	}
	return it.err
}

// SelectingEvdevIterator reduces an iterator output to a selection of top-level
// event types.
type SelectingEvdevIterator struct {
	Wrapped   EvdevIterator
	Selection []uint16
}

// Next continually calls the wrapped Next() until it either returns a value
// that matches the selection criteria or it returns a false.
func (it *SelectingEvdevIterator) Next() bool {
	for it.Wrapped.Next() {
		t := it.Wrapped.Current().Type
		for _, selection := range it.Selection {
			if t == selection {
				return true
			}
		}
	}
	return false
}

// Current returns the active element.
func (it *SelectingEvdevIterator) Current() EvdevEvent {
	return it.Wrapped.Current()
}

// Close proxies to the wrapped instance.
func (it *SelectingEvdevIterator) Close() error {
	return it.Wrapped.Close()
}

// FilteringEvdevIterator reduces an iterator output to all but a selection of
// top-level event types.
type FilteringEvdevIterator struct {
	Wrapped EvdevIterator
	Filter  []uint16
}

// Next continually calls the wrapped Next() until it either returns a value
// that matches the filter criteria or it returns a false.
func (it *FilteringEvdevIterator) Next() bool {
	for it.Wrapped.Next() {
		t := it.Wrapped.Current().Type
		ok := true
		for _, filter := range it.Filter {
			if t == filter {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

// Current returns the active element.
func (it *FilteringEvdevIterator) Current() EvdevEvent {
	return it.Wrapped.Current()
}

// Close proxies to the wrapped instance.
func (it *FilteringEvdevIterator) Close() error {
	return it.Wrapped.Close()
}
