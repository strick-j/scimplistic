package views

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/internal/types"
)

// Generate Struct for the forms required by User functions
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

	ob := Object{
		Type:   "users",
		Method: "GET",
	}

	// Retrieve byte based object via ScimApiGetObject function
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UsersHandler"}).Info("Attempting to obtain User Data from SCIM API")
	//res, err := BuildUrl("Users", "GET")
	res, _, err := ob.ScimType1Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UsersHandler"}).Error(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UsersHandler"}).Info("Retrieved All User Data")
	}

	// Establish context for populating allinfo template
	context := types.Context{
		Navigation: "Users",
		CreateForm: userFormData,
		Users:      *res,
	}

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

///////////////////////// User Action Handlers /////////////////////////

// UserHandler Obtains all details about a specific user
func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse Vars from URL Variables
	vars := mux.Vars(r)

	// Add "method" and "id" to userObject

	ob := Object{
		Type:   "users",
		Method: "GET",
		Id:     vars["id"],
	}

	// Retrieve byte based object via BuildUrl function
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserHandler"}).Info("Attempting to retrieve User Info for User: ", ob.Id)
	res, _, err := ob.ScimType1Api()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserHandler"}).Info("Retrieved User Information")
	}

	fmt.Println(res)

	//http.Redirect(w, r, redirectURL, http.StatusFound)
}

// UserAddHandler reads in form data from the modal that is triggered
// when the add user button is pressed. This function calls the SCIM
// function which submits the ADD action to the SCIM Endpoint.
func UserAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserAddHandler"}).Info("Starting User Add Process")
	redirectURL := GetRedirectUrl(r.Referer())

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserAddHandler"}).Trace("Reading User Add Form")
	uname := r.FormValue("FormUserName")
	fname := r.FormValue("FormGivenName")
	lname := r.FormValue("FormFamilyName")
	udname := r.FormValue("FormDisplayName")
	uemail := r.FormValue("FormEmail")
	upassword := r.FormValue("FormPassword")
	scimschema := []string{"urn:ietf:params:scim:schemas:core:2.0:User"}

	addUserData := &types.Type1Resources{
		UserName: uname,
		Name: types.Name{
			FamilyName: lname,
			GivenName:  fname,
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

	// Add "method" and "payload" to userObject
	ob := Object{
		Type:           "users",
		Method:         "POST",
		Type1Resources: *addUserData,
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserAddHandler"}).Trace("Calling ScimType1Api with User Add data")
	_, res, err := ob.ScimType1Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserAddHandler"}).Error(err)
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserAddHandler"}).Info("User added: ", res.DisplayName)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// UserDelHandler reads in form data from the modal that is triggered
// when the delete user button is pressed. This function calls the SCIM
// function which submits the DELETE action to the SCIM Endpoint.
func UserDelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserDelHandler"}).Info("Starting User Delete Process")
	redirectURL := GetRedirectUrlNoId(r.Referer())

	// Retrieve UserID from URL to send to Del Function.
	vars := mux.Vars(r)
	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserDelHandler"}).Info("User ID to Delete: ", vars["id"])

	ob := Object{
		Type:   "users",
		Method: "DELETE",
		Id:     vars["id"],
	}

	// Call ScimUserApi with userObject Details
	_, _, err := ob.ScimType1Api()
	if err != nil {
		log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserDelHandler"}).Error(err)
	}

	log.WithFields(log.Fields{"Category": "SCIM API Request", "Function": "UserDelHandler"}).Info("User Deleted", vars["id"])
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

//////////////////////// User Form Handlers /////////////////////////

// UserAddForm is the form utilized to build the Modal when the
// Add button is pressed from the Users page.
func UserAddForm(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"Category": "Form Handler", "Function": "UserAddForm"}).Trace("Initializing Add User Form")

	// Establish context for populating add user template
	context := types.Context{
		Navigation: "Add User",
		CreateForm: userFormData,
	}

	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

// UserUpdateForm is the form utilized to build the Modal when the
// Update/Edit button is pressed from the Users page.
/*func UserUpdateForm(w http.ResponseWriter, r *http.Request) {
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
}
*/
