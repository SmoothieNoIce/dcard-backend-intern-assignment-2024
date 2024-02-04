package main

import (
	"database/sql"
	"fmt"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.Setup("../../config.json")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(err.Error())
	}

	file, err := os.ReadFile("create_schema.sql")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = db.Exec(string(file))
	if err != nil {
		fmt.Println(err.Error())
	}
}
