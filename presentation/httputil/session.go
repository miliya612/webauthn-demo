package httputil

import (
	"fmt"
	"sync"
	"time"
)

const KeySessionID = "sid"

type Manager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxlifetime int64
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

var provides = make(map[string]Provider)

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provider
}

type SessionManagement struct {
	sid          string
	lastAccessed time.Time
	store        map[interface{}]interface{}
}

func (m *SessionManagement) Set(key, value interface{}) error {
	m.store[key] = value
	touch(m)
	return nil
}

func (m *SessionManagement) Get(key interface{}) interface{} {
	touch(m)
	if v, ok := m.store[key]; ok {
		return v
	}
	return nil
}

func (m *SessionManagement) Delete(key interface{}) error {
	delete(m.store, key)
	touch(m)
	return nil
}

func (m *SessionManagement) SessionID() string {
	return m.sid
}

func touch(m *SessionManagement) *SessionManagement {
	m.lastAccessed = time.Now()
	return m
}
