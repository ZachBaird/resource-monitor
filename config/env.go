package config

import (
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
