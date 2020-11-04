package services

import (
	"fmt"
	"github.com/revel/revel"
	"paltronus-backend/app/models"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB instance connects to the Database
var DB *gorm.DB

func getParamString(param string, defaultValue string) string {
	p, found := revel.Config.String(param)
	if !found {
		if defaultValue == "" {
			revel.AppLog.Fatal("Cound not find parameter: " + param)
		} else {
			return defaultValue
		}
	}
	return p
}

func getConnectionString() string {
	host := getParamString("db.host", "")
	port := getParamString("db.port", "3306")
	user := getParamString("db.user", "")
	pass := getParamString("db.password", "")
	dbname := getParamString("db.name", "auction")
	protocol := getParamString("db.protocol", "tcp")
	dbargs := getParamString("dbargs", " ")

	if strings.Trim(dbargs, " ") != "" {
		dbargs = "?" + dbargs
	} else {
		dbargs = ""
	}
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, protocol, host, port, dbname, dbargs)
}

// InitDB initialises DB instance
func InitDB() {
	var err error
	prod := os.Getenv("PROD")
	if prod == "True" {
		println("Using Postgres")
		//dsn := "host=db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
		connectionString := getConnectionString()
		DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	} else {
		println("Using Sqlite")
		DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	}

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(&models.File{}, &models.Version{})
}
