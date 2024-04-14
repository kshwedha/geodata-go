package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
}

// ParseConfig parses the configuration file and returns a Config struct.
func ParseConfig(filePath string) (Config, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	// Unmarshal the JSON data into a Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

var config *viper.Viper

// this will read yaml file in common/configdata/dev.yaml
func Init() {
	env := "dev"
	config = viper.New()

	log.Print("This is the environment: ", env)

	// Set the name of the config file (without extension)
	config.SetConfigName(fmt.Sprint(env))

	// Set the path to look for the config file, can set multiple
	config.AddConfigPath(path.Join("src", "common", "configdata"))
	// config.AddConfigPath(".")             // used for docker
	// config.AddConfigPath("../../config/") // used for unit tests

	// Set the config type (YAML in this case)
	config.SetConfigType("yaml")

	// Read the config file
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("!! error on parsing configuration file: %s\n", err)
		return
	}
}

func GetConfig() *viper.Viper {
	return config
}
