package types

import (
	"time"
)

// ServiceProviderResponse handles the response when the SCIM ServiceProvider is queried
type ServiceProvider struct {
	Schemas               []string                `json:"schemas"`
	Patch                 Patch                   `json:"patch"`
	Bulk                  Bulk                    `json:"bulk"`
	Filter                Filter                  `json:"filter"`
	ChangePassword        ChangePassword          `json:"changePassword"`
	Sort                  Sort                    `json:"sort"`
	Etag                  Etag                    `json:"etag"`
	AuthenticationSchemes []AuthenticationSchemes `json:"authenticationSchemes"`
	Meta                  Meta                    `json:"meta"`
}

// ScimType2 should be used for the following:
//   - Users
//   - Groups
type ScimType1 struct {
	Schemas      []string         `json:"schemas"`
	TotalResults int              `json:"totalResults"`
	ItemsPerPage int              `json:"itemsPerPage"`
	StartIndex   int              `json:"startIndex"`
	Resources    []Type1Resources `json:"Resources"`
}

type Type1Resources struct {
	Active                                            bool                                              `json:"active,omitempty"`
	Container                                         Container                                         `json:"container"`
	Description                                       string                                            `json:"description,omitempty"`
	DisplayName                                       string                                            `json:"displayName,omitempty"`
	Emails                                            []Emails                                          `json:"emails,omitempty"`
	Entitlements                                      []string                                          `json:"entitlements,omitempty"`
	Group                                             Group                                             `json:"group,omitempty"`
	Groups                                            []Groups                                          `json:"groups,omitempty"`
	ID                                                string                                            `json:"id,omitempty"`
	ItemsPerPage                                      int                                               `json:"itemsPerPage,omitempty"`
	Members                                           []Members                                         `json:"members,omitempty"`
	Meta                                              Meta                                              `json:"meta,omitempty"`
	Name                                              Name                                              `json:"name,omitempty"`
	Owner                                             Owner                                             `json:"owner,omitempty"`
	Password                                          string                                            `json:"password,omitempty"`
	PrivilegedData                                    []PrivilegedData                                  `json:"privilegedData,omitempty"`
	Rights                                            []string                                          `json:"rights,omitempty"`
	Schemas                                           []string                                          `json:"schemas,omitempty"`
	TotalResults                                      int                                               `json:"totalResults,omitempty"`
	Type                                              string                                            `json:"type,omitempty"`
	User                                              User                                              `json:"user,omitempty"`
	UserName                                          string                                            `json:"userName,omitempty"`
	UserType                                          string                                            `json:"userType,omitempty"`
	StartIndex                                        int                                               `json:"startIndex,omitempty"`
	UrnIetfParamsScimSchemasCyberark10Safe            UrnIetfParamsScimSchemasCyberark10Safe            `json:"urn:ietf:params:scim:schemas:cyberark:1.0:Safe,omitempty"`
	UrnIetfParamsScimSchemasCyberark10PrivilegedData  UrnIetfParamsScimSchemasCyberark10PrivilegedData  `json:"urn:ietf:params:scim:schemas:cyberark:1.0:PrivilegedData,omitempty"`
	UrnIetfParamsScimSchemasPam10LinkedObject         UrnIetfParamsScimSchemasPam10LinkedObject         `json:"urn:ietf:params:scim:schemas:pam:1.0:LinkedObject,omitempty"`
	UrnIetfParamsScimSchemasExtensionEnterprise20User UrnIetfParamsScimSchemasExtensionEnterprise20User `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
	UrnIetfParamsScimSchemasCyberark10SafeMember      UrnIetfParamsScimSchemasCyberark10SafeMember      `json:"urn:ietf:params:scim:schemas:cyberark:1.0:SafeMember"`
}

// ScimType2 should be used for the following:
//   - Containers
//   - Privileged Data
//   - Schemas
//   - ResourceTypes
type ScimType2 struct {
	Schemas      []string         `json:"schemas"`
	TotalResults int              `json:"totalResults"`
	ItemsPerPage int              `json:"itemsPerPage"`
	StartIndex   int              `json:"startIndex"`
	Resources    []Type2Resources `json:"Resources"`
}

type Type2Resources struct {
	Attributes                                        []Attributes                                      `json:"attributes,omitempty"`
	DisplayName                                       string                                            `json:"displayName,omitempty"`
	Description                                       string                                            `json:"description,omitempty"`
	Endpoint                                          string                                            `json:"endpoint,omitempty"`
	ID                                                string                                            `json:"id,omitempty"`
	Meta                                              Meta                                              `json:"meta,omitempty"`
	Name                                              string                                            `json:"name,omitempty"`
	Owner                                             Owner                                             `json:"owner,omitempty"`
	PrivilegedData                                    []PrivilegedData                                  `json:"privilegedData,omitempty"`
	SchemaExtensions                                  []SchemaExtensions                                `json:"schemaExtensions,omitempty"`
	Operations                                        []Operations                                      `json:"Operations"`
	Schemas                                           []string                                          `json:"schemas,omitempty"`
	Type                                              string                                            `json:"type,omitempty"`
	UniqueSafeId                                      string                                            `json:"uniqueSafeId,omitempty"`
	UrnIetfParamsScimSchemasCyberark10Safe            UrnIetfParamsScimSchemasCyberark10Safe            `json:"urn:ietf:params:scim:schemas:cyberark:1.0:Safe,omitempty"`
	UrnIetfParamsScimSchemasCyberark10PrivilegedData  UrnIetfParamsScimSchemasCyberark10PrivilegedData  `json:"urn:ietf:params:scim:schemas:cyberark:1.0:PrivilegedData,omitempty"`
	UrnIetfParamsScimSchemasPam10LinkedObject         UrnIetfParamsScimSchemasPam10LinkedObject         `json:"urn:ietf:params:scim:schemas:pam:1.0:LinkedObject,omitempty"`
	UrnIetfParamsScimSchemasExtensionEnterprise20User UrnIetfParamsScimSchemasExtensionEnterprise20User `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
	UrnIetfParamsScimSchemasCyberark10SafeMember      UrnIetfParamsScimSchemasCyberark10SafeMember      `json:"urn:ietf:params:scim:schemas:cyberark:1.0:SafeMember"`
}

