package kmgBootstrap

import (
	"fmt"

	"github.com/bronze1man/kmg/kmgView"
	"github.com/bronze1man/kmg/kmgXss"
	"strconv"
	"strings"
)

type Panel struct {
	Title string
	Body  kmgView.HtmlRenderer
}

func (p Panel) HtmlRender() string {
	return tplPanel(p)
}

type HorizontalDescription struct {
	Key   string
	Value string
}

func (p HorizontalDescription) HtmlRender() string {
	return tplHorizontalDescription(p)
}

type Table struct {
	Caption   kmgView.HtmlRenderer
	TitleList []kmgView.HtmlRenderer
	DataList  [][]kmgView.HtmlRenderer
}

func (p *Table) SetTitleListString(TitleList []string){
	p.TitleList = make([]kmgView.HtmlRenderer,len(TitleList))
	for i, title := range TitleList {
		p.TitleList[i] = kmgView.String(title)
	}
}

func (p Table) HtmlRender() string {
	return tplTable(p)
}

func NewTable(TitleList []string, DataList [][]string) Table {
	t := Table{}
	t.SetTitleListString(TitleList)
	for _, row := range DataList {
		renderRow := []kmgView.HtmlRenderer{}
		for _, value := range row {
			renderRow = append(renderRow, kmgView.String(value))
		}
		t.DataList = append(t.DataList, renderRow)
	}
	return t
}

// 显示键值对显示信息的 表格
func NewNoTitleTable(DataList [][]kmgView.HtmlRenderer) Table {
	return Table{
		DataList: DataList,
	}
}

// 1.这个函数顺序会变,基本不靠谱.
// 2.直接从数据库取的数据,字段信息没有翻译.
func NewTableFromMapList(m []map[string]string) Table {
	t := Table{}
	if len(m) == 0 {
		return t
	}
	i := 0
	for k := range m[0] {
		// TODO 排序
		t.TitleList = append(t.TitleList, kmgView.String(k))
		i++
	}
	for _, row := range m {
		viewRow := []kmgView.HtmlRenderer{}
		for _, k := range t.TitleList {
			viewRow = append(viewRow, kmgView.String(row[string(k.(kmgView.String))]))
		}
		t.DataList = append(t.DataList, viewRow)
	}
	return t
}

type HorizontalHtmlRenderList []kmgView.HtmlRenderer

func (p HorizontalHtmlRenderList) HtmlRender() string {
	return tplHorizontalHtmlRenderList(p)
}

type NavTabList struct {
	ActiveName  string //选中的名称
	OptionList  []NavTabOption
	CustomClass string // 两种样式 nav-pills / nav-tabs，不填将使用 nav-pills
}

type NavTabOption struct {
	Name string
	Url  string
}

func (p NavTabList) HtmlRender() string {
	return tplNavTabList(p)
}

type Image struct {
	Src       string
	MaxHeight int
	MaxWidth  int
}

func (p Image) HtmlRender() string {
	return tplImage(p)
}

func (p Image) getStyle() string {
	if p.MaxHeight == 0 {
		p.MaxHeight = 100
	}
	if p.MaxWidth == 0 {
		p.MaxWidth = 100
	}
	return fmt.Sprintf("max-height: %dpx;max-width: %dpx;", p.MaxHeight, p.MaxWidth)
}

type A struct {
	Href        string
	Title       string
	IsNewWindow bool
	Id          string
	ClassName   string
}

func (a A) HtmlRender() string {
	s := `<a href="` + kmgXss.H(a.Href) + `" class="` + kmgXss.H(a.ClassName) + `"`
	if a.IsNewWindow {
		s += ` target="_blank" `
	}
	if a.Id != "" {
		s += ` id="` + kmgXss.H(a.Id) + `" `
	}
	return s + `>` + kmgXss.H(a.Title) + `</a>`
}

type TextColor string

const (
	TextMuted   TextColor = "text-muted"
	TextPrimary TextColor = "text-primary"
	TextSuccess TextColor = "text-success" // 绿色
	TextInfo    TextColor = "text-info"
	TextWarning TextColor = "text-warning"
	TextDanger  TextColor = "text-danger"
)

