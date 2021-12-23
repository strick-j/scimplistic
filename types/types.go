package types

import "time"

// ServiceProviderResponse handles the response when the SCIM ServiceProvider is queried
type ServiceProviderResponse struct {
	Schemas []string `json:"schemas"`
	Patch   struct {
		Supported bool `json:"supported"`
	} `json:"patch"`
	Bulk struct {
		Supported      bool `json:"supported"`
		MaxOperations  int  `json:"maxOperations"`
		MaxPayloadSize int  `json:"maxPayloadSize"`
	} `json:"bulk"`
	Filter struct {
		Supported  bool `json:"supported"`
		MaxResults int  `json:"maxResults"`
	} `json:"filter"`
	ChangePassword struct {
		Supported bool `json:"supported"`
	} `json:"changePassword"`
	Sort struct {
		Supported bool `json:"supported"`
	} `json:"sort"`
	Etag struct {
		Supported bool `json:"supported"`
	} `json:"etag"`
	AuthenticationSchemes []struct {
		Type        string `json:"type"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"authenticationSchemes"`
	Meta struct {
		ResourceType string    `json:"resourceType"`
		Created      time.Time `json:"created"`
		LastModified time.Time `json:"lastModified"`
		Location     string    `json:"location"`
	} `json:"meta"`
}

// ResourceTypes handles the response when the SCIM ResourceTypes are queried
type ResourceTypes struct {
	Schemas      []string `json:"schemas"`
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	StartIndex   int      `json:"startIndex"`
	Resources    []struct {
		Name             string `json:"name"`
		Endpoint         string `json:"endpoint"`
		Schema           string `json:"schema"`
		SchemaExtensions []struct {
			Schema   string `json:"schema"`
			Required bool   `json:"required"`
		} `json:"schemaExtensions,omitempty"`
		Schemas []string `json:"schemas"`
		ID      string   `json:"id"`
		Meta    struct {
			ResourceType string    `json:"resourceType"`
			Created      time.Time `json:"created"`
			LastModified time.Time `json:"lastModified"`
			Location     string    `json:"location"`
		} `json:"meta"`
	} `json:"Resources"`
}

// Schemas handles the response when the SCIM Schemas are queried
type Schemas struct {
	Schemas      []string `json:"schemas"`
	TotalResults int      `json:"totalResults"`
	ItemsPerPage int      `json:"itemsPerPage"`
	StartIndex   int      `json:"startIndex"`
	Resources    []struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Attributes  []struct {
			Name          string `json:"name"`
			Type          string `json:"type"`
			MultiValued   bool   `json:"multiValued"`
			Required      bool   `json:"required"`
			CaseExact     bool   `json:"caseExact,omitempty"`
			SubAttributes []struct {
				Name        string `json:"name"`
				Type        string `json:"type"`
				MultiValued bool   `json:"multiValued"`
				Required    bool   `json:"required"`
				CaseExact   bool   `json:"caseExact"`
			} `json:"subAttributes,omitempty"`
			Mutability string `json:"mutability,omitempty"`
			Returned   string `json:"returned,omitempty"`
		} `json:"attributes,omitempty"`
		ID   string `json:"id"`
		Meta struct {
			ResourceType string    `json:"resourceType"`
			Created      time.Time `json:"created"`
			LastModified time.Time `json:"lastModified"`
			Location     string    `json:"location"`
		} `json:"meta"`
	} `json:"Resources"`
}

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

// CreateForm holds the initial fields required to setup a form
// Action: "DestinationURL"
// Method: "Get, Post, etc..."
// Legend: "Form Title"
// Role:   "adduser,addgroup,addsafe,etc..."
type CreateForm struct {
	FormAction string       `json:"formAction"`
	FormMethod string       `json:"formMethod"`
	FormLegend string       `json:"formLegend"`
	FormRole   string       `json:"formRole,omitempty"`
	FormFields []FormFields `json:"formFields"`
}

