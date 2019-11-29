package cfg

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// Flag - dynamic flags structure
type Flag struct {
	Key   string
	Value string
}

// DatabaseKeyword - struct for database keywords
type DatabaseKeyword struct {
	Flag
}

// Endpoint - endpoint struct
type Endpoint struct {
	ID      string
	Name    string
	Address string
}

// NotificationInfo - notification information on connecting to Notify API
type NotificationInfo struct {
	ID            string
	APIHost       string
	APIPath       string
	Type          string
	Login         string
	Password      string
	Active        bool
	SenderAddress string
	SenderName    string
	Recipients    []NotificationRecipient
}

//DomainInfo - domain information for LDAP authentication
type DomainInfo struct {
	Name               string
	Host               string
	Port               uint16
	Path               string
	AuthorizedUser     string
	AuthorizedPassword string
	Filter             string
}

//SequenceGeneratorInfo - sequence generator query
type SequenceGeneratorInfo struct {
	UpsertQuery     string
	ResultQuery     string
	NamePlaceHolder string
}

//DatabaseInfo - database configuration setting
type DatabaseInfo struct {
	ID                    string                // A unique ID that will identify the connection to a database
	ConnectionString      string                // ConnectionString specific to the database
	DriverName            string                // DriverName needs to be specified depending on the driver id used by the Go database driver
	StorageType           string                // FILE for filebased database such as Access, SQlite or LocalDB. SERVER for SQL Server, MySQL etc
	ParameterPlaceholder  string                // Parameter place holder for prepared statements. Default is '?'
	ParameterInSequence   bool                  // Parameter place holder is in sequence. Default is false
	GroupID               string                // GroupID allows us to get groups of connection
	SequenceGenerator     SequenceGeneratorInfo // Sequence generator configuration
	IdentityQuery         string                // A query to get the generated identity
	DateFunction          string                // The date function of each SQL database driver
	UTCDateFunction       string                // The UTC date function of each SQL database driver
	MaxOpenConnection     int                   // Maximum open connection
	MaxIdleConnection     int                   // Maximum idle connection
	MaxConnectionLifetime int                   // Max connection lifetime
	Ping                  bool                  // Ping connection
	KeywordMap            []DatabaseKeyword     // various keyword equivalents
}

//NotificationRecipient - notification standard recipients
type NotificationRecipient struct {
	ID             string
	ContactName    string
	ContactAddress string
}

//Configuration - for various configuration settings. This struct can be modified depending on the requirement.
type Configuration struct {
	APIEndpoints          []Endpoint         `json:"APIEndpoints,omitempty"`
	APIKey                string             `json:"APIKey,omitempty"`
	ApplicationID         string             `json:"ApplicationID,omitempty"`
	ApplicationName       string             `json:"ApplicationName,omitempty"`
	ApplicationTheme      string             `json:"ApplicationTheme,omitempty"`
	CookieDomain          string             `json:"CookieDomain,omitempty"`
	DefaultDatabaseID     string             `json:"DefaultDatabaseID,omitempty"`
	DefaultEndpointID     string             `json:"DefaultEndpointID,omitempty"`
	DefaultNotificationID string             `json:"DefaultNotificationID,omitempty"`
	HostInternalURL       string             `json:"HostInternalURL,omitempty"`
	HostExternalURL       string             `json:"HostExternalURL,omitempty"`
	HostPort              int                `json:"HostPort,omitempty"`
	HMAC                  string             `json:"HMAC,omitempty"`
	LicenseSerial         string             `json:"LicenseSerial,omitempty"`
	Databases             []DatabaseInfo     `json:"Databases,omitempty"`
	Domains               []DomainInfo       `json:"Domains,omitempty"`
	Notifications         []NotificationInfo `json:"Notifications,omitempty"`
	Flags                 []Flag             `json:"Flags,omitempty"`
	fileName              string
	errorText             string
}

