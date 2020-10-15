package cfg

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

// CacheInfo connection information
type CacheInfo struct {
	Provider string
	Address  string
	Password string
	DB       int
}

// DomainInfo - domain information for LDAP authentication
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
	ID                     string                // A unique ID that will identify the connection to a database
	ConnectionString       string                // ConnectionString specific to the database
	DriverName             string                // DriverName needs to be specified depending on the driver id used by the Go database driver
	StorageType            string                // FILE for filebased database such as Access, SQlite or LocalDB. SERVER for SQL Server, MySQL etc
	ParameterPlaceholder   string                // Parameter place holder for prepared statements. Default is '?'
	ParameterInSequence    bool                  // Parameter place holder is in sequence. Default is false
	GroupID                string                // GroupID allows us to get groups of connection
	SequenceGenerator      SequenceGeneratorInfo // Sequence generator configuration
	StringEnclosingChar    string                // Gets or sets the character that encloses a string in the query
	StringEscapeChar       string                // Gets or Sets the character that escapes a reserved character such as the character that encloses a s string
	Schema                 string                // Schema for any of the database operations
	IdentityQuery          string                // A query to get the generated identity
	DateFunction           string                // The date function of each SQL database driver
	UTCDateFunction        string                // The UTC date function of each SQL database driver
	MaxOpenConnection      int                   // Maximum open connection
	MaxIdleConnection      int                   // Maximum idle connection
	MaxConnectionLifetime  int                   // Max connection lifetime
	Ping                   bool                  // Ping connection
	ReservedWordEscapeChar string                // Reserved word escape chars. For escaping with different opening and closing characters, just set to both. Example. `[]` for SQL server
	KeywordMap             []DatabaseKeyword     // various keyword equivalents
}

// NotificationRecipient - notification standard recipients
type NotificationRecipient struct {
	ID             string
	ContactName    string
	ContactAddress string
}

// QueueInfo - queue info connector
type QueueInfo struct {
	ID            string `json:"ID,omitempty"`            // ID of the setting
	ServerAddress string `json:"ServerAddress,omitempty"` // Queue server address
	Cluster       string `json:"Cluster,omitempty"`       // Cluster name
	ClientID      string `json:"ClientID,omitempty"`      // ClientID of the service
}

// SourceInfo - file sources for configuration
type SourceInfo struct {
	ID            string `json:"id,omitempty"`        // ID of the source for quick reference
	Type          string `json:"type,omitempty"`      // Type of Inbound file. Supported types are ORDER and SNAPSHOT
	FolderSource  string `json:"source,omitempty"`    // Source folder of the source
	FolderError   string `json:"error,omitempty"`     // Error folder of the source
	FolderSuccess string `json:"success,omitempty"`   // Success folder of the source
	FileExtension string `json:"extension,omitempty"` // Extension of the file to pickup
}

// Configuration - for various configuration settings. This struct can be modified depending on the requirement.
type Configuration struct {
	APIEndpoints          []Endpoint         `json:"APIEndpoints,omitempty"`          // External API endpoints that this application can communicate
	APIKey                string             `json:"APIKey,omitempty"`                // Registration key
	ApplicationID         string             `json:"ApplicationID,omitempty"`         // ID of this application
	ApplicationName       string             `json:"ApplicationName,omitempty"`       // Name of this application
	ApplicationTheme      string             `json:"ApplicationTheme,omitempty"`      // Theme of this application
	Cache                 CacheInfo          `json:"Cache,omitempty"`                 // Cache info of this application
	CertificateFile       string             `json:"CertificateFile,omitempty"`       // Certificate file
	CertificateKey        string             `json:"CertificateKey,omitempty"`        // Certificate private key
	CookieDomain          string             `json:"CookieDomain,omitempty"`          // The domain of the cookie that this application will send
	CrossOriginDomains    []string           `json:"CrossOriginDomains,omitempty"`    // Domains or endpoints that this application will allow
	Databases             []DatabaseInfo     `json:"Databases,omitempty"`             // Configured databases for this application use
	DefaultDatabaseID     string             `json:"DefaultDatabaseID,omitempty"`     // The default database id that this application will find on the database configuration
	DefaultEndpointID     string             `json:"DefaultEndpointID,omitempty"`     // The default endpoint that this application will find on the API endpoints configuration
	DefaultNotificationID string             `json:"DefaultNotificationID,omitempty"` // The default notification id that this application will find on the notification configuration
	Domains               []DomainInfo       `json:"Domains,omitempty"`               // Configured domains for this application use
	FileName              string             `json:"FileName,omitempty"`              // Filename of the current configuration
	Flags                 []Flag             `json:"Flags,omitempty"`                 // Miscellaneous flags for this application use
	HostInternalURL       string             `json:"HostInternalURL,omitempty"`       // The internal host URL that this application will use to set returned resources and assets
	HostExternalURL       string             `json:"HostExternalURL,omitempty"`       // The external host URL that this application will use to set returned resources and assets
	HostPort              int                `json:"HostPort,omitempty"`              // The network port for the application
	JWTSecret             string             `json:"JWTSecret,omitempty"`             // Application wide JSON Web Token (JWT) secret
	LicenseSerial         string             `json:"LicenseSerial,omitempty"`         // License serial of this application
	Notifications         []NotificationInfo `json:"Notifications,omitempty"`         // Configured notifications for this application use
	Queue                 QueueInfo          `json:"Queue,omitempty"`                 // Queue or message queue
	ReadTimeout           int                `json:"ReadTimeout,omitempty"`           // Default network timeout setting for reading data uploaded to this application
	Secure                bool               `json:"Secure,omitempty"`                // Flags if secure
	Sources               []SourceInfo       `json:"Sources,omitempty"`               // Folder sources
	WriteTimeout          int                `json:"WriteTimeout,omitempty"`          // Default network timeout setting for writing data downloaded from this application
	errorText             string
}

