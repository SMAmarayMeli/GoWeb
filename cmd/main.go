package main

import (
	"GoWeb/cmd/routes"
	"GoWeb/internal/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
)

func readJson(dbP *[]domain.Producto) {
	data, err := ioutil.ReadFile("../GoWeb/internal/product/products.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, dbP)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var StorageDb *sql.DB

func initDB() {
	databaseConfig := mysql.Config{
		User:   "root",
		Passwd: "abc",
		Addr:   "127.0.0.1:3306",
		DBName: "my_db",
	}

	StorageDB, err := sql.Open("mysql", databaseConfig.FormatDSN())
	if err != nil {
		panic(err)
	}
	err = StorageDB.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("database Configured")
}

func main() {
	err := godotenv.Load("go.env")
	if err != nil {
		panic("env not loadable")
	}

	initDB()
	var dbP []domain.Producto
	readJson(&dbP)

	en := gin.Default()
	rt := routes.NewRouter(StorageDb, en, &dbP)
	rt.SetRoutes()

	if err := en.Run(); err != nil {
		log.Fatal(err)
	}

}
