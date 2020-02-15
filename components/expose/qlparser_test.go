package expose

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestExpose_Parse(t *testing.T) {
	b, _ := ioutil.ReadFile(`C:\Users\mreza\go\src\github.com\iesreza\foundation\components\expose\graphql.test`)
	fmt.Println(ParseQueryLanguage(string(b)))
}