// LoadConfig - load configuration file and return a configuration
func LoadConfig(Source string) (*Configuration, error) {
	return load(Source)
}

func load(Source string) (*Configuration, error) {

	var (
		err       error
		b         []byte
		islocfile bool
	)

	const def string = `DEFAULT`

	config := &Configuration{
		errorText: ``,
	}

	if !(strings.HasPrefix(Source, `http://`) || strings.HasPrefix(Source, `https://`)) {
		islocfile = true
	}

	if islocfile {
		if b, err = ioutil.ReadFile(Source); err != nil {
			config.errorText = err.Error()
			return config, err
		}
	} else {
		if b, err = func() ([]byte, error) {

			var ob []byte

			nr, err := http.Get(Source)
			if err != nil {
				return ob, err
			}
			defer nr.Body.Close()

			ob, err = ioutil.ReadAll(nr.Body)
			if err != nil {
				return ob, err
			}

			return ob, nil
		}(); err != nil {
			config.errorText = err.Error()
			return config, err
		}
	}

	if len(b) == 0 {
		err = errors.New(`No data from source for configuration`)
		config.errorText = err.Error()
		return config, err
	}

	config.FileName = Source

	err = json.Unmarshal(b, config)
	if err != nil {
		config.errorText = err.Error()
		return nil, err
	}

	if config.DefaultDatabaseID == "" {
		config.DefaultDatabaseID = def
	}

	if config.DefaultEndpointID == "" {
		config.DefaultEndpointID = def
	}

	if config.DefaultNotificationID == "" {
		config.DefaultNotificationID = def
	}

	// Default setting for database
	for i, cd := range config.Databases {

		if cd.StringEnclosingChar == "" {
			cd.StringEnclosingChar = `'`
		}

		if cd.StringEscapeChar == "" {
			cd.StringEscapeChar = `\`
		}

		if cd.ParameterPlaceholder == "" {
			config.Databases[i].ParameterPlaceholder = `?`
		}

		if cd.ReservedWordEscapeChar == "" {
			cd.ReservedWordEscapeChar = `"`
		}

		if cd.StorageType == "" {
			config.Databases[i].StorageType = `SERVER`
		}
		cd.StorageType = strings.ToUpper(cd.StorageType)

		drivern := strings.ToLower(cd.DriverName)
		if cd.StorageType == `SERVER` && (drivern == `sqlserver` || drivern == `mssql`) {

			if config.Databases[i].IdentityQuery == "" {
				config.Databases[i].IdentityQuery = `SELECT SCOPE_IDENTITY();`
			}

			if config.Databases[i].UTCDateFunction == "" {
				config.Databases[i].UTCDateFunction = `GETUTCDATE()`
			}

			if config.Databases[i].DateFunction == "" {
				config.Databases[i].DateFunction = `GETDATE()`
			}

		}
	}

	// check if there is a notification
	defnum := ""
	for i, cn := range config.Notifications {
		if i > 0 {
			defnum = strconv.Itoa(i)
		}

		if cn.ID == "" {
			config.Notifications[i].ID = def + defnum
		}
	}

	return config, nil
}

// GetDatabaseInfo - get a database info by ID
func (c *Configuration) GetDatabaseInfo(ConnectionID string) *DatabaseInfo {
	for _, v := range c.Databases {
		if v.ID == ConnectionID {
			return &v
		}
	}

	return nil
}

// GetDatabaseInfoGroup - get database infos based on the group id
func (c *Configuration) GetDatabaseInfoGroup(GroupID string) *[]DatabaseInfo {
	dbgi := make([]DatabaseInfo, 0)
	for _, v := range c.Databases {
		if v.GroupID == GroupID {
			dbgi = append(dbgi, v)
		}
	}

	return &dbgi
}

// GetDomainInfo - get a domain info by name
func (c *Configuration) GetDomainInfo(DomainName string) *DomainInfo {
	for _, v := range c.Domains {
		if v.Name == DomainName {
			return &v
		}
	}

	return nil
}

// GetEndpoint - get an endpoint value
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

// GetSourceInfo - get source by id
func (c *Configuration) GetSourceInfo(SourceID string) (*SourceInfo, error) {
	for _, v := range c.Sources {
		if v.ID == SourceID {
			return &v, nil
		}
	}

	return &SourceInfo{}, nil
}

// Save - save configuration file
func (c *Configuration) Save() bool {
	var b []byte
	var err error

	c.errorText = ""
	if b, err = json.MarshalIndent(c, "", "\t"); err != nil {
		c.errorText = err.Error()
		return false
	}

	if err = ioutil.WriteFile(c.FileName, b, os.ModePerm); err != nil {
		c.errorText = err.Error()
		return false
	}
	return true
}

// Reload configuration
func (c *Configuration) Reload() error {
	_, err := load(c.FileName)
	return err
}

// LastErrorText - gets the last error
func (c *Configuration) LastErrorText() string {
	return c.errorText
}
