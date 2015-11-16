package kmgJpush
import (
	"testing"
)

func TestJpush(*testing.T){
	client := NewClient(NewClientRequest{
		Name:         "Android 快喵（测试）", // 蜗牛企业版
		AppKey:       "4dce413e3eec454ea0284247",
		Secret:       "4be60187fd42f0a8328064bc",
		IsIosProduct: false,
		Platform:     Android,
		IsActive:     true,
	})
	client.PushToOne("110","这是测试消息")
	client.PushToTag("test","这是测试消息")
}
