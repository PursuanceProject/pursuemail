package main

import (
	"database/sql"
	"strings"

	log "github.com/Sirupsen/logrus"
	emailLib "github.com/jordan-wright/email"
	_ "github.com/lib/pq"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	// TODO - Handle basic signals
	addr := "127.0.0.1:8080"

	pgUser := "pursuemail"
	pgHost := "127.0.0.1:5432"
	database := "pursuemail"
	dbUrl := BuildPGUrl(pgUser, pgHost, database)
	db := MustGetDb(dbUrl)
	defer db.Close()

	// TODO - Setup configuration for Mailgun or similar
	emailPool := emailLib.NewPool("localhost:1025", 5, nil)
	defer emailPool.Close()

	srv := NewServer(addr, db, emailPool)
	log.Infof("Listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func BuildPGUrl(pgUser, pgHost, databaseName string) string {
	return strings.Join([]string{"postgres://", pgUser, "@", pgHost, "/", databaseName, "?sslmode=disable"}, "")
}

func MustGetDb(dbUrl string) *sql.DB {
	log.Debugf("Creating new db postgres instance @ `%s`", dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Could not open postgres db @ `%s`", dbUrl)
	}
	if err = db.Ping(); err != nil {
		defer db.Close()
		log.Fatalf("Error connecting to db. Err: %s", err)
	}
	return db
}
