package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

//MailServer - mail server setting
type MailServer struct {
	Host       string
	Port       int
	User       string
	Password   string
	SenderName string
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

//NotificationInfo - notification configuration setting
type NotificationInfo struct {
	ID            string
	FullName      string
	EmailAddress  string
	MessengerName string
}

//Configuration - for various configuration settings. This struct can be modified depending on the requirement.
type Configuration struct {
	ApplicationID     string             `json:"ApplicationID,omitempty"`
	ApplicationTheme  string             `json:"ApplicationTheme,omitempty"`
	APIKey            string             `json:"APIKey,omitempty"`
	CookieDomain      string             `json:"CookieDomain,omitempty"`
	DefaultDatabaseID string             `json:"DefaultDatabaseID,omitempty"`
	HostPort          int                `json:"HostPort,omitempty"`
	HMAC              string             `json:"HMAC,omitempty"`
	LicenseSerial     string             `json:"LicenseSerial,omitempty"`
	MailServer        MailServer         `json:"MailServer,omitempty"`
	Databases         []DatabaseInfo     `json:"Databases,omitempty"`
	Domains           []DomainInfo       `json:"Domains,omitempty"`
	NotifyRecipients  []NotificationInfo `json:"NotifyRecipients,omitempty"`
	Flags             []Flag             `json:"Flags,omitempty"`
}

// LoadConfig - load configuration file and return a configuration
func LoadConfig(fileName string) (*Configuration, error) {
	config := &Configuration{}
	b, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, config)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return nil, err
	}

	for i := range config.Databases {
		if config.Databases[i].ParameterPlaceholder == "" {
			config.Databases[i].ParameterPlaceholder = "?"
		}

		if config.Databases[i].StorageType == "" {
			config.Databases[i].StorageType = "SERVER"
		}

		drivern := strings.ToLower(config.Databases[i].DriverName)
		if config.Databases[i].StorageType == "SERVER" && (drivern == "sqlserver" || drivern == "mssql") {
			config.Databases[i].IdentityQuery = "SELECT SCOPE_IDENTITY();"
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

//GetNotifyRecipient - get a notification recipient info by name
func (c *Configuration) GetNotifyRecipient(RecipientID string) *NotificationInfo {
	for _, v := range c.NotifyRecipients {
		if v.ID == RecipientID {
			return &v
		}
	}

	return nil
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
