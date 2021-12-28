package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/strick-j/scimplistic/config"
	"github.com/strick-j/scimplistic/types"
)

func BuildUrl(target string, apimethod string) ([]byte, error) {
	values, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}

	finu := &url.URL{
		Scheme:      "https",
		Opaque:      "",
		Host:        values.ScimURL,
		Path:        "/scim/" + target,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	// Pull API Method from function call (Get, Patch, Del, etc..)
	method := apimethod

	finurl := finu.String()

	client := &http.Client{}
	req, err := http.NewRequest(method, finurl, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "Bearer "+values.AuthToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return body, err
}

func ScimAPI(target string, apimethod string, data types.PostObjectRequest, userdata types.PostUserRequest) ([]byte, int, error) {
	values, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}

	// Generate json from passed struct, check if user or other (Safe/Group)
	var payload []byte
	if target == "Users" {
		payload, err = json.Marshal(userdata)
		if err != nil {
			log.Println("ERROR ScimAPI: Error reading user data: ", err)
			return nil, 0, err
		}
	} else {
		payload, err = json.Marshal(data)
		if err != nil {
			log.Println("ERROR ScimAPI: Error reading object data: ", err)
			return nil, 0, err
		}
	}
	log.Println("INFO ScimAPI: Object data converted to JSON for Post.")

	// Generate target URL from passed target
	finu := &url.URL{
		Scheme:      "https",
		Opaque:      "",
		Host:        values.ScimURL,
		Path:        "/scim/" + target,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	// Pull API Method from function call (Get, Patch, Del, etc..)
	method := apimethod
	finurl := finu.String()
	client := &http.Client{}
	// Build Request
	req, err := http.NewRequest(method, finurl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
	}

	// Add authorization for request
	req.Header.Add("Authorization", "Bearer "+values.AuthToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.Body)

	return body, res.StatusCode, err
}

// ScimApiDel is used by Users, Groups, and Safes to perform
// an object deletion. Note: Users and Groups return no response
// Safes returns a "1" upon successful deletion.
func ScimApiDel(data types.DelObjectRequest) (int, error) {
	// Read in Target and Token from Configuration
	values, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}

	// Generate target URL from passed target
	finu := &url.URL{
		Scheme:      "https",
		Opaque:      "",
		Host:        values.ScimURL,
		Path:        "/scim/" + data.ResourceType + "/" + data.ID,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	// Pull API Method from function call (Get, Patch, Del, etc..)
	method := "DELETE"
	finurl := finu.String()
	client := &http.Client{}
	// Build Request
	req, err := http.NewRequest(method, finurl, nil)
	if err != nil {
		log.Println("ERROR ScimApiDel:", err)
	}

	// Add authorization for request
	req.Header.Add("Authorization", "Bearer "+values.AuthToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println("ERROR ScimApiDel:", err)
	}
	defer res.Body.Close()

	return res.StatusCode, err
}

func ScimApiReq(data types.DelObjectRequest) ([]byte, error) {
	values, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}

	rawQ := ""
	if data.DisplayName != "" {
		rawQ = url.PathEscape("filter=DisplayName eq " + data.DisplayName)
		log.Println("ScimApiReq: Display Name Provided. Query -" + rawQ)
	} else {
		log.Println("ScimApiReq: No Display Name provided. No filter applied")
	}

	finu := &url.URL{
		Scheme:      "https",
		Opaque:      "",
		Host:        values.ScimURL,
		Path:        "/scim/" + data.ResourceType,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    rawQ,
		Fragment:    "",
		RawFragment: "",
	}

	// Pull API Method from function call (Get, Patch, Del, etc..)
	method := "GET"

	finurl := finu.String()
	client := &http.Client{}
	req, err := http.NewRequest(method, finurl, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "Bearer "+values.AuthToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return body, err
}
