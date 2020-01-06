package system

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iesreza/foundation/lib/log"
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

	RegisterCLI("sql", &struct{}{}, func(command string, data interface{}) {
		fmt.Println("Type exit to exit")
		for {
			cmd := WaitForConsole(GetConfig().App.Title + " SQL>")
			lower := strings.ToLower(cmd)
			if strings.TrimSpace(lower) == "exit" {
				break
			}
			if strings.TrimSpace(cmd) != "" {

				if strings.HasPrefix(lower, "select") {
					res, err := Database.Query(cmd)
					if err != nil {
						log.Warning(err)
					} else {
						if len(res) > 0 {
							mk := make([]string, len(res[0]))
							c := 0
							for key, _ := range res[0] {
								mk[c] = key
								c++
							}

							for i, row := range res {
								fmt.Printf("\r\n%d- ", i+1)
								for _, key := range mk {
									fmt.Printf("%s:%s\t", key, string(row[key]))
								}
							}
						}
						fmt.Printf("\r\nTotal %d records\r\n", len(res))
					}
				} else {
					res, err := Database.Exec(cmd)
					if err != nil {
						log.Warning(err)
					} else {
						if strings.HasPrefix(cmd, "insert into") {
							insertid, _ := res.LastInsertId()
							fmt.Printf("last inserted id: %d", insertid)
						} else {
							affected, _ := res.RowsAffected()
							fmt.Printf("%d rows affected", affected)
						}
					}
				}
			}

		}
	}, "Run SQL command on database and return result")
}

func GetDBO() *xorm.Engine {
	return Database
}
