package util

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var EnvMap map[string]string

func ConfigureEnv(envFileVarName string) {
	if envFile, isOk := os.LookupEnv(envFileVarName); !isOk {
		slog.Error("Environment variable not found", slog.String("env_var", envFileVarName))
		os.Exit(1)
	} else {
		slog.Info("Loading env file", slog.String("file", envFile))
		if env, err := godotenv.Read(envFile); err != nil {
			slog.Error("Failed to load environment file", slog.String("file", envFile), slog.Any("err", err))
			os.Exit(1)
		} else {
			EnvMap = env
		}
	}
}
