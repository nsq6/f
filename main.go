package main

import (
	"flag"
	"fmt"
)

var (
	configFilePath = flag.String("config", "client_secret.live.json", "Path to client configuration file")
)

func main() {
	config, err := createConfig(configFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.ClientID)
}
