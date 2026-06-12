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
	config, err := Load("samples/config.json")
	if err != nil {
		t.Fail()
		t.Fatalf(`Error %v`, err)
	}

	ddf := config.GetDefault("DATABASE")
	v := config.GetDatabaseInfo(ddf)
	if v != nil {
		fmt.Println(`Parameter PlaceHolder: `, v.ParameterPlaceholder)
	}

	vi := config.Flag("Joan").String()
	fmt.Println(vi)

	b, _ := json.MarshalIndent(config, "", "\t")

	fmt.Printf("%v+", string(b))

	dei := config.GetDefault("ENDPOINT")
	apiInfo := config.GetEndpointInfo(dei)
	_ = apiInfo

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

	// ok := config.Save()
	// if !ok {
	// 	fmt.Printf("%s", config.LastErrorText())
	// }
}

func TestGetField(t *testing.T) {
	config, err := Load("samples/config.mssql.json")
	if err != nil {
		t.Fail()
		t.Fatalf(`Error %v`, err)
	}

	f := GetField[CacheInfo](config, "cache")
	_ = f

	// ok := config.Save()
	// if !ok {
	// 	fmt.Printf("%s", config.LastErrorText())
	// }
}
