package kmgTime

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestMustFromMysqlFormatDefaultTZ(ot *testing.T) {
	t := MustFromMysqlFormatDefaultTZ("2001-01-01 00:00:00")
	kmgTest.Equal(t.Hour(), 0)
	kmgTest.Equal(t.Day(), 1)

	t = MustFromMysqlFormatDefaultTZ("0000-00-00 00:00:00")
	kmgTest.Equal(t.IsZero(), true)
}

const localTime string = "2015-11-12 14:15:55"


func TestFixLocalTimeToOffsetSpecifiedZoneTime(ot *testing.T){
	otherZoneTime := FixLocalTimeToOffsetSpecifiedZoneTime(3600, localTime)
	kmgTest.Equal(otherZoneTime, "2015-11-12 07:15:55")
	otherZoneTime = FixLocalTimeToOffsetSpecifiedZoneTime(7200, localTime)
	kmgTest.Equal(otherZoneTime, "2015-11-12 08:15:55")
	otherZoneTime = FixLocalTimeToOffsetSpecifiedZoneTime(-18000, localTime)
	kmgTest.Equal(otherZoneTime, "2015-11-12 01:15:55")
}