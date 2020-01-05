package system

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iesreza/foundation/log"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"xorm.io/xorm"
)

var Database *xorm.Engine

func SetupDatabase() {
	config := GetConfig().Database
	var err error

	switch strings.ToLower(config.Type) {
	case "mysql":
		connectionString := fmt.Sprintf("%s:%s@%s/%s", config.Username, config.Password, config.Server, config.Database)
		Database, err = xorm.NewEngine("mysql", connectionString)
	case "postgres":
		connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s", config.Username, config.Password, config.Server, config.Database, config.SSLMode)
		Database, err = xorm.NewEngine("postgres", connectionString)
	case "mssql":
		connectionString := fmt.Sprintf("user id=%s;password=%s;server=%s;database:%s;", config.Username, config.Password, config.Server, config.Database)
		Database, err = xorm.NewEngine("mssql", connectionString)
	default:
		Database, err = xorm.NewEngine("sqlite3", config.Database)
	}
	if err != nil {
		log.Critical(err)
	}

}

func GetDBO() *xorm.Engine {
	return Database
}
