package db

import (
	"Stardigi-Policy/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	// DBsqler = make(map[string]Sqler)
	DBsqler *utils.SyncMap
)

func init() {
	DBsqler = utils.NewSyncMap()
}

type DBclient struct {
	object interface{}
	Db     *sql.DB
	StrSql string
}

func SetDBsqler(name string, sqler Sqler) {
	DBsqler.Set(name, sqler)
}

func GetDBsqker(name string) Sqler {

	sqler, _ := DBsqler.Get(name)
	return sqler.(Sqler)
}

func NewDBClient(dbType, dataSourceName string) {

	dbClient := new(DBclient)
	dbClient = dbClient.newDBClient(dbType, dataSourceName)

	go func(dbType, dataSourceName string, dbClient *DBclient) {
		for {
			select {
			case <-time.After(1 * time.Second):
				if dbClient.Db.Ping() != nil {
					dbClient = dbClient.newDBClient(dbType, dataSourceName)
				}
			}
		}
	}(dbType, dataSourceName, dbClient)

}

func (dbClient *DBclient) newDBClient(dbType, dataSourceName string) *DBclient {
	db, err := sql.Open(dbType, dataSourceName)
	if err != nil {
		fmt.Println("数据库初始化失败，错误信息：", err)
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)

	dbClient = &DBclient{
		Db: db,
	}
	SetDBsqler(dbType, dbClient)
	return dbClient
}

func (dbClient *DBclient) Insert(object interface{}) {

	name := utils.NameParaseWithReflect(object)
	table := utils.StringParaseWith_(name)
	stat, value := utils.ParaseInterface(object)
	fmt.Println("-----------------value---------------",value)
	con := utils.StringJoin("INSERT INTO ", table, " SET ", stat)
	db := dbClient.Db
	stmt, err := db.Prepare(con)
	if err != nil {
		fmt.Println("插入数据看失败", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(value...)
	if err != nil {
		fmt.Println("插入数据看失败", err)
	}
}

func (dbClient *DBclient) QueryTable(object interface{}) QuerySet {

	// dbClient.StrChan = make(chan string, 1)

	name := utils.NameParaseWithReflect(object)
	table := utils.StringParaseWith_(name)

	str := utils.StringJoin("SELECT * from ", table)

	dbClient.StrSql = str
	dbClient.object = object
	starQuerySet := NewStarQuerySet(dbClient)
	return starQuerySet
}
