package kmgTime

import (
	"testing"
//	"time"
//
//	"github.com/bronze1man/kmg/kmgTest"
	"time"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestGetPeriodFromSortedList(ot *testing.T) {
	SortedList := []Period{
		{
			Start: MustFromMysqlFormat("2001-01-01 00:00:00"),
			End:   MustFromMysqlFormat("2001-01-01 01:00:00"),
		},
		{
			Start: MustFromMysqlFormat("2001-01-01 02:00:00"),
			End:   MustFromMysqlFormat("2001-01-01 03:00:00"),
		},
		{
			Start: MustFromMysqlFormat("2001-01-01 03:00:00"),
			End:   MustFromMysqlFormat("2001-01-01 04:00:00"),
		},
	}
	for _, testcase := range []struct {
		t  time.Time
		i  int
		ok bool
	}{
		{MustFromMysqlFormat("2001-01-00 23:30:00"), 0, false},
		{MustFromMysqlFormat("2001-01-01 00:30:00"), 0, true},
		{MustFromMysqlFormat("2001-01-01 03:00:00"), 2, true},
		{MustFromMysqlFormat("2001-01-01 04:30:00"), 0, false},
	} {
		i, ok := GetPeriodFromSortedList(testcase.t, SortedList)
		kmgTest.Equal(i, testcase.i)
		kmgTest.Equal(ok, testcase.ok)
	}
}
