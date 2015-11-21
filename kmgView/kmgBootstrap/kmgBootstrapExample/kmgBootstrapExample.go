package main

import (
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/kmgControllerRunner"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgView"
	"github.com/bronze1man/kmg/kmgView/kmgBootstrap"
)

func main() {
	kmgControllerRunner.RegisterController(Example{})
	kmgControllerRunner.EnterPointApiName = "main.Example.Demo1"
	kmgControllerRunner.StartServerCommand()
}

type Example struct{}

func (e Example) Demo1(ctx *kmgHttp.Context) {
	ctx.WriteString(kmgBootstrap.NewWrap("kmg bootstrap demo",
		kmgBootstrap.NavBar{
			Title: kmgView.String("kmgBootstrap DEMO"),
			OptionList: []kmgBootstrap.NavBarNode{
				{Name: "主页", Url: "/"},
				{
					Name: "测试菜单1",
					Url:  "/?menu=1",
					ChildList: []kmgBootstrap.NavBarNode{
						{
							Name: "测试菜单-子菜单1-1",
							Url:  "/?menu=1-1",
							ChildList: []kmgBootstrap.NavBarNode{
								{
									Name: "测试菜单-子菜单1-1-1",
									Url:  "/?menu=1-1-1",
									ChildList: []kmgBootstrap.NavBarNode{
										{Name: "测试菜单-子菜单1-1-1-1", Url: "/?menu=1-1-1-1"},
									},
								},
								{Name: "测试菜单-子菜单1-1-2", Url: "/?menu=1-1-2"},
							},
						},
						{Name: "测试菜单-子菜单1-2", Url: "/?menu=1-2"},
					},
				},
			},
			RightOptionList: []kmgBootstrap.NavBarNode{
				{
					Name: "测试菜单1右",
					Url:  "/",
					ChildList: []kmgBootstrap.NavBarNode{
						{Name: "测试菜单-子菜单1-1右边", Url: "/"},
					},
				},
			},
		},
		kmgBootstrap.Panel{
			Title: "Panel",
			Body: kmgBootstrap.Panel{
				Title: "Panel",
				Body:  kmgView.String("Body"),
			},
		},
		kmgBootstrap.Panel{
			Title: "Table",
			Body: kmgBootstrap.Table{
				Caption: kmgView.String("Caption"),
				TitleList: []kmgView.HtmlRenderer{
					kmgView.String("title1"),
					kmgView.String("title2"),
				},
				DataList: [][]kmgView.HtmlRenderer{
					{
						kmgView.String("r1c1"),
						kmgView.String("r1c2"),
					},
					{
						kmgView.String("r2c1"),
						kmgView.String("r2c2"),
					},
				},
			},
		},
		kmgBootstrap.Panel{
			Title: "DropDown",
			Body: kmgView.HtmlRendererList{
				kmgView.String("使用 DropDown"),
				kmgBootstrap.DropDown{
					Title: kmgBootstrap.Button{
						Size:  kmgBootstrap.ButtonSizeSmall,
						Color: kmgBootstrap.ButtonColorSuccess,
						Content: kmgView.HtmlRendererList{
							kmgView.String("更多"),
							kmgBootstrap.Blank(1),
							kmgBootstrap.NewCaret(),
						},
					},
					OptionList: []kmgBootstrap.NavBarNode{
						{
							Name: "Say",
							Url:  "/",
							ChildList: []kmgBootstrap.NavBarNode{
								{
									Name: "你好",
								},
								{
									Name: "Hello",
								},
								{
									Name: "もしもし",
								},
								{
									Name: "Ciao",
								},
							},
						},
						{
							Name: "Do",
							ChildList: []kmgBootstrap.NavBarNode{
								{
									Name: "Walk",
								},
								{
									Name: "Sleep",
								},
							},
						},
					},
				},
				kmgView.String("将操作集中起来,节省空间"),
				kmgBootstrap.Br(1),
				kmgView.String("还可以使用快捷方法 kmgBootstrap.NewMoreButton,获得默认样式"),
				kmgBootstrap.NewMoreButton([]kmgBootstrap.NavBarNode{
					{
						Name: "新增",
					},
					{
						Name: "删除",
					},
					{
						Name: "编辑",
					},
				}),
			},
		},
		kmgBootstrap.Panel{
			Title: "Icon",
			Body: kmgBootstrap.Icon{
				IconName:      "exchange",
				IconColor:     kmgBootstrap.TextDanger,
				AttributeNode: kmgBootstrap.Popover{Title: "这里应该有一个图标"},
			},
		},
		kmgBootstrap.Panel{
			Title: "NavTabList",
			Body: kmgView.HtmlRendererList{
				kmgBootstrap.NavTabList{
					ActiveName: "状态1",
					OptionList: []kmgBootstrap.NavTabOption{
						{Name: "状态1", Url: "/#1"},
						{Name: "状态2", Url: "/#2"},
					},
				},
				kmgBootstrap.Br(1),
				kmgBootstrap.NavTabList{
					ActiveName:  "状态3",
					CustomClass: "nav-tabs",
					OptionList: []kmgBootstrap.NavTabOption{
						{Name: "状态3", Url: "/#3"},
						{Name: "状态4", Url: "/#4"},
					},
				},
			},
		},
		kmgBootstrap.Panel{
			Title: "Form",
			Body: kmgView.HtmlRendererList{
				kmgBootstrap.Form{
					Url: "/",
					InputList: []kmgView.HtmlRenderer{
						kmgBootstrap.InputVerticalString{
							Name:     "UserName",
							ShowName: "用户名",
							Value:    "enter your username here",
							Comment:  "必填",
							Need:     true,
						},
						kmgBootstrap.TextAreaVerticalString{
							ShowName: "服务条款",
							Value:    "服务条款",
							ReadOnly: true,
						},
						kmgBootstrap.SelectVerticalString{
							ShowName: "性别",
							Value:    "girl",
							Name:     "Gender",
							OptionList: []kmgBootstrap.SelectOption{
								{ShowName: "男", Value: "boy"},
								{ShowName: "女", Value: "girl"},
							},
						},
					},
				},
				kmgBootstrap.Br(2),
				kmgBootstrap.Pre(`从 URL 直接发送 POST 的链接`),
				kmgBootstrap.NewPostButton("/?n=main.Example.DemoPostAction&Name=kmg&Age=12", "POST 请求"),
				kmgBootstrap.Blank(2),
				kmgBootstrap.NewGetButton("/?n=main.Example.DemoPostAction&Name=kmg&Age=12", "非 POST 请求"),
			},
		},
	).HtmlRender())
}

func (e Example) DemoPostAction(ctx *kmgHttp.Context) {
	ctx.MustPost()
	b := kmgYaml.MustMarshal(ctx.GetInMap())
	ctx.WriteString(kmgBootstrap.NewWrap("DemoPostAction", kmgBootstrap.Pre(`
You send me a HTTP POST Request
`+string(b))).HtmlRender())
}
