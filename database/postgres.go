package database

import (
	"fmt"
	"github.com/kaleido-io/kaleido-fabric-go/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"runtime"
	"time"
)

var (
	clientpostgres *gorm.DB
)

func PostgresConnect() {

	var host string

	user := utils.GetEnvVariables("DB_USER")
	pass := utils.GetEnvVariables("DB_PASSWORD")
	dbname := utils.GetEnvVariables("DB_NAME")
	port := utils.GetEnvVariables("DB_PORT")

	os := runtime.GOOS
	switch os {
	case "darwin":
		host = utils.GetEnvVariables("DB_HOST")
	default:
		host = utils.GetEnvVariables("DB_HOST")
	}

	Conn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", host, port, user, dbname, pass)
	db, err := gorm.Open(postgres.Open(Conn), &gorm.Config{SkipDefaultTransaction: true, DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic("Failed to connect to database!")
	} else {
		log.Println("Success connect to postgres")
	}

	db.Use(
		dbresolver.Register(dbresolver.Config{ /* xxx */ }).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)

	//db.AutoMigrate(&domain.Customer{})
	//db.AutoMigrate(&domain.Auth{})

	utils.GenerateKeyFiles("jwtRS256")
	SetUpPostgresdb(db)

}

func SetUpPostgresdb(DB *gorm.DB) {
	clientpostgres = DB
}

func GetPostgresdb() *gorm.DB {
	return clientpostgres
}
