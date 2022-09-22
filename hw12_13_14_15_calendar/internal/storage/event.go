package storage

import (
	"crypto/rand"
	"fmt"
	"time"
)

func GenUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

type Event struct {
	ID        string
	Title     string
	StartDate time.Time
	Details   string
	UserID    int
}

type EventList struct {
	Events   []*Event
	EventMap map[string]*Event
}
