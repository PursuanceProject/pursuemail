package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// TODO - Make these configurable
var (
	GPG_DIR                  = os.Getenv("HOME") + "/.gnupg/"
	PUBLIC_KEYRING_FILENAME  = GPG_DIR + "pubring.gpg"
	PRIVATE_KEYRING_FILENAME = GPG_DIR + "secring.gpg"
)

type Buffer []byte

func (buf *Buffer) Write(p []byte) (int, error) {
	*buf = append(*buf, p...)
	return len(p), nil
}

func encryptEmailBody(from, to, body string) (enc []byte, err error) {
	var buf Buffer

	privKey, err := GetEntityFrom(from, PRIVATE_KEYRING_FILENAME)
	if err != nil {
		return nil, fmt.Errorf("Error getting private key for %s from %s: %v",
			from, PRIVATE_KEYRING_FILENAME, err)
	}

	pubKey, err := GetEntityFrom(to, PUBLIC_KEYRING_FILENAME)
	if err != nil {
		return nil, fmt.Errorf("Error getting public key for %s from %s: %v",
			to, PRIVATE_KEYRING_FILENAME, err)
	}

	// Produce new writer to... write encrypted messages to?
	w, err := armor.Encode(&buf, "PGP MESSAGE", nil)
	if err != nil {
		return nil, fmt.Errorf("Error from armor.Encode: %v", err)
	}
	defer w.Close()

	// Encrypt message from ME to recipient
	plaintext, err := openpgp.Encrypt(w, []*openpgp.Entity{pubKey},
		privKey, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("Error from openpgp.Encrypt: %v", err)
	}
	defer plaintext.Close()

	// Write message to `plaintext` WriteCloser
	_, err = fmt.Fprintf(plaintext, body)
	if err != nil {
		return nil, fmt.Errorf("Error writing to plaintext: %v", err)
	}

	return []byte(buf), nil
}

// TODO: We can make this a lot better. Memoization, etc
func GetEntityFrom(email, sourceFile string) (*openpgp.Entity, error) {
	ringFile, err := os.Open(sourceFile)
	if err != nil {
		return nil, err
	}
	defer ringFile.Close()

	ring, err := openpgp.ReadKeyRing(ringFile)
	if err != nil {
		return nil, err
	}

	var key *openpgp.Entity
	for _, entity := range ring {
		for _, ident := range entity.Identities {
			if ident.UserId.Email == email {
				key = entity
			}
		}
	}

	if key == nil {
		e := fmt.Errorf("Couldn't find key for user %s: %v", email, err)
		return nil, e
	}

	return key, nil
}
