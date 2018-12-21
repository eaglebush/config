package configuration

import (
	"encoding/json"
	"os"
)

//Configuration - for various configuration settings
type Configuration struct {
	HostPort         int
	ConnectionString string
}

// LoadConfig - load configuration file and return a configuration
func LoadConfig(fileName string) (*Configuration, error) {
	config := &Configuration{}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
