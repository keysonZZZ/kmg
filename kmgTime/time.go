package kmgTime

import (
	"fmt"
	"time"
	"math"
)

const (
	Day = 24 * time.Hour
)

var DefaultTimeZone = BeijingZone

//auto guess format from ParseFormatGuessList
func ParseAutoInLocation(sTime string, loc *time.Location) (t time.Time, err error) {
	for _, format := range ParseFormatGuessList {
		t, err = time.ParseInLocation(format, sTime, loc)
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("[ParseAutoInLocation] time: %s can not parse", sTime)
	return
}

func FixLocalTimeToOffsetSpecifiedZoneTime(timeOffset int, localTime string) string {
	clientTimeZone := time.FixedZone("ClientTimeZone", timeOffset)
	serverTime := MustParseAutoInDefault(localTime)
	return serverTime.In(clientTimeZone).Format(FormatMysql)
}

func MustParseAutoInDefault(sTime string) (t time.Time) {
	t, err := ParseAutoInLocation(sTime, DefaultTimeZone)
	if err != nil {
		panic(err)
	}
	return t
}

func ParseAutoInDefault(sTime string) (t time.Time, err error) {
	return ParseAutoInLocation(sTime, DefaultTimeZone)
}

//utc time
func MustFromMysqlFormat(timeString string) time.Time {
	t, err := time.Parse(FormatMysql, timeString)
	if err != nil {
		panic(err)
	}
	return t
}

func MustFromMysqlFormatInLocation(timeString string, loc *time.Location) time.Time {
	t, err := time.ParseInLocation(FormatMysql, timeString, loc)
	if err != nil {
		panic(err)
	}
	return t
}

//使用默认时区解释mysql 数据结构
func MustFromMysqlFormatDefaultTZ(timeString string) time.Time {
	if timeString == "0000-00-00 00:00:00" {
		return time.Time{}
	}
	t, err := time.ParseInLocation(FormatMysql, timeString, DefaultTimeZone)
	if err != nil {
		panic(err)
	}
	return t
}

//local time
func MustFromLocalMysqlFormat(timeString string) time.Time {
	t, err := time.ParseInLocation(FormatMysql, timeString, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func ToLocal(t time.Time) time.Time {
	return t.Local()
}

func ToDateString(t time.Time) string {
	return t.Format(FormatDateMysql)
}

//规整到日期,去掉时分秒
func ToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

//规整到日期然后相减
func DateSub(t1 time.Time, t2 time.Time, loc *time.Location) time.Duration {
	return ToDate(t1.In(loc)).Sub(ToDate(t2.In(loc)))
}

//规整到日期然后相减,返回天数
func DateSubToDay(t1 time.Time, t2 time.Time, loc *time.Location) int {
	dur := ToDate(t1.In(loc)).Sub(ToDate(t2.In(loc)))
	return int(dur.Hours() / 24)
}

func DateSubLocal(t1 time.Time, t2 time.Time) time.Duration {
	return DateSub(t1, t2, time.Local)
}

//是否同一天
func IsSameDay(t1 time.Time, t2 time.Time, loc *time.Location) bool {
	return DateSub(t1, t2, loc) == 0
}

//规则到秒,去掉毫秒什么的
func ModBySecond(t1 time.Time) time.Time {
	return t1.Round(time.Second)
}

func GetUnixFloat(t1 time.Time) float64{
	return (float64(t1.Nanosecond())/1e9)+float64(t1.Unix())
}

func FromUnixFloat(f float64) time.Time{
	s,ns:=math.Modf(f)
	return time.Unix(int64(s),int64(ns*1e9)).In(DefaultTimeZone)
}