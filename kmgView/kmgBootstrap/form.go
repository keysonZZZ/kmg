package kmgBootstrap

import (
	"github.com/bronze1man/kmg/kmgStrconv"
	"github.com/bronze1man/kmg/kmgView"
)

// 自带一个提交按钮的From表单.
// TODO 允许调用者配置,去掉自带的按钮.
type Form struct {
	Id        string
	Url       string
	IsGet     bool // 默认是POST
	InputList []kmgView.HtmlRenderer
	NoSubmit  bool // 默认有一个提交按钮.
	IsHidden  bool
}

func (f Form) HtmlRender() string {
	return tplForm(f)
}

// 纵向输入框外壳
type InputWrapVertical struct {
	ShowName  string
	Comment   string
	Body      kmgView.HtmlRenderer
	AppendTpl kmgView.HtmlRenderer
	Need      bool
}

func (f InputWrapVertical) HtmlRender() string {
	return tplInputWrapVertical(f)
}

// html的输入框
type Input struct {
	Type     string
	Name     string
	Value    string
	ReadOnly bool
}

func (f Input) HtmlRender() string {
	return tplInput(f)
}

type Select struct {
	Name       string
	Value      string
	ReadOnly   bool
	OptionList []SelectOption
}

type SelectOption struct {
	Value    string
	ShowName string
	Disable  bool
}

func (s Select) HtmlRender() string {
	return tplSelect(s)
}

type SelectVerticalString struct {
	Name       string
	Value      string
	ReadOnly   bool
	OptionList []SelectOption
	ShowName   string
	Comment    string
	Need       bool
}

func (s SelectVerticalString) HtmlRender() string {
	if s.ShowName == "" {
		s.ShowName = s.Name
	}
	return InputWrapVertical{
		ShowName: s.ShowName,
		Comment:  s.Comment,
		Need:     s.Need,
		Body: Select{
			Name:       s.Name,
			Value:      s.Value,
			ReadOnly:   s.ReadOnly,
			OptionList: s.OptionList,
		},
	}.HtmlRender()
}

type TextArea struct {
	Name     string
	Value    string
	ReadOnly bool
}

func (f TextArea) HtmlRender() string {
	return tplTextArea(f)
}

// 纵向输入框 (横向占满)
type InputVerticalString struct {
	Name      string
	Value     string
	ShowName  string
	Comment   string
	Need      bool
	ReadOnly  bool
	AppendTpl kmgView.HtmlRenderer
}

func (f InputVerticalString) HtmlRender() string {
	if f.ShowName == "" {
		f.ShowName = f.Name
	}
	return InputWrapVertical{
		ShowName: f.ShowName,
		Comment:  f.Comment,
		Need:     f.Need,
		Body: Input{
			Type:     "text",
			Name:     f.Name,
			Value:    f.Value,
			ReadOnly: f.ReadOnly,
		},
		AppendTpl: f.AppendTpl,
	}.HtmlRender()
}

type InputVerticalInt struct {
	Name     string
	Value    int
	ShowName string
	Comment  string
	Need     bool
	ReadOnly bool
}

func (f InputVerticalInt) HtmlRender() string {
	if f.ShowName == "" {
		f.ShowName = f.Name
	}
	return InputWrapVertical{
		ShowName: f.ShowName,
		Comment:  f.Comment,
		Need:     f.Need,
		Body: Input{
			Type:     "number",
			Name:     f.Name,
			Value:    kmgStrconv.FormatInt(f.Value),
			ReadOnly: f.ReadOnly,
		},
	}.HtmlRender()
}

type TextAreaVerticalString struct {
	Name     string
	Value    string
	ShowName string
	Comment  string
	Need     bool
	ReadOnly bool
}

func (f TextAreaVerticalString) HtmlRender() string {
	if f.ShowName == "" {
		f.ShowName = f.Name
	}
	return InputWrapVertical{
		ShowName: f.ShowName,
		Comment:  f.Comment,
		Need:     f.Need,
		Body: TextArea{
			Name:     f.Name,
			Value:    f.Value,
			ReadOnly: f.ReadOnly,
		},
	}.HtmlRender()
}

type InputHidden struct {
	Name  string
	Value string
}

func (f InputHidden) HtmlRender() string {
	return Input{
		Type:  "hidden",
		Name:  f.Name,
		Value: f.Value,
	}.HtmlRender()
}
