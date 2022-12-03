package config

import "os"

var (
	//Service
	LogLevel = os.Getenv("LOG_LEVEL")

	//Auth
	PrivateKeyPath = os.Getenv("PRIVATE_KEY_PATH")

	//Database

)
