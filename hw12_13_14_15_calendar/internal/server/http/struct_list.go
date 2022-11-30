package internalhttp

type StartDate struct {
	StartDateStr string
}

type ForCreate struct {
	Title     string
	StartDate string
	Details   string
	UserID    uint32
}

type ForUpdate struct {
	EventID   string
	Title     string
	StartDate string
	Details   string
	UserID    uint32
}

type ForDelete struct {
	EventID string
}

type OneEvent struct {
	EventID   string
	Title     string
	StartDate string
	Details   string
	UserID    uint32
}
