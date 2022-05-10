package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
)

var ErrDate = errors.New("invalid Date format")
var IdNotFound = errors.New("EventID not found")
var ErrNotBeginMonth = errors.New("Date is not Begin Month")
var ErrNotBeginWeek = errors.New("Date is not Begin Week")

type Storage struct {
	DbName     string
	DbUserName string
	DbPassword string
	DbConnect  *sql.DB
}

func New(dbName, dbUserName, dbPassword string) *Storage {
	return &Storage{
		DbName: dbName, DbUserName: dbUserName, DbPassword: dbPassword,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.DbConnect, err = sql.Open("postgres", "user="+s.DbUserName+" dbname="+s.DbName+" password="+s.DbPassword+" sslmode=disable")
	if err == nil {
		err = s.DbConnect.PingContext(ctx)
	}
	return err
}

func (s *Storage) CreateEvent(title, startDateStr, details string, userID int) error {
	myId := st.GenUUID()
	query := `
		insert into events(ID, Title, StartDate, Details,UserID)
		values($1, $2, $3, $4, $5)
	`
	result, err := s.DbConnect.Exec(query, myId, title, startDateStr, details, userID)
	fmt.Println(result, err)
	return err
}

func (s *Storage) UpdateEvent(eventID, title, startDateStr, details string, userID int) error {
	query := `
		update events
		set 
		Title=$1,
		StartDate=$2, 
		Details=$3,
		UserID=$4
		where id=ID=$5	
    `
	result, err := s.DbConnect.Exec(query, title, startDateStr, details, userID, eventID)
	fmt.Println(result, err)
	return err
}

func (s *Storage) DeleteEvent(eventID string) error {
	query := `
		delete 
		from events
		where StartDate=$1	
    `
	result, err := s.DbConnect.Exec(query, eventID)
	fmt.Println(result, err)
	return err
}

func rowsToStruct(rows *sql.Rows) ([]st.Event, error) {
	var myEventList []st.Event
	var eventID, title, details string
	var userID int
	var startDateStr time.Time

	for rows.Next() {
		if err := rows.Scan(&eventID, &title, &startDateStr, &details, &userID); err != nil {
			return nil, err
		}
		myEventList = append(myEventList, st.Event{ID: eventID, Title: title, StartDate: startDateStr, Details: details, UserID: userID})
	}
	return myEventList, nil
}

func (s *Storage) GetEventByDate(startDateStr string) ([]st.Event, error) {
	query := `
		select ID, Title, StartDate, Details,UserID
		from events
		where StartDate=$1
    `

	rows, err := s.DbConnect.Query(query, startDateStr)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	myEventList, newErr := rowsToStruct(rows)
	return myEventList, newErr
}

func (s *Storage) GetEventMonth(startDateStr string) ([]st.Event, error) {
	myTime, err := time.Parse("2006-1-2", startDateStr)
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
		where date_trunc('month',StartDate)=$1	
	`
	rows, err := s.DbConnect.Query(query, startDateStr)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	myEventList, newErr := rowsToStruct(rows)
	return myEventList, newErr
}

func (s *Storage) GetEventWeek(startDateStr string) ([]st.Event, error) {
	myTime, err := time.Parse("2006-1-2", startDateStr)
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
		where date_trunc('week',StartDate)=$1	
	`
	rows, err := s.DbConnect.Query(query, startDateStr)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	myEventList, newErr := rowsToStruct(rows)
	return myEventList, newErr
}

func (s *Storage) CreateTable() error {
	query := `
	create table events (
		ID text primary key,
		Title text,
		StartDate date,
		Details text,
		UserID bigint
	);
	`
	result, err := s.DbConnect.Exec(query)
	fmt.Println(result, err)
	return err
}

func (s *Storage) Close(ctx context.Context) error {
	err := s.DbConnect.Close()
	return err
}

func dbConect(dsn string) (*sql.DB, error) {
	//	db, err := sql.Open("pgx", dsn)
	db, err := sql.Open("postgres", dsn)
	return db, err
}
