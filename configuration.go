package cfg

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type (
	// DirectoryInfo contains a directory info configuration
	DirectoryInfo struct {
		GroupID     string // Group id of the directory
		Description string // Description of the directory
		Items       []Flag // Item flags of this directory
	}

	// EndpointInfo contains an endpoint info configuration
	EndpointInfo struct {
		ID      string        // Endpoint ID for quick access
		Name    string        // Endpoint Name to show
		Address string        // The absolute URL to the resource
		GroupID *string       // A group id to get certain endpoint set
		Token   *string       // A static JWT token for instant access
		APIKey  *string       // An API key for the endpoint
		Secrets *[]SecretInfo // Secrets for any or each part of an API
		Flags   *[]Flag       // Miscellaneous flags inclusive to this endpoint

		cfgAddress string
		cfgToken   *string
		cfgAPIKey  *string
	}

	// OAuthProviderInfo for OAuth configuration
	OAuthProviderInfo struct {
		ID             string // OAuth provider info id for quick access
		Name           string // OAuth name for miscellaneous purposes
		IconUrl        string // OAuth icon image for miscellaneous purposes
		EmbedText      string // OAuth embed options
		Label          string // OAuth label for visual controls
		ClientID       string // Represents the application id registered in an OAuth provider
		ProviderHost   string // The host name of the provider. This is useful to get assets from the provider.
		ProviderWebUri string // The web URI to get authorization and access keys
		ProviderApiUri string // The API URI to get authorization and access keys
		ResponseType   string // The type of response that the application needs from the OAuth provider
		Scope          string // The scope of access to resources

		cfgIconUrl        string
		cfgProviderHost   string
		cfgProviderWebUri string
		cfgProviderApiUri string
	}

	// NotificationInfo - notification information on connecting to Notify API
	NotificationInfo struct {
		ID            string                  // ID of the notification application
		APIHost       string                  // API host of the notification application
		APIPath       string                  // API path of the notification application
		Type          string                  // Notification type (E-mail or messaging)
		Login         string                  // Login credential
		Password      string                  // Password of the login credential
		Active        bool                    // Tags if the notification configuration is active
		SenderAddress string                  // Sender address or id
		SenderName    string                  // Sender name
		ReplyTo       string                  // Reply to address
		Recipients    []NotificationRecipient // Recipients

		cfgAPIHost       string
		cfgLogin         string
		cfgPassword      string
		cfgSenderAddress string
		cfgReplyTo       string
	}

	// CacheInfo connection information
	CacheInfo struct {
		Provider string
		Address  string
		Password string
		DB       int

		cfgAddress  string
		cfgPassword string
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

		cfgConnStr string
	}

	// NotificationRecipient - notification standard recipients
	NotificationRecipient struct {
		ID      string // ID of the recipient
		Name    string // Name of the recipient
		Address string // E-mail address or any identity for notification
	}

	// QueueInfo - queue info connector
	QueueInfo struct {
		ID                 string   // ID of the setting
		ServerAddressGroup []string // Queue server address group
		Cluster            string   // Cluster name
		ClientID           string   // ClientID of the service
		StreamName         string   // Stream name
	}

	// SecretInfo identifies the secrets for application use
	SecretInfo struct {
		GroupID *string // GroupID allows to get secrets in group
		ID      string  // ID of the secret
		Name    string  // Name of the secret
		Value   string  // Value of the secret
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
		// External API endpoints that this application can communicate
		APIEndpoints *[]EndpointInfo

		// ID of this application
		ApplicationID *string

		// Name of this application
		ApplicationName *string

		// Theme of this application
		ApplicationTheme *string

		// Cache info of this application
		Cache *CacheInfo

		// Certificate file
		CertificateFile *string

		// Certificate private key
		CertificateKey *string

		// The domain of the cookie that this application will send
		CookieDomain *string

		// Domains or endpoints that this application will allow
		CrossOriginDomains *[]string

		// Configured databases for this application use
		Databases *[]DatabaseInfo

		// Configured directory for this application use
		Directories *[]DirectoryInfo

		// The default database id that this application will find on the database configuration
		DefaultDatabaseID *string

		// The default endpoint that this application will find on the API endpoints configuration
		DefaultEndpointID *string

		// The default notification id that this application will find on the notification configuration
		DefaultNotificationID *string

		// Configured domains for this application use
		Domains *[]DomainInfo

		// Filename of the current configuration
		FileName string

		// Miscellaneous flags for this application use
		Flags *[]Flag

		// The internal host URL that this application will use to set returned resources and assets
		HostInternalURL *string

		// The external host URL that this application will use to set returned resources and assets
		HostExternalURL *string

		// The network port for the application
		HostPort *int

		// Application wide JSON Web Token (JT) secret
		//
		// Deprecated: Use Secrets instead
		JWTSecret *string

		// License serial of this application
		LicenseSerial *string

		// Configured notifications for this application use
		Notifications *[]NotificationInfo

		// OAuth definitions
		OAuths *[]OAuthProviderInfo

		// Queue or message queue
		Queue *QueueInfo

		// Default network timeout setting for reading data uploaded to this application
		ReadTimeout *int

		// Configured secrets for this application
		Secrets *[]SecretInfo

		// Flags if secure
		Secure *bool

		// Folder sources
		Sources *[]SourceInfo

		// Default network timeout setting for writing data downloaded from this application
		WriteTimeout *int

		// Local file
		local bool
	}
)

