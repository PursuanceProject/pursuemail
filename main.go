package main

import (
	"database/sql"
	"net/smtp"
	"os"
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
	addr := "127.0.0.1:9080"

	pgUser := "pursuemail"
	pgHost := "127.0.0.1:5432"
	database := "pursuemail"
	dbUrl := BuildPGUrl(pgUser, pgHost, database)
	db := MustGetDb(dbUrl)
	defer db.Close()

	// TODO - Setup configuration for Mailgun or similar
	emailPool, err := emailLib.NewPool(os.Getenv("SMTP_SERVER"), 5,
		smtp.PlainAuth("", os.Getenv("SMTP_LOGIN"),
			os.Getenv("SMTP_PASSWORD"), strings.SplitN(os.Getenv("SMTP_SERVER"), ":", 2)[0]))
	if err != nil {
		log.Fatalf("Error from email.NewPool: %v", err)
	}
	defer emailPool.Close()

	srv := NewServer(addr, db, emailPool)
	log.Infof("Listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func BuildPGUrl(pgUser, pgHost, databaseName string) string {
	password := os.Getenv("PGPASSWORD")
	dbUrl := "postgres://" + pgUser + ":" + password + "@" +
		pgHost + "/" + databaseName + "?sslmode=disable"
	log.Debugf("Building new Postgres URL `%s`",
		strings.Replace(dbUrl, password, "${PGPASSWORD}", -1))
	return dbUrl
}

func MustGetDb(dbUrl string) *sql.DB {
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
