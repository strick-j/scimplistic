package views

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/internal/types"
)

// Generate Struct for the forms required by Safe Functions
// Form Data is used by several functions (add, get, update, etc...)
var safeFormData = types.CreateForm{
	FormAction: "/safes/add",
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

var safeScimSchema = []string{"urn:ietf:params:scim:schemas:pam:1.0:Container"}

///////////////////////// Safe Default Handler /////////////////////////

// SafesHandler is the function for displaying the basic Safes page.
func SafesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafesHandler"}).Info("Starting Safe retrieval Process")

	ob := Object{
		Type:   "containers",
		Method: "GET",
	}

	res, _, err := ob.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafesHandler"}).Error(err)
		return
	}

	for i := 0; i < len(res.Resources); i++ {
		res.Resources[i].UniqueSafeId = strconv.Itoa(i)
	}

	// Establish context for populating allinfo template
	context := types.Context{
		Navigation: "Safes",
		CreateForm: safeFormData,
		Safes:      *res,
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafesHandler"}).Info("Safe retrieval process completed without error")

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

///////////////////////// Safe Action Handlers /////////////////////////

func SafeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	// Parse Vars from URL Variables
	vars := mux.Vars(r)

	ob := Object{
		Type:   "containers",
		Method: "GET",
		Id:     vars["id"],
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeHandler"}).Trace("Derived Request Id: ", ob.Id)

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeHandler"}).Trace("Calling ScimType2Api to retrieve Safe info")
	_, res, err := ob.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeHandler"}).Error(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeHandler"}).Info("Safe retrieved: ", res.DisplayName)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// SafeAddHandler reads in form data from the modal that is triggered
// when the add safe button is pressed. This function calls the SCIM
// function which submits the Delete action to the SCIM Endpoint.
func SafeAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	redirectURL := GetRedirectUrl(r.Referer())

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeAddHandler"}).Info("Starting Add Safe Process")

	ob := Object{
		Type:   "containers",
		Method: "POST",
	}

	ob.Type2Resources = types.Type2Resources{
		Name:        r.FormValue("FormSafeName"),
		DisplayName: r.FormValue("FormSafeDisplayName"),
		Description: r.FormValue("FormSafeDescription"),
		Schemas:     safeScimSchema,
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeAddHandler"}).Trace("Calling ScimType2Api to Add Safe")
	_, res, err := ob.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeAddHandler"}).Error(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeAddHandler"}).Info("Safe add process completed. Safe Added: ", res.Name)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// SafeDelHandler reads in form data from the modal that is triggered
// when the delete button for a particular safe is pressed. This function
// calls the SCIM function which submits the Delete action to the SCIM Endpoint.
// Note: Safe deletions are based off of safe Name not ID like users or groups
func SafeDelHandler(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := GetRedirectUrlNoId(r.Referer())

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeDelHandler"}).Info("Starting Safe Delete Process")

	// Retrieve SafeID from URL to send to Del Function
	vars := mux.Vars(r)

	ob := Object{
		Type:   "containers",
		Method: "DELETE",
		Id:     vars["id"],
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeDelHandler"}).Trace("Safe Id: ", ob.Id)

	// Delete Safe and recieve response from Delete Function
	_, _, err := ob.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeDelHandler"}).Error(err)
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "SafeDelHandler"}).Info("Safe Delete Process completed")

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// SafeUpdateHandler reads in form data from the modal that is triggered
// when the Update button for a particular safe is pressed. This function
// calls the SCIM function which submits the UPDATE action to the SCIM Endpoint.
func SafeUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from Update Form

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}

///////////////////////// Safe Form Handlers /////////////////////////

// SafeAddForm is the form for collecting data to add a new Safe
func SafeAddForm(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"Category": "Form Handler", "Function": "SafeAddForm"}).Trace("Initializing Add Safe Form")

	// Establish context for populating add safe template
	context := types.Context{
		Navigation: "Add Safe",
		CreateForm: safeFormData,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}
