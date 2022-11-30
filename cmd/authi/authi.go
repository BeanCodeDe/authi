package main

import (
	"github.com/BeanCodeDe/authi/internal/app/authi/api"
	"github.com/BeanCodeDe/authi/internal/app/authi/config"
	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/BeanCodeDe/authi/pkg/parser"
	log "github.com/sirupsen/logrus"
)

func main() {
	setLogLevel(config.LogLevel)
	log.Info("Start Server")

	authAdapter := adapter.NewAuthiAdapter()
	tokenParser, err := parser.NewJWTParser()
	if err != nil {
		log.Fatalf("Error while initializing token parser: %v", err)
	}

	_, err = api.NewUserApi(authAdapter, tokenParser)
	if err != nil {
		log.Fatalf("Error while initializing user api: %v", err)
	}

}

func setLogLevel(logLevel string) {
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	}
}
