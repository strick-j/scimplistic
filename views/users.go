package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	types "github.com/strick-j/go-form-webserver/types"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

//UserAllReq is the function for requesting user info for collecting data to add a new user
func UserAllReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res := BuildUrl("Users", "GET")

	var bodyObject types.User
	json.Unmarshal(res, &bodyObject)

	tpl.ExecuteTemplate(w, "userallinfo.html", bodyObject)
}

//UserAddForm is the form for collecting data to add a new user
func UserAddForm(w http.ResponseWriter, r *http.Request) {

	userFormData := types.CreateForm{
		FormAction: "/useraddreq",
		FormMethod: "GET",
		FormLegend: "Add User Form",
		FormRole:   "adduser",
		FormFields: []types.FormFields{
			{
				FieldLabel:     "Username",
				FieldLabelText: "Username",
				FieldInputType: "Text",
				FieldRequired:  true,
				FieldInputName: "FormUserName",
				FieldIdNum:     1,
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
				FieldLabel:     "inputPassword",
				FieldLabelText: "User Password",
				FieldInputType: "password",
				FieldRequired:  true,
				FieldInputName: "FormPassword",
				FieldDescBy:    "PasswordHelp",
				FieldHelp:      "User password must be 8-20 characters long, contain letters and numbers, and a special character. It must not contain spaces, special characters, or emoji.",
				FieldIdNum:     5,
			},
		},
	}

	tpl.ExecuteTemplate(w, "objectaddform.html", userFormData)
}

//UserDelForm is the form for deleting a user
func UserDelForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "userdelform.html", nil)
}

//UserAddReq is used to add users from the /useraddreq URL
func UserAddReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println("made it to form read in UserAddReq")
	uname := r.FormValue("FormUserName")
	//fname := r.FormValue("FormGivenName")
	//lname := r.FormValue("FormFamilyName")
	//uemail := r.FormValue("FormEmail")
	upassword := r.FormValue("FormPassword")
	scimschema := []string{"urn:ietf:params:scim:schemas:core:2.0:User"}

	addUserData := types.PostUserRequest{
		UserName: uname,
		Password: upassword,
		Schemas:  scimschema,
	}

	url := "https://identity.strlab.us/scim/v2/Users"
	method := "POST"
	payload, err := json.Marshal(addUserData)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Make request with marshalled JSON as the POST body
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
		return
	}
	// removed header auth token

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode == 201 {

		http.Redirect(w, r, "/success", http.StatusCreated)
	} //else if http.StatusText(res.StatusCode) == "Unauthorized" {
	//	fmt.Println("User does not have the appropriate permissions to utilize the REST API")
	//}
	fmt.Println(res)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))

}