type BgColor string

const (
	BgPrimary TextColor = "bg-primary"
	BgSuccess TextColor = "bg-success"
	BgInfo    TextColor = "bg-info"
	BgWarning TextColor = "bg-warning"
	BgDanger  TextColor = "bg-danger"
)

func IconGreenCircle() kmgView.HtmlRenderer {
	return Icon{IconColor: TextSuccess, IconName: "circle"}
}

type Icon struct {
	IconName      string
	IconColor     TextColor
	AttributeNode kmgView.HtmlRenderer
}

func (icon Icon) HtmlRender() string {
	return tplIcon(icon)
}

type PlacementType string

const (
	PlacementTypeLeft   PlacementType = "left"
	PlacementTypeRight  PlacementType = "right"
	PlacementTypeTop    PlacementType = "top"
	PlacementTypeBottom PlacementType = "bottom"
	PlacementTypeAuto   PlacementType = "auto"
)

type PopoverType string

const (
	PopoverTypePopover PopoverType = "popover"
	PopoverTypeTooltip PopoverType = "tooltip"
)

//注意
//Bootstrap 有两种提示小窗口，一种叫 Tooltip，一种叫 Popover，两者数据结构基本一致，这里为了简单，直接合并为一种 kmgBootstrap 的结构 Popover
type Popover struct {
	Type      PopoverType
	Title     string
	Content   string
	Placement PlacementType
}

func (p Popover) HtmlRender() string {
	if p.Type == "" {
		p.Type = PopoverTypeTooltip
	}
	if p.Placement == "" {
		p.Placement = PlacementTypeAuto
	}
	return tplPopover(p)
}

type LabelColor string

const (
	LabelColorDefault LabelColor = "label-default"
	LabelColorPrimary LabelColor = "label-primary"
	LabelColorSuccess LabelColor = "label-success"
	LabelColorInfo    LabelColor = "label-info"
	LabelColorWarning LabelColor = "label-warning"
	LabelColorDanger  LabelColor = "label-danger"
)

type Label struct {
	Color   LabelColor
	Content kmgView.HtmlRenderer
}

func (l Label) HtmlRender() string {
	if l.Color == "" {
		l.Color = LabelColorDefault
	}
	if l.Content == nil {
		return ""
	}
	return `<span class="label ` + kmgXss.H(string(l.Color)) + `">` + l.Content.HtmlRender() + `</span>`
}

func Blank(num int) kmgView.HtmlRenderer {
	return kmgView.Html(strings.Repeat("&nbsp;", num))
}

func BlankChinese(num int) kmgView.HtmlRenderer {
	return kmgView.Html(strings.Repeat("&#12288;", num))
}

func H(index int, content kmgView.HtmlRenderer) kmgView.HtmlRenderer {
	if content == nil {
		return kmgView.String("")
	}
	if index == 0 {
		index = 1
	}
	iStr := strconv.Itoa(index)
	return kmgView.Html(`<h` + iStr + `>` + content.HtmlRender() + `</h` + iStr + `>`)
}

func HString(index int,content string) kmgView.HtmlRenderer{
	if index==0{
		index=1
	}
	return kmgView.Html(`<h` + strconv.Itoa(index) + `>` + kmgXss.H(content) + `</h` + strconv.Itoa(index) + `>`)
}

func Pre(content string) kmgView.HtmlRenderer {
	return kmgView.Html(`<pre>` + kmgXss.H(content) + `</pre>`)
}

func Br(num int) kmgView.HtmlRenderer {
	return kmgView.Html(strings.Repeat("<br />", num))
}

func Hr(num int) kmgView.HtmlRenderer {
	return kmgView.Html(strings.Repeat("<hr />", num))
}

type NavBarNode struct {
	Name      string
	Url       string
	ChildList []NavBarNode
}

type NavBar struct {
	Title           kmgView.HtmlRenderer //可以是 Log 或者文字
	OptionList      []NavBarNode
	RightOptionList []NavBarNode
}

func (navBar NavBar) HtmlRender() string {
	return tplNavBar(navBar)
}
