package main

import (
	"fmt"
	"longest-sum/handlers"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

func main() {

	// get file path from config
	initConfig()
	path := viper.GetString("file.hard")
	if path == "" {
		log.Panic("Cannot read file path from config")
	}

	// sum path
	result, err := handlers.SumMaxPath(path)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Result ->", result)

}

func initConfig() {
	viper.SetConfigName("config") // Name of config file (without extension)
	viper.AddConfigPath(".")      // Path to look for the config file in
	viper.ReadInConfig()          // Read the config file
}
