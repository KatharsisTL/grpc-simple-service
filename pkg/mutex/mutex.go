package mutex

import "sync"

type BoolMutex interface {
	BoolValue() bool
	SetBoolValue(boolValue bool)
}

type Int64Mutex interface {
	Int64Value() int64
	SetInt64Value(int64Value int64)
}

type mutex struct {
	mu         sync.RWMutex
	boolValue  bool
	int64Value int64
}

func NewBoolMutex(boolValue bool) BoolMutex {
	return &mutex{
		mu:        sync.RWMutex{},
		boolValue: boolValue,
	}
}

func NewInt64Mutex(int64Value int64) Int64Mutex {
	return &mutex{
		mu:         sync.RWMutex{},
		int64Value: int64Value,
	}
}

func (m *mutex) BoolValue() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.boolValue
}

func (m *mutex) SetBoolValue(boolValue bool) {
	m.mu.Lock()
	m.boolValue = boolValue
	m.mu.Unlock()
}

func (m *mutex) Int64Value() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.int64Value
}

func (m *mutex) SetInt64Value(int64Value int64) {
	m.mu.Lock()
	m.int64Value = int64Value
	m.mu.Unlock()
}
