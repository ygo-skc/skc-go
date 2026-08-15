package util

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var EnvMap map[string]string

func ConfigureEnv(envFileVarName string) {
	envFile, isOk := os.LookupEnv(envFileVarName)
	if !isOk {
		slog.Error("Environment variable not found", slog.String("env_var", envFileVarName))
		os.Exit(1)
	}

	slog.Info("Loading env file", slog.String("file", envFile))
	env, err := godotenv.Read(envFile)
	if err != nil {
		slog.Error("Failed to load environment file", slog.String("file", envFile), slog.Any("err", err))
		os.Exit(1)
	}
	EnvMap = env
}
