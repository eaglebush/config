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
	config, err := Load("samples/config.mssql.json")
	if err != nil {
		t.Fail()
		t.Fatalf(`Error %v`, err)
	}

	v := config.GetDatabaseInfo("DEFAULT")

	vi := config.Flag("Joan").String()
	fmt.Println(vi)

	fmt.Println(`Parameter PlaceHolder: `, v.ParameterPlaceholder)

	b, _ := json.MarshalIndent(config, "", "\t")

	fmt.Printf("%v+", string(b))

	config.LicenseSerial = new_string("12345678")

	// ok := config.Save()
	// if !ok {
	// 	fmt.Printf("%s", config.LastErrorText())
	// }
}

func TestLoadURLConfig(t *testing.T) {
	config, err := Load("http://valkyrie.vdimdci.com.ph/xtest/config.json")
	if err != nil {
		t.Fail()
		t.Fatalf(`Error %v`, err)
	}

	v := config.GetDatabaseInfo("DEFAULT")

	fmt.Println(`Parameter PlaceHolder: `, v.ParameterPlaceholder)

	b, _ := json.MarshalIndent(config, "", "\t")

	fmt.Printf("%v+", string(b))

	config.LicenseSerial = new_string("12345678")

	// ok := config.Save()
	// if !ok {
	// 	fmt.Printf("%s", config.LastErrorText())
	// }
}
