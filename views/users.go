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
var userFormData = types.CreateForm{
	FormAction: "/users/add",
	FormMethod: "POST",
	FormLegend: "Add User Form",
	FormRole:   "adduser",
	FormFields: []types.FormFields{
		{
			FieldLabel:      "Username",
			FieldLabelText:  "Username",
			FieldInputType:  "Text",
			FieldRequired:   true,
			FieldInputName:  "FormUserName",
			FieldInFeedback: "Username is required.",
			FieldIdNum:      1,
		},
		{
			FieldLabel:     "givenName",
			FieldLabelText: "First Name",
			FieldInputType: "Text",
			FieldRequired:  false,
			FieldInputName: "FormGivenName",
			FieldDescBy:    "givenHelp",
			FieldHelp:      "Optional",
			FieldIdNum:     2,
		},
		{
			FieldLabel:     "familyName",
			FieldLabelText: "Last Name",
			FieldInputType: "Text",
			FieldRequired:  false,
			FieldInputName: "FormFamilyName",
			FieldDescBy:    "familyHelp",
			FieldHelp:      "Optional",
			FieldIdNum:     3,
		},
		{
			FieldLabel:     "inputEmail",
			FieldLabelText: "Email Address",
			FieldInputType: "email",
			FieldRequired:  false,
			FieldInputName: "FormEmail",
			FieldPlaceHold: "username@acme.com",
			FieldIdNum:     4,
		},
		{
			FieldLabel:      "inputPassword",
			FieldLabelText:  "User Password",
			FieldInputType:  "password",
			FieldRequired:   true,
			FieldInputName:  "FormPassword",
			FieldDescBy:     "PasswordHelp",
			FieldHelp:       "User password must be 8-20 characters long, contain letters and numbers, and a special character. It must not contain spaces, special characters, or emoji.",
			FieldInFeedback: "Password is required.",

			FieldIdNum: 5,
		},
	},
}

///////////////////////// User Default Handler /////////////////////////

//UsersHandler is the function for requesting user info for collecting data to add a new user
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve byte based object via BuildUrl function
	log.Println("INFO UserAllReq: Attempting to obtain User Data from SCIM API.")
	res, err := BuildUrl("Users", "GET")
	if err != nil {
		log.Println("ERROR UserAllReq:", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.Println("INFO User AllReq: User Information Recieved")
	}

	// Declare and unmarshal byte based response
	var bodyObject types.User
	json.Unmarshal(res, &bodyObject)

	// Establish context for populating allinfo template
	context := types.Context{
		Navigation: "Users",
		CreateForm: userFormData,
		Users:      bodyObject,
	}

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

///////////////////////// User Action Handlers /////////////////////////

// UsersActionHandler will decide what is required based on the action provided.
// add: Proceed to UserAddHandler
// del: Proceed to UserDelHandler
// update: Proceed to UserUpdateHandler
// review: Proceed to UserReviewHandler
func UsersActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Extract action from URL using mux.Vars.
	vars := mux.Vars(r)
	action := vars["action"]
	log.Println("INFO UsersActionHandler: Action = ", action)

	// Switch to appropriate handler based on action type.
	switch action {
	case "add":
		log.Println("INFO UsersActionHandler: Calling to UserAddHandler")
		UserAddHandler(w, r)
	case "del":
		log.Println("INFO UsersActionHandler: Calling to UserDelHandler")
		UserDelHandler(w, r)
	case "update":
		log.Println("INFO UsersActionHandler: Calling to UserUpdateHandler")
		UserUpdateHandler(w, r)
	case "review":
		log.Println("INFO UsersActionHandler: Calling to UserReviewHandler")
		UserReviewHandler(w, r)
	}
}

