package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"candy-machine/internal/appconfig"
	"candy-machine/internal/artmachine"

	"github.com/spf13/viper"
)

func main() {

	// init viper
	initViper()

	// read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("configuration file not found")
		} else {
			log.Fatalf("failed to read configuration file, error: %v", err)
		}
	}

	// config file found and successfully parsed
	allSettings := viper.AllSettings()

	// convert map to json string
	jsonSettings, err := json.Marshal(allSettings)
	if err != nil {
		log.Fatalf("failed to marshal configuration, error: %v", err)
	}

	// convert json string to struct
	var config appconfig.Config
	if err := json.Unmarshal(jsonSettings, &config); err != nil {
		log.Fatalf("failed to unmarshal configuration, error: %v", err)
	}

	generateArt := flag.Bool("generate-art", false, "Generate art")
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Parse()

	if *generateArt {
		log.Println("generating art ...")
		log.Println("please wait, this will take some time")
		// init and run art machine
		f, err := artmachine.New(&config)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = f.Run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func initViper() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
}
