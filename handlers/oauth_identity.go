package handlers

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/db"
	"github.com/strick-j/scimplistic/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type tokenSource struct {
	ctx  context.Context
	conf *clientcredentials.Config
}

// OauthCredClient returns a client (http.Client) that has authenticated via the oauth2 client_credentials
// method.
func OauthCredClient() *oauth2.Token {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "OauthCredClient"}).Trace("Starting Oauth Token process")
	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "OauthCredClient"}).Error(err)
	}

	// Establish credential config based on user settings
	var credConfig = clientcredentials.Config{
		ClientID:     values.ClientID,
		ClientSecret: values.ClientSecret,
		TokenURL:     "https://" + values.ScimURL + "/oauth2/token/" + values.ClientAppID,
		AuthStyle:    0,
		Scopes:       []string{"scim"},
	}

	ts := &tokenSource{
		ctx:  context.Background(),
		conf: &credConfig,
	}

	ctx := context.Background()

	// Call dbTokenSource. dbTokenSource checks token in database
	// 	  - Returns database token if unexpired
	//    - Returns new token if expired
	authToken, err := ts.dbTokenSource(ctx)
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "OauthCredClient"}).Fatal("error retrieving token: ", err)
	}

	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "OauthCredClient"}).Trace("Oauth Token process completed")
	return authToken
}

// dbTokenSource attempts to retrieve token from Database
// If database is not present, table does not exist or, token does not exist
// dbTokenSource uses the clientcredentials.Config to retrieve an
// updated token. After a new token a retrieve dbTokenSource attempts to write
// new token to Database.
func (c tokenSource) dbTokenSource(ctx context.Context) (*oauth2.Token, error) {
	authToken, err := db.GetToken()
	switch {
	// No token should only occur after the database is configured but no tokens exist.
	case err != nil:
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DbTokenSource"}).Trace("No token returned, retrieving new Token")
		// No token found in database, request new token based on "Client Credentials"
		authToken, err := c.conf.Token(ctx)
		if err != nil {
			log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DbTokenSource"}).Error(err)
		}
		// AddToken as new row in database
		err = db.AddToken(authToken)
		if err != nil {
			log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DbTokenSource"}).Error(err)
		}
		return authToken, nil
	default:
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DbTokenSource"}).Trace("Token Struct returned from Database")
		// Check if token returned from Database is expired
		if time.Now().After(authToken.Expiry) {
			log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DbTokenSource"}).Trace("Token Expired, requesting new token")
			// Request new token from SCIM server using Client Credentials
			authToken, err := c.conf.Token(ctx)
			if err != nil {
				log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DbTokenSource"}).Error(err)
			}
			// AddToken as new row in database
			err = db.AddToken(authToken)
			if err != nil {
				log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DbTokenSource"}).Error(err)
			}

			return authToken, nil
		} else {
			log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DbTokenSource"}).Trace("Token is still valid. Using exiting token")
			return authToken, nil
		}
	}
}
