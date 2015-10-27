package event

import (
	"sync"
)

// 定义监听器的接口对象
type Listener interface {
	Trigger(event string, data interface{})
}

// 定义事件管理器对象
type EventManager struct {
	listeners map[string][]Listener
}

// 创建新的事件管理器对象
// 返回值：
// 事件管理器对象
func New() *EventManager {
	return &EventManager{
		listeners: make(map[string][]Listener, 0),
	}
}

// 附加监听器
// event：事件名称
// listener：监听器对象
// 返回值：无
func (em *EventManager) AttachListener(event string, listener Listener) {
	if em.listeners[event] == nil {
		em.listeners[event] = make([]Listener, 0)
	}

	em.listeners[event] = append(em.listeners[event], listener)
}

// 分离监听器
// event：事件名称
// listener：监听器对象
// 返回值：无
func (em *EventManager) DetachListener(event string, listener Listener) {
	if em.listeners[event] == nil {
		return
	}

	for k, v := range em.listeners[event] {
		if v == listener {
			em.listeners[event] = append(em.listeners[event][:k], em.listeners[event][k+1:]...)
		}
	}
}

// 触发监听器事件
// event：事件名称
// listener：监听器对象
// 返回值：无
func (em *EventManager) Trigger(event string, data interface{}) {
	if em.listeners[event] == nil {
		return
	}

	for _, listener := range em.listeners[event] {
		go listener.Trigger(event, data)
	}
}

// 同步触发监听器事件
// event：事件名称
// listener：监听器对象
// 返回值：无
func (em *EventManager) TriggerAndWait(event string, data interface{}) {
	wg := new(sync.WaitGroup)
	if em.listeners[event] == nil {
		return
	}

	for _, listener := range em.listeners[event] {
		wg.Add(1)
		go func(listener Listener) {
			defer wg.Done()
			listener.Trigger(event, data)
		}(listener)
	}

	wg.Wait()
}