// FormFields builds out the individual fields within a form
type FormFields struct {
	FieldType       string `json:"fieldType"`
	FieldLabel      string `json:"fieldLabel"`
	FieldLabelText  string `json:"fieldLabelText,omitempty"`
	FieldInputType  string `json:"fieldInputType"`
	FieldRequired   bool   `json:"fieldReuired,omitempty"`
	FieldInputName  string `json:"fieldInputName,omitempty"`
	FieldDescBy     string `json:"fieldDescBy,omitempty"`
	FieldHelp       string `json:"fieldHelp,omitempty"`
	FieldPlaceHold  string `json:"fieldPlaceHold,omitempty"`
	FieldIdNum      int    `json:"fieldId,omitempty"`
	FieldInFeedback string `json:"fieldInFeedback,omitempty"`
	FieldVaFeedback string `json:"fieldVaFeedback,omitempty"`
	FieldDisabled   bool   `json:"fieldDisabled,omitempty"`
}

// Context is utilized for sending data to templates / forms
type Context struct {
	Navigation         string     `json:"navigation,omitempty"`
	Message            string     `json:"message,omitempty"`
	SettingsConfigured bool       `json:"settingsConfigured,omitempty"`
	Token              string     `json:"authToken,omitempty"`
	CreateForm         CreateForm `json:"createForm,omitempty"`
	Safes              Safes      `json:"safes,omitempty"`
	Users              User       `json:"users,omitempty"`
	Groups             Group      `json:"groups,omitempty"`
	Members            Members    `json:"members,omitempty"`
}

// ConfigSettings is the struct to hold settings information
// used in conjunction with MarshallIndent to write configuration
type ConfigSettings struct {
	ScimURL   string `json:"scimURL"`
	AuthToken string `json:"AuthToken"`
	PrevConf  bool   `json:"prevConf"`
}

// PostUser is the struct created for adding users
type PostUserRequest struct {
	UserName    string   `json:"userName"`
	Name        Name     `json:"fullName,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Password    string   `json:"password"`
	UserType    string   `json:"userType,omitempty"`
	Active      bool     `json:"active,omitempty"`
	Emails      []Emails `json:"emails,omitempty"`
	Schemas     []string `json:"schemas"`
}

type Emails struct {
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
	Value   string `json:"value,omitempty"`
}

type Name struct {
	FamilyName string `json:"familyName,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
}

// PostObjectRequest contains the required fields for a Posts for
// adding Groups, and Safes
type PostObjectRequest struct {
	Name        string    `json:"name,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	Description string    `json:"description,omitempty"`
	Members     []Members `json:"members,omitempty"`
	Schemas     []string  `json:"schemas"`
}

type Members struct {
	Value   string `json:"value,omitempty"`
	Ref     string `json:"$ref,omitempty"`
	Display string `json:"display,omitempty"`
}

// PostGroupResponse contains the fields returned when a group is added
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

type DelObjectRequest struct {
	ResourceType string    `json:"resourceType"`
	ID           string    `json:"id"`
	DisplayName  string    `json:"displayName"`
	Members      []Members `json:"members,omitempty"`
	Schemas      []string  `json:"schemas"`
}

// PostUserResponse contains the fields returned when a user is added
type PostUserResponse struct {
	UserName string `json:"userName"`
	Name     struct {
		Formatted string `json:"formatted"`
		GivenName string `json:"givenName"`
	} `json:"name"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	Emails      []struct {
		Type    string `json:"type"`
		Primary bool   `json:"primary"`
		Value   string `json:"value"`
	} `json:"emails"`
	Schemas []string `json:"schemas"`
	ID      string   `json:"id"`
	Meta    struct {
		ResourceType string    `json:"resourceType"`
		Created      time.Time `json:"created"`
		LastModified time.Time `json:"lastModified"`
		Location     string    `json:"location"`
	} `json:"meta"`
}

type PostSafeResponse struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Owner       struct {
		Value   string `json:"value"`
		Display string `json:"display"`
	} `json:"owner"`
	Schemas []string `json:"schemas"`
	ID      string   `json:"id"`
	Meta    struct {
		ResourceType string    `json:"resourceType"`
		Created      time.Time `json:"created"`
		LastModified time.Time `json:"lastModified"`
		Location     string    `json:"location"`
	} `json:"meta"`
	UrnIetfParamsScimSchemasCyberark11Safe struct {
		NumberOfDaysRetention int `json:"NumberOfDaysRetention"`
	} `json:"urn:ietf:params:scim:schemas:cyberark:1.1:Safe"`
}
