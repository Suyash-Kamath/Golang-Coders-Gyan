package config

import (
	
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	
)


type HTTPServer struct{
	Addr string `yaml:"address" env-required:"true" env-default:":8082"`
}

type Config struct{
	Env string   `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	 HTTPServer `yaml:"http_server" env-required:"true"`
}

// This is the most essential function , don't return error from this function, if there is an error, log it and exit the program becuase , if this not works then nothing can work
func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == ""{
		flags:= flag.String("config","","path to the configuration file")
		flag.Parse()

		configPath = *flags
	}

	if configPath==" " {
		log.Fatal("config path is required")
	}

	if _,err:= os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file does not exist at path: %s",configPath)
	}

	var cfg Config

	err:=cleanenv.ReadConfig(configPath,&cfg)
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}


	return &cfg

}