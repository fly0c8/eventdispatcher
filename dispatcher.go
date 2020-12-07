package main

import (
	"reflect"
	"sync"
)

const (
	UnknownEvent = iota
	CommitEvent
	MessageEvent
	FileOpenedEvent
)

var eventTypeStrings = [...]string{
	"unkown", "entryCommitted", "messageReceived", "fileOpened",
}

type EventType uint16

func (t EventType) String() string {
	return eventTypeStrings[t]
}

type Event interface {
	Type() EventType
	Source() interface{}
	Value() interface{}
}

// event is an internal implementation of the Event interface
type event struct {
	etype  EventType
	source interface{}
	value  interface{}
}

func (e *event) Type() EventType {
	return e.etype
}
func (e *event) Source() interface{} {
	return e.source
}
func (e *event) Value() interface{} {
	return e.value
}

// Callback is a function that can receive events
type Callback func(Event) error

// Dispatcher objects can register callbacks for specific events,
// then when those events occur, dispatch them to all callback functions
type Dispatcher struct {
	sync.RWMutex
	source    interface{}
	callbacks map[EventType][]Callback
}

func (d *Dispatcher) Init(source interface{}) {
	d.source = source
	d.callbacks = make(map[EventType][]Callback)
}
func (d *Dispatcher) Register(etype EventType, callback Callback) {
	d.Lock()
	defer d.Unlock()
	d.callbacks[etype] = append(d.callbacks[etype], callback)
}
func (d *Dispatcher) Remove(etype EventType, callback Callback) {
	d.Lock()
	defer d.Unlock()

	// grab a ref to the function pointer
	ptr := reflect.ValueOf(callback).Pointer()
	callbacks := d.callbacks[etype]
	for idx, cb := range callbacks {
		if reflect.ValueOf(cb).Pointer() == ptr {
			d.callbacks[etype] = append(callbacks[:idx], callbacks[idx+1:]...)
		}
	}
}
func (d *Dispatcher) Dispatch(etype EventType, value interface{}) error {
	d.RLock()
	defer d.RUnlock()
	e := &event{
		etype:  etype,
		source: d.source,
		value:  value,
	}
	for _, cb := range d.callbacks[etype] {
		if err := cb(e); err != nil {
			return err
		}
	}
	return nil
}