// LoadConfig - load configuration file and return a configuration
func LoadConfig(fileName string) (*Configuration, error) {
	config := &Configuration{}
	config.errorText = ""

	b, err := ioutil.ReadFile(fileName)

	if err != nil {
		config.errorText = err.Error()
		return nil, err
	}

	config.fileName = fileName

	err = json.Unmarshal(b, config)
	if err != nil {
		config.errorText = err.Error()
		log.Fatal(err)
	}

	if err != nil {
		return nil, err
	}

	if config.DefaultDatabaseID == "" {
		config.DefaultDatabaseID = "DEFAULT"
	}

	if config.DefaultEndpointID == "" {
		config.DefaultEndpointID = "DEFAULT"
	}

	if config.DefaultNotificationID == "" {
		config.DefaultNotificationID = "DEFAULT"
	}

	// Default setting for database
	for i, cd := range config.Databases {

		if cd.ParameterPlaceholder == "" {
			config.Databases[i].ParameterPlaceholder = "?"
		}

		if cd.StorageType == "" {
			config.Databases[i].StorageType = "SERVER"
		}

		drivern := strings.ToLower(cd.DriverName)
		if cd.StorageType == "SERVER" && (drivern == "sqlserver" || drivern == "mssql") {
			config.Databases[i].IdentityQuery = "SELECT SCOPE_IDENTITY();"
		}
	}

	// check if there is a notification
	defnum := ""
	for i, cn := range config.Notifications {
		if i > 0 {
			defnum = string(i)
		}
		if cn.ID == "" {
			config.Notifications[i].ID = "DEFAULT" + defnum
		}
	}

	return config, nil
}

//GetDatabaseInfo - get a database info by ID
func (c *Configuration) GetDatabaseInfo(ConnectionID string) *DatabaseInfo {
	for _, v := range c.Databases {
		if v.ID == ConnectionID {
			return &v
		}
	}

	return nil
}

//GetDatabaseInfoGroup - get database infos based on the group id
func (c *Configuration) GetDatabaseInfoGroup(GroupID string) *[]DatabaseInfo {
	dbgi := make([]DatabaseInfo, 0)
	for _, v := range c.Databases {
		if v.GroupID == GroupID {
			dbgi = append(dbgi, v)
		}
	}

	return &dbgi
}

//GetDomainInfo - get a domain info by name
func (c *Configuration) GetDomainInfo(DomainName string) *DomainInfo {
	for _, v := range c.Domains {
		if v.Name == DomainName {
			return &v
		}
	}

	return nil
}

//GetEndpoint - get an endpoint value
func (c *Configuration) GetEndpoint(id ...string) string {
	var k string
	if len(id) > 0 {
		k = strings.ToLower(id[0])
	}

	if k == "" {
		k = strings.ToLower(c.DefaultEndpointID)
	}

	if k == "" {
		return ""
	}

	for i := range c.APIEndpoints {
		k2 := strings.TrimSpace(strings.ToLower(c.APIEndpoints[i].ID))
		if k == k2 {
			return c.APIEndpoints[i].Address
		}
	}

	return ""
}

// GetNotificationInfo - get notification info
func (c *Configuration) GetNotificationInfo(id ...string) (NotificationInfo, error) {
	ni := NotificationInfo{}

	var k string
	if len(id) > 0 {
		k = strings.ToLower(id[0])
	}

	if k == "" {
		k = strings.ToLower(c.DefaultNotificationID)
	}

	if len(c.Notifications) == 0 {
		return ni, errors.New("No notification configuration could be found")
	}

	for i := range c.Notifications {
		k2 := strings.TrimSpace(strings.ToLower(c.Notifications[i].ID))
		if k == k2 {
			return c.Notifications[i], nil
		}
	}

	return ni, errors.New("Notification could not be found")
}

// Save - save configuration file
func (c *Configuration) Save() bool {
	c.errorText = ""
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		c.errorText = err.Error()
		return false
	}

	err = ioutil.WriteFile(c.fileName, b, os.ModePerm)
	if err != nil {
		c.errorText = err.Error()
		return false
	}
	return true
}

// LastErrorText - gets the last error
func (c *Configuration) LastErrorText() string {
	return c.errorText
}

//Flag - get a flag value
func (c *Configuration) Flag(key string) Flag {
	k := strings.ToLower(key)

	for i := range c.Flags {
		k2 := strings.TrimSpace(strings.ToLower(c.Flags[i].Key))
		if k == k2 {
			return c.Flags[i]
		}
	}

	return Flag{}
}

// Bool - return a boolean from flag value
func (f Flag) Bool() bool {
	v := strings.TrimSpace(f.Value)
	v = strings.ToLower(v)
	switch v {
	case "1", "on", "yes", "enabled", "true":
		return true
	}
	return false
}

// Int64 - return an int64 from flag value
func (f Flag) Int64() int64 {
	v := strings.TrimSpace(f.Value)
	vi, _ := strconv.ParseInt(v, 0, 64)
	return vi
}

// Int - return an int from flag value
func (f Flag) Int() int {
	v := strings.TrimSpace(f.Value)
	vi, _ := strconv.Atoi(v)
	return vi
}

// String - return a string from flag value
func (f Flag) String() string {
	return strings.TrimSpace(f.Value)
}
