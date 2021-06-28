package cfg

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/eaglebush/stdutil"
)

// DatabaseKeyword - struct for database keywords
type DatabaseKeyword struct {
	Flag
}

// APIKeys
type APIKeyInfo struct {
	ID   string
	Name string
	Key  string
}

// Endpoint - endpoint struct
type EndpointInfo struct {
	ID      string  // Endpoint ID for quick access
	Name    string  // Endpoint Name to show
	GroupID *string // A group id to get certain endpoint set
	Address string  // The absolute URL to the resource
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
	ID                     string                 `json:"ID"`                               // A unique ID that will identify the connection to a database
	ConnectionString       string                 `json:"ConnectionString"`                 // ConnectionString specific to the database
	DriverName             string                 `json:"DriverName"`                       // DriverName needs to be specified depending on the driver id used by the Go database driver
	StorageType            string                 `json:"StorageType"`                      // FILE for filebased database such as Access, SQlite or LocalDB. SERVER for SQL Server, MySQL etc
	ParameterPlaceholder   string                 `json:"ParameterPlaceholder,omitempty"`   // Parameter place holder for prepared statements. Default is '?'
	ParameterInSequence    bool                   `json:"ParameterInSequence,omitempty"`    // Parameter place holder is in sequence. Default is false
	Schema                 string                 `json:"Schema,omitempty"`                 // Schema for any of the database operations
	InterpolateTables      *bool                  `json:"InterpolateTables,omitempty"`      // Enables the tables to be interpolated with schema
	GroupID                *string                `json:"GroupID,omitempty"`                // GroupID allows us to get groups of connection
	SequenceGenerator      *SequenceGeneratorInfo `json:"SequenceGenerator,omitempty"`      // Sequence generator configuration
	StringEnclosingChar    *string                `json:"StringEnclosingChar,omitempty"`    // Gets or sets the character that encloses a string in the query
	StringEscapeChar       *string                `json:"StringEscapeChar,omitempty"`       // Gets or Sets the character that escapes a reserved character such as the character that encloses a s string
	IdentityQuery          *string                `json:"IdentityQuery,omitempty"`          // A query to get the generated identity
	DateFunction           *string                `json:"DateFunction,omitempty"`           // The date function of each SQL database driver
	UTCDateFunction        *string                `json:"UTCDateFunction,omitempty"`        // The UTC date function of each SQL database driver
	MaxOpenConnection      *int                   `json:"MaxOpenConnection,omitempty"`      // Maximum open connection
	MaxIdleConnection      *int                   `json:"MaxIdleConnection,omitempty"`      // Maximum idle connection
	MaxConnectionLifetime  *int                   `json:"MaxConnectionLifetime,omitempty"`  // Max connection lifetime
	Ping                   *bool                  `json:"Ping,omitempty"`                   // Ping connection
	ReservedWordEscapeChar *string                `json:"ReservedWordEscapeChar,omitempty"` // Reserved word escape chars. For escaping with different opening and closing characters, just set to both. Example. `[]` for SQL server
	KeywordMap             *[]DatabaseKeyword     `json:"KeywordMap,omitempty"`             // various keyword equivalents
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
	APIEndpoints          *[]EndpointInfo     `json:"APIEndpoints,omitempty"`          // External API endpoints that this application can communicate
	APIKeys               *[]APIKeyInfo       `json:"APIKeys,omitempty"`               // API Keys
	ApplicationID         *string             `json:"ApplicationID,omitempty"`         // ID of this application
	ApplicationName       *string             `json:"ApplicationName,omitempty"`       // Name of this application
	ApplicationTheme      *string             `json:"ApplicationTheme,omitempty"`      // Theme of this application
	Cache                 *CacheInfo          `json:"Cache,omitempty"`                 // Cache info of this application
	CertificateFile       *string             `json:"CertificateFile,omitempty"`       // Certificate file
	CertificateKey        *string             `json:"CertificateKey,omitempty"`        // Certificate private key
	CookieDomain          *string             `json:"CookieDomain,omitempty"`          // The domain of the cookie that this application will send
	CrossOriginDomains    *[]string           `json:"CrossOriginDomains,omitempty"`    // Domains or endpoints that this application will allow
	Databases             *[]DatabaseInfo     `json:"Databases,omitempty"`             // Configured databases for this application use
	DefaultDatabaseID     *string             `json:"DefaultDatabaseID,omitempty"`     // The default database id that this application will find on the database configuration
	DefaultEndpointID     *string             `json:"DefaultEndpointID,omitempty"`     // The default endpoint that this application will find on the API endpoints configuration
	DefaultNotificationID *string             `json:"DefaultNotificationID,omitempty"` // The default notification id that this application will find on the notification configuration
	Domains               *[]DomainInfo       `json:"Domains,omitempty"`               // Configured domains for this application use
	FileName              string              `json:"FileName,omitempty"`              // Filename of the current configuration
	Flags                 *[]Flag             `json:"Flags,omitempty"`                 // Miscellaneous flags for this application use
	HostInternalURL       *string             `json:"HostInternalURL,omitempty"`       // The internal host URL that this application will use to set returned resources and assets
	HostExternalURL       *string             `json:"HostExternalURL,omitempty"`       // The external host URL that this application will use to set returned resources and assets
	HostPort              *int                `json:"HostPort,omitempty"`              // The network port for the application
	JWTSecret             *string             `json:"JWTSecret,omitempty"`             // Application wide JSON Web Token (JWT) secret
	LicenseSerial         *string             `json:"LicenseSerial,omitempty"`         // License serial of this application
	Notifications         *[]NotificationInfo `json:"Notifications,omitempty"`         // Configured notifications for this application use
	Queue                 *QueueInfo          `json:"Queue,omitempty"`                 // Queue or message queue
	ReadTimeout           *int                `json:"ReadTimeout,omitempty"`           // Default network timeout setting for reading data uploaded to this application
	Secure                *bool               `json:"Secure,omitempty"`                // Flags if secure
	Sources               *[]SourceInfo       `json:"Sources,omitempty"`               // Folder sources
	WriteTimeout          *int                `json:"WriteTimeout,omitempty"`          // Default network timeout setting for writing data downloaded from this application
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

	if config.DefaultDatabaseID == nil || *config.DefaultDatabaseID == "" {
		config.DefaultDatabaseID = stdutil.NewString(def)
	}

	if config.DefaultEndpointID == nil || *config.DefaultEndpointID == "" {
		config.DefaultEndpointID = stdutil.NewString(def)
	}

	if config.DefaultNotificationID == nil || *config.DefaultNotificationID == "" {
		config.DefaultNotificationID = stdutil.NewString(def)
	}

	// Default setting for database
	if config.Databases != nil {
		dbs := *config.Databases
		for i, cd := range dbs {

			if cd.InterpolateTables == nil {
				cd.InterpolateTables = new(bool)
				*cd.InterpolateTables = true
			}

			if cd.StringEnclosingChar == nil || *cd.StringEnclosingChar == "" {
				cd.StringEnclosingChar = new(string)
				*cd.StringEnclosingChar = `'`
			}

			if cd.StringEscapeChar == nil || *cd.StringEscapeChar == "" {
				cd.StringEscapeChar = new(string)
				*cd.StringEscapeChar = `\`
			}

			if cd.ReservedWordEscapeChar == nil || *cd.ReservedWordEscapeChar == "" {
				cd.ReservedWordEscapeChar = new(string)
				*cd.ReservedWordEscapeChar = `"`
			}

			if cd.ParameterPlaceholder == "" {
				cd.ParameterPlaceholder = `?`
			}

			if cd.StorageType == "" {
				cd.StorageType = `SERVER`
			}

			cd.StorageType = strings.ToUpper(cd.StorageType)

			drivern := strings.ToLower(cd.DriverName)
			if cd.StorageType == `SERVER` && (drivern == `sqlserver` || drivern == `mssql`) {

				if cd.IdentityQuery == nil || *cd.IdentityQuery == "" {
					cd.IdentityQuery = new(string)
					*cd.IdentityQuery = `SELECT SCOPE_IDENTITY();`
				}

				if cd.UTCDateFunction == nil || *cd.UTCDateFunction == "" {
					cd.UTCDateFunction = new(string)
					*cd.UTCDateFunction = `GETUTCDATE()`
				}

				if cd.DateFunction == nil || *cd.DateFunction == "" {
					cd.DateFunction = new(string)
					*cd.DateFunction = `GETDATE()`
				}

			}

			dbs[i] = cd
		}

		config.Databases = &dbs
	}

	// check if there is a notification
	defnum := ""
	if config.Notifications != nil {
		nfs := *config.Notifications

		for i, cn := range nfs {
			if i > 0 {
				defnum = strconv.Itoa(i)
			}

			if cn.ID == "" {
				nfs[i].ID = def + defnum
			}
		}

		config.Notifications = &nfs
	}

	return config, nil
}

