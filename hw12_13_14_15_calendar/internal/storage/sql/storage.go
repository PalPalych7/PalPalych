package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" // justifying
	_ "github.com/lib/pq"
)

var (
	ErrDate          = errors.New("invalid Date format")
	ErrNotBeginMonth = errors.New("date is not Begin Month")
	ErrNotBeginWeek  = errors.New("date is not Begin Week")
)

type Storage struct {
	DBName     string
	DBUserName string
	DBPassword string
	DBHost     string
	DBPort     string
	DBConnect  *sql.DB
}

func New(dbName, dbUserName, dbPassword, dbHost, dbPort string) *Storage {
	return &Storage{
		DBName: dbName, DBUserName: dbUserName, DBPassword: dbPassword, DBHost: dbHost, DBPort: dbPort,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	//	myStr := "user=" + s.DBUserName + " dbname=" + s.DBName + " password=" + s.DBPassword + " sslmode=disable"
	//	s.DBConnect, err = sql.Open("postgres", myStr)
	myStr := "postgres://" + s.DBUserName + ":" + s.DBPassword + "@"
	myStr += s.DBHost + ":" + s.DBPort + "/" + s.DBName + "?sslmode=disable"
	s.DBConnect, err = sql.Open("postgres", myStr)

	if err == nil {
		err = s.DBConnect.PingContext(ctx)
	}
	return err
}

func (s *Storage) CreateEvent(title, startDateStr, details string, userID int) error {
	myID := st.GenUUID()
	query := `
		insert into events(ID, Title, StartDate, Details,UserID)
		values($1, $2, to_date($3,'DD.MM.YYYY'), $4, $5)
	`
	result, err := s.DBConnect.Exec(query, myID, title, startDateStr, details, userID)
	fmt.Println(result, err)
	return err
}

func (s *Storage) UpdateEvent(eventID, title, startDateStr, details string, userID int) error {
	query := `
		update events
		set 
		Title=$1,
		StartDate=to_date($2,'DD.MM.YYYY'), 
		Details=$3,
		UserID=$4
		where ID=$5	
    `
	result, err := s.DBConnect.Exec(query, title, startDateStr, details, userID, eventID)
	fmt.Println(result, err)
	return err
}

func (s *Storage) DeleteEvent(eventID string) error {
	query := `
		delete 
		from events
		where ID=$1
    `
	result, err := s.DBConnect.Exec(query, eventID)
	fmt.Println(result, err)
	return err
}

func rowsToStruct(rows *sql.Rows) ([]st.Event, error) {
	var myEventList []st.Event
	var eventID, title, details string
	var userID int
	var startDateStr time.Time
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&eventID, &title, &startDateStr, &details, &userID); err != nil {
			return nil, err
		}
		myEventList = append(myEventList, st.Event{ID: eventID, Title: title, StartDate: startDateStr, Details: details, UserID: userID}) //nolint
	}
	return myEventList, nil
}

func (s *Storage) GetEventByDate(startDateStr string) ([]st.Event, error) {
	query := `
		select ID, Title, StartDate, Details,UserID
		from events
		where StartDate=to_date($1,'DD.MM.YYYY')
    `

	rows, err := s.DBConnect.Query(query, startDateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	myEventList, newErr := rowsToStruct(rows)
	return myEventList, newErr
}

func (s *Storage) GetEventMonth(startDateStr string) ([]st.Event, error) {
	myTime, err := time.Parse("2.1.2006", startDateStr)
	fmt.Println("myT=", myTime)
	if err != nil {
		return nil, ErrDate
	}
	if myTime.Day() != 1 {
		return nil, ErrNotBeginMonth
	}

	query := `
		select ID, Title, StartDate, Details,UserID
		from events
		where date_trunc('month',StartDate)=to_date($1,'DD.MM.YYYY')
	`
	rows, err := s.DBConnect.Query(query, startDateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	myEventList, newErr := rowsToStruct(rows)
	return myEventList, newErr
}

func (s *Storage) GetEventWeek(startDateStr string) ([]st.Event, error) {
	myTime, err := time.Parse("2.1.2006", startDateStr)
	fmt.Println("myT=", myTime)
	if err != nil {
		return nil, ErrDate
	}
	if myTime.Weekday() != time.Monday {
		return nil, ErrNotBeginWeek
	}

	query := `
		select ID, Title, StartDate, Details,UserID
		from events
		where date_trunc('week',StartDate)=to_date($1,'DD.MM.YYYY')
	`
	rows, err := s.DBConnect.Query(query, startDateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	myEventList, newErr := rowsToStruct(rows)
	return myEventList, newErr
}

func (s *Storage) GetNotSendEventByDate(startDateStr string) ([]st.Event, error) {
	query := `
		select ID, Title, StartDate, Details,UserID
		from events
		where StartDate=to_date($1,'DD.MM.YYYY')
		  and id not in (select event_id from shed_send_id)
    `
	rows, err := s.DBConnect.Query(query, startDateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	myEventList, newErr := rowsToStruct(rows)
	return myEventList, newErr
}

func (s *Storage) SetSendMessID(messID string) error {
	query := `
	   insert into shed_send_id(event_id)  values($1)
	`
	_, err := s.DBConnect.Exec(query, messID)
	return err
}

func (s *Storage) SendMessStat(messID string) error {
	query := `
		insert into send_events_stat(send_mess) values($1)
	`
	_, err := s.DBConnect.Exec(query, messID)
	return err
}

func (s *Storage) Close(ctx context.Context) error {
	err := s.DBConnect.Close()
	return err
}
