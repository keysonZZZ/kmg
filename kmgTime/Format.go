package kmgTime

import "time"

const (
	FormatMysql      = "2006-01-02 15:04:05"
	FormatFileName   = "2006-01-02_15-04-05" //适合显示在文件上面的日期格式 @deprecated
	FormatFileNameV2 = "2006-01-02-15-04-05" //版本2,更规整,方便使用正则取出
	FormatDateMysql  = "2006-01-02"
	Iso3339Hour      = "2006-01-02T15"
	Iso3339Minute    = "2006-01-02T15:04"
	Iso3339Second    = "2006-01-02T15:04:05"
	AppleJsonFormat  = "2006-01-02 15:04:05 Etc/MST" //仅解决GMT的这个特殊情况.其他不管,如果苹果返回的字符串换时区了就悲剧了
)

var ParseFormatGuessList = []string{
	FormatMysql,
	FormatDateMysql,
	Iso3339Hour,
	Iso3339Minute,
	Iso3339Second,
}

//输出成mysql的格式,并且使用默认时区,并且在0值的时候输出空字符串
func DefaultFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(DefaultTimeZone).Format(FormatMysql)
}
