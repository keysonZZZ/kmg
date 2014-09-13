package command

import (
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/kmgProjectConfig"
	"github.com/bronze1man/kmg/kmgCmd"
	"os"
)

type GoFmt struct {
}

func (command *GoFmt) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "GoFmt", Short: `format all golang code in current project,or current dir`}
}
func (command *GoFmt) Execute(context *console.Context) (err error) {
	var fmtDir string
	kmgc, err := kmgProjectConfig.FindFromWd()
	if err == nil {
		fmtDir = kmgc.ProjectPath
	} else {
		fmtDir, err = os.Getwd()
		if err != nil {
			return
		}
	}
	cmd := kmgCmd.NewStdioCmd(context, "gofmt", "-w=true", ".")
	cmd.Dir = fmtDir
	return cmd.Run()
}
