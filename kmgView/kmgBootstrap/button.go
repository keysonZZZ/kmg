package kmgBootstrap

import (
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgView"
	"net/url"
)

type ButtonSize string

const (
	ButtonSizeLarge      ButtonSize = "btn-lg"
	ButtonSizeDefault    ButtonSize = ""
	ButtonSizeSmall      ButtonSize = "btn-sm"
	ButtonSizeExtraSmall ButtonSize = "btn-xs"
)

type ButtonColor string

const (
	ButtonColorDefault ButtonColor = "btn-default"
	ButtonColorPrimary ButtonColor = "btn-primary"
	ButtonColorSuccess ButtonColor = "btn-success"
	ButtonColorInfo    ButtonColor = "btn-info"
	ButtonColorWarning ButtonColor = "btn-warning"
	ButtonColorDanger  ButtonColor = "btn-danger"
	ButtonColorLink    ButtonColor = "btn-link"
)

type ButtonType string

const (
	ButtonTypeA      ButtonType = "a"
	ButtonTypeButton ButtonType = "button"
)

type Button struct {
	Type          ButtonType
	Url           string
	Color         ButtonColor
	Size          ButtonSize
	Content       kmgView.HtmlRenderer
	AttributeNode kmgView.HtmlRenderer
	ClassName     string
	Id            string
	FormId        string //HTML5 Attribute
	Name          string
	Value         string
}

func (b Button) HtmlRender() string {
	if b.Color == "" {
		b.Color = ButtonColorDefault
	}
	if b.Type == "" {
		b.Type = ButtonTypeButton
	}
	if b.Content == nil {
		panic("kmgBootstrap.Button.Content must not be empty")
	}
	return tplButton(b)
}

// 一个发起Get请求的按钮
// url: 点击后的跳转链接
// title: 按钮上面的字
// 在当前页面打开连接
func NewGetButton(url string, title string) Button {
	return Button{
		Type:    ButtonTypeA,
		Url:     url,
		Content: kmgView.String(title),
		Size:    ButtonSizeExtraSmall,
		Color:   ButtonColorInfo,
	}
}

// 一个发起Post请求的按钮
// url: 将把url转化为表单并使用Post方法提交
// title: 按钮上面的字
// 在当前页面打开连接
//不支持一个 key 里塞多个 value 的 url
func NewPostButton(urlStr string, title string) kmgView.HtmlRenderer {
	formAction, err := url.Parse(urlStr)
	kmgErr.PanicIfError(err)
	query, err := url.ParseQuery(formAction.RawQuery)
	kmgErr.PanicIfError(err)
	formAction.RawQuery = ""
	id := kmgRand.MustCryptoRandToReadableAlphaNum(10)
	form := Form{
		Id:        id,
		NoSubmit:  true,
		Url:       formAction.String(),
		InputList: []kmgView.HtmlRenderer{},
		IsHidden:  true,
	}
	for k, v := range query {
		form.InputList = append(form.InputList, InputHidden{
			Name:  k,
			Value: v[0],
		})
	}
	button := Button{
		FormId:  id,
		Type:    ButtonTypeButton,
		Content: kmgView.String(title),
		Size:    ButtonSizeExtraSmall,
		Color:   ButtonColorInfo,
	}
	return kmgView.Html(form.HtmlRender() + button.HtmlRender())
}
