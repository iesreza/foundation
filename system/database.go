package system

import (
	"fmt"
	"github.com/iesreza/foundation/lib/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strings"
)

var Database *gorm.DB

func SetupDatabase() {
	config := GetConfig().Database
	log.Debug(config)
	var err error
	if config.Enabled == false {
		return
	}

	switch strings.ToLower(config.Type) {
	case "mysql":
		connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?"+config.Params, config.Username, config.Password, config.Server, config.Database)
		Database, err = gorm.Open("mysql", connectionString)
	case "postgres":
		connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s "+config.Params, config.Username, config.Password, config.Server, config.Database, config.SSLMode)
		Database, err = gorm.Open("postgres", connectionString)
	case "mssql":
		connectionString := fmt.Sprintf("user id=%s;password=%s;server=%s;database:%s;"+config.Params, config.Username, config.Password, config.Server, config.Database)
		Database, err = gorm.Open("mssql", connectionString)
	default:
		Database, err = gorm.Open("sqlite3", config.Database+config.Params)
	}
	Database.LogMode(true)
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
					var res []map[string]interface{}
					obj := Database.Raw(cmd)
					obj.Scan(&res)
					err := obj.GetErrors()
					if len(err) > 0 {
						for _, item := range err {
							log.Error(item.Error())
						}
					}
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
								fmt.Printf("%s:%v\t", key, row[key])
							}
						}
					}
					fmt.Printf("\r\nTotal %d records\r\n", len(res))

				} else {
					obj := Database.Exec(cmd)
					err := obj.GetErrors()
					if len(err) > 0 {
						for _, item := range err {
							log.Error(item.Error())
						}
					} else {
						if strings.HasPrefix(cmd, "insert into") {
							fmt.Printf("successfuly inserted")
						} else {

							fmt.Printf("%d rows affected", obj.RowsAffected)
						}
					}
				}

			}

		}
	}, "Run SQL command on database and return result")
}

func GetDBO() *gorm.DB {
	if Database == nil {
		SetupDatabase()
	}
	return Database
}
