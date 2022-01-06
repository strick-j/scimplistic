package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	config "github.com/strick-j/scimplistic/config"
	types "github.com/strick-j/scimplistic/types"
	utils "github.com/strick-j/scimplistic/utils"
)

func SettingsForm(w http.ResponseWriter, r *http.Request) {
	values, err := config.ReadConfig("config.json")
	if err != nil {
		log.Println("ERROR SettingsForm:", err)
	}

	settingsFormData := types.CreateForm{
		FormEncType: "multipart/form-data",
		FormAction:  "/configuresettings/",
		FormMethod:  "POST",
		FormLegend:  "Configure Settings",
		FormRole:    "configuresettings",
		FormFields: []types.FormFields{
			{
				FieldLabel:      "scimurl",
				FieldLabelText:  "SCIM Endpoint URL",
				FieldInputType:  "Text",
				FieldRequired:   true,
				FieldInputName:  "FormSCIMURL",
				FieldInFeedback: "SCIM Endpoint URL is Required.",
				FieldPlaceHold:  values.ScimURL,
				FieldIdNum:      1,
			},
			{
				FieldLabel:     "serverHostname",
				FieldLabelText: "Scimplistic Server URL",
				FieldInputType: "Text",
				FieldRequired:  false,
				FieldInputName: "FormServerURL",
				FieldPlaceHold: values.ServerURL,
				FieldIdNum:     2,
			},
			{
				FieldLabel:     "serverCertFile",
				FieldLabelText: "Server Certificate for TLS",
				FieldInputType: "file",
				FieldRequired:  false,
				FieldInputName: "FormServerCert",
				FieldIdNum:     3,
			},
			{
				FieldLabel:     "serverCertKey",
				FieldLabelText: "Server Private Key for TLS",
				FieldInputType: "file",
				FieldRequired:  false,
				FieldInputName: "FormServerKey",
				FieldIdNum:     4,
			},
		},
	}

	context := types.Context{
		Navigation: "Settings",
		CreateForm: settingsFormData,
		Token:      values.AuthToken,
	}

	if values.ScimURL != "" {
		context.SettingsConfigured = true
	} else {
		context.SettingsConfigured = false
	}

	if values.EnableHTTPS {
		context.HTTPSEnabled = true
	} else {
		context.HTTPSEnabled = false
	}

	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

func ConfigureSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Read data from Configure Settings Form
	log.Println("INFO ConfigureSettings: Reading Data from Configure Settings Form")
	r.ParseMultipartForm(10 << 20)
	fmt.Printf("%+v\n", r.Form)
	scimURL := r.FormValue("FormSCIMURL")
	scimToken := r.FormValue("FormOathToken")
	serverURL := r.FormValue("FormServerURL")
	enableHTTPS := r.Form["FormEnableHTTPS"][0]

	certFile, handler, err := r.FormFile("FormServerCert")
	if err != nil {
		fmt.Println("ERROR ConfigureSettings:Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer certFile.Close()
	log.Printf("INFO ConfigureSettings: Uploaded File: %+v\n", handler.Filename)

	tempCert, err := ioutil.TempFile("files", "cert-*.crt")
	if err != nil {
		fmt.Println(err)
	}
	defer tempCert.Close()

	certBytes, err := ioutil.ReadAll(certFile)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempCert.Write(certBytes)

	// Repeat above for Private Key
	keyFile, handler, err := r.FormFile("FormServerKey")
	if err != nil {
		fmt.Println("ERROR ConfigureSettings: Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer certFile.Close()
	log.Printf("INFO ConfigureSettings: Uploaded File: %+v\n", handler.Filename)

	tempKey, err := ioutil.TempFile("files", "key-*.key")
	if err != nil {
		fmt.Println(err)
	}
	defer tempKey.Close()

	keyBytes, err := ioutil.ReadAll(keyFile)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempKey.Write(keyBytes)

	// Set configuration info in json file
	configSettings := types.ConfigSettings{
		ScimURL:   scimURL,
		AuthToken: scimToken,
		PrevConf:  true,
		ServerURL: serverURL,
		KeyName:   tempKey.Name(),
		CertName:  tempCert.Name(),
	}

	// Check HTTPS Setting and set in config file
	if enableHTTPS == "1" {
		configSettings.EnableHTTPS = true
		log.Println("INFO ConfigureSettings: Enable HTTPS set to true")
		// Enable listening on 8443 and redirect to 8443 from 8080

		certPath := tempCert.Name()
		keyPath := tempKey.Name()

		utils.startTlsListen(certPath, keyPath)

	}

	file, err := json.MarshalIndent(configSettings, "", "   ")
	if err != nil {
		log.Println("ERROR ConfigureSettings:", err)
	}
	err = ioutil.WriteFile("config.json", file, 0644)
	if err != nil {
		log.Println("ERROR ConfigureSettings:", err)
	}

	log.Println("INFO ConfigureSettings: Configuration File written.")

	// Redirect back to settings.
	http.Redirect(w, r, "/settings/", http.StatusFound)
}
