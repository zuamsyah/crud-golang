package db

import (
	"crud-golang/config"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func New() {
	DBMS := "mysql"
	mySqlConfig := &mysql.Config{
		User:                 config.Get("DB_USER").String(),
		Passwd:               config.Get("DB_PASSWORD").String(),
		Net:                  "tcp",
		Addr:                 config.Get("DB_HOST").String(),
		DBName:               config.Get("DB_NAME").String(),
		AllowNativePasswords: true,
		Params: map[string]string{
			"parseTime": "true",
		},
	}

	var err error
	db, err = gorm.Open(DBMS, mySqlConfig.FormatDSN())
	if err != nil {
		panic("failed to connect database")
	}

	if config.Get("DB_IS_DEBUG").Bool() {
		db = db.Debug()
	}
}

func CloseDB() {
	db.Close()
}

func GetDB() *gorm.DB {
	return db
}
