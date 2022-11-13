package db

import (
	"context"
	"fmt"
	"os"

	"github.com/BeanCodeDe/authi/internal/app/authi/config"
	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbpool *pgxpool.Pool

func Init() {

	user := config.PostgresUser
	name := config.PostgresDB
	password := config.PostgresPassword
	host := config.PostgresHost
	port := config.PostgresPort

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)
	var err error
	dbpool, err = pgxpool.Connect(context.Background(), psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func Close() {
	dbpool.Close()
}

func getConnection() *pgxpool.Pool {
	return dbpool
}
