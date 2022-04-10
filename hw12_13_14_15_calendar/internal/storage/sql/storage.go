package sqlstorage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
)

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

func (s *Storage) Close(ctx context.Context) error {
	err := s.DbConnect.Close()
	return err
}

func dbConect(dsn string) (*sql.DB, error) {
	//	db, err := sql.Open("pgx", dsn)
	db, err := sql.Open("postgres", dsn)
	//	fmt.Println(db, err)
	return db, err
}
