package conf

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env = &env{}

type env struct {
	Host string
	Port string
}

var IsLoaded = false

func LoadEnv() {
	if IsLoaded {
		return
	}

	// load the .env file from /
	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	loadErr := godotenv.Load(curDir + "/.env")
	if loadErr != nil {
		log.Fatalln("can't load env file from current directory: " + curDir)
	}

	// Env
	Env.Host = os.Getenv("HOST")
	Env.Port = os.Getenv("PORT")

	IsLoaded = true
}
