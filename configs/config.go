package configs

import (
	"api/constant"
	"log"

	"github.com/spf13/viper"
)
type AppConfig struct {
	ServerPort string
}

func GetAppConfig() AppConfig {
	log.Println("App configuration executed...")
	viper.SetConfigFile(constant.ENV_FILE_PATH)
	err := viper.ReadInConfig()
    if err != nil {
        log.Fatal("Error while reading .env file ", err)
    }
	appconfig := AppConfig {
		ServerPort: viper.GetString(constant.PORT),
	}
	log.Println("Server Port : ",appconfig.ServerPort)
	return appconfig
}

