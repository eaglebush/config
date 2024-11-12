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

type (
	// DatabaseKeyword for database keywords
	DatabaseKeyword struct {
		Flag
	}

	// APIKeyInfo contains an API info configuration
	APIKeyInfo struct {
		ID    string
		Name  string
		Key   string
		Token *string
	}

	// DirectoryInfo contains a directory info configuration
	DirectoryInfo struct {
		GroupID     string
		Description string
		Items       []Flag
	}

	// Endpoint contains an endpoint info configuration
	EndpointInfo struct {
		ID      string  // Endpoint ID for quick access
		Name    string  // Endpoint Name to show
		Address string  // The absolute URL to the resource
		GroupID *string // A group id to get certain endpoint set
		Token   *string
	}

	// OAuthProviderInfo for OAuth configuration
	OAuthProviderInfo struct {
		ID             string // OAuth provider info id for quick access
		ClientID       string // Represents the application id registered in an OAuth provider
		ProviderWebUri string // The web URI to get authorization and access keys
		ProviderApiUri string // The API URI to get authorization and access keys
		ResponseType   string // The type of response that the application needs from the OAuth provider
		Scope          string // The scope of access to resources
	}

	// NotificationInfo - notification information on connecting to Notify API
	NotificationInfo struct {
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
	CacheInfo struct {
		Provider string
		Address  string
		Password string
		DB       int
	}

	// DomainInfo - domain information for LDAP authentication
	DomainInfo struct {
		Name               string
		Host               string
		Port               uint16
		Path               string
		AuthorizedUser     string
		AuthorizedPassword string
		Filter             string
	}

	// SequenceGeneratorInfo - sequence generator query
	SequenceGeneratorInfo struct {
		UpsertQuery     string
		ResultQuery     string
		NamePlaceHolder string
	}

	// DatabaseInfo - database configuration setting
	DatabaseInfo struct {
		GroupID                *string                // GroupID allows us to get groups of connection
		ID                     string                 // A unique ID that will identify the connection to a database
		ConnectionString       string                 // ConnectionString specific to the database
		DriverName             string                 // DriverName needs to be specified depending on the driver id used by the Go database driver
		StorageType            string                 // FILE for filebased database such as Access, SQlite or LocalDB. SERVER for SQL Server, MySQL etc
		HelperID               string                 // When using github.com/NarsilWorks-Inc/datahelperlite, this is needed in the configuration file
		ParameterPlaceholder   string                 // Parameter place holder for prepared statements. Default is '?'
		ParameterInSequence    bool                   // Parameter place holder is in sequence. Default is false
		Schema                 string                 // Schema for any of the database operations
		InterpolateTables      *bool                  // Enables the tables to be interpolated with schema
		SequenceGenerator      *SequenceGeneratorInfo // Sequence generator configuration
		StringEnclosingChar    *string                // Gets or sets the character that encloses a string in the query
		StringEscapeChar       *string                // Gets or Sets the character that escapes a reserved character such as the character that encloses a s string
		MaxOpenConnection      *int                   // Maximum open connection
		MaxIdleConnection      *int                   // Maximum idle connection
		MaxConnectionLifetime  *int                   // Max connection lifetime
		MaxConnectionIdleTime  *int                   // Max idle connection lifetime
		Ping                   *bool                  // Ping connection
		ReservedWordEscapeChar *string                // Reserved word escape chars. For escaping with different opening and closing characters, just set to both. Example. `[]` for SQL server
	}

	// NotificationRecipient - notification standard recipients
	NotificationRecipient struct {
		ID      string
		Name    string
		Address string
	}

	// QueueInfo - queue info connector
	QueueInfo struct {
		ID                 string   // ID of the setting
		ServerAddressGroup []string // Queue server address group
		Cluster            string   // Cluster name
		ClientID           string   // ClientID of the service
		StreamName         string   // Stream name
	}

	// SourceInfo - file sources for configuration
	SourceInfo struct {
		ID        string // ID of the source for quick reference
		Type      string // Type of Inbound file. Supported types are ORDER and SNAPSHOT
		Source    string // Source folder of the source
		Relative  bool   // Indicates that the Error and Success folders are relative to Source
		Error     string // Error folder of the source
		Success   string // Success folder of the source
		Extension string // Extension of the file to pickup
	}

	// Configuration
	Configuration struct {
		APIEndpoints          *[]EndpointInfo      // External API endpoints that this application can communicate
		APIKeys               *[]APIKeyInfo        // API Keys
		ApplicationID         *string              // ID of this application
		ApplicationName       *string              // Name of this application
		ApplicationTheme      *string              // Theme of this application
		Cache                 *CacheInfo           // Cache info of this application
		CertificateFile       *string              // Certificate file
		CertificateKey        *string              // Certificate private key
		CookieDomain          *string              // The domain of the cookie that this application will send
		CrossOriginDomains    *[]string            // Domains or endpoints that this application will allow
		Databases             *[]DatabaseInfo      // Configured databases for this application use
		Directories           *[]DirectoryInfo     // Configured directory for this application use
		DefaultDatabaseID     *string              // The default database id that this application will find on the database configuration
		DefaultEndpointID     *string              // The default endpoint that this application will find on the API endpoints configuration
		DefaultNotificationID *string              // The default notification id that this application will find on the notification configuration
		Domains               *[]DomainInfo        // Configured domains for this application use
		FileName              string               // Filename of the current configuration
		Flags                 *[]Flag              // Miscellaneous flags for this application use
		HostInternalURL       *string              // The internal host URL that this application will use to set returned resources and assets
		HostExternalURL       *string              // The external host URL that this application will use to set returned resources and assets
		HostPort              *int                 // The network port for the application
		JWTSecret             *string              // Application wide JSON Web Token (JT) secret
		LicenseSerial         *string              // License serial of this application
		Notifications         *[]NotificationInfo  // Configured notifications for this application use
		OAuths                *[]OAuthProviderInfo // OAuth definitions
		Queue                 *QueueInfo           // Queue or message queue
		ReadTimeout           *int                 // Default network timeout setting for reading data uploaded to this application
		Secure                *bool                // Flags if secure
		Sources               *[]SourceInfo        // Folder sources
		WriteTimeout          *int                 // Default network timeout setting for writing data downloaded from this application
		local                 bool                 // Local file
	}
)

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
func (c *Configuration) GetDatabaseInfo(id string) *DatabaseInfo {
	if c.Databases == nil {
		return nil
	}
	for _, v := range *c.Databases {
		if v.ID == id {
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
func (c *Configuration) GetEndpointInfo(id string) *EndpointInfo {
	if c.APIEndpoints == nil || (len(id) == 0 && (c.DefaultEndpointID == nil || *c.DefaultEndpointID == "")) {
		return nil
	}
	k := strings.ToLower(*c.DefaultEndpointID)
	if len(id) > 0 {
		k = strings.ToLower(id)
	}
	eps := *c.APIEndpoints
	for _, ep := range eps {
		if strings.EqualFold(k, ep.ID) {
			return &ep
		}
	}
	return nil
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
func (c *Configuration) GetNotificationInfo(id string) *NotificationInfo {
	if c.Notifications == nil || (len(id) == 0 && (c.DefaultNotificationID == nil || *c.DefaultNotificationID == "")) {
		return nil
	}
	k := strings.ToLower(*c.DefaultNotificationID)
	if len(id) > 0 {
		k = strings.ToLower(id)
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
func (c *Configuration) GetSourceInfo(id string) *SourceInfo {
	if c.Sources == nil || id == "" {
		return nil
	}
	for _, v := range *c.Sources {
		if strings.EqualFold(v.ID, id) {
			return &v
		}
	}
	return nil
}

// GetOAuthInfo gets an OAuth info by id
func (c *Configuration) GetOAuthInfo(id string) *OAuthProviderInfo {
	if c.OAuths == nil || len(*c.OAuths) == 0 || len(id) == 0 {
		return nil
	}
	for _, oa := range *c.OAuths {
		if strings.EqualFold(id, oa.ID) {
			return &oa
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
