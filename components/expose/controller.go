package expose

import (
	"github.com/iesreza/foundation/system"
	"reflect"
	"strings"
)

var db = system.GetDBO()
var objects = map[string]internalQueryItem{}

type internalQueryItem struct {
	name  string
	pkg   string
	path  string
	table string
}

func Bind(name string, object interface{}) {

	ref := reflect.TypeOf(object)
	item := internalQueryItem{}
	item.name = ref.Name()
	item.pkg = strings.Split(ref.String(), ".")[0]
	item.path = string(ref.PkgPath())
	item.table = db.NewScope(object).TableName()
	objects[item.name] = item

}

func Test() {
	type Human struct {
		Name   string  `expose:"selectable filterable"`
		Height int     `expose:"selectable filterable"`
		Age    int     `expose:"filterable"`
		Extra  float64 `expose:"selectable filterable"`
	}
	db.AutoMigrate(&Human{})
	Bind("", &Human{})
}
