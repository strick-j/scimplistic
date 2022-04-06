package types

import (
	"github.com/gorilla/mux"
)

// CreateForm holds the initial fields required to setup a form
// Action: "DestinationURL"
// Method: "Get, Post, etc..."
// Legend: "Form Title"
// Role:   "adduser,addgroup,addsafe,etc..."
type CreateForm struct {
	FormEncType string       `json:"forEncyType"`
	FormAction  string       `json:"formAction"`
	FormMethod  string       `json:"formMethod"`
	FormLegend  string       `json:"formLegend"`
	FormRole    string       `json:"formRole,omitempty"`
	FormFields  []FormFields `json:"formFields"`
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
	Navigation         string         `json:"navigation,omitempty"`
	Message            string         `json:"message,omitempty"`
	SettingsConfigured bool           `json:"settingsConfigured,omitempty"`
	CreateForm         CreateForm     `json:"createForm,omitempty"`
	SecretForm         CreateForm     `json:"secretForm,omitempty"`
	Safes              ScimType2      `json:"safes,omitempty"`
	Users              ScimType1      `json:"users,omitempty"`
	Groups             ScimType1      `json:"groups,omitempty"`
	Accounts           ScimType2      `json:"accounts,omitempty"`
	Members            Members        `json:"members,omitempty"`
	Settings           ConfigSettings `json:"settings,omitempty"`
}

// ConfigSettings is the struct to hold settings information
// used in conjunction with MarshallIndent to write configuration
type ConfigSettings struct {
	ScimURL         string      `json:"scimURL,omitempty"`     // e.g. <identity tenant id>.my.idaptive.app
	ApiEndpoint     string      `json:"apiEndpoint,omitempty"` // e.g. "scim"
	ApiVersion      string      `json:"apiVersion,omitempty"`  //e.g. "v2"
	AuthToken       string      `json:"authToken,omitempty"`
	DatabaseIP      string      `json:"databaseIp,omitempty"`   // IP for Postgres Database
	DatabasePort    int         `json:"databasePort,omitempty"` // Port used for Postgres Database
	DatabaseName    string      `json:"databaseName,omitempty"` // Name for Postgres Database
	DatabaseUser    string      `json:"databaseUser,omitempty"` // Port used for Postgres Database
	DatabasePass    string      `json:"databasePass,omitempty"` // Name for Postgres Database
	DatabaseEnabled bool        `json:"databaseEnabled,omitempty"`
	AuthDisabled    bool        `json:"authDisabled,omitempty"` // Default setting when no SCIM Authentication is configured
	TokenEnabled    bool        `json:"tokenEnabled,omitempty"` // True when Oauth Bearer Token Provided
	CredEnabled     bool        `json:"credEnabled,omitempty"`  // True when Client Credentials Provided
	PrevConf        bool        `json:"prevConf,omitempty"`
	ServerURL       string      `json:"serverUrl,omitempty"`
	EnableHTTPS     bool        `json:"enableHTTPS,omitempty"`
	Schema          string      `json:"schema,omitempty"`
	ServerName      string      `json:"serverName,omitempty"`
	MaxConnections  int         `json:"maxConnections"`        // The maximum amount of concurrent connections the server will accept. Setting this to 0 means infinite.
	HostName        string      `json:"hostName"`              // Server's host name. Use 'https://' for TLS connections. (ex: 'https://example.com') (Required)
	HostAlias       string      `json:"hostAlias"`             // Server's host alias name. Use 'https://' for TLS connections. (ex: 'https://www.example.com')
	IP              string      `json:"ip"`                    // Server's IP address. (Required)
	Port            int         `json:"port"`                  // Server's port. (Required)
	TLS             bool        `json:"tls"`                   // Enables TLS/SSL connections.
	DB              bool        `json:"db,omitempty"`          // Enables Database
	CertFile        string      `json:"certfile,omitempty"`    // SSL/TLS certificate file location (starting from system's root folder). (Required for TLS)
	PrivKeyFile     string      `json:"privKeyFile,omitempty"` // SSL/TLS private key file location (starting from system's root folder). (Required for TLS)
	OriginOnly      bool        `json:"originOnly,omitempty"`  // When enabled, the server declines connections made from outside the origin server (Admin logins always check origin). IMPORTANT: Enable this for web apps and LAN servers.
	Router          *mux.Router `json:"router,omitempty"`
	ClientID        string      `json:"clientID,omitempty"`
	ClientSecret    string      `json:"clientSecret,omitempty"`
	ClientAppID     string      `json:"clientAppID,omitempty"`
	ServiceLogging  bool        `json:"serviceLogging,omitempty"` // Sets log level for API Service functions
	LogLevel        string      `json:"logLevel"`                 // Server Log Level - Trace, Debug, Info, Warn, Error, Fatal, Panic
}
