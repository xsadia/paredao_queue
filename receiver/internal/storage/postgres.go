package storage

import (
	"database/sql"
	"fmt"
)

func NewDB(host, user, password, dbname string) (*sql.DB, error) {
	connectionString :=
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
