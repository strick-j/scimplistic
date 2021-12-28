package views

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/strick-j/scimplistic/types"
	"github.com/strick-j/scimplistic/utils"
)

// Generate Struct for the forms required by GroupFunctions
// Form Data is used by several functions (add, get, update, etc...)
var userFormData = types.CreateForm{
	FormAction: "/useraddreq/",
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

//UserAllReq is the function for requesting user info for collecting data to add a new user
func UserAllReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve byte based object via BuildUrl function
	log.Println("INFO UserAllReq: Attempting to obtain User Data from SCIM API.")
	res, err := BuildUrl("Users", "GET")
	if err != nil {
		log.Println("ERROR UserAllReq:", err)
		return
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

//UserAddForm is the form for collecting data to add a new user
func UserAddForm(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO GroupAddForm: Initializing Add Group Form")

	// Establish context for populating add group template
	context := types.Context{
		Navigation: "Add User",
		CreateForm: userFormData,
	}

	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

//UserDelFunc is the form for deleting a user
func UserDelFunc(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := utils.GetRedirectUrl(r.Referer())

	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Println("INFO UserDelFunc: Starting User Delete Process")

	// Retrieve UserID from URL to send to Del Function
	id, err := strconv.Atoi(r.URL.Path[len("/userdel/"):])
	if err != nil {
		log.Println("ERROR UserDelFunc:", err)
		return
	} else {
		log.Println("INFO UserDelFunc: User ID to Delete:", id)
	}

	// Create Struct for passing data to SCIM API Delete Function
	delObjectData := types.DelObjectRequest{
		ResourceType: "users",
		ID:           strconv.Itoa(id),
	}

	// Delete User and recieve response from Delete Function
	res, err := ScimApiDel(delObjectData)
	if res == 204 {
		log.Println("INFO UserDelFunc: UserDeleted:", id)
		log.Println("INFO UserDelFunc: Valid HTTP StatusCode Recieved:", res)
		log.Println("SUCCESS UserDelFunc: UserDelete Process Complete.")
	} else {
		log.Println(err)
		log.Println("ERROR UserDelFunc: Invalid Http StatusCode Recieved:", res)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

//UserAddReq is used to add users from the /useraddreq URL
func UserAddReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	redirectURL := utils.GetRedirectUrl(r.Referer())

	log.Println("INFO UserAddReq: Reading data from UserAddReq Form")
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

	// Required as placeholder
	blankstruct := types.PostObjectRequest{}

	log.Println("INFO UserAddReq: Sending POST to SCIM server for user add.")
	res, code, err := ScimAPI("Users", "POST", blankstruct, addUserData)
	if code != 201 {
		log.Println("ERROR UserAddReq: Error Adding User - Response StatusCode:", code)
		log.Println("ERROR UserAddReq: Error Adding User", err)
		return
	} else {
		log.Println("INFO UserAddReq: Recieved SCIM Response - Valid HTTP StatusCode:", code)
	}

	// Declare and unmarshal byte based response
	log.Println("INFO UserAddReq: Parsing User Add SCIM Reponse.")
	var bodyObject types.PostUserResponse
	err = json.Unmarshal(res, &bodyObject)
	if err != nil {
		log.Println("ERROR UserAddReq:", err)
		return
	} else {
		log.Println("SUCCESS UserAddReq: User Display Name and ID:", bodyObject.DisplayName, "-", bodyObject.ID)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
