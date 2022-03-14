package handlers

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// OauthCredClient returns a client (http.Client)that has authenticated via the oauth2 client_credentials
// method
func OauthCredClient() *oauth2.Token {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "OauthCredClient"}).Trace("Starting Oauth Token process")
	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "OauthCredClient"}).Error(err)
	}

	var credConfig = clientcredentials.Config{
		ClientID:     values.ClientID,
		ClientSecret: values.ClientSecret,
		TokenURL:     "https://" + values.ScimURL + "/oauth2/token/" + values.ClientAppID,
		AuthStyle:    0,
		Scopes:       []string{"scim"},
	}

	ctx := context.Background()

	authToken, err := credConfig.Token(ctx)
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "OauthCredClient"}).Error(err)
	}

	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "OauthCredClient"}).Trace("Oauth Token process completed")
	return authToken
}
