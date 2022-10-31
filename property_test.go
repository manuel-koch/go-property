package property

import (
	"testing"
	"time"
)

func TestBasePropertyInterfaces(t *testing.T) {
	// check that BaseProperty inplements required interfaces
	var _ Comparable = (*BaseProperty[int])(nil)
}

func TestSimplePropertyInterfaces(t *testing.T) {
	// check that SimpleProperty inplements required interfaces
	var _ Comparable = (*BasicProperty[int])(nil)
	var _ Property[int] = (*BasicProperty[int])(nil)
}

// check that ComparableProperty inplements required interfaces
type dummycomparable struct {
	foo int
}

func (*dummycomparable) Equals(v interface{}) bool {
	return false
}

func TestComparablePropertyInterfaces(t *testing.T) {
	var _ Comparable = (*ComparableProperty[*dummycomparable])(nil)
	var _ Property[*dummycomparable] = (*ComparableProperty[*dummycomparable])(nil)
}

func TestSimpleIntProperty(t *testing.T) {
	p := NewBasicProperty(0)
	newValue := 42
	calledWith := -1

	p.ChangedSignal().Subscribe(func(v int) {
		calledWith = v
	})

	go func() {
		time.Sleep(500 * time.Millisecond)
		p.Set(newValue)
	}()

	time.Sleep(time.Second)
	currentValue := p.Get()

	if currentValue != newValue {
		t.Errorf("Value mismatch: want %d, got %d", newValue, currentValue)
	}
	if calledWith != newValue {
		t.Errorf("Value mismatch: want %d, got %d", calledWith, newValue)
	}
}

func TestSimpleIntPropertyMultipleSubscriptions(t *testing.T) {
	p := NewBasicProperty(0)
	newValue := 42
	callCount := 0

	p.ChangedSignal().Subscribe(func(v int) {
		callCount++
	})

	p.ChangedSignal().Subscribe(func(v int) {
		callCount++
	})

	go func() {
		time.Sleep(500 * time.Millisecond)
		p.Set(newValue)
	}()

	time.Sleep(time.Second)
	currentValue := p.Get()

	if currentValue != newValue {
		t.Errorf("Value mismatch: want %d, got %d", newValue, currentValue)
	}
	if callCount != 2 {
		t.Errorf("Callcount mismatch: want %d, got %d", 2, callCount)
	}
}

func TestSimpleIntPropertyUnchanged(t *testing.T) {
	p := NewBasicProperty(0)
	newValue := 42
	calledWith := -1
	callCount := 0

	p.ChangedSignal().Subscribe(func(v int) {
		t.Logf("Got %d", v)
		calledWith = v
		callCount++
	})

	go func() {
		time.Sleep(500 * time.Millisecond)
		t.Logf("Set %d", newValue)
		p.Set(newValue)
		time.Sleep(10 * time.Millisecond)
		t.Logf("Set %d", newValue)
		p.Set(newValue)
		time.Sleep(10 * time.Millisecond)
		t.Logf("Set %d", newValue)
		p.Set(newValue)
		time.Sleep(10 * time.Millisecond)
		t.Logf("Set %d", newValue)
		p.Set(newValue)
	}()

	time.Sleep(time.Second)
	currentValue := p.Get()

	if currentValue != newValue {
		t.Errorf("Value mismatch: want %d, got %d", newValue, currentValue)
	}
	if calledWith != newValue {
		t.Errorf("Value mismatch: want %d, got %d", calledWith, newValue)
	}
	if callCount != 1 {
		t.Errorf("Callcount mismatch: want %d, got %d", 1, callCount)
	}
}

func TestSimpleIntPropertySubscribeOnce(t *testing.T) {
	p := NewBasicProperty(0)
	firstValue := 42
	secondValue := 12
	thirdValue := 66
	calledWith := -1
	callCount := 0

	p.ChangedSignal().SubscribeOnce(func(v int) {
		t.Logf("Got %d", v)
		calledWith = v
		callCount++
	})

	go func() {
		time.Sleep(500 * time.Millisecond)
		t.Logf("Set %d", firstValue)
		p.Set(firstValue)
		time.Sleep(20 * time.Millisecond)
		t.Logf("Set %d", secondValue)
		p.Set(secondValue)
		time.Sleep(20 * time.Millisecond)
		t.Logf("Set %d", thirdValue)
		p.Set(thirdValue)
	}()

	time.Sleep(time.Second)
	currentValue := p.Get()

	if currentValue != thirdValue {
		t.Errorf("Value mismatch: want %d, got %d", thirdValue, currentValue)
	}
	if calledWith != firstValue {
		t.Errorf("Value mismatch: want %d, got %d", calledWith, firstValue)
	}
	if callCount != 1 {
		t.Errorf("Callcount mismatch: want %d, got %d", 1, callCount)
	}
}

