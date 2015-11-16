package kmgJpush

import (
	"fmt"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/isdamir/jpush"
)

type Client struct {
	c            *jpush.PushClient
	IsIosProduct bool
	name         string
	Platform     SystemPlatform
	IsActive     bool
}

type NewClientRequest struct {
	Name         string //这个客户端的名字
	AppKey       string
	Secret       string
	Platform     SystemPlatform //iOS、Android平台
	IsIosProduct bool           //如果是否false表示向测试设备推送,如果是true表示向正式设备推送,后台的那个开发与正式似乎没有作用.
	IsActive     bool
}
type JpushConfig struct {
	Content string
	Alias string
	Tag string
	Badge string
}
func (c *Client) EasyPush(config *JpushConfig)(err error){
	if c.IsActive == false{
		kmgLog.Log("jpush", "Jpush Client is not active,please checkout your configure",c.name,c.IsIosProduct,c.IsActive,config)
		return
	}
	if config.Badge == ""{
		config.Badge = "1"
	}
	nb := jpush.NewNoticeBuilder()
	nb.SetPlatform(jpush.AllPlatform())
	if config.Alias == ""||config.Tag == ""{
		nb.SetAudience(jpush.AllAudience())
	}else{
		au := &jpush.Audience{}
		if config.Alias != ""{
			au.SetAlias([]string{config.Alias})
		}
		if config.Tag != ""{
			au.SetTag([]string{config.Tag})
		}
		nb.SetAudience(au)
	}
	//Android配置
	notice := jpush.NewNoticeAndroid()
	notice.Alert = config.Content
	nb.SetNotice(notice)

	//iOS配置
	iosNotice := jpush.NewNoticeIos()
	iosNotice.Sound = "default"
	iosNotice.Badge = "1"
	iosNotice.Alert = config.Content
	nb.SetNotice(iosNotice)

	op := jpush.NewOptions()
	op.SetApns_production(c.IsIosProduct)
	nb.SetOptions(op)
	ret, err := c.c.Send(nb)
	if err != nil {
		return err
	}
	if ret.Error.Code == 0 {
		kmgLog.Log("jpush", "Push success", c.name, config.Content)
		return nil
	}
	if ret.Error.Code == 1011 {
		kmgLog.Log("jpush","Not Found User",c.name,config)
		return NotFoundUser
	}
	return fmt.Errorf("code:%d err: %s", ret.Error.Code, ret.Error.Message)
}

type SystemPlatform string

const (
	Android     SystemPlatform = "Android"
	IOS         SystemPlatform = "iOS"
	AllPlatform SystemPlatform = "AllPlatform"
)

func NewClient(req NewClientRequest) *Client {
	return &Client{
		c:            jpush.NewPushClient(req.Secret, req.AppKey),
		IsIosProduct: req.IsIosProduct,
		name:         req.Name,
		Platform:     req.Platform,
		IsActive:     req.IsActive,
	}
}

var NotFoundUser error = errors.New("[kmgJpush] user not exist")

func (c *Client) PushToOne(alias string, content string) (err error) {
	err = c.EasyPush(&JpushConfig{
		Alias: alias,
		Content: content,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PushToTag(tag string, content string) (err error) {
	err = c.EasyPush(&JpushConfig{
		Tag: tag,
		Content: content,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PushToAll(content string) (err error) {
	err = c.EasyPush(&JpushConfig{
		Content: content,
	})
	if err != nil {
		return err
	}
	return nil
}
