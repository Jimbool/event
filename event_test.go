package event

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

type eventListener struct {}

const eventName = "foobar"

type eventData struct {
    t *testing.T
    c chan bool
}

func (l *eventListener) Trigger(event string, data interface{}) {
    if d, ok := data.(*eventData); ok {
        d.c <- true
        return
    }
    panic("Type assertion failed")
}

type dontListen struct {}

func (l *dontListen) Trigger(event string, data interface{}) {
    if c, ok := data.(*eventData); ok {
        c.t.FailNow()
        return
    }
    panic("Type assertion failed")
}

func TestTrigger(t *testing.T) {
    data := &eventData{t, make(chan bool)}
    em := New()
    el := new(eventListener)

    em.AttachListener(eventName, el)
    em.Trigger(eventName, data)

    assert.True(t, <-data.c)
}

func TestDetach(t *testing.T) {
    data := &eventData{t, make(chan bool)}
    em := New()
    el := new(dontListen)

    em.AttachListener(eventName, el)
    em.DetachListener(eventName, el)
    em.Trigger(eventName, data)

    assert.True(t, true)
}