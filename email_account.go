package main

import (
	"database/sql"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
	emailLib "github.com/jordan-wright/email"
)

type EmailAccount struct {
	Id        string    `json:"id,omitempty"`
	Email     string    `json:"email"`
	PubKey    string    `json:"pubkey,omitempty"`
	HasPubKey bool      `json:"has_pubkey,omitempty"`
	Created   time.Time `json:"created,omitempty"`
}

func GetEmailAccount(db *sql.DB, id string) (*EmailAccount, error) {
	emailAccount := &EmailAccount{}
	err := db.QueryRow(`
		SELECT
			id, email, has_pubkey, created
		FROM
			email_account
		WHERE
		id = $1
	`, id).Scan(
		&emailAccount.Id, &emailAccount.Email,
		&emailAccount.HasPubKey, &emailAccount.Created,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("Error getting email_account. Err: %s", err)
		}
		return nil, err
	}

	return emailAccount, nil
}

func GetEmailAccounts(db *sql.DB, ids []string) ([]*EmailAccount, error) {
	rows, err := db.Query(`
		SELECT
			id, email, has_pubkey, created
		FROM
			email_account
		WHERE
			id in ($1)
	`, ids)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("Error getting email_accounts. Err: %s", err)
		}
		return nil, err
	}
	defer rows.Close()

	emailAccounts := []*EmailAccount{}
	for rows.Next() {
		var ea EmailAccount

		err := rows.Scan(&ea.Id, &ea.Email, &ea.HasPubKey, &ea.Created)

		if err != nil {
			log.Errorf("Error with scan. Err: %v", err)
			return nil, err
		}

		emailAccounts = append(emailAccounts, &ea)
	}

	return emailAccounts, nil
}

// Save Email and PubKey, attach Id that is returned.
func (e *EmailAccount) Save(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("Error beginning transaction. Err: %s", err)
		return err
	}

	if e.PubKey != "" {
		err = importPublicKey(e.PubKey)
		if err != nil {
			log.Errorf("Error importing public key. Err: %s", err)
			return err
		}
		e.HasPubKey = true
	}

	err = tx.QueryRow(`
		INSERT INTO email_account(email, has_pubkey)
		VALUES ($1, $2)
		RETURNING id, created
	`, e.Email, e.HasPubKey,
	).Scan(&e.Id, &e.Created)

	if err != nil {
		log.Errorf("Error adding email_account. Err: %s", err)
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Errorf("Got error rolling back transaction. Err: %s", rollbackErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("Error committing transaction. Err: %s", err)
		return err
	}
	return nil
}

func (e *EmailAccount) Send(emailData EmailData, emailPool *emailLib.Pool) error {
	sendableEmail := emailData.toSendableEmail()

	if e.HasPubKey {
		encryptedMsg, err := encryptEmailBody(sendableEmail.From, e.Email, string(sendableEmail.Text))
		if err != nil {
			log.Errorf("Error encrypting message: %v\n", err)
			return err
		}
		sendableEmail.Text = encryptedMsg
	}

	sendableEmail.To = []string{e.Email}

	// TODO - Make timeout configurable?
	return emailPool.Send(sendableEmail, 15*time.Second)
}

type EmailData struct {
	// TODO: Have a default from email
	From    string `json:"from,omitempty"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (ed EmailData) toSendableEmail() *emailLib.Email {
	em := emailLib.NewEmail()

	em.From = ed.From
	em.Subject = ed.Subject
	em.Text = []byte(ed.Body)

	return em
}

func importPublicKey(pubkey string) error {
	// TODO - make tempfile directory configurable
	tmpfile, err := ioutil.TempFile("", "pubkey-import")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	// Save pubkey to temp file
	if _, err := tmpfile.Write([]byte(pubkey)); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	cmd := "gpg"
	args := []string{"--import", tmpfile.Name()}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return err
	}

	return nil
}
