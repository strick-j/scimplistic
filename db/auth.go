package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

////// Oauth Token Functions

// AddToken adds a new token to the "access" database
func AddToken(t *oauth2.Token) error {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "AddToken"}).Trace("Starting AddToken process")
	err = genQuery(`insert into "auth"("access_token", "token_type", "expiry") values($1,$2,$3)`, t.AccessToken, t.TokenType, t.Expiry)
	if err != nil {
		return err
	}
	return nil
}

// GetToken returns the latest token from the "access" database
func GetToken() (*oauth2.Token, error) {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GetToken"}).Trace("Starting GetToken process")
	var (
		id           int
		access_token string
		token_type   string
		expiry       time.Time
	)

	// Build query for latest Auth Token, select last row
	q := `
	select * from auth
	order by "id"
	desc limit 1
	`

	// User QueryRow to to return single row
	err := database.db.QueryRow(q).Scan(&id, &access_token, &token_type, &expiry)
	switch {
	case err == sql.ErrNoRows:
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GetToken"}).Error(err)
		return nil, err
	case err != nil:
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GetToken"}).Error(err)
		return nil, err
	default:
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GetToken"}).Trace("Token retrieved, recreating token struct")
	}

	// Reconstruct Oauth2.Token for return
	t := oauth2.Token{
		AccessToken: access_token,
		TokenType:   token_type,
		Expiry:      expiry,
	}

	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GetToken"}).Trace("GetToken process completed")
	return &t, nil
}