var (
	ErrNoDataFromSource = errors.New(`no data from source for configuration`)
	ErrSaveNotLocalFile = errors.New("configuration file is not local")
)

var envPattern = regexp.MustCompile(`\$\{[A-Z0-9_]+\}`)

const def string = `DEFAULT`

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

	if config.DefaultDatabaseID == nil || *config.DefaultDatabaseID == "" {
		config.DefaultDatabaseID = newString(def)
	}
	if config.DefaultEndpointID == nil || *config.DefaultEndpointID == "" {
		config.DefaultEndpointID = newString(def)
	}
	if config.DefaultNotificationID == nil || *config.DefaultNotificationID == "" {
		config.DefaultNotificationID = newString(def)
	}
	if config.CookieDomain == nil {
		config.CookieDomain = newString(`localhost`)
	}
	if config.JWTSecret == nil {
		config.JWTSecret = newString(`defaultsecretkey`)
	}
	// Default setting for database
	if config.Databases != nil {
		dbs := *config.Databases
		for i, cd := range dbs {
			// A database configuration without the connection string is invalid
			if cd.ConnectionString == "" {
				continue
			}
			// Store loaded connection string
			// If there are any environment variables, get the system value and set it
			// Note: If the ConnectionString was modified directly, the stored cfgConnStr
			// will put the value back if Save() is called
			cd.cfgConnStr = cd.ConnectionString
			cd.ConnectionString = interpolateEnvVars(cd.ConnectionString)

			if cd.InterpolateTables == nil {
				cd.InterpolateTables = new(bool)
				*cd.InterpolateTables = true
			}
			if cd.StringEnclosingChar == nil || *cd.StringEnclosingChar == "" {
				cd.StringEnclosingChar = newString(`'`)
			}
			if cd.StringEscapeChar == nil || *cd.StringEscapeChar == "" {
				cd.StringEscapeChar = newString(`\`)
			}
			if cd.ReservedWordEscapeChar == nil || *cd.ReservedWordEscapeChar == "" {
				cd.ReservedWordEscapeChar = newString(`"`)
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

	// Load and parse environment variables set into the values of Address, APIKey and Token
	if config.APIEndpoints != nil {
		aep := *config.APIEndpoints
		for i, ep := range aep {
			if ep.Address == "" {
				continue
			}
			interpolateEndpoint(&aep[i])
		}
		config.APIEndpoints = &aep
	}

	// Load and parse environment variables set into the values
	if config.OAuths != nil {
		oas := *config.OAuths
		for i := range oas {
			interpolateOAuth(&oas[i])
		}
		config.OAuths = &oas
	}

	// Load notifications
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
			interpolateNotifications(&nfs[i])
		}
		config.Notifications = &nfs
	}

	if config.Cache != nil {
		cac := *config.Cache
		cac.cfgAddress = cac.Address
		cac.cfgPassword = cac.Password
		cac.Address = interpolateEnvVars(cac.cfgAddress)
		cac.Password = interpolateEnvVars(cac.cfgPassword)
		config.Cache = &cac
	}

	config.FileName = source
	return config, nil
}

// GetDatabaseInfo get a database info by its ID
func (c *Configuration) GetDatabaseInfo(id string) *DatabaseInfo {
	return findByID(c.Databases, func(v DatabaseInfo) string { return v.ID }, id)
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
	return findByID(c.Directories, func(v DirectoryInfo) string { return v.GroupID }, groupId)
}

// GetDirectoryItem retrieves a directory item under a group
func (c *Configuration) GetDirectoryItem(groupId, key string) *Flag {
	dir := c.GetDirectory(groupId)
	if dir == nil {
		return nil
	}
	for i := range dir.Items {
		if strings.EqualFold(dir.Items[i].Key, key) {
			return &dir.Items[i]
		}
	}
	return nil
}

// GetDomainInfo gets a domain info by name
func (c *Configuration) GetDomainInfo(domainName string) *DomainInfo {
	return findByID(c.Domains, func(v DomainInfo) string { return v.Name }, domainName)
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
	return findByID(c.APIEndpoints, func(v EndpointInfo) string { return v.ID }, k)
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
	for i := range nfs {
		if strings.EqualFold(k, nfs[i].ID) {
			return &nfs[i]
		}
	}
	return nil
}

// GetSourceInfo gets source by id
func (c *Configuration) GetSourceInfo(id string) *SourceInfo {
	return findByID(c.Sources, func(v SourceInfo) string { return v.ID }, id)
}

// GetOAuthInfo gets an OAuth info by id
func (c *Configuration) GetOAuthInfo(id string) *OAuthProviderInfo {
	return findByID(c.OAuths, func(v OAuthProviderInfo) string { return v.ID }, id)
}

// Save saves configuration file
func (c *Configuration) Save() error {
	if !c.local {
		return ErrSaveNotLocalFile
	}
	// Put back the raw connection string to save loaded
	if c.Databases != nil {
		for i, db := range *c.Databases {
			if db.ConnectionString == "" {
				continue
			}
			(*c.Databases)[i].ConnectionString = db.cfgConnStr
		}
	}
	// Put back endpoint info
	if c.APIEndpoints != nil {
		aep := *c.APIEndpoints
		for i, ep := range aep {
			if ep.Address == "" {
				continue
			}
			// Store loaded values to retrieve later
			ep.Address = ep.cfgAddress
			ep.APIKey = ep.cfgAPIKey
			ep.Token = ep.cfgToken
			aep[i] = ep
		}
		c.APIEndpoints = &aep
	}
	if c.OAuths != nil {
		oas := *c.OAuths
		for i, oa := range oas {
			oa.IconUrl = oa.cfgIconUrl
			oa.ProviderHost = oa.cfgProviderHost
			oa.ProviderWebUri = oa.cfgProviderWebUri
			oa.ProviderApiUri = oa.cfgProviderApiUri
			oas[i] = oa
		}
		c.OAuths = &oas
	}
	if c.Notifications != nil {
		nfs := *c.Notifications
		for i, cn := range nfs {
			cn.APIHost = cn.cfgAPIHost
			cn.Login = cn.cfgLogin
			cn.Password = cn.cfgPassword
			cn.SenderAddress = cn.cfgSenderAddress
			cn.ReplyTo = cn.cfgReplyTo
			nfs[i] = cn
		}
		c.Notifications = &nfs
	}
	if c.Cache != nil {
		cac := *c.Cache
		cac.Address = cac.cfgAddress
		cac.Password = cac.cfgPassword
		c.Cache = &cac
	}

	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	if err = os.WriteFile(c.FileName, b, os.ModePerm); err != nil {
		return err
	}
	// Re-interpret the environment variables back
	if c.Databases != nil {
		for i, db := range *c.Databases {
			if db.ConnectionString == "" {
				continue
			}
			(*c.Databases)[i].ConnectionString = interpolateEnvVars(db.cfgConnStr)
		}
	}
	if c.APIEndpoints != nil {
		aep := *c.APIEndpoints
		for i, ep := range aep {
			if ep.Address == "" {
				continue
			}
			interpolateEndpoint(&aep[i])
		}
		c.APIEndpoints = &aep
	}

	if c.OAuths != nil {
		oas := *c.OAuths
		for i := range oas {
			interpolateOAuth(&oas[i])
		}
		c.OAuths = &oas
	}

	if c.Notifications != nil {
		defnum := ""
		nfs := *c.Notifications
		for i, cn := range nfs {
			if i > 0 {
				defnum = strconv.Itoa(i)
			}
			if cn.ID == "" {
				nfs[i].ID = def + defnum
			}
			interpolateNotifications(&nfs[i])
		}

		c.Notifications = &nfs
	}

	if c.Cache != nil {
		cac := *c.Cache
		cac.cfgAddress = cac.Address
		cac.cfgPassword = cac.Password
		cac.Address = interpolateEnvVars(cac.cfgAddress)
		cac.Password = interpolateEnvVars(cac.cfgPassword)
		c.Cache = &cac
	}

	return nil
}

// Load loads configuration file and return a configuration
func Load(source string) (*Configuration, error) {
	return load(source)
}

// Reload configuration
func (c *Configuration) Reload() error {
	newConfig, err := load(c.FileName)
	if err != nil {
		return err
	}

	// Copy each field explicitly to avoid overwriting the struct
	// and breaking references (e.g. pointers held elsewhere).
	c.APIEndpoints = newConfig.APIEndpoints
	c.ApplicationID = newConfig.ApplicationID
	c.ApplicationName = newConfig.ApplicationName
	c.Cache = newConfig.Cache
	c.CookieDomain = newConfig.CookieDomain
	c.Databases = newConfig.Databases
	c.Directories = newConfig.Directories
	c.Flags = newConfig.Flags
	c.JWTSecret = newConfig.JWTSecret
	c.LicenseSerial = newConfig.LicenseSerial
	c.Notifications = newConfig.Notifications
	c.OAuths = newConfig.OAuths
	c.Queue = newConfig.Queue
	c.ReadTimeout = newConfig.ReadTimeout
	c.Secure = newConfig.Secure
	c.Sources = newConfig.Sources
	c.Secrets = newConfig.Secrets
	c.WriteTimeout = newConfig.WriteTimeout

	return nil
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

// GetFlag retrieves a flag value and return it converted to type indicated
func GetFlag[T FlagTypes](flgs *[]Flag, key string) T {
	var zero T
	if key == "" || flgs == nil || len(*flgs) == 0 {
		return zero
	}
	for _, flg := range *flgs {
		if !strings.EqualFold(flg.Key, key) {
			continue
		}
		if flg.Value == nil {
			return zero
		}

		switch any(*new(T)).(type) {
		case string:
			return any(*flg.Value).(T)
		case int:
			v, _ := strconv.Atoi(*flg.Value)
			return any(v).(T)
		case int32:
			v, _ := strconv.ParseInt(*flg.Value, 10, 32)
			return any(v).(T)
		case int64:
			v, _ := strconv.ParseInt(*flg.Value, 10, 64)
			return any(v).(T)
		case bool:
			v, _ := strconv.ParseBool(*flg.Value)
			return any(v).(T)
		case float32:
			v, _ := strconv.ParseFloat(*flg.Value, 32)
			return any(v).(T)
		case float64:
			v, _ := strconv.ParseFloat(*flg.Value, 64)
			return any(v).(T)
		default:
			return zero
		}

	}
	return zero
}

// normalizeFieldName normalizes a field name for lookup:
// - lowercases
// - removes spaces, underscores, and dashes
func normalizeFieldName(name string) string {
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, "-", "")
	return name
}

// GetField safely retrieve a field in the configuration value and return it converted to type indicated.
func GetField[T any](config *Configuration, fieldName string) T {
	var zv T
	if config == nil || fieldName == "" {
		return zv
	}

	var v any
	switch normalizeFieldName(fieldName) {
	case "apiendpoints":
		v = config.APIEndpoints
	case "applicationid":
		v = config.ApplicationID
	case "applicationname":
		v = config.ApplicationName
	case "applicationtheme":
		v = config.ApplicationTheme
	case "cache":
		v = config.Cache
	case "certificatefile":
		v = config.CertificateFile
	case "certificatekey":
		v = config.CertificateKey
	case "cookiedomain":
		v = config.CookieDomain
	case "crossorigindomains":
		v = config.CrossOriginDomains
	case "databases":
		v = config.Databases
	case "directories":
		v = config.Directories
	case "defaultdatabaseid":
		v = config.DefaultDatabaseID
	case "defaultendpointid":
		v = config.DefaultEndpointID
	case "defaultnotificationid":
		v = config.DefaultNotificationID
	case "domains":
		v = config.Domains
	case "flags":
		v = config.Flags
	case "hostinternalurl":
		v = config.HostInternalURL
	case "hostexternalurl":
		v = config.HostExternalURL
	case "hostport":
		v = config.HostPort
	case "licenseserial":
		v = config.LicenseSerial
	case "notifications":
		v = config.Notifications
	case "oauths":
		v = config.OAuths
	case "queue":
		v = config.Queue
	case "readtimeout":
		v = config.ReadTimeout
	case "secrets":
		v = config.Secrets
	case "secure":
		v = config.Secure
	case "sources":
		v = config.Sources
	case "writetimeout":
		v = config.WriteTimeout
	default:
		return zv
	}

	// Now enforce that v matches T
	if cast, ok := v.(T); ok {
		return cast
	}

	// Wrong type requested → return zero value
	return zv
}

// GetSecretInfo get a secret info by its ID
func (c *Configuration) GetSecretInfo(id string) *SecretInfo {
	return findByID(c.Secrets, func(v SecretInfo) string { return v.ID }, id)
}

// GetSecretInfoGroup gets secret infos based on the group id
func (c *Configuration) GetSecretInfoGroup(groupId string) []SecretInfo {
	scts := make([]SecretInfo, 0)
	if c.Secrets == nil || groupId == "" {
		return scts
	}
	for _, v := range *c.Secrets {
		if v.GroupID == nil {
			continue
		}
		if strings.EqualFold(*v.GroupID, groupId) {
			scts = append(scts, v)
		}
	}
	return scts
}

// GetSecretInfo get a secret info by its ID
func (c *EndpointInfo) GetSecretInfo(id string) *SecretInfo {
	return findByID(c.Secrets, func(v SecretInfo) string { return v.ID }, id)
}

// GetSecretInfoGroup gets secret infos based on the group id
func (c *EndpointInfo) GetSecretInfoGroup(groupId string) []SecretInfo {
	scts := make([]SecretInfo, 0)
	if c.Secrets == nil || groupId == "" {
		return scts
	}
	for _, v := range *c.Secrets {
		if v.GroupID == nil {
			continue
		}
		if strings.EqualFold(*v.GroupID, groupId) {
			scts = append(scts, v)
		}
	}
	return scts
}

func newString(initial string) (init *string) {
	init = new(string)
	*init = initial
	return
}

func interpolateEnvVars(text string) string {
	if text == "" {
		return text
	}

	// Replace all matches with actual env values
	result := envPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the env var name from ${ENV_NAME}
		varName := match[2 : len(match)-1]
		value := os.Getenv(varName)
		return value // If not set, will return empty string
	})

	return result
}

