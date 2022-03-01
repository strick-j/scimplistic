package views

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/types"
	"github.com/strick-j/scimplistic/utils"
)

//////////////////////// Settings Default Handler /////////////////////////

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"Category": "Configuration",
		"Function": "SettingsHandler",
	})

	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		logger.Error(err)
	}

	settingsFormData := types.CreateForm{
		FormEncType: "multipart/form-data",
		FormAction:  "/settings/general",
		FormMethod:  "POST",
		FormLegend:  "General Settings",
		FormRole:    "configuresettings",
		FormFields: []types.FormFields{
			{
				FieldLabel:     "serverHostname",
				FieldLabelText: "Scimplistic Server URL",
				FieldInputType: "Text",
				FieldRequired:  false,
				FieldInputName: "FormServerURL",
				FieldPlaceHold: values.ServerName,
				FieldIdNum:     1,
			},
			{
				FieldLabel:     "serverCertFile",
				FieldLabelText: "Server Certificate for TLS",
				FieldInputType: "file",
				FieldRequired:  false,
				FieldInputName: "ServerCert",
				FieldIdNum:     2,
			},
			{
				FieldLabel:     "serverCertKey",
				FieldLabelText: "Server Private Key for TLS",
				FieldInputType: "file",
				FieldRequired:  false,
				FieldInputName: "ServerKey",
				FieldIdNum:     3,
			},
			{
				FieldLabel:      "scimurl",
				FieldLabelText:  "SCIM Endpoint URL",
				FieldInputType:  "Text",
				FieldRequired:   true,
				FieldInputName:  "FormSCIMURL",
				FieldInFeedback: "SCIM Endpoint URL is Required.",
				FieldPlaceHold:  values.ScimURL,
				FieldIdNum:      4,
			},
		},
	}

	secretFormData := types.CreateForm{
		FormAction: "/settings/secrets",
		FormMethod: "POST",
		FormLegend: "Secrets Settings",
		FormRole:   "secretsettings",
	}

	context := types.Context{
		Navigation: "Settings",
		CreateForm: settingsFormData,
		SecretForm: secretFormData,
		Token:      values.AuthToken,
	}

	if values.PrevConf {
		context.SettingsConfigured = true
	} else {
		context.SettingsConfigured = false
	}

	if values.TLS {
		context.HTTPSEnabled = true
	} else {
		context.HTTPSEnabled = false
	}

	tpl.ExecuteTemplate(w, "settings.html", context)
}

///////////////////////// Settings Type Handlers /////////////////////////

func GeneralSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	// Initialize logger for GeneralSettinsHandler
	logger := log.WithFields(log.Fields{
		"Category": "Server Processes",
		"Function": "SettingsGenHandler",
	})

	logger.Info("Starting General Settings Process")

	// Read in current config settings
	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		logger.Error(err)
	}

	// initialize  Variables
	var ext string
	var fileUpload [2]string

	// Read data from Configure Settings Form
	logger.Info("Reading Data from General Settings Form")

	// Parse Form Data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user updated Server Port and set it
	if len(r.FormValue("FormServerIP")) != 0 {
		values.IP = r.FormValue("FormServerIP")
	}

	// Check if user updated Server Port and set it
	if len(r.FormValue("FormServerPort")) != 0 {
		values.Port, _ = strconv.Atoi(r.FormValue("FormServerPort"))
	}

	// Read in standard strings, these fields are required in the form
	values.ScimURL = r.FormValue("FormSCIMURL")    // Required
	values.HostName = r.FormValue("FormServerURL") // Required

	// Check for Files using range to iterate upload fields
	index := 0
	for fileHeader := range r.MultipartForm.File {
		// If a file header is found, process the upload
		if fileHeader != "" {
			File, handler, err := r.FormFile(fileHeader)
			if err != nil {
				logger.Error(err)
				return
			}
			defer File.Close()

			logger.Trace("Uploading File: %+v\n", handler.Filename)
			// Catch filename file extenstion for filename
			fileExt := strings.Split(handler.Filename, ".")
			if ext = "pem"; fileExt[1] == "pem" || fileExt[1] == "crt" {
				logger.Debug("Detected File Extension:", ext)
			} else if ext = "key"; fileExt[1] == "key" {
				logger.Debug("Detected File Extension:", ext)
			} else {
				logger.Error("Could not detect file extension")
				http.Error(w, "The provided file format is not allowed. Valid Formats are .crt, .pem, and .key", http.StatusBadRequest)
				return
			}

			tempFile, err := ioutil.TempFile("files", fileHeader+"-*."+ext)
			if err != nil {
				logger.Error(err)
			}
			fileUpload[index] = tempFile.Name()
			defer tempFile.Close()

			fileBytes, err := ioutil.ReadAll(File)
			if err != nil {
				logger.Error(err)
			}
			// write this byte array to our temporary file
			tempFile.Write(fileBytes)

			logger.Trace("Uploaded Filename: %+v\n", fileUpload[index])
			index++
		}
	}
	values.CertFile = fileUpload[0]
	values.PrivKeyFile = fileUpload[1]

	// Set Log Level - Default is "info"
	values.LogLevel = r.FormValue("FormLogLevel")
	values.OriginOnly = false
	values.PrevConf = true

	// Check if user enabled TLS.
	// If user enabled TLS validate cert and key file
	if len(r.Form.Get("FormEnableHTTPS")) != 0 {
		logger.Trace("Enable TLS Selected, checking for required cert and private key")
		if _, err := CheckTLS(fileUpload[0], fileUpload[1]); err != nil {
			logger.Error(err)
			values.TLS = false
		} else {
			logger.Trace("Verified Certificate and PrivateKey are available")
			values.TLS = true
		}
	}

	file, err := json.MarshalIndent(values, "", "   ")
	if err != nil {
		logger.Error(err)
	}
	err = ioutil.WriteFile("settings.json", file, 0644)
	if err != nil {
		logger.Error(err)
	}

	log.Println("INFO ConfigureSettings: Configuration File written.")

	if values.TLS {
		logger.Info("TLS configured - Restart to enable")
	} else {
		logger.Info("TLS not configured")
	}

	//if sdErr := utils.ShutDown(); sdErr != nil {
	//	log.Println(sdErr.Error())
	//}

	logger.Info("General Settings Process completed")
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// SettingsSecretHandler reads in data and sets configurations for the
// secrets page.
func SecretSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	logger := log.WithFields(log.Fields{
		"Category": "Server Processes",
		"Function": "SecretSettingsHandler",
	})

	logger.Info("Starting Secret Settings Process")

	// Read data from Secret Settings Form
	logger.Trace("Reading Data from Secret Settings Form")

	// Parse Form Data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read in current config settings
	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		logger.Error(err)
	}

	// Check if user updated Client ID and set it
	if len(r.FormValue("FormClientID")) != 0 {
		values.ClientID = r.FormValue("FormClientID")
	}

	// Check if user updated Client Secret and set it
	if len(r.FormValue("FormClientSecret")) != 0 {
		values.ClientSecret = r.FormValue("FormClientSecret")
	}

	// Check if user updated Client Secret and set it
	if len(r.FormValue("FormClientAppID")) != 0 {
		values.ClientAppId = r.FormValue("FormClientAppID")
	}

	// Check if user updated Oauth Token and set it
	if len(r.FormValue("FormOauthToken")) != 0 {
		values.AuthToken = r.FormValue("FormOauthToken")
	}

	file, err := json.MarshalIndent(values, "", "   ")
	if err != nil {
		logger.Error(err)
	}
	err = ioutil.WriteFile("settings.json", file, 0644)
	if err != nil {
		logger.Error(err)
	}

	logger.Info("Configuration File written.")

	//if sdErr := utils.ShutDown(); sdErr != nil {
	//	log.Println(sdErr.Error())

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
