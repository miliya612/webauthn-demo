package model

import "time"

type Session struct {
	ID           string
	UserID       []byte
	Challenge    []byte
	LastAccessed time.Time
}
