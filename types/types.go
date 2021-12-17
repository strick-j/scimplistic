package types

import "time"

type User struct {
	Schemas      []string `json:"schemas"`
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	StartIndex   int      `json:"startIndex"`
	Resources    []struct {
		UserName string `json:"userName"`
		Name     struct {
			FamilyName string `json:"familyName"`
			GivenName  string `json:"givenName"`
		} `json:"name"`
		DisplayName string `json:"displayName"`
		UserType    string `json:"userType"`
		Active      bool   `json:"active"`
		Groups      []struct {
			Type    string `json:"type"`
			Display string `json:"display"`
			Value   string `json:"value"`
			Ref     string `json:"$ref"`
		} `json:"groups"`
		Entitlements []string `json:"entitlements"`
		Schemas      []string `json:"schemas"`
		ID           string   `json:"id"`
		Meta         struct {
			ResourceType string    `json:"resourceType"`
			Created      time.Time `json:"created"`
			LastModified time.Time `json:"lastModified"`
			Location     string    `json:"location"`
		} `json:"meta"`
		UrnIetfParamsScimSchemasPam11LinkedObject struct {
			Source           string `json:"source"`
			NativeIdentifier string `json:"nativeIdentifier"`
		} `json:"urn:ietf:params:scim:schemas:pam:1.1:LinkedObject"`
		UrnIetfParamsScimSchemasExtensionEnterprise21User struct {
			Organization string `json:"organization"`
		} `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.1:User"`
	}
}

type Group struct {
	Schemas      []string `json:"schemas"`
	TotalResults int      `json:"totalResults"`
	Resources    []struct {
		DisplayName string `json:"displayName"`
		Members     []struct {
			Value   string `json:"value"`
			Ref     string `json:"$ref"`
			Display string `json:"display"`
		} `json:"members,omitempty"`
		Schemas []string `json:"schemas"`
		ID      string   `json:"id"`
		Meta    struct {
			ResourceType string    `json:"resourceType"`
			Created      time.Time `json:"created"`
			LastModified time.Time `json:"lastModified"`
			Location     string    `json:"location"`
		} `json:"meta"`
		ExternalID string `json:"externalId,omitempty"`
	} `json:"Resources"`
}

type Groups []Group

type AddUser struct {
	UserName string `json:"userName"`
	Name     struct {
		FamilyName string `json:"familyName,omitempty"`
		GivenName  string `json:"givenName,omitempty"`
	} `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Password    string `json:"password,omitempty"`
	UserType    string `json:"userType,omitempty"`
	Active      bool   `json:"active,omitempty"`
	//Emails      []struct {
	//	Type    string `json:"type,omitempty"`
	//	Primary bool   `json:"primary,omitempty"`
	//	Value   string `json:"value,omitempty"`
	//} `json:"emailsomitempty"`
	Schemas []string `json:"schemas"`
}

type Safes struct {
	Schemas      []string `json:"schemas"`
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	StartIndex   int      `json:"startIndex"`
	Resources    []struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Owner       struct {
			Value   string `json:"value"`
			Ref     string `json:"$ref"`
			Display string `json:"display"`
		} `json:"owner"`
		PrivilegedData []struct {
			Value   string `json:"value"`
			Ref     string `json:"$ref"`
			Display string `json:"display"`
		} `json:"privilegedData"`
		Schemas []string `json:"schemas"`
		ID      string   `json:"id"`
		Meta    struct {
			ResourceType string    `json:"resourceType"`
			Created      time.Time `json:"created"`
			LastModified time.Time `json:"lastModified"`
			Location     string    `json:"location"`
		} `json:"meta"`
		UrnIetfParamsScimSchemasCyberark11Safe struct {
			NumberOfDaysRetention int    `json:"NumberOfDaysRetention"`
			ManagingCPM           string `json:"ManagingCPM"`
		} `json:"urn:ietf:params:scim:schemas:cyberark:1.1:Safe"`
		UniqueSafeId string `json:"uniqueSafeId"`
	} `json:"Resources"`
}

type CreateForm struct {
	FormAction string       `json:"formAction"`
	FormMethod string       `json:"formMethod"`
	FormLegend string       `json:"formLegend"`
	FormFields []FormFields `json:"formFields"`
}

type FormFields struct {
	FieldLabel     string `json:"fieldLabel"`
	FieldLabelText string `json:"fieldLabelText,omitempty"`
	FieldInputType string `json:"fieldInputType,omitempty"`
	FieldRequired  bool   `json:"fieldReuired,omitempty"`
	FieldInputName string `json:"fieldInputName,omitempty"`
	FieldIdNum     int    `json:"fieldId,omitempty"`
}

type PostGroupRequest struct {
	DisplayName string `json:"displayName"`
	Members     []struct {
		Value   string `json:"value"`
		Ref     string `json:"$ref"`
		Display string `json:"display"`
	} `json:"members"`
	Schemas []string `json:"schemas"`
}

type AddSafe struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description"`
	Schemas     []string `json:"schemas"`
}

type PostObjectRequest struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
	Members     []struct {
		Value   string `json:"value,omitempty"`
		Ref     string `json:"$ref,omitempty"`
		Display string `json:"display,omitempty"`
	} `json:"members,omitempty"`
	Schemas []string `json:"schemas"`
}

type PostGroupResponse struct {
	DisplayName string   `json:"displayName"`
	Schemas     []string `json:"schemas"`
	ID          string   `json:"id"`
	Meta        struct {
		ResourceType string    `json:"resourceType"`
		Created      time.Time `json:"created"`
		LastModified time.Time `json:"lastModified"`
		Location     string    `json:"location"`
	} `json:"meta"`
}
