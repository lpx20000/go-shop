package models

import (
	"fmt"
	"log"
	"shop/pkg/setting"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Db       *gorm.DB
	HostInfo Info
)

type Model struct {
	CreateTime int `json:"create_time"`
	UpdateTime int `json:"update_time"`
}

type Info struct {
	Host string
}

func init() {
	var (
		err         error
		dbType      string
		dbName      string
		user        string
		password    string
		host        string
		tablePrefix string
		logMode     bool
	)

	sec, err := setting.Cfg.GetSection("databases")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()
	logMode, _ = sec.Key("LOG_MODE").Bool()

	Db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	Db.LogMode(logMode)
	Db.SingularTable(true)
	Db.DB().SetMaxIdleConns(0)
	Db.DB().SetConnMaxLifetime(time.Second)
	Db.DB().SetMaxOpenConns(100)
}

func SetInfo(host string) {
	HostInfo.Host = host
}

func CloseDB() {
	defer Db.Close()
}
