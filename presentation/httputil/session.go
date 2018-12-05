package httputil

import "time"

type SessionManager interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
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