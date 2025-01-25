package util

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var EnvMap map[string]string

func ConfigureEnv(envFileVarName string) {
	if envFile, isOk := os.LookupEnv(envFileVarName); !isOk {
		log.Fatalf("Could not find environment variable %s in path", envFileVarName)
	} else {
		slog.Info(fmt.Sprintf("Loading env from file %s", envFile))
		if env, err := godotenv.Read(envFile); err != nil {
			log.Fatalln("Could not load environment file (does it exist?). Terminating program.")
		} else {
			EnvMap = env
		}
	}
}
