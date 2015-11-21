package kmgTime

import (
	"testing"
	"time"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestFormat(ot *testing.T) {
//	tc := kmgTest.NewTestTools(ot)
//	t, err := time.Parse(AppleJsonFormat, "2014-04-16 18:26:18 Etc/GMT")
//	tc.Equal(err, nil)
//	tc.Ok(t.Equal(MustFromMysqlFormat("2014-04-16 18:26:18")))
	t, err := time.Parse(AppleJsonFormat, "2014-04-16 18:26:18 Etc/GMT")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(t, MustFromMysqlFormatInLocation("2014-04-16 18:26:18", time.FixedZone("GMT", 0)))

	kmgTest.Equal(DefaultFormat(t), "2014-04-17 02:26:18")
	kmgTest.Equal(MonthAndDayFormat(t), "04-17")
	kmgTest.Equal(NowWithFileNameFormatV2(), time.Now().Format(FormatFileNameV2))
}