// UserAddHandler reads in form data from the modal that is triggered
// when the add user button is pressed. This function calls the SCIM
// function which submits the ADD action to the SCIM Endpoint.
func UserAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	log.Println("INFO UserAddHandler: Reading data from UserAddReq Form")
	uname := r.FormValue("FormUserName")
	fname := r.FormValue("FormGivenName")
	lname := r.FormValue("FormFamilyName")
	udname := r.FormValue("FormDisplayName")
	uemail := r.FormValue("FormEmail")
	upassword := r.FormValue("FormPassword")
	scimschema := []string{"urn:ietf:params:scim:schemas:core:2.0:User"}

	addUserData := types.PostUserRequest{
		UserName: uname,
		Name: types.Name{
			FamilyName: fname,
			GivenName:  lname,
		},
		DisplayName: udname,
		Emails: []types.Emails{
			{
				Type:    "Primary",
				Value:   uemail,
				Primary: true,
			},
		},
		Password: upassword,
		UserType: "EPVUser",
		Active:   true,
		Schemas:  scimschema,
	}

	// Required as placeholder.
	blankstruct := types.PostObjectRequest{}

	log.Println("INFO UserAddHandler: Sending POST to SCIM server for user add.")
	res, code, err := ScimAPI("Users", "POST", blankstruct, addUserData)
	if code != 201 {
		log.Println("ERROR UserAddHandler: Error Adding User - Response StatusCode:", code)
		log.Println("ERROR UserAddHandler: Error Adding User", err)
		return
	} else {
		log.Println("INFO UserAddHandler: Recieved SCIM Response - Valid HTTP StatusCode:", code)
	}

	// Declare and unmarshal byte based response.
	log.Println("INFO UserAddHandler: Parsing User Add SCIM Reponse.")
	var bodyObject types.PostUserResponse
	err = json.Unmarshal(res, &bodyObject)
	if err != nil {
		log.Println("ERROR UserAddHandler:", err)
		return
	} else {
		log.Println("SUCCESS UserAddHandler: User Display Name and ID:", bodyObject.DisplayName, "-", bodyObject.ID)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// UserDelHandler reads in form data from the modal that is triggered
// when the delete user button is pressed. This function calls the SCIM
// function which submits the DELETE action to the SCIM Endpoint.
func UserDelHandler(w http.ResponseWriter, r *http.Request) {
	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Println("INFO UserDelHandler: Starting User Delete Process")

	// Retrieve UserID from URL to send to Del Function.
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("INFO UserDelHandler: User ID to Delete:", id)

	// Create Struct for passing data to SCIM API Delete Function.
	delObjectData := types.DelObjectRequest{
		ResourceType: "users",
		ID:           id,
	}

	// Delete User and recieve response from Delete Function.
	res, err := ScimApiDel(delObjectData)
	if res == 204 {
		log.Println("INFO UserDelHandler: UserDeleted:", id)
		log.Println("INFO UserDelHandler: Valid HTTP StatusCode Recieved:", res)
		log.Println("SUCCESS UserDelHandler: UserDelete Process Complete.")
	} else {
		log.Println(err)
		log.Println("ERROR UserDelHandler: Invalid Http StatusCode Recieved:", res)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// UserUpdateHandler reads in form data from the modal that is triggered
// when the Update button for a particular user is pressed. This function
// calls the SCIM function which submits the UPDATE action to the SCIM Endpoint.
func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from Update Form

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}

// UserReviewHandler reads in safe id information when the Review button
// for a particular user is pressed. This function then provides the user with
// more information about a particular user.
func UserReviewHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read information from review (safe id)

	// TODO: Parse info and form PUT/POST for Group Update

	// TODO: Execute PUT/POST

	// TODO: Return to Groups page after success
	tpl.ExecuteTemplate(w, "/", nil)
}

//////////////////////// User Form Handlers /////////////////////////

// UserAddForm is the form utilized to build the Modal when the
// Add button is pressed from the Users page.
func UserAddForm(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO UserAddForm: Initializing Add User Form")

	// Establish context for populating add user template
	context := types.Context{
		Navigation: "Add User",
		CreateForm: userFormData,
	}

	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

// UserUpdateForm is the form utilized to build the Modal when the
// Update/Edit button is pressed from the Users page.
func UserUpdateForm(w http.ResponseWriter, r *http.Request) {
	// For best UX we want the user to be returned to the page making
	// the delete transaction, we use the r.Referer() function to get the link.
	redirectURL := GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, redirectURL, http.StatusBadRequest)
		return
	}

	// Retrieve User ID from URL to send to request function
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
