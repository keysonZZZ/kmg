package kmgStrings

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
	"strings"
)

func TestIsAllAphphabet(ot *testing.T) {
	kmgTest.Equal(IsAllAphphabet("abc"), true)
	kmgTest.Equal(IsAllAphphabet(""), true)
	kmgTest.Equal(IsAllAphphabet("123"), false)
}

func TestIsAllNum(ot *testing.T) {
	kmgTest.Equal(IsAllNum("abc"), false)
	kmgTest.Equal(IsAllNum(""), true)
	kmgTest.Equal(IsAllNum("123"), true)
	kmgTest.Equal(IsAllNum("123.1"), false)
}

func TestFirstLetterToUpper(t *testing.T) {
	s := FirstLetterToUpper("title")
	kmgTest.Equal("Title", s)

	s1 := FirstLetterToUpper("")
	kmgTest.Equal("", s1)

	s2 := FirstLetterToUpper("123")
	kmgTest.Equal("123", s2)

	s3 := FirstLetterToUpper("中文")
	kmgTest.Equal("中文", s3)
}
func TestSubStr(t *testing.T) {
	kmgTest.Equal("abc", SubStr("abcdefg", 0, 3))
	kmgTest.Equal("defg", SubStr("abcdefg", 3, 0))
	kmgTest.Equal("abcdef", SubStr("abcdefg", 0, -1))
	kmgTest.Equal("", SubStr("abcdefg", 0, -10))
}
func TestStartWith(t *testing.T) {
	kmgTest.Ok(IsStartWith("abc", "a"))
	kmgTest.Ok(!IsStartWith("abc", "d"))
}

func TestSplit(ot *testing.T){
	kmgTest.Equal(strings.Split("",""),[]string{})
	kmgTest.Equal(strings.Split("","1"),[]string{""})
	kmgTest.Equal(strings.Split("123",","),[]string{"123"})
}