// GetDatabaseInfo - get a database info by ID
func (c *Configuration) GetDatabaseInfo(ConnectionID string) *DatabaseInfo {

	if c.Databases == nil {
		return &DatabaseInfo{}
	}

	for _, v := range *c.Databases {
		if v.ID == ConnectionID {
			return &v
		}
	}

	return nil
}

// GetDatabaseInfoGroup - get database infos based on the group id
func (c *Configuration) GetDatabaseInfoGroup(GroupID string) *[]DatabaseInfo {

	dbgi := make([]DatabaseInfo, 0)

	if c.Databases == nil {
		return &dbgi
	}

	for _, v := range *c.Databases {

		if v.GroupID == nil {
			continue
		}

		if *v.GroupID == GroupID {
			dbgi = append(dbgi, v)
		}
	}

	return &dbgi
}

// GetDomainInfo - get a domain info by name
func (c *Configuration) GetDomainInfo(DomainName string) *DomainInfo {

	if c.Domains == nil {
		return &DomainInfo{}
	}

	for _, v := range *c.Domains {
		if v.Name == DomainName {
			return &v
		}
	}

	return nil
}

// GetEndpointAddress - get an endpoint value
func (c *Configuration) GetEndpointAddress(id ...string) string {

	if c.APIEndpoints != nil {
		return ""
	}

	var k string
	if len(id) > 0 {
		k = strings.ToLower(id[0])
	}

	if k == "" {
		k = strings.ToLower(*c.DefaultEndpointID)
	}

	if k == "" {
		return ""
	}

	eps := *c.APIEndpoints

	for i := range eps {
		k2 := strings.TrimSpace(strings.ToLower(eps[i].ID))
		if k == k2 {
			return eps[i].Address
		}
	}

	return ""
}

