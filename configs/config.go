package configs

import (
	"api/constant"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	ServerPort string
	BasePath   string
	ApiVersion string
	Timeout    time.Duration
	Client     *http.Client
}

var (
	config *AppConfig
	once   sync.Once
)

// GetConfig returns the singleton AppConfig instance.
func GetConfig() *AppConfig {
	// ensure the thread-safe singleton loading
	once.Do(func() {
		config = loadConfig()
	})
	return config
}

// loadConfig initializes and returns a configured AppConfig instance.
func loadConfig() *AppConfig {
	log.Println("Loading configuration...")

	if os.Getenv(constant.TEST_ENV) == "true" {
		viper.SetConfigFile(constant.ENV_TEST_PATH)
	} else {
		viper.SetConfigFile(constant.ENV_FILE_PATH)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error while reading .env file ", err)
	}

	timeout := viper.GetDuration(constant.TIMEOUT_IN_MS) * time.Second

	return &AppConfig{
		ServerPort: viper.GetString(constant.PORT),
		BasePath:   viper.GetString(constant.BASE_PATH),
		ApiVersion: viper.GetString(constant.API_VERSION),
		Timeout:    timeout,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}
