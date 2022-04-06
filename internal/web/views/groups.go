package views

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/internal/types"
)

// Generate Struct for the forms required by GroupFunctions
// Form Data is used by several functions (add, get, update, etc...)
var groupFormData = types.CreateForm{
	FormAction: "/groups/add",
	FormMethod: "POST",
	FormLegend: "Add Group Form",
	FormFields: []types.FormFields{
		{
			FieldLabel:      "DisplayName",
			FieldLabelText:  "Group Display Name",
			FieldInputType:  "Text",
			FieldRequired:   true,
			FieldInputName:  "FormGroupDisplayName",
			FieldInFeedback: "Group Display Name is Required",
			FieldIdNum:      1,
		},
	},
}

//////////////////////// Group Default Handler /////////////////////////

// GroupHandler is the function for displaying the basic Groups page
func GroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ob := Object{
		Type:   "groups",
		Method: "GET",
	}

	// Retrieve byte based object via ScimGroupApifunction
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupsHandler"}).Info("Attempting to obtain Groups from SCIM API.")
	res, _, err := ob.ScimType1Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupsHandler"}).Error(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupsHandler"}).Info("Retrieved All User Data")
	}

	// Establish context for populating allinfo template
	context := types.Context{
		Navigation: "Groups",
		CreateForm: groupFormData,
		Groups:     *res,
	}

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

//////////////////////// Group Action Handlers /////////////////////////

func GroupHandler(w http.ResponseWriter, r *http.Request) {
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
		Type:   "groups",
		Method: "GET",
		Id:     vars["id"],
	}

	// Retrieve byte based object via ScimGroupApifunction
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupsHandler"}).Info("Attempting to obtain Groups from SCIM API.")
	res, _, err := ob.ScimType1Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupsHandler"}).Error(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupsHandler"}).Info("Retrieved All User Data")
	}

	fmt.Println(res)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// GroupAddHandler reads in form data from the modal that is triggered
// when the add group button is pressed. This function calls the SCIM
// function which submits the ADD action to the SCIM Endpoint.
func GroupAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupAddHandler"}).Info("Reading Data from Group Add Form")
	displayName := r.FormValue("FormGroupDisplayName")
	scimschema := []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}

	ob := Object{
		Type:   "groups",
		Method: "POST",
	}

	ob.Type1Resources = types.Type1Resources{
		DisplayName: displayName,
		Schemas:     scimschema,
	}

	// Utilize the SCIM API funciton to POST a new Group
	// Use the group data retrieved from the Add Group Form.
	_, res, err := ob.ScimType1Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupDelHandler"}).Error(err)
	}
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupDelHandler"}).Trace(fmt.Sprintf("Group Display Name %s, Group Id %s ", res.DisplayName, res.ID))

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// GroupDelHandler reads in form data from the modal that is triggered
// when the delete button for a particular group is pressed. This function
// calls the SCIM function which submits the DELETE action to the SCIM Endpoint.
func GroupDelHandler(w http.ResponseWriter, r *http.Request) {
	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrlNoId(r.Referer())

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupDelHandler"}).Trace("Starting Group Delete Process")

	vars := mux.Vars(r)
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupDelHandler"}).Trace("Attempting to delete Group: ", vars["id"])

	ob := Object{
		Type:   "groups",
		Method: "DELETE",
		Id:     vars["id"],
	}

	// Delete Group. No response should be returned unless there is an error.
	_, _, err := ob.ScimType1Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupDelHandler"}).Error(err)
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "GroupDelHandler"}).Trace("Group Delete Process completed without error")
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// GroupUpdateHandler reads in form data from the modal that is triggered
// when the Update button for a particular group is pressed. This function
// calls the SCIM function which submits the UPDATE action to the SCIM Endpoint.
func GroupUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from Update Form

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}

//////////////////////// Group Form Handlers /////////////////////////

// GroupAddForm is the form utilized to build the Modal when the
// Add button is pressed from the Groups page.
func GroupAddForm(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"Category": "Form Handler", "Function": "GroupAddForm"}).Trace("Initializing Add Group Form")
	// Establish context for populating add group template
	context := types.Context{
		Navigation: "Add Group",
		CreateForm: groupFormData,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

/*
// GroupUpdateForm is the form utilized to build the Modal when the
// Update/Edit button is pressed from the Groups page.
func GroupUpdateForm(w http.ResponseWriter, r *http.Request) {
	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, redirectURL, http.StatusBadRequest)
		return
	}

	// Retrieve Group ID from URL to send to request function
	// required to retrieve latest Group Details
	displayName := r.URL.Path[len("/groupupdate/"):]
	log.Println("GroupUpdateFunc: Group DisplayName Obtained -", displayName)

	// Create Struct for passing data to SCIM API Request Function
	log.Println("GroupUpdateFunc: Querying SCIM Server for Group Details for -", displayName)
	reqObjectData := types.DelObjectRequest{
		ResourceType: "groups",
		DisplayName:  displayName,
	}

	// Request Group Information Function
	res, err := ScimApiReq(reqObjectData)
	if res != nil {
		log.Println("GroupUpdateFunc: Group Details Obtained for:", displayName)
	} else {
		log.Println(err)
	}

	// Declare and unmarshal byte based response
	var bodyObject types.Resources
	json.Unmarshal(res, &bodyObject)

	// Generate Struct for the Update Form that will be created
	groupFormData := types.CreateForm{
		FormAction: "/groupupdatereq/",
		FormMethod: "POST",
		FormLegend: "Update Group Form",
		FormRole:   "updategroup",
		FormFields: []types.FormFields{
			{
				FieldLabel:      "DisplayName",
				FieldInputType:  "Text",
				FieldInputName:  "FormGroupDisplayName",
				FieldInFeedback: "Group Display Name is Required",
				FieldIdNum:      1,
				FieldPlaceHold:  bodyObject.DisplayName,
				FieldDisabled:   true,
			},
		},
	}

	// Pass the context for the Update Form Page. Includes:
	// Navigation Information
	// Create Form Struct for Creating Form Layout
	// Response from Group Query (bodyObject) for populating form
	context := types.Context{
		Navigation: "Update Group",
		CreateForm: groupFormData,
		Groups:     bodyObject.Group,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectupdateform.html", context)
}*/
