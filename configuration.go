package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

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
	ID                string //A unique ID that will identify the connection to a database
	ConnectionString  string //ConnectionString specific to the database
	DriverName        string //DriverName needs to be specified depending on the driver id used by the Go database driver
	StorageType       string //StorageType: FILE for filebased database such as Access, SQlite or LocalDB. SERVER for SQL Server, MySQL etc
	GroupID           string `json:"GroupID,omitempty"` //GroupID allows us to get groups of connection
	SequenceGenerator SequenceGeneratorInfo
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
	ApplicationID     string
	CookieDomain      string
	HostPort          int
	HMAC              string
	MailServer        MailServer         `json:"MailServer,omitempty"`
	DefaultDatabaseID string             `json:"DefaultDatabaseID,omitempty"`
	Databases         []DatabaseInfo     `json:"Databases,omitempty"`
	Domains           []DomainInfo       `json:"Domains,omitempty"`
	NotifyRecipients  []NotificationInfo `json:"NotifyRecipients,omitempty"`
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
	/*
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
	*/
	if err != nil {
		return nil, err
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
