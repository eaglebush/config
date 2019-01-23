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

//Configuration - for various configuration settings. This struct can be modified depending on the requirement.
type Configuration struct {
	HostPort         int
	ConnectionString string
	ConvergenceHost  string
	DriverName       string
	HMAC             string
	MailServer       MailServer
	SequenceInfo     SequenceGeneratorInfo
	Domains          []DomainInfo
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
