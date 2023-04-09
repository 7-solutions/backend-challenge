package main

import (
	"fmt"
	"word-count/routers"

	"github.com/spf13/viper"
)

func main() {

	initConfig()

	e := routers.InitRouter()

	err := e.Start(fmt.Sprintf(":%s", viper.GetString("app.port")))
	if err != nil {
		panic(err)
	}
}

func initConfig() {
	viper.SetConfigName("config") // Name of config file (without extension)
	viper.AddConfigPath(".")      // Path to look for the config file in
	viper.ReadInConfig()          // Read the config file
}
