package InstallCmd

import (
	"github.com/bronze1man/kmg/kmgConsole"

	"github.com/bronze1man/kmg/kmgCmd"
	//"strings"
	"fmt"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgPlatform"
	"os"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "install",
		Desc:   "install tool",
		Runner: installCmd,
	})
}

var toolList = []kmgConsole.Command{}

func addTool(name string, f func()) {
	toolList = append(toolList, kmgConsole.Command{
		Name:   name,
		Runner: f,
	})
}

func installCmd() {
	addTool("golang", installGolang)
	for _, cmd := range toolList {
		if cmd.Name == os.Args[1] {
			cmd.Runner()
			return
		}
	}
}

/*
本次实现遇到下列问题:
	* 不同操作系统判断
	* root权限判断
	* 如果同名 wget会修改下载组件的名称,使用临时文件夹处理
	* 如果已经安装 cp -rf go /usr/local/go 会再创建一个/usr/local/go/go 目录,而不是更新它
	* 如果已经在/usr/local/bin/go处放了一个执行,又在/bin/go处放了一个执行文件, /bin/go的版本不会被使用,使用函数专门判断这种情况,并且把多余的去除掉
	* 如果上一种情况发生,当前bash不能执行go version,因为当前bash有路径查询缓存(暂时无解了..)
*/
func installGolang() {
	kmgFile.MustChangeToTmpPath()
	if !kmgCmd.MustIsRoot() {
		fmt.Println("you need to be root to install golang")
		return
	}

	p := kmgPlatform.GetCompiledPlatform()
	packageName := ""

	switch {
	case p.Compatible(kmgPlatform.LinuxAmd64):
		packageName = "go1.4.2.linux-amd64.tar.gz"
	case p.Compatible(kmgPlatform.DarwinAmd64):
		packageName = "go1.4.2.darwin-amd64-osx10.8.tar.gz"
	default:
		kmgConsole.ExitOnErr(fmt.Errorf("not support platform [%s]", p))
	}
	kmgCmd.ProxyRun("wget http://kmgtools.qiniudn.com/v1/" + packageName)
	kmgCmd.ProxyRun("tar -xf " + packageName)
	kmgCmd.ProxyRun("cp -rf go /usr/local")
	kmgFile.MustDeleteFile("/bin/go")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/go /bin/go")
	kmgFile.MustDeleteFile("/bin/godoc")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/godoc /bin/godoc")
	kmgFile.MustDeleteFile("/bin/gofmt")
	kmgCmd.ProxyRun("ln -s /usr/local/go/bin/gofmt /bin/gofmt")
	kmgCmd.MustEnsureBinPath("/bin/go")
	kmgCmd.MustEnsureBinPath("/bin/godoc")
	kmgCmd.MustEnsureBinPath("/bin/gofmt")

}
