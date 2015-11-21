package kmgTime

import (
	"testing"
	"github.com/bronze1man/kmg/kmgTest"
	"time"
)

func TestTimeRecoverInt(ot *testing.T) {
	getTri := func() *TimeRecoverInt {
		return &TimeRecoverInt{
			Num:             1,
			Max:             10,
			LastRecoverTime: MustFromMysqlFormat("2001-01-01 01:01:01"),
			AddDuration:     time.Hour,
		}
	}
	tri := getTri()
	tri.Full(MustFromMysqlFormat("2001-01-01 01:02:01"))
	kmgTest.Equal(tri.Num, 10)

	tri = getTri()
	tri.Sync(MustFromMysqlFormat("2001-01-01 02:01:01"))
	kmgTest.Equal(tri.Num, 2)

	tri = getTri()
	tri.Sync(MustFromMysqlFormat("2001-01-01 13:01:01"))
	kmgTest.Equal(tri.Num, 10)
}
