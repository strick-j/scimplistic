package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	types "github.com/strick-j/go-form-webserver/types"
)

//SafeAllReq is the function for requesting user info for collecting data to add a new user
func SafeAllReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res := BuildUrl("Containers", "GET")

	var bodyObject types.Safes
	json.Unmarshal(res, &bodyObject)

	for i := 0; i < len(bodyObject.Resources); i++ {
		bodyObject.Resources[i].UniqueSafeId = strconv.Itoa(i)
	}

	safeFormData := types.CreateForm{
		FormAction: "/safeaddreq/",
		FormMethod: "POST",
		FormLegend: "Add Safe Form",
		FormFields: []types.FormFields{
			{
				FieldLabel:      "SafeName",
				FieldLabelText:  "Safe Name",
				FieldInputType:  "Text",
				FieldRequired:   true,
				FieldInputName:  "FormSafeName",
				FieldInFeedback: "Safe Name is Required.",
				FieldIdNum:      1,
			},
			{
				FieldLabel:     "DisplayName",
				FieldLabelText: "Safe Display Name",
				FieldInputType: "Text",
				FieldRequired:  false,
				FieldInputName: "FormSafeDisplayName",
				FieldDescBy:    "displayHelp",
				FieldHelp:      "Optional",
				FieldIdNum:     2,
			},
			{
				FieldLabel:     "SafeDescription",
				FieldLabelText: "Description",
				FieldInputType: "Text",
				FieldRequired:  false,
				FieldInputName: "FormSafeDescription",
				FieldDescBy:    "descHelp",
				FieldHelp:      "Optional",
				FieldIdNum:     3,
			},
		},
	}

	context := types.Context{
		Navigation: "Safes",
		CreateForm: safeFormData,
		Safes:      bodyObject,
	}
	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

//UserAddForm is the form for collecting data to add a new user
func SafeAddForm(w http.ResponseWriter, r *http.Request) {

	log.Printf("Initializing Add Safe Form")

	safeFormData := types.CreateForm{
		FormAction: "/safeaddreq/",
		FormMethod: "POST",
		FormLegend: "Add Safe Form",
		FormFields: []types.FormFields{
			{
				FieldLabel:      "SafeName",
				FieldLabelText:  "Safe Name",
				FieldInputType:  "Text",
				FieldRequired:   true,
				FieldInputName:  "FormSafeName",
				FieldInFeedback: "Safe Name is Required.",
				FieldIdNum:      1,
			},
			{
				FieldLabel:     "DisplayName",
				FieldLabelText: "Safe Display Name",
				FieldInputType: "Text",
				FieldRequired:  false,
				FieldInputName: "FormSafeDisplayName",
				FieldDescBy:    "displayHelp",
				FieldHelp:      "Optional",
				FieldIdNum:     2,
			},
			{
				FieldLabel:     "SafeDescription",
				FieldLabelText: "Description",
				FieldInputType: "Text",
				FieldRequired:  false,
				FieldInputName: "FormSafeDescription",
				FieldDescBy:    "descHelp",
				FieldHelp:      "Optional",
				FieldIdNum:     3,
			},
		},
	}

	context := types.Context{
		Navigation:         "Add Safe",
		SettingsConfigured: false,
		CreateForm:         safeFormData,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

//UserDelForm is the form for deleting a user
func SafeDelForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "safedelform.html", nil)
}

//SafeAddReq is used to add users from the /useraddreq URL
func SafeAddReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Println("Reading Data from Safe Add Form")
	safeName := r.FormValue("FormSafeName")
	displayName := r.FormValue("FormSafeDisplayName")
	description := r.FormValue("FormSafeDescription")
	scimschema := []string{"urn:ietf:params:scim:schemas:pam:1.0:Container"}

	addSafeData := types.PostObjectRequest{
		Name:        safeName,
		DisplayName: displayName,
		Description: description,
		Schemas:     scimschema,
	}

	// Required as placeholder
	blankstruct := types.PostUserRequest{}

	res := ScimAPI("Containers", "POST", addSafeData, blankstruct)

	fmt.Println(string(res))
	// Redirect back to all safes
	http.Redirect(w, r, "/allsafes/", http.StatusFound)
}
