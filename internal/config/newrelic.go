package config

import (
	"log"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var NRApp *newrelic.Application

func InitNewRelic() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("ride-hailing-app"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		log.Println("New Relic init failed:", err)
		return
	}

	NRApp = app
	log.Println("New Relic initialized ✅")
}