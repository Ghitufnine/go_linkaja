package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	DBConn *gorm.DB
)

func ConnectDb() {
	// Update these values with your local MySQL credentials
	username := "root"
	password := ""
	host := "localhost"
	port := "3306"
	dbName := "go_linkaja"

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC", username, password, host, port, dbName)

	db, err := gorm.Open("mysql", dbConnectionString)
	if err != nil {
		panic(err)
	}

	log.Println("Connected to local MySQL")

	DBConn = db
}
