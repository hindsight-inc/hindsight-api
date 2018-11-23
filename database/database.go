package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"hindsight/config"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	cfg := config.Shared()
	//db, err := gorm.Open("sqlite3", "./../gorm.db")
	db, err := gorm.Open("mysql", cfg.MySQLDatabase + ":" + cfg.MySQLPassword + "@/golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("db err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	DB = db
	return DB
}

func Shared() *gorm.DB {
	return DB
}
