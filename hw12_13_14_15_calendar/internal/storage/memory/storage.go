package memorystorage

import (
	"errors"
	"sync"
	"time"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrDate           = errors.New("invalid Date format")
	EventIDIsNotFound = errors.New("eventID not found") // nolint
	ErrNotBeginMonth  = errors.New("date is not Begin Month")
	ErrNotBeginWeek   = errors.New("date is not Begin Week")
)

type Storage struct {
	Events map[string]st.Event
	mu     sync.RWMutex
}

func (s *Storage) CreateEvent(title, startDateStr, details string, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	myTime, err := time.Parse("2.1.2006", startDateStr)
	if err != nil {
		return ErrDate
	}
	myID := st.GenUUID()
	myEvent := st.Event{ID: myID, Title: title, StartDate: myTime, Details: details, UserID: userID}
	s.Events[myID] = myEvent
	return nil
}

func (s *Storage) UpdateEvent(eventID, title, startDateStr, details string, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Events[eventID]; !ok {
		return EventIDIsNotFound
	}
	myTime, err := time.Parse("2.1.2006", startDateStr)
	if err != nil {
		return ErrDate
	}
	myEvent := st.Event{ID: eventID, Title: title, StartDate: myTime, Details: details, UserID: userID}
	s.Events[eventID] = myEvent
	return nil
}

func (s *Storage) DeleteEvent(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Events[eventID]; !ok {
		return EventIDIsNotFound
	}
	delete(s.Events, eventID)
	return nil
}

func (s *Storage) GetEventByDate(startDateStr string) ([]st.Event, error) {
	myTime, err := time.Parse("2.1.2006", startDateStr)
	if err != nil {
		return nil, ErrDate
	}
	var myEventList []st.Event
	for _, Val := range s.Events {
		if myTime == Val.StartDate {
			myEventList = append(myEventList, Val)
		}
	}
	return myEventList, nil
}

func (s *Storage) GetEventMonth(startDateStr string) ([]st.Event, error) {
	myTime, err := time.Parse("2.1.2006", startDateStr)
	if err != nil {
		return nil, ErrDate
	}
	if myTime.Day() != 1 {
		return nil, ErrNotBeginMonth
	}
	var myEventList []st.Event
	for _, Val := range s.Events {
		if Val.StartDate.Year() == myTime.Year() && Val.StartDate.Month() == myTime.Month() {
			myEventList = append(myEventList, Val)
		}
	}
	return myEventList, nil
}

func (s *Storage) GetEventWeek(startDateStr string) ([]st.Event, error) {
	myTime, err := time.Parse("2.1.2006", startDateStr)
	if err != nil {
		return nil, ErrDate
	}
	if myTime.Weekday() != time.Monday {
		return nil, ErrNotBeginWeek
	}
	myYer, myWeek := myTime.ISOWeek()
	var myEventList []st.Event
	for _, Val := range s.Events {
		calYer, calWeek := Val.StartDate.ISOWeek()
		if myYer == calYer && myWeek == calWeek {
			myEventList = append(myEventList, Val)
		}
	}
	return myEventList, nil
}

func New() *Storage {
	return &Storage{
		map[string]st.Event{},
		sync.RWMutex{},
	}
}
