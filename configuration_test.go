package cfg

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, _ := LoadConfig("config.mssql.json")

	v := config.GetDatabaseInfo("DEFAULT").KeywordMap[0]
	fmt.Println(v.Key)
	fmt.Println(v.Value)
	vi := config.Flag("Joan").String()
	fmt.Println(vi)
}
