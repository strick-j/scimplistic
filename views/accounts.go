package views

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/types"
)

// Generate Struct for the forms required by Account Functions
// Form Data is used by several functions (add, get, update, etc...)
var accountFormData = types.CreateForm{
	FormAction: "/accounts/add",
	FormMethod: "POST",
	FormLegend: "Add Account Form",
	FormFields: []types.FormFields{
		{
			FieldLabel:      "AccountName",
			FieldLabelText:  "Account Name",
			FieldInputType:  "Text",
			FieldRequired:   true,
			FieldInputName:  "FormAccountName",
			FieldInFeedback: "Account Name is Required.",
			FieldIdNum:      1,
		},
		{
			FieldLabel:     "DisplayName",
			FieldLabelText: "Account Display Name",
			FieldInputType: "Text",
			FieldRequired:  false,
			FieldInputName: "FormAccountDisplayName",
			FieldDescBy:    "displayHelp",
			FieldHelp:      "Optional",
			FieldIdNum:     2,
		},
		{
			FieldLabel:     "AccountDescription",
			FieldLabelText: "Description",
			FieldInputType: "Text",
			FieldRequired:  false,
			FieldInputName: "FormAccountDescription",
			FieldDescBy:    "descHelp",
			FieldHelp:      "Optional",
			FieldIdNum:     3,
		},
	},
}

// Generate Struct for actions required by User functions
var accountObject = Object{
	Type: "privilegedData",
}

///////////////////////// Account Default Handler /////////////////////////

// AccountsHandler is the function for displaying the basic Safes page.
func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountsHandler"}).Info("Starting Accounts retrieval Process")

	accountObject.Method = "GET"

	res, _, err := accountObject.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountsHandler"}).Error(err)
		return
	}
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountsHandler"}).Trace("Accounts data retrieved")

	// Establish context for populating allinfo template
	context := types.Context{
		Navigation: "Accounts",
		CreateForm: accountFormData,
		Accounts:   *res,
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountsHandler"}).Info("Accounts retrieval process completed without error")

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

///////////////////////// Account Action Handlers /////////////////////////

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountHandler"}).Info("Starting Account retrieval Process")

	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	// Parse Vars from URL Variables
	vars := mux.Vars(r)

	// Initialize Object Struct for SCIM Request
	accountObject.Id = vars["id"]
	accountObject.Method = "GET"
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountHandler"}).Trace("Request Id: ", accountObject.Id)

	// Retrieve byte based object via BuildUrl function
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountHandler"}).Trace("Calling ScimType2Api to retrieve Account info")
	_, res, err := accountObject.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountHandler"}).Error(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountHandler"}).Trace(fmt.Sprintf("Retrieved Account Name: %s, Account Id: %s", res.DisplayName, res.ID))
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// SafeAddHandler reads in form data from the modal that is triggered
// when the add safe button is pressed. This function calls the SCIM
// function which submits the Delete action to the SCIM Endpoint.
func AccountAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := GetRedirectUrl(r.Referer())

	// TODO

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// AccountDelHandler reads in form data from the modal that is triggered
// when the delete button for a particular safe is pressed. This function
// calls the SCIM function which submits the Delete action to the SCIM Endpoint.
// Note: Safe deletions are based off of safe Name not ID like users or groups
func AccountDelHandler(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := GetRedirectUrl(r.Referer())

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountDelHandler"}).Info("Starting Account Delete Process")

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Retrieve USerID from URL to send to Del Function
	vars := mux.Vars(r)
	accountObject.Id = vars["id"]
	accountObject.Method = "DELETE"

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountDelHandler"}).Trace("Request Id: ", accountObject.Id)

	// Delete Account - No response unless error
	_, _, err := accountObject.ScimType2Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountDelHandler"}).Error("Account Delete Process completed finished Error")
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "AccountDelHandler"}).Info("Account Delete Process completed finished Error")

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// AccountUpdateHandler reads in form data from the modal that is triggered
// when the Update button for a particular safe is pressed. This function
// calls the SCIM function which submits the UPDATE action to the SCIM Endpoint.
func AccountUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from Update Form

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}

///////////////////////// Account Form Handlers /////////////////////////

// SafeAddForm is the form for collecting data to add a new Safe
func AccountAddForm(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO SafeAddForm: Initializing Add Safe Form")

	// Establish context for populating add safe template
	context := types.Context{
		Navigation: "Add Safe",
		CreateForm: safeFormData,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}
