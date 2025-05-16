package configs

import (
	"api/constant"
	"log"

	"github.com/spf13/viper"
)
type AppConfig struct {
	ServerPort string
	BasePath string
	ApiVersion string
}

var config *AppConfig

func GetConfig() *AppConfig {
	if config == nil {
		config = &AppConfig{}
	}
	return config
}

func (a *AppConfig) LoadConfig() {
	log.Println("Load configuration executed...")
	viper.SetConfigFile(constant.ENV_FILE_PATH)
	err := viper.ReadInConfig()
    if err != nil {
        log.Fatal("Error while reading .env file ", err)
    }
	a.ServerPort = viper.GetString(constant.PORT)
	a.BasePath = viper.GetString(constant.BASE_PATH)
	a.ApiVersion = viper.GetString(constant.API_VERSION)
}

