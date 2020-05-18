package envutil

import (
	"fmt"
	"os"
	"strconv"
)

const (
	EnvDbUrl          = "DB_URL"
	EnvDbHost         = "DB_HOST"
	EnvDbPort         = "DB_PORT"
	EnvDbUser         = "DB_USER"
	EnvDbPassword     = "DB_PASSWORD"
	EnvDbName         = "DB_NAME"
	EnvGitlabUrl      = "GITLAB_URL"
	EnvGitlabToken    = "GITLAB_TOKEN"
	EnvUpdateInterval = "UPDATE_INTERVAL"
	EnvLogLevel       = "LOG_LEVEL"
	EnvLogMode        = "LOG_MODE"
)

func MustGetEnvStr(env string) string {
	if val := GetEnvStr(env); len(val) > 0 {
		return val
	}
	panic(fmt.Errorf("env %s is not set", env))
}

func GetEnvStr(env string) string {
	return os.Getenv(env)
}

func GetEnvUintOrDefault(env string, defaultVal uint64) uint64 {
	valStr := os.Getenv(env)
	if val, err := strconv.ParseUint(valStr, 10, 64); err == nil {
		return val
	}
	return defaultVal
}