func interpolateEndpoint(ep *EndpointInfo) {
	ep.cfgAddress = ep.Address
	ep.cfgAPIKey = ep.APIKey
	ep.cfgToken = ep.Token
	ep.Address = interpolateEnvVars(ep.cfgAddress)
	if ep.cfgAPIKey != nil && ep.APIKey != nil {
		*ep.APIKey = interpolateEnvVars(*ep.cfgAPIKey)
	}
	if ep.cfgToken != nil && ep.Token != nil {
		*ep.Token = interpolateEnvVars(*ep.cfgToken)
	}
}

func interpolateOAuth(oa *OAuthProviderInfo) {
	oa.cfgIconUrl = oa.IconUrl
	oa.cfgProviderHost = oa.ProviderHost
	oa.cfgProviderWebUri = oa.ProviderWebUri
	oa.cfgProviderApiUri = oa.ProviderApiUri

	oa.IconUrl = interpolateEnvVars(oa.cfgIconUrl)
	oa.ProviderHost = interpolateEnvVars(oa.cfgProviderHost)
	oa.ProviderWebUri = interpolateEnvVars(oa.cfgProviderWebUri)
	oa.ProviderApiUri = interpolateEnvVars(oa.cfgProviderApiUri)
}

func interpolateNotifications(cn *NotificationInfo) {
	cn.cfgAPIHost = cn.APIHost
	cn.cfgLogin = cn.Login
	cn.cfgPassword = cn.Password
	cn.cfgSenderAddress = cn.SenderAddress
	cn.cfgReplyTo = cn.ReplyTo

	cn.APIHost = interpolateEnvVars(cn.cfgAPIHost)
	cn.Login = interpolateEnvVars(cn.cfgLogin)
	cn.Password = interpolateEnvVars(cn.cfgPassword)
	cn.SenderAddress = interpolateEnvVars(cn.cfgSenderAddress)
	cn.ReplyTo = interpolateEnvVars(cn.cfgReplyTo)
}

func findByID[T any](slice *[]T, getID func(T) string, id string) *T {
	if slice == nil || id == "" {
		return nil
	}
	for i := range *slice {
		if strings.EqualFold(getID((*slice)[i]), id) {
			return &(*slice)[i]
		}
	}
	return nil
}
