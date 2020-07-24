package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type persistenseStorage struct {
	db *sql.DB
}

const initDbSql = `CREATE TABLE users (
	id TEXT PRIMARY KEY,
	country VARCHAR(10) NOT NULL,
	level INT NOT NULL,
	properties JSON NOT NULL,
	candidate TEXT NOT NULL,
	context VARCHAR(254) NOT NULL 
);`

func newPersistenseStorageSqllite() (*persistenseStorage, error) {
	db, err := sql.Open("sqlite3", "file:db.sqlite?cache=shared")
	if err != nil {
		return nil, err
	}

	batch := []string{initDbSql}

	for _, b := range batch {
		_, _ = db.Exec(b)
	}

	return &persistenseStorage{
		db: db,
	}, nil
}

func newPersistenseStoragePq() (*persistenseStorage, error) {
	connStr := os.Getenv("DB_CONNECTION")
	if connStr == "" {
		return nil, errors.New("DB_CONNECTION is empty")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)

	return &persistenseStorage{
		db: db,
	}, nil
}

func (s *persistenseStorage) load(id string) (*StorageUser, error) {
	sqlStatement := `SELECT id, country, context, level, properties, candidate FROM users WHERE id = $1;`
	var user StorageUser
	row := s.db.QueryRow(sqlStatement, id)
	var properties string
	err := row.Scan(&user.Id, &user.Country, &user.Context, &user.Level, &properties, &user.Candidate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.Properties = map[string]string{}
	err = json.Unmarshal([]byte(properties), &user.Properties)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *persistenseStorage) clear(id string) error {
	sqlStatement := `DELETE FROM users WHERE id = $1;`
	_, err := s.db.Exec(sqlStatement, id)
	return err
}

func (s *persistenseStorage) count() (int, error) {
	sqlStatement := `SELECT COUNT(*) FROM users;`
	row := s.db.QueryRow(sqlStatement)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *persistenseStorage) save(user *StorageUser) error {
	sqlStatement := `SELECT COUNT(*) FROM users WHERE id = $1;`
	row := s.db.QueryRow(sqlStatement, user.Id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	properties, err := json.Marshal(user.Properties)
	if err != nil {
		return err
	}

	if count == 0 {
		sqlStatement = `INSERT INTO users (id, country, context, level, properties, candidate) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = s.db.Exec(sqlStatement, user.Id, user.Country, user.Context, user.Level, string(properties), user.Candidate)
		return err
	}
	sqlStatement = `UPDATE users SET country=$2, context=$3, level=$4, properties=$5, candidate=$6 WHERE id = $1`
	_, err = s.db.Exec(sqlStatement, user.Id, user.Country, user.Context, user.Level, string(properties), user.Candidate)
	return err
}