// GetEndpointInfo - get an endpoint by id
func (c *Configuration) GetEndpointInfo(id ...string) *EndpointInfo {

	if c.APIEndpoints != nil {
		return nil
	}

	var k string
	if len(id) > 0 {
		k = strings.ToLower(id[0])
	}

	if k == "" {
		k = strings.ToLower(*c.DefaultEndpointID)
	}

	if k == "" {
		return nil
	}

	eps := *c.APIEndpoints

	for i := range eps {
		k2 := strings.TrimSpace(strings.ToLower(eps[i].ID))
		if k == k2 {
			return &eps[i]
		}
	}

	return nil
}

// GetDatabaseInfoGroup - get database infos based on the group id
func (c *Configuration) GetEndpointInfoGroup(GroupID string) *[]EndpointInfo {

	ee := make([]EndpointInfo, 0)

	if c.APIEndpoints != nil {
		return &ee
	}

	eps := *c.APIEndpoints

	for _, v := range eps {

		if v.GroupID == nil {
			continue
		}

		if *v.GroupID == GroupID {
			ee = append(ee, v)
		}
	}

	return &ee
}

// GetNotificationInfo - get notification info
func (c *Configuration) GetNotificationInfo(id ...string) (NotificationInfo, error) {

	ni := NotificationInfo{}

	if c.Notifications == nil {
		return ni, nil
	}

	var k string
	if len(id) > 0 {
		k = strings.ToLower(id[0])
	}

	if k == "" {
		k = strings.ToLower(*c.DefaultNotificationID)
	}

	nfs := *c.Notifications

	if len(nfs) == 0 {
		return ni, errors.New("No notification configuration could be found")
	}

	for i := range nfs {
		k2 := strings.TrimSpace(strings.ToLower(nfs[i].ID))
		if k == k2 {
			return nfs[i], nil
		}
	}

	return ni, errors.New("Notification could not be found")
}

// GetSourceInfo - get source by id
func (c *Configuration) GetSourceInfo(SourceID string) (*SourceInfo, error) {

	if c.Sources == nil {
		return nil, nil
	}

	srcs := *c.Sources

	for _, v := range srcs {
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
