package models

import (
	"fmt"
	"log"
	"shop/pkg/logging"
	"shop/pkg/setting"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	CreateTime int `json:"-"`
	UpdateTime int `json:"-"`
}

var (
	Db *gorm.DB
	Host string
)

// Setup initializes the database instance
func Setup() {
	var err error
	Db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	Db.SingularTable(true)
	Db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	Db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)
	Db.DB().SetConnMaxLifetime(time.Second)
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer Db.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreateTime"); ok {
			if createTimeField.IsBlank {
				err := createTimeField.Set(nowTime)
				if err != nil {
					logging.LogError(err.Error())
				}
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdateTime"); ok {
			if modifyTimeField.IsBlank {
				err := modifyTimeField.Set(nowTime)
				if err != nil {
					logging.LogError(err.Error())
				}
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_time"); !ok {
		err := scope.SetColumn("UpdateTime", time.Now().Unix())
		if err != nil {
			logging.LogError(err)
		}
	}
}