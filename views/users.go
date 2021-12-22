package views

import (
	"encoding/json"
	"fmt"
	"net/http"

	types "github.com/strick-j/go-form-webserver/types"
)

//UserAllReq is the function for requesting user info for collecting data to add a new user
func UserAllReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res := BuildUrl("Users", "GET")

	var bodyObject types.User
	json.Unmarshal(res, &bodyObject)

	userFormData := types.CreateForm{
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

	context := types.Context{
		Navigation: "Users",
		CreateForm: userFormData,
		Users:      bodyObject,
	}

	tpl.ExecuteTemplate(w, "objectallinfo.html", context)
}

//UserAddForm is the form for collecting data to add a new user
func UserAddForm(w http.ResponseWriter, r *http.Request) {

	userFormData := types.CreateForm{
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

	context := types.Context{
		Navigation: "Add User",
		CreateForm: userFormData,
	}

	tpl.ExecuteTemplate(w, "objectaddform.html", context)
}

//UserDelForm is the form for deleting a user
func UserDelForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "userdelform.html", nil)
}

//UserAddReq is used to add users from the /useraddreq URL
func UserAddReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println("Reading data from UserAddReq Form")
	uname := r.FormValue("FormUserName")
	fname := r.FormValue("FormGivenName")
	lname := r.FormValue("FormFamilyName")
	udname := r.FormValue("FormDisplayName")
	//uemail := r.FormValue("FormEmail")
	upassword := r.FormValue("FormPassword")
	scimschema := []string{"urn:ietf:params:scim:schemas:core:2.0:User"}

	addUserData := types.PostUserRequest{
		UserName: uname,
		Name: types.FullName{
			FamilyName: fname,
			GivenName:  lname,
		},
		DisplayName: udname,
		Password:    upassword,
		UserType:    "EPVUser",
		Active:      true,
		Schemas:     scimschema,
	}

	// Required as placeholder
	blankstruct := types.PostObjectRequest{}

	res := ScimAPI("Containers", "POST", blankstruct, addUserData)

	fmt.Println(string(res))

}
