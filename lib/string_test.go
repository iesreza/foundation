package lib

import (
	"fmt"
	"testing"
)

func TestString_String(t *testing.T) {
	s := String("Soft kitty warm kitty little ball of fur. happy kitty sleepy kitty pur pur pur")
	fmt.Println(s.TruncateWord(10,"..."))
	fmt.Println(RandomString(10,ALPHANUM_SIGNS))

	fmt.Println(FormatInt(15000000))
	fmt.Println(FormatInt64(4500000000000))
}
