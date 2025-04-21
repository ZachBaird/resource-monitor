package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func GetDryRunConfig() bool {
	_ = godotenv.Load()
	s := os.Getenv("RESOURCE_MONITOR_DRY_RUN")

	dryRun, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal(err)
	}
	return dryRun
}

func GetSecretConfig(a string) string {
	_ = godotenv.Load()
	env := os.Getenv(fmt.Sprintf("%s_SECRET", a))
	return env
}
