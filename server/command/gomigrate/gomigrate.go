package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	testPtr := flag.Bool("test", false, "using test databse")
	flag.Parse()

	if testPtr != nil {
		if *testPtr {
			config.Setup("../../config.json")
		} else {
			config.Setup("../../config.json.test")
		}
	}

	dbName := config.AppConfig.Database.DBName

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		dbName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(err.Error())
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../database/migration",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if m != nil {
		v1, d1, err := m.Version()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if d1 {
				fmt.Printf("Current version is %d and is dirty\n", v1)
			} else {
				fmt.Printf("Current version is %d\n", v1)
			}
		}

		errUP := m.Up()
		if errUP == nil {
			v, dirty, err := m.Version()
			if err != nil {
				fmt.Print(err.Error())
				fmt.Println(err.Error())
			}
			if dirty {
				fmt.Printf("Current version is %d and is dirty\n", v)
			} else {
				fmt.Printf("Current version is %d\n", v)
			}
		} else {
			fmt.Println(errUP.Error())
		}
	} else {
		fmt.Println("nil migrate.Migrate")
	}
}
