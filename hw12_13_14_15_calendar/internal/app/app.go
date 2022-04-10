package app

import (
	"context"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	Storage Storage
	Logg    Logger
}

type Logger interface { // TODO
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}

type Storage interface {
	CreateEvent(title, startDateStr, details string, userID int) error
	UpdateEvent(eventID, title, startDateStr, details string, userID int) error
	DeleteEvent(eventID string) error
	GetEventByDate(startDateStr string) ([]st.Event, error)
	GetEventMonth(startDateStr string) ([]st.Event, error)
	GetEventWeek(startDateStr string) ([]st.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{Storage: storage, Logg: logger}
}

func (a *App) CreateEvent(ctx context.Context, title, startDateStr, details string, userID int) error {
	return a.Storage.CreateEvent(title, startDateStr, details, userID)
}
func (a *App) UpdateEvent(ctx context.Context, eventID, title, startDateStr, details string, userID int) error {
	return a.Storage.UpdateEvent(eventID, title, startDateStr, details, userID)
}
func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return a.Storage.DeleteEvent(eventID)
}
func (a *App) GetEventByDate(ctx context.Context, startDateStr string) ([]st.Event, error) {
	return a.Storage.GetEventByDate(startDateStr)
}
func (a *App) GetEventMonth(ctx context.Context, startDateStr string) ([]st.Event, error) {
	return a.Storage.GetEventMonth(startDateStr)
}
func (a *App) GetEventWeek(ctx context.Context, startDateStr string) ([]st.Event, error) {
	return a.Storage.GetEventWeek(startDateStr)
}

func (a *App) Trace(args ...interface{}) {
	a.Logg.Trace(args)
}
func (a *App) Debug(args ...interface{}) {
	a.Logg.Debug(args)
}
func (a *App) Info(args ...interface{}) {
	a.Logg.Info(args)
}
func (a *App) Print(args ...interface{}) {
	a.Logg.Print(args)
}
func (a *App) Warning(args ...interface{}) {
	a.Logg.Warning(args)
}
func (a *App) Error(args ...interface{}) {
	a.Logg.Error(args)
}
