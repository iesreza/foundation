package lib

import (
	"fmt"
	"testing"
	"time"
)

func TestDate_Calculate(t *testing.T) {

	d, _ := Now().DiffExpr("+1 day midnight")
	fmt.Println(fmtDuration(d))
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
