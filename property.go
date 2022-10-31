package property

import (
	"sync"
)

//-- Interfaces ------------------------------------------

type Comparable interface {
	Equals(interface{}) bool
}

type Property[T any] interface {
	Get() T
	Set(T)
	Equals(interface{}) bool
	ChangedSignal() *Signal[T]
}

//-- Structs ---------------------------------------------

type Signal[T any] struct {
	mutex     sync.RWMutex
	listeners []*Listener[T]
}

type Listener[T any] struct {
	channel chan T
	signal  *Signal[T]
}

// BaseProperty holds common data for any property implementation
type BaseProperty[T any] struct {
	value      T
	changed    Signal[T]
	comparator Comparable
}

// BasicProperty can be used for all built-in types that support `comparable`
type BasicProperty[T comparable] struct {
	BaseProperty[T]
}

// ComparableProperty can be used for anything that is not a built-in type and
// requires a custom equality implementation
type ComparableProperty[T Comparable] struct {
	BaseProperty[T]
}

//-- Constructors -----------------------------------------

func NewBasicProperty[T comparable](value T) Property[T] {
	p := &BasicProperty[T]{BaseProperty[T]{value: value}}
	p.comparator = p
	return p
}

func NewComparableProperty[T Comparable](value T) Property[T] {
	p := &ComparableProperty[T]{BaseProperty[T]{value: value}}
	p.comparator = p
	return p
}

//-- Methods ----------------------------------------------

func (p *BaseProperty[T]) Equals(v interface{}) bool {
	return p.comparator.Equals(v)
}

func (p *BaseProperty[T]) Get() T {
	return p.value
}

func (p *BaseProperty[T]) Set(value T) {
	changed := !p.Equals(value)
	p.value = value
	if changed {
		p.changed.Emit(p.value)
	}
}

func (p *BasicProperty[T]) Equals(v interface{}) bool {
	if v_, ok := v.(T); ok {
		return p.value == v_
	}
	return false
}

func (p *BasicProperty[T]) ChangedSignal() *Signal[T] {
	return &p.changed
}

func (p *ComparableProperty[T]) Equals(v interface{}) bool {
	return p.value.Equals(v)
}

func (p *ComparableProperty[T]) ChangedSignal() *Signal[T] {
	return &p.changed
}

func (s *Signal[T]) Subscribe(callback func(T)) *Listener[T] {
	listener := &Listener[T]{channel: make(chan T), signal: s}
	s.addListener(listener)
	cb := callback
	go func() {
		for {
			v, ok := <-listener.channel
			if !ok || listener.signal == nil {
				break
			}
			cb(v)
		}
		close(listener.channel)
		listener.Unsubscribe()
	}()
	return listener
}

func (s *Signal[T]) SubscribeOnce(callback func(T)) *Listener[T] {
	var listener *Listener[T]
	listener = s.Subscribe(func(v T) {
		callback(v)
		listener.Unsubscribe()
	})
	return listener
}

func (s *Signal[T]) addListener(listener *Listener[T]) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.listeners = append(s.listeners, listener)
}

func (s *Signal[T]) removeListener(listener *Listener[T]) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, l := range s.listeners {
		if l == listener {
			s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
			break
		}
	}
}

func (s *Signal[T]) Emit(value T) {
	s.mutex.RLock()
	listeners := make([]*Listener[T], len(s.listeners))
	copy(listeners, s.listeners)
	s.mutex.RUnlock()

	for _, l := range listeners {
		l.channel <- value
	}
}

func (l *Listener[T]) Unsubscribe() {
	if l.signal != nil {
		l.signal.removeListener(l)
		l.signal = nil
	}
}

func main() {
}
