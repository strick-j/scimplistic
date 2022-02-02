package views

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/strick-j/scimplistic/types"
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

	// Retrieve byte based object via BuildUrl function
	log.Println("INFO GroupsHandler: Attempting to obtain Group Data from SCIM API.")
	res, err := BuildUrl("Groups", "GET")
	if err != nil {
		log.Println("ERROR GroupsHandler:", err)
		return
	} else {
		log.Println("INFO GroupsHandler: Group Information Recieved")
	}

	// Declare and unmarshal byte based response
	var bodyObject types.Group
	json.Unmarshal(res, &bodyObject)

	// Establish context for populating allinfo template
	context := types.Context{
		Navigation: "Groups",
		CreateForm: groupFormData,
		Groups:     bodyObject,
	}

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

//////////////////////// Group Action Handlers /////////////////////////

// GroupActionHandler will decide what is required based on the action provided.
// add: Proceed to GroupAddHandler
// del: Proceed to GroupDelHandler
// update: Proceed to GroupUpdateHandler
// review: Proceed to GroupReviewHandler
func GroupsActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Extract action from URL using mux.Vars.
	vars := mux.Vars(r)
	action := vars["action"]

	// Switch to appropriate handler based on action type.
	switch action {
	case "add":
		log.Println("INFO GroupsActionHandler: Calling GroupAddHandler")
		GroupAddHandler(w, r)
	case "del":
		log.Println("INFO GroupsActionHandler: Calling GroupDelHandler")
		GroupDelHandler(w, r)
	case "update":
		log.Println("INFO GroupsActionHandler: Calling GroupUpdateHandler")
		GroupUpdateHandler(w, r)
	case "review":
		log.Println("INFO GroupsActionHandler: Calling GroupReviewHandler")
		GroupReviewHandler(w, r)
	}
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

	log.Println("INFO GroupAddHandler: Reading Data from Group Add Form")
	displayName := r.FormValue("FormGroupDisplayName")
	scimschema := []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}

	addGroupData := types.PostObjectRequest{
		DisplayName: displayName,
		Schemas:     scimschema,
	}

	// Required as placeholder.
	blankstruct := types.PostUserRequest{}

	// Utilize the SCIM API funciton to POST a new Group
	// Use the group data retrieved from the Add Group Form.
	res, code, err := ScimAPI("Groups", "POST", addGroupData, blankstruct)
	if code != 201 {
		log.Println("ERROR GroupAddHandler:", err)
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		log.Println("INFO GroupAddHandler: Group Added - Response StatusCode:", code)
	}

	// Declare and unmarshal byte based response.
	var bodyObject types.PostGroupResponse
	err = json.Unmarshal(res, &bodyObject)
	if err != nil {
		log.Println("ERROR GroupAddHandler:", err)
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		log.Println("INFO GroupAddHandler: Group Display Name and ID:", bodyObject.DisplayName, "-", bodyObject.ID)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// GroupDelHandler reads in form data from the modal that is triggered
// when the delete button for a particular group is pressed. This function
// calls the SCIM function which submits the DELETE action to the SCIM Endpoint.
func GroupDelHandler(w http.ResponseWriter, r *http.Request) {
	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Println("INFO GroupDelFunc: Starting Group Delete Process")

	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("INFO GroupDelFunc: Group ID to Delete:", id)

	// Create Struct for passing data to SCIM API Delete Function.
	delObjectData := types.DelObjectRequest{
		ResourceType: "groups",
		ID:           id,
	}

	// Delete Group and recieve response from Delete Function.
	res, err := ScimApiDel(delObjectData)
	if res == 204 {
		log.Println("INFO GroupDelFunc: Group Deleted:", id)
		log.Println("INFO GroupDelFunc: Valid HTTP StatusCode Recieved:", res)
		log.Println("SUCCESS GroupDelFunc: Group Delete Process Complete.")
	} else {
		log.Println(err)
		log.Println("ERROR GroupDelFunc: Invalid Http StatusCode Recieved:", res)
	}

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

// GroupReviewHandler reads in safe id information when the Review button
// for a particular group is pressed. This function then provides the user with
// more information about a particular group.
func GroupReviewHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from review (safe id)

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}

//////////////////////// Group Form Handlers /////////////////////////

// GroupAddForm is the form utilized to build the Modal when the
// Add button is pressed from the Groups page.
func GroupAddForm(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO GroupAddForm: Initializing Add Group Form")

	// Establish context for populating add group template
	context := types.Context{
		Navigation: "Add Group",
		CreateForm: groupFormData,
	}

	// Pass form data to form template to dynamically build form
	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

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
