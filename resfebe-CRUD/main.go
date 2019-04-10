package main

import (
	"fmt"

	"github.com/spf13/viper"
)

var err error

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config/")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	a := App{}
	a.Initialize()

	a.Run(viper.GetString("Server.port"))
}