// All structs below are utilized in the above top level structs
type Container struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Name    string `json:"name"`
	Display string `json:"display"`
}

type User struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Display string `json:"display"`
}

type Group struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Display string `json:"display"`
}

type Groups struct {
	Type    string `json:"type"`
	Display string `json:"display"`
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
}

type Owner struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Display string `json:"display"`
}

type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Location     string    `json:"location"`
}

type Properties struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Members struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Display string `json:"display"`
}

type PrivilegedData struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
	Display string `json:"display"`
}

type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

type Emails struct {
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
	Value   string `json:"value,omitempty"`
}

type Patch struct {
	Supported bool `json:"supported"`
}

type Bulk struct {
	Supported      bool `json:"supported"`
	MaxOperations  int  `json:"maxOperations"`
	MaxPayloadSize int  `json:"maxPayloadSize"`
}

type Filter struct {
	Supported  bool `json:"supported"`
	MaxResults int  `json:"maxResults"`
}

// Used in Privileged Data PATCH functions
type Value struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Used in Privileged Data PATCH functions
type Operations struct {
	Op    string  `json:"op"`
	Path  string  `json:"path"`
	Value []Value `json:"value"`
}

type ChangePassword struct {
	Supported bool `json:"supported"`
}

type Sort struct {
	Supported bool `json:"supported"`
}

type Etag struct {
	Supported bool `json:"supported"`
}

type AuthenticationSchemes struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SchemaExtensions struct {
	Schema   string `json:"schema"`
	Required bool   `json:"required"`
}

type SubAttributes struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	MultiValued bool   `json:"multiValued"`
	Required    bool   `json:"required"`
	CaseExact   bool   `json:"caseExact"`
}

type Attributes struct {
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	MultiValued   bool            `json:"multiValued"`
	Required      bool            `json:"required"`
	CaseExact     bool            `json:"caseExact,omitempty"`
	SubAttributes []SubAttributes `json:"subAttributes,omitempty"`
	Mutability    string          `json:"mutability,omitempty"`
	Returned      string          `json:"returned,omitempty"`
}

type UrnIetfParamsScimSchemasCyberark10SafeMember struct {
}

type UrnIetfParamsScimSchemasPam10LinkedObject struct {
	Source           string `json:"source"`
	NativeIdentifier string `json:"nativeIdentifier"`
}

type UrnIetfParamsScimSchemasExtensionEnterprise20User struct {
	Organization string `json:"organization"`
}

type UrnIetfParamsScimSchemasCyberark10Safe struct {
	NumberOfDaysRetention int    `json:"NumberOfDaysRetention"`
	ManagingCPM           string `json:"ManagingCPM"`
}

type UrnIetfParamsScimSchemasCyberark10PrivilegedData struct {
	Safe       string       `json:"safe"`
	Properties []Properties `json:"properties"`
}
