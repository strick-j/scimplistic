package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	config "github.com/strick-j/scimplistic/config"
	types "github.com/strick-j/scimplistic/types"
	utils "github.com/strick-j/scimplistic/utils"
)

//////////////////////// Settings Default Handler /////////////////////////

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	values, err := config.ReadConfig("settings.json")
	if err != nil {
		log.Println("ERROR SettingsHandler:", err)
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
		FormEncType: "multipart/form-data",
		FormAction:  "/settings/secret",
		FormMethod:  "POST",
		FormLegend:  "Secrets Settings",
		FormRole:    "secretsettings",
		FormFields: []types.FormFields{
			{
				FieldLabel:     "oauthToken",
				FieldLabelText: "Path to Auth Token Secret",
				FieldInputType: "Text",
				FieldRequired:  false,
				FieldInputName: "FormOAUTHPath",
				FieldIdNum:     1,
			},
			{
				FieldLabel:     "postgresUsername",
				FieldLabelText: "Username for PostreSQL server",
				FieldInputType: "Text",
				FieldRequired:  false,
				FieldInputName: "psqlUname",
				FieldIdNum:     2,
			},
			{
				FieldLabel:     "postgresPassword",
				FieldLabelText: "Path to PostgreSQL Password",
				FieldInputType: "text",
				FieldRequired:  false,
				FieldInputName: "psqlPwd",
				FieldIdNum:     3,
			},
		},
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

// SettingsTypeHandler will decide what is required based on the type assigned
// gen: Proceed to SettingsGenHandler
// sec: Proceed to SettingsSecretlHandler
func SettingsTypeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Extract action from URL using mux.Vars.
	vars := mux.Vars(r)
	action := vars["action"]
	log.Println("INFO SettingsActionHandler: Action = ", action)

	// Switch to appropriate handler based on action type.
	switch action {
	case "general":
		log.Println("INFO SettingsActionHandler: Calling SettingsGenHandler")
		SettingsGenHandler(w, r)
	case "secret":
		log.Println("INFO UsersActionHandler: Calling SettingSecretHandler")
		SettingsSecretHandler(w, r)
	}
}

func SettingsGenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// initialize  Variables
	var serverIP, ext string
	var serverPort int
	var fileUpload [2]string

	// Read data from Configure Settings Form
	log.Println("INFO ConfigureSettings: Reading Data from Configure Settings Form")

	// Parse Form Data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user updated Server Port and set it
	if len(r.FormValue("FormServerIP")) != 0 {
		serverIP = r.FormValue("FormServerIP")
	}

	// Check if user updated Server Port and set it
	if len(r.FormValue("FormServerPort")) != 0 {
		serverPort, _ = strconv.Atoi(r.FormValue("FormServerPort"))
	}

	// Read in standard strings, these fields are required in the form
	scimURL := r.FormValue("FormSCIMURL")     // Required
	scimToken := r.FormValue("FormOathToken") // Required
	serverURL := r.FormValue("FormServerURL") // Required

	// Check for Files using range to iterate upload fields
	index := 0
	for fileHeader := range r.MultipartForm.File {
		// If a file header is found, process the upload
		if fileHeader != "" {
			File, handler, err := r.FormFile(fileHeader)
			if err != nil {
				fmt.Println("ERROR ConfigureSettings: Error Retrieving the File")
				fmt.Println(err)
				return
			}
			defer File.Close()

			log.Printf("INFO ConfigureSettings: Uploading File: %+v\n", handler.Filename)
			// Catch filename file extenstion for filename
			fileExt := strings.Split(handler.Filename, ".")
			if ext = "pem"; fileExt[1] == "pem" || fileExt[1] == "crt" {
				log.Println("INFO ConfigureSettings: Detected File Extension:", ext)
			} else if ext = "key"; fileExt[1] == "key" {
				log.Println("INFO ConfigureSettings: Detected File Extension:", ext)
			} else {
				log.Printf("Error ConfigurationSettings: Could not detect file extension")
				http.Error(w, "The provided file format is not allowed. Valid Formats are .crt, .pem, and .key", http.StatusBadRequest)
				return
			}

			tempFile, err := ioutil.TempFile("files", fileHeader+"-*."+ext)
			if err != nil {
				fmt.Println(err)
			}
			fileUpload[index] = tempFile.Name()
			defer tempFile.Close()

			fileBytes, err := ioutil.ReadAll(File)
			if err != nil {
				fmt.Println(err)
			}
			// write this byte array to our temporary file
			tempFile.Write(fileBytes)

			log.Printf("INFO ConfigureSettings: Uploaded Filename: %+v\n", fileUpload[index])
			index++
		}
	}

	// Read in current config settings
	configSettings, err := config.ReadConfig("settings.json")
	if err != nil {
		log.Println("ERROR Main:", err)
	}

	// Update configuration info in json file
	configSettings = types.ConfigSettings{
		ScimURL:     scimURL,
		AuthToken:   scimToken,
		PrevConf:    true,
		IP:          serverIP,
		Port:        serverPort,
		ServerName:  serverURL,
		HostName:    serverURL,
		CertFile:    fileUpload[0],
		PrivKeyFile: fileUpload[1],
		OriginOnly:  false,
		TLS:         false,
	}

	// Check if user enabled TLS.
	// If user enabled TLS validate cert and key file
	if len(r.Form.Get("FormEnableHTTPS")) != 0 {
		log.Println("INFO ConfigureSettings: Enable TLS Selected, checking for required cert and private key")
		if _, err := CheckTLS(fileUpload[0], fileUpload[1]); err != nil {
			log.Println("ERROR ConfigureSettings: Error Verifying Certificate and PrivateKey")
			configSettings.TLS = false
		} else {
			log.Println("INFO ConfigureSettings: Verified Certificate and PrivateKey are available")
			configSettings.TLS = true
		}
	}

	file, err := json.MarshalIndent(configSettings, "", "   ")
	if err != nil {
		log.Println("ERROR ConfigureSettings:", err)
	}
	err = ioutil.WriteFile("settings.json", file, 0644)
	if err != nil {
		log.Println("ERROR ConfigureSettings:", err)
	}

	log.Println("INFO ConfigureSettings: Configuration File written.")

	if configSettings.TLS {
		log.Printf("INFO ConfigureSettings: Attempting to stop HTTP Listener and restart via HTTPS")
	} else {
		log.Printf("INFO ConfigureSettings: Attempting to stop HTTP Listener and restart")
	}

	if sdErr := utils.ShutDown(); sdErr != nil {
		log.Println(sdErr.Error())
	}
}

func SettingsSecretHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Update application to pull secrets from secret store
	tpl.ExecuteTemplate(w, "/", nil)
}
