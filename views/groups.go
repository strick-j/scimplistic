package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	types "github.com/strick-j/go-form-webserver/types"
	"github.com/strick-j/go-form-webserver/utils"
)

//GroupAllReq is the function for requesting user info for collecting data to add a new user
func GroupAllReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve byte based object via BuildUrl function
	res := BuildUrl("Groups", "GET")

	// Declare and unmarshal byte based response
	var bodyObject types.Group
	json.Unmarshal(res, &bodyObject)

	groupFormData := types.CreateForm{
		FormAction: "/groupaddreq/",
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

	context := types.Context{
		Navigation: "Groups",
		CreateForm: groupFormData,
		Groups:     bodyObject,
	}

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

//GroupAddForm is the form for deleting a user
func GroupAddForm(w http.ResponseWriter, r *http.Request) {

	log.Printf("Initializing Add Group Form")

	groupFormData := types.CreateForm{
		FormAction: "/groupaddreq/",
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

	context := types.Context{
		Navigation:         "Add Group",
		SettingsConfigured: false,
		CreateForm:         groupFormData,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

//GroupAddReq is used to add users from the /useraddreq URL
func GroupAddReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := utils.GetRedirectUrl(r.Referer())

	log.Println("Reading Data from Group Add Form")
	displayName := r.FormValue("FormGroupDisplayName")
	scimschema := []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}

	addGroupData := types.PostObjectRequest{
		DisplayName: displayName,
		Schemas:     scimschema,
	}

	// Required as placeholder
	blankstruct := types.PostUserRequest{}

	res := ScimAPI("Groups", "POST", addGroupData, blankstruct)

	fmt.Println(string(res))

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

//GroupDel is the function for deleting a group
func GroupDelFunc(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := utils.GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Retrieve Group ID from URL to send to Del Function
	id, err := strconv.Atoi(r.URL.Path[len("/groupdel/"):])
	if err != nil {
		log.Println("GroupDelFunc", err)
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		log.Println("Group ID Obtained: ", id)
	}

	// Create Struct for passing data to SCIM API Delete Function
	delObjectData := types.DelObjectRequest{
		ResourceType: "groups",
		ID:           strconv.Itoa(id),
	}

	// Delete Group and recieve response from Delete Function
	res, err := ScimApiDel(delObjectData)
	if res == 204 {
		log.Println("Group Deleted:", id)
		log.Println("Valid http.StatusCode Recieved:", res)
	} else {
		log.Println(err)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// GroupUpdateFunc is responsible for displaying group
// properties for user update.
func GroupUpdateForm(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := utils.GetRedirectUrl(r.Referer())

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
	var bodyObject types.Group
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
				FieldLabelText:  "Group Display Name",
				FieldInputType:  "Text",
				FieldInputName:  "FormGroupDisplayName",
				FieldInFeedback: "Group Display Name is Required",
				FieldIdNum:      1,
				FieldPlaceHold:  bodyObject.Resources[0].DisplayName,
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
		Groups:     bodyObject,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectupdateform.html", context)
}

// GroupUpdateFunc is the function for updating a group based on information
// returned from GroupUpdateForm
func GroupUpdateFunc(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from Update Form

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}
