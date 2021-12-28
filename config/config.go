package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/strick-j/scimplistic/types"
)

var err error
var config types.ConfigSettings

// ReadConfig will read the configuration json file to read the parameters
// which will be passed in the config file
func ReadConfig(fileName string) (types.ConfigSettings, error) {
	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print("Unable to read config file, switching to flag mode")
		return types.ConfigSettings{}, err
	}
	//log.Print(configFile)
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Print("Invalid JSON")
		return types.ConfigSettings{}, err
	}
	return config, nil
}
