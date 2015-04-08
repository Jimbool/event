package event
import "sync"

type Listener interface {
    Trigger(event string, data interface{})
}

type eventManager struct {
    listeners map[string][]Listener
}

func New() *eventManager {
    return &eventManager{
        listeners: make(map[string][]Listener, 0),
    }
}

func (em *eventManager) AttachListener(event string, listener Listener) {
    if em.listeners[event] == nil {
        em.listeners[event] = make([]Listener,0)
    }
    em.DetachListener(event, listener)
    em.listeners[event] = append(em.listeners[event], listener)
}

func (em *eventManager) DetachListener(event string, listener Listener) {
    if em.listeners[event] == nil {
        return
    }
    for k, v := range em.listeners[event] {
        if v == listener {
            em.listeners[event] = append(em.listeners[event][:k], em.listeners[event][k+1:]...)
        }
    }
}

func (em *eventManager) Trigger(event string, data interface{}) {
    if em.listeners[event] == nil {
        return
    }
    for _, listener := range em.listeners[event] {
        go listener.Trigger(event, data)
    }
}


func (em *eventManager) TriggerAndWait(event string, data interface{}) {
    wg := new(sync.WaitGroup)
    if em.listeners[event] == nil {
        return
    }
    for _, listener := range em.listeners[event] {
        wg.Add(1)
        go func() {
            defer wg.Done()
            listener.Trigger(event, data)
        }()
    }
    wg.Wait()
}
