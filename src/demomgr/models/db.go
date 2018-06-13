package models

import (
	"os"
	"strings"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //mysqldriver
)

func init() {
	//Register Tables
	orm.RegisterModel("mysql", new(DemoTest))
}
func connect(config map[string]string) {
	//Setup DB connection
	sqlDriver := "mysql"
	if strings.ToLower(os.Getenv("debug")) == "true" {
		orm.Debug = true
	}
	dbURL := config["DatasourceURL"]
	if dbURL == "" {
		panic("dbURL can not be empty")
	}
	// Register DB drivers
	err := orm.RegisterDriver(sqlDriver, orm.DR_MySQL)
	if err != nil {
		panic(err)
	}
	err = orm.RegisterDataBase("default", sqlDriver, dbURL, 30, 30, 300)
	if err != nil {
		panic(err)
	}
}

//Setup database
func Setup(config map[string]string) {
	connect(config)
	name := "default"
	force := false
	verbose := true
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		panic(err)
	}
}
