package main

import (
	"log"

	"crud-golang/db"
	"crud-golang/handler"
	"crud-golang/routes"
)

func main() {
	db.New()
	defer db.CloseDB()

	migrate()

	e := routes.Init()

	err := e.Start(":3000")
	logFatal(err)
}

func logFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func migrate() {
	db.GetDB().AutoMigrate(&handler.Student{})
	db.GetDB().AutoMigrate(&handler.Course{})
	db.GetDB().AutoMigrate(&handler.User{})
}
