package utils

import (
	"fmt"
	"reflect"

	"github.com/tkanos/gonfig"
	"github.com/wyrdnixx/KBZOZeit/models"
)

func GetConfig(params ...string) models.Configuration {
	configuration := models.Configuration{}
	env := "dev"
	if len(params) > 0 {
		env = params[0]
	}
	fileName := fmt.Sprintf("./%s_config.json", env)
	Log(1, "GetConfig()", "using configfile: "+fileName)

	gonfig.GetConf(fileName, &configuration)

	v := reflect.ValueOf(configuration)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		cfg := fmt.Sprintf("Field: %s\tValue: %v", typeOfS.Field(i).Name, v.Field(i).Interface())
		//fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		Log(1, "GetConfig()", "config: "+cfg)
	}

	return configuration
}
