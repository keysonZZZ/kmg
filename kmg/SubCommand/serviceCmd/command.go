package serviceCmd

import (
	"fmt"
	"os"

	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgConsole"
)

//TODO 统一进程输出的输出位置,默认不区分err和out.
//TODO 添加stopanduninstall命令(停掉,并且卸载),删除或隐藏其他无用命令

//TODO 添加至少一个平台的测试
//TODO kmg gorun 可以在service运行的进程里面,目前是restart 存在bug
//TODO 完整描述使用过程
//TODO 在osx和linux上达到一致的行为

func AddCommandList() {
	cmdGroup := kmgConsole.NewCommandGroup().
		AddCommand(kmgConsole.Command{
		Name:   "setAndRestart",
		Desc:   "install the service,and restart the service,uninstall and stop if need",
		Runner: setAndRestartCmd,
	}).AddCommandWithName(
		"setAndRestartV1",
		setAndRestartCmdV1,
	).AddCommand(kmgConsole.Command{
		Name:   "install",
		Desc:   "install the service",
		Runner: installCmd,
	}).AddCommand(kmgConsole.Command{
		Name:   "uninstall",
		Desc:   "uninstall the serivce",
		Runner: newNameCmd(Uninstall),
	}).AddCommand(kmgConsole.Command{
		Name:   "start",
		Desc:   "start the service",
		Runner: newNameCmd(Start),
	}).AddCommand(kmgConsole.Command{
		Name:   "stop",
		Desc:   "stop the service",
		Runner: newNameCmd(Stop),
	}).AddCommand(kmgConsole.Command{
		Name:   "restart",
		Desc:   "restart the service",
		Runner: newNameCmd(Restart),
	})

	kmgConsole.AddCommand(kmgConsole.Command{
		Name:   "Service",
		Runner: cmdGroup.Main,
	})
	kmgConsole.AddCommand(kmgConsole.Command{
		Name:   "Service.Process",
		Runner: processCmd,
		Hidden: true,
	})
}

type installRequest struct {
	Name             string   //名字
	ExecuteArgs      []string //执行的命令,第一个是执行命令的进程地址
	WorkingDirectory string   //工作目录(默认是当前目录) 不能在upstart系统上面运行
}

func setAndRestartCmd() {
	s, err := parseInstallRequest()
	kmgConsole.ExitOnStdErr(err)
	err = Install(s)
	if err != nil {
		if err != ErrServiceExist {
			kmgConsole.ExitOnStdErr(err)
		}
		err = Uninstall(s.Name)
		kmgConsole.ExitOnStdErr(err)
		err = Install(s)
		kmgConsole.ExitOnStdErr(err)
	}
	err = Restart(s.Name)
	kmgConsole.ExitOnStdErr(err)
}

func setAndRestartCmdV1() {
	s, err := parseInstallRequest()
	kmgConsole.ExitOnStdErr(err)
	err = Install(s)
	if err != nil {
		if err != ErrServiceExist {
			kmgConsole.ExitOnStdErr(err)
		}
		err = Uninstall(s.Name)
		kmgConsole.ExitOnStdErr(err)
		err = Install(s)
		kmgConsole.ExitOnStdErr(err)
	}
	err = RestartV1(s.Name)
	kmgConsole.ExitOnStdErr(err)
}

func installCmd() {
	s, err := parseInstallRequest()
	kmgConsole.ExitOnStdErr(err)
	err = Install(s)
	kmgConsole.ExitOnStdErr(err)
}

func parseInstallRequest() (s *Service, err error) {
	if len(os.Args) < 3 {
		return nil, errors.New("require name,exec args")
	}
	s = &Service{}
	s.Name = os.Args[1]
	s.CommandLineSlice = os.Args[2:]
	s.WorkingDirectory, err = os.Getwd()
	if err != nil {
		return
	}
	return
}

func newNameCmd(fn func(name string) error) func() {
	return func() {
		if len(os.Args) <= 1 {
			fmt.Println("require name args")
			return
		}
		name := os.Args[1]
		err := fn(name)
		kmgConsole.ExitOnStdErr(err)
	}
}
