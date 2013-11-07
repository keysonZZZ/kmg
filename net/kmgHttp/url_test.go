package kmgHttp

import (
	"kmg/test"
	"testing"
)

func TestNewUrlByString(T *testing.T) {
	url, err := NewUrlByString("http://www.google.com")
	test.Assert(T, nil, err)
	test.Assert(T, "http://www.google.com", url.String())

}
