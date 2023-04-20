package main

import (
	"strings"

	"github.com/BeanCodeDe/authi/internal/app/authi/api"
	"github.com/BeanCodeDe/authi/internal/app/authi/util"
	"github.com/BeanCodeDe/authi/pkg/parser"
	log "github.com/sirupsen/logrus"
)

const banner = `
 ______           __    __             
/\  _  \         /\ \__/\ \      __    
\ \ \L\ \  __  __\ \ ,_\ \ \___ /\_\   
 \ \  __ \/\ \/\ \\ \ \/\ \  _  \/\ \  
  \ \ \/\ \ \ \_\ \\ \ \_\ \ \ \ \ \ \ 
   \ \_\ \_\ \____/ \ \__\\ \_\ \_\ \_\
	\/_/\/_/\/___/   \/__/ \/_/\/_/\/_/
____________________________________O/_______
                                    O\
`

func main() {
	setLogLevel()
	if log.GetLevel() > log.WarnLevel {
		println(banner)
	}
	log.Debug("Start Server")
	tokenParser, err := parser.NewJWTParser()
	if err != nil {
		log.Fatalf("Error while initializing token parser: %v", err)
	}

	_, err = api.NewUserApi(tokenParser)
	if err != nil {
		log.Fatalf("Error while initializing user api: %v", err)
	}

}

func setLogLevel() {
	logLevel := util.GetEnvWithFallback("LOG_LEVEL", "info")
	switch strings.ToLower(logLevel) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	}
}
