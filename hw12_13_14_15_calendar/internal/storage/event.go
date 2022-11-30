package storage

import (
	"time"

	"github.com/google/uuid"
)

func GenUUID() string {
	return uuid.New().String()
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
