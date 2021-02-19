package cfg

import (
	"encoding/json"
	"fmt"
	"testing"
)

type EmbeddedConfiguration struct {
	Configuration
	ID   string
	Name string
}

func TestLoadConfig(t *testing.T) {
	config, _ := LoadConfig("config.mssql.json")

	v := config.GetDatabaseInfo("DEFAULT")

	var km []DatabaseKeyword
	if v.KeywordMap != nil || len(*v.KeywordMap) > 0 {
		km = *v.KeywordMap
	}

	fmt.Println(km[0].Key)
	fmt.Println(km[0].Value)
	vi := config.Flag("Joan").String()
	fmt.Println(vi)

	fmt.Println(`Parameter PlaceHolder: `, v.ParameterPlaceholder)

	b, _ := json.MarshalIndent(config, "", "\t")

	fmt.Printf("%v+", string(b))

	config.LicenseSerial = "12345678"

	// ok := config.Save()
	// if !ok {
	// 	fmt.Printf("%s", config.LastErrorText())
	// }
}

func TestLoadURLConfig(t *testing.T) {
	config, err := LoadConfig("http://valkyrie.vdimdci.com.ph/xtest/config.json")
	if err != nil {
		t.Fail()
	}

	v := config.GetDatabaseInfo("DEFAULT")

	if v.KeywordMap != nil || len(*v.KeywordMap) > 0 {
		km := *v.KeywordMap

		fmt.Println(km[0].Key)
		fmt.Println(km[0].Value)
		vi := config.Flag("Joan").String()
		fmt.Println(vi)

	}

	fmt.Println(`Parameter PlaceHolder: `, v.ParameterPlaceholder)

	b, _ := json.MarshalIndent(config, "", "\t")

	fmt.Printf("%v+", string(b))

	config.LicenseSerial = "12345678"

	// ok := config.Save()
	// if !ok {
	// 	fmt.Printf("%s", config.LastErrorText())
	// }
}
