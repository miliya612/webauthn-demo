package model

import "time"

type Session struct {
	ID           string
	UserID       string
	Challenge    string
	LastAccessed time.Time
}
