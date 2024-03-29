package cfg

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// DatabaseKeyword - struct for database keywords
type DatabaseKeyword struct {
	Flag
}

// APIKeyInfo contains an API info
type APIKeyInfo struct {
	ID    string
	Name  string
	Key   string
	Token *string
}

// DirectoryInfo contains a directory info
type DirectoryInfo struct {
	GroupID     string
	Description string
	Items       []Flag
}

// Endpoint - endpoint struct
type EndpointInfo struct {
	ID      string  // Endpoint ID for quick access
	Name    string  // Endpoint Name to show
	Address string  // The absolute URL to the resource
	GroupID *string // A group id to get certain endpoint set
	Token   *string
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
	ReplyTo       string
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

// SequenceGeneratorInfo - sequence generator query
type SequenceGeneratorInfo struct {
	UpsertQuery     string
	ResultQuery     string
	NamePlaceHolder string
}

// DatabaseInfo - database configuration setting
type DatabaseInfo struct {
	GroupID                *string                `json:"GroupID,omitempty"`                // GroupID allows us to get groups of connection
	ID                     string                 `json:"ID,omitempty"`                     // A unique ID that will identify the connection to a database
	ConnectionString       string                 `json:"ConnectionString,omitempty"`       // ConnectionString specific to the database
	DriverName             string                 `json:"DriverName,omitempty"`             // DriverName needs to be specified depending on the driver id used by the Go database driver
	StorageType            string                 `json:"StorageType,omitempty"`            // FILE for filebased database such as Access, SQlite or LocalDB. SERVER for SQL Server, MySQL etc
	HelperID               string                 `json:"HelperID,omitempty"`               // When using github.com/NarsilWorks-Inc/datahelperlite, this is needed in the configuration file
	ParameterPlaceholder   string                 `json:"ParameterPlaceholder,omitempty"`   // Parameter place holder for prepared statements. Default is '?'
	ParameterInSequence    bool                   `json:"ParameterInSequence,omitempty"`    // Parameter place holder is in sequence. Default is false
	Schema                 string                 `json:"Schema,omitempty"`                 // Schema for any of the database operations
	InterpolateTables      *bool                  `json:"InterpolateTables,omitempty"`      // Enables the tables to be interpolated with schema
	SequenceGenerator      *SequenceGeneratorInfo `json:"SequenceGenerator,omitempty"`      // Sequence generator configuration
	StringEnclosingChar    *string                `json:"StringEnclosingChar,omitempty"`    // Gets or sets the character that encloses a string in the query
	StringEscapeChar       *string                `json:"StringEscapeChar,omitempty"`       // Gets or Sets the character that escapes a reserved character such as the character that encloses a s string
	IdentityQuery          *string                `json:"IdentityQuery,omitempty"`          // A query to get the generated identity
	DateFunction           *string                `json:"DateFunction,omitempty"`           // The date function of each SQL database driver
	UTCDateFunction        *string                `json:"UTCDateFunction,omitempty"`        // The UTC date function of each SQL database driver
	MaxOpenConnection      *int                   `json:"MaxOpenConnection,omitempty"`      // Maximum open connection
	MaxIdleConnection      *int                   `json:"MaxIdleConnection,omitempty"`      // Maximum idle connection
	MaxConnectionLifetime  *int                   `json:"MaxConnectionLifetime,omitempty"`  // Max connection lifetime
	MaxConnectionIdleTime  *int                   `json:"MaxConnectionIdleTime,omitempty"`  // Max idle connection lifetime
	Ping                   *bool                  `json:"Ping,omitempty"`                   // Ping connection
	ReservedWordEscapeChar *string                `json:"ReservedWordEscapeChar,omitempty"` // Reserved word escape chars. For escaping with different opening and closing characters, just set to both. Example. `[]` for SQL server
	KeywordMap             *[]DatabaseKeyword     `json:"KeywordMap,omitempty"`             // various keyword equivalents
}

// NotificationRecipient - notification standard recipients
type NotificationRecipient struct {
	ID      string
	Name    string
	Address string
}

// QueueInfo - queue info connector
type QueueInfo struct {
	ID                 string   `json:"ID,omitempty"`                 // ID of the setting
	ServerAddressGroup []string `json:"ServerAddressGroup,omitempty"` // Queue server address group
	Cluster            string   `json:"Cluster,omitempty"`            // Cluster name
	ClientID           string   `json:"ClientID,omitempty"`           // ClientID of the service
	StreamName         string   `json:"StreamName,omitempty"`         // Stream name
}

// SourceInfo - file sources for configuration
type SourceInfo struct {
	ID        string `json:"id,omitempty"`        // ID of the source for quick reference
	Type      string `json:"type,omitempty"`      // Type of Inbound file. Supported types are ORDER and SNAPSHOT
	Source    string `json:"source,omitempty"`    // Source folder of the source
	Relative  bool   `json:"relative"`            // Indicates that the Error and Success folders are relative to Source
	Error     string `json:"error,omitempty"`     // Error folder of the source
	Success   string `json:"success,omitempty"`   // Success folder of the source
	Extension string `json:"extension,omitempty"` // Extension of the file to pickup
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
	Directories           *[]DirectoryInfo    `json:"Directories,omitempty"`           // Configured directory for this application use
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
	local                 bool                // Local file
}

var (
	ErrNoDataFromSource = errors.New(`no data from source for configuration`)
	ErrSaveNotLocalFile = errors.New("configuration file is not local")
)

func load(source string) (*Configuration, error) {

	config := &Configuration{}
	if !(strings.HasPrefix(source, `http://`) || strings.HasPrefix(source, `https://`)) {
		config.local = true
	}

	var (
		err error
		b   []byte
	)
	if config.local {
		b, err = os.ReadFile(source)
	} else {
		b, err =
			func() ([]byte, error) {
				var ob []byte
				nr, err := http.Get(source)
				if err != nil {
					return ob, err
				}
				defer nr.Body.Close()

				ob, err = io.ReadAll(nr.Body)
				if err != nil {
					return ob, err
				}
				return ob, nil
			}()
	}
	if err != nil {
		return config, err
	}

	if len(b) == 0 {
		return config, ErrNoDataFromSource
	}

	err = json.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}

	const def string = `DEFAULT`
	if config.DefaultDatabaseID == nil || *config.DefaultDatabaseID == "" {
		config.DefaultDatabaseID = new_string(def)
	}
	if config.DefaultEndpointID == nil || *config.DefaultEndpointID == "" {
		config.DefaultEndpointID = new_string(def)
	}
	if config.DefaultNotificationID == nil || *config.DefaultNotificationID == "" {
		config.DefaultNotificationID = new_string(def)
	}
	if config.CookieDomain == nil {
		config.CookieDomain = new_string(`localhost`)
	}
	if config.JWTSecret == nil {
		config.JWTSecret = new_string(`defaultsecretkey`)
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
				cd.StringEnclosingChar = new_string(`'`)
			}
			if cd.StringEscapeChar == nil || *cd.StringEscapeChar == "" {
				cd.StringEscapeChar = new_string(`\`)
			}
			if cd.ReservedWordEscapeChar == nil || *cd.ReservedWordEscapeChar == "" {
				cd.ReservedWordEscapeChar = new_string(`"`)
			}
			if cd.ParameterPlaceholder == "" {
				cd.ParameterPlaceholder = `?`
			}
			if cd.StorageType == "" {
				cd.StorageType = `SERVER`
			} else {
				cd.StorageType = strings.ToUpper(cd.StorageType)
			}
			drivern := strings.ToLower(cd.DriverName)
			if cd.StorageType == `SERVER` && (drivern == `sqlserver` || drivern == `mssql`) {
				if cd.IdentityQuery == nil || *cd.IdentityQuery == "" {
					cd.IdentityQuery = new_string(`SELECT SCOPE_IDENTITY();`)
				}
				if cd.UTCDateFunction == nil || *cd.UTCDateFunction == "" {
					cd.UTCDateFunction = new_string(`GETUTCDATE()`)
				}
				if cd.DateFunction == nil || *cd.DateFunction == "" {
					cd.DateFunction = new_string(`GETDATE()`)
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

	config.FileName = source
	return config, nil
}

// GetDatabaseInfo get a database info by its ID
func (c *Configuration) GetDatabaseInfo(connectionId string) *DatabaseInfo {
	if c.Databases == nil {
		return nil
	}
	for _, v := range *c.Databases {
		if v.ID == connectionId {
			return &v
		}
	}
	return nil
}

// GetDatabaseInfoGroup gets database infos based on the group id
func (c *Configuration) GetDatabaseInfoGroup(groupId string) []DatabaseInfo {
	dbgi := make([]DatabaseInfo, 0)
	if c.Databases == nil || groupId == "" {
		return dbgi
	}
	for _, v := range *c.Databases {
		if v.GroupID == nil {
			continue
		}
		if strings.EqualFold(*v.GroupID, groupId) {
			dbgi = append(dbgi, v)
		}
	}
	return dbgi
}

// GetDirectory retrieves a directory under a group
func (c *Configuration) GetDirectory(groupId string) *DirectoryInfo {
	if c.Directories == nil || len(*c.Directories) == 0 {
		return nil
	}
	for _, dir := range *c.Directories {
		if strings.EqualFold(dir.GroupID, groupId) {
			return &dir
		}
	}
	return nil
}

// GetDirectoryItem retrieves a directory item under a group
func (c *Configuration) GetDirectoryItem(groupId, key string) *Flag {
	dir := c.GetDirectory(groupId)
	if dir == nil {
		return nil
	}
	for _, item := range dir.Items {
		if strings.EqualFold(item.Key, key) {
			return &item
		}
	}
	return nil
}

// GetDomainInfo gets a domain info by name
func (c *Configuration) GetDomainInfo(domainName string) *DomainInfo {
	if c.Domains == nil || domainName == "" {
		return nil
	}
	for _, v := range *c.Domains {
		if strings.EqualFold(v.Name, domainName) {
			return &v
		}
	}
	return nil
}

// GetEndpointInfo - get an endpoint by id
func (c *Configuration) GetEndpointInfo(id ...string) *EndpointInfo {
	if c.APIEndpoints == nil || (len(id) == 0 && (c.DefaultEndpointID == nil || *c.DefaultEndpointID == "")) {
		return nil
	}
	k := strings.ToLower(*c.DefaultEndpointID)
	if len(id) > 0 {
		k = strings.ToLower(id[0])
	}
	eps := *c.APIEndpoints
	for _, ep := range eps {
		if strings.EqualFold(k, ep.ID) {
			return &ep
		}
	}
	return nil
}

// GetEndpointAddress gets an endpoint address value
func (c *Configuration) GetEndpointAddress(id ...string) string {
	if ep := c.GetEndpointInfo(id...); ep != nil {
		return ep.Address
	}
	return ""
}

// GetDatabaseInfoGroup gets database infos based on the group id
func (c *Configuration) GetEndpointInfoGroup(groupId string) []EndpointInfo {
	eps := make([]EndpointInfo, 0)
	if c.APIEndpoints == nil {
		return eps
	}
	for _, ep := range *c.APIEndpoints {
		if ep.GroupID == nil {
			continue
		}
		if strings.EqualFold(*ep.GroupID, groupId) {
			eps = append(eps, ep)
		}
	}
	return eps
}

// GetNotificationInfo gets notification info
func (c *Configuration) GetNotificationInfo(id ...string) *NotificationInfo {
	if c.Notifications == nil || (len(id) == 0 && (c.DefaultNotificationID == nil || *c.DefaultNotificationID == "")) {
		return nil
	}

	k := strings.ToLower(*c.DefaultNotificationID)
	if len(id) > 0 {
		k = strings.ToLower(id[0])
	}

	nfs := *c.Notifications
	for _, nf := range nfs {
		if strings.EqualFold(k, nf.ID) {
			return &nf
		}
	}
	return nil
}

// GetSourceInfo gets source by id
func (c *Configuration) GetSourceInfo(sourceId string) *SourceInfo {
	if c.Sources == nil || sourceId == "" {
		return nil
	}
	for _, v := range *c.Sources {
		if strings.EqualFold(v.ID, sourceId) {
			return &v
		}
	}
	return nil
}

// Save saves configuration file
func (c *Configuration) Save() error {
	if c.local {
		return ErrSaveNotLocalFile
	}
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	if err = os.WriteFile(c.FileName, b, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// Load loads configuration file and return a configuration
func Load(source string) (*Configuration, error) {
	return load(source)
}

// Reload configuration
func (c *Configuration) Reload() error {
	_, err := load(c.FileName)
	return err
}

// Flag gets a flag value
func (c *Configuration) Flag(key string) Flag {
	key = strings.TrimSpace(key)
	ret := Flag{
		Key:   key,
		Value: nil,
	}
	if c.Flags == nil {
		return ret
	}
	// get flags to loop from
	// also loop from variations
	// of convention, like underscore
	// and dash
	for _, f := range *c.Flags {
		for _, v := range []string{"_", "-"} {
			ki := strings.ReplaceAll(f.Key, v, "")
			if strings.EqualFold(key, ki) {
				return f
			}
		}
	}

	return ret
}

func new_string(initial string) (init *string) {
	init = new(string)
	*init = initial
	return
}
