package cfg

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, _ := LoadConfig("config.mssql.json")

	v := config.GetDatabaseInfo("DEFAULT")
	km := v.KeywordMap[0]

	fmt.Println(km.Key)
	fmt.Println(km.Value)
	vi := config.Flag("Joan").String()
	fmt.Println(vi)

	fmt.Println(`Parameter PlaceHolder: `, v.ParameterPlaceHolder)
}
