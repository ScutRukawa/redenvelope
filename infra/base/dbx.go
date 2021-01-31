package base

import (
	"fmt"
	"redenvelope/infra"

	"github.com/tietang/props/kvs"

	"github.com/tietang/dbx"
	//mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var database *dbx.Database

//DbxDatabase get database instance
func DbxDatabase() *dbx.Database {
	return database
}

//DbxDatabaseStarter starter
type DbxDatabaseStarter struct {
	infra.BaseStarter
}

//Init Dbx
func (dbx *DbxDatabaseStarter) Init(ctx infra.StarterContext) {
}

//Setup dbx
func (dbStarter *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {

	conf := ctx.Props()
	setting := dbx.Settings{}
	err := kvs.Unmarshal(conf, &setting, "mysql")
	if err != nil {
		panic(err)
	}
	fmt.Println("mysql:", setting.ShortDataSourceName())
	db, err := dbx.Open(setting)
	if err != nil {
		panic(err)
	}
	fmt.Println("mysqlxxxxxxxxxxxxx:", db)
	database = db
}
