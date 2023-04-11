package conf

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

var Env = &env{}

type env struct {
	AppMode string
	GinMode string
	Host    string
	Port    string
}

var IsLoaded = false

func LoadEnv() {

	once.Do(func() {

		// if not "prod", load the env vars from the .env file
		if Env.AppMode != "prod" {
			curDir, err := os.Getwd()
			if err != nil {
				log.Fatal(err, "conf", "LoadEnv", "error loading os.Getwd()")
			}
			// load the /.env file
			loadErr := godotenv.Load(curDir + "/.env")
			if loadErr != nil {
				log.Fatal(loadErr, "conf", "LoadEnv", "can't load env file from current directory: "+curDir)
			}
			Env.GinMode = "debug"
		} else {
			Env.GinMode = "release"
		}

		// load the env vars
		Env.AppMode = os.Getenv("APP_MODE")
		Env.GinMode = os.Getenv("GIN_MODE")
		Env.Host = os.Getenv("HOST")
		Env.Port = os.Getenv("PORT")
	})
}