func TestSimpleStringProperty(t *testing.T) {
	p := NewBasicProperty("0")
	newValue := "42"
	calledWith := ""

	p.ChangedSignal().Subscribe(func(v string) {
		calledWith = v
	})

	go func() {
		time.Sleep(500 * time.Millisecond)
		p.Set(newValue)
	}()

	time.Sleep(time.Second)
	currentValue := p.Get()

	if currentValue != newValue {
		t.Errorf("Value mismatch: want %s, got %s", newValue, currentValue)
	}
	if calledWith != newValue {
		t.Errorf("Value mismatch: want %s, got %s", calledWith, newValue)
	}
}

func TestSimplePropertyUnsubscribe(t *testing.T) {
	p := NewBasicProperty(0)
	unsubscripeAt := 42
	callCount := 0
	finalValue := 99

	var l *Listener[int]
	l = p.ChangedSignal().Subscribe(func(v int) {
		callCount++
		time.Sleep(10 * time.Millisecond)
		if v >= unsubscripeAt {
			l.Unsubscribe()
			t.Log("Unsubscriped")
		}
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		for i := 0; i < finalValue+1; i++ {
			p.Set(i)
		}
		t.Log("Done setting")
	}()

	time.Sleep(1 * time.Second)
	currentValue := p.Get()

	if currentValue != finalValue {
		t.Errorf("Value mismatch: want %d, got %d", finalValue, currentValue)
	}
	if callCount != unsubscripeAt {
		t.Errorf("Call count mismatch: want %d, got %d", unsubscripeAt, callCount)
	}
}

type StructValue struct {
	A int
	B string
}

func TestStructProperty(t *testing.T) {
	p := NewBasicProperty(StructValue{})
	newValue := StructValue{A: 42, B: "FOO"}
	calledWith := StructValue{}

	p.ChangedSignal().Subscribe(func(v StructValue) {
		calledWith = v
	})

	go func() {
		time.Sleep(500 * time.Millisecond)
		p.Set(newValue)
	}()

	time.Sleep(time.Second)
	currentValue := p.Get()

	if currentValue != newValue {
		t.Errorf("Value mismatch: want %v, got %v", newValue, currentValue)
	}
	if calledWith != newValue {
		t.Errorf("Value mismatch: want %v, got %v", calledWith, newValue)
	}
}

type ComplexValue struct {
	A []int
	B [3]string
}

var _ Comparable = (*ComplexValue)(nil)

func (c ComplexValue) Equals(other interface{}) bool {
	if o, ok := other.(ComplexValue); !ok {
		return false
	} else {
		if c.B != o.B {
			return false
		}
		if len(c.A) != len(o.A) {
			return false
		}
		for i, v := range o.A {
			if c.A[i] != v {
				return false
			}
		}
	}
	return true
}

func TestComplexProperty(t *testing.T) {
	p := NewComparableProperty(ComplexValue{A: []int{1, 2, 3}, B: [3]string{"hello", "world", "!"}})
	newValue := ComplexValue{A: []int{3, 2, 1, 0}, B: [3]string{"foo", "bar", "!"}}
	otherValue := ComplexValue{A: []int{0}, B: [3]string{"foo", "bar", "!"}}
	calledWith := ComplexValue{}
	callCount := 0

	p.ChangedSignal().Subscribe(func(v ComplexValue) {
		t.Logf("Got %v", v)
		calledWith = v
		callCount++
	})

	go func() {
		time.Sleep(500 * time.Millisecond)
		t.Logf("Set %v", newValue)
		p.Set(newValue)
		time.Sleep(10 * time.Millisecond)
		t.Logf("Set %v", otherValue)
		p.Set(otherValue)
		time.Sleep(10 * time.Millisecond)
		t.Logf("Set %v", newValue)
		p.Set(newValue)
		time.Sleep(10 * time.Millisecond)
		t.Logf("Set %v", newValue)
		p.Set(newValue)
	}()

	time.Sleep(time.Second)
	currentValue := p.Get()

	if !currentValue.Equals(newValue) {
		t.Errorf("Value mismatch: want %v, got %v", newValue, currentValue)
	}
	if !calledWith.Equals(newValue) {
		t.Errorf("Value mismatch: want %v, got %v", calledWith, newValue)
	}
	if callCount != 3 {
		t.Errorf("Callcount mismatch: want %v, got %v", 3, callCount)
	}
}
