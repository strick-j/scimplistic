package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	config "github.com/strick-j/scimplistic/config"
	types "github.com/strick-j/scimplistic/types"
)

func SettingsForm(w http.ResponseWriter, r *http.Request) {
	values, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}

	settingsFormData := types.CreateForm{
		FormAction: "/configuresettings/",
		FormMethod: "POST",
		FormLegend: "Configure Settings",
		FormRole:   "configuresettings",
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

	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

func ConfigureSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Read data from Configure Settings Form
	log.Println("Reading Data from Configure Settings Form")
	scimURL := r.FormValue("FormSCIMURL")
	scimToken := r.FormValue("FormOathToken")

	configSettings := types.ConfigSettings{
		ScimURL:   scimURL,
		AuthToken: scimToken,
		PrevConf:  true,
	}

	file, err := json.MarshalIndent(configSettings, "", "   ")
	if err != nil {
		log.Println("ConfigureSettings:", err)
	}
	err = ioutil.WriteFile("config.json", file, 0644)
	if err != nil {
		log.Println("ConfigureSettings:", err)
	}

	log.Println("ConfigureSettings: Configuration File written.")

	// Redirect back to settings.
	http.Redirect(w, r, "/settings/", http.StatusFound)

}
