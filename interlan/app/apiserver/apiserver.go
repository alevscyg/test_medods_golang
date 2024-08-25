package apiserver

import (
	"database/sql"
	"medos/interlan/app/apiserver/config"
	sqlstorage "medos/storage/sqlStorage"
	"net/http"
)

// Start ...
func Start(config *config.Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	storage := sqlstorage.New(db)
	srv := newServer(storage)

	return http.ListenAndServe(config.PORT, srv.router)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
