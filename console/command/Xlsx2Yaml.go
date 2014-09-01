package command

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/encoding/kmgExcel"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
)

type Xlsx2Yaml struct {
	filePath         string
	format           string
	isOutputAllSheet bool
}

func (command *Xlsx2Yaml) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Xlsx2Yaml",
		Short: "convert from xlsx(Microsoft excel 2007) to yaml",
	}
}
func (command *Xlsx2Yaml) ConfigFlagSet(f *flag.FlagSet) {
	f.StringVar(&command.filePath, "input", "", "input file path")
	f.StringVar(&command.format, "format", "grid", "output yaml format(grid,raw)")
	f.BoolVar(&command.isOutputAllSheet, "outputAllSheet", false, "is output all sheet(default just out first one)?")
}
func (command *Xlsx2Yaml) Execute(context *console.Context) error {
	if *command.filePath == "" {
		if context.FlagSet().NArg() == 1 {
			*command.filePath = context.FlagSet().Arg(0)
		} else {
			return fmt.Errorf("need input file")
		}
	}
	rawArray, err := kmgExcel.XlsxFile2Array(command.filePath)
	if err != nil {
		return err
	}
	output, err := command.formatOutput(rawArray)
	if err != nil {
		return err
	}
	outByte, err := kmgYaml.Marshal(output)
	if err != nil {
		return err
	}
	_, err = context.Stdout.Write(outByte)
	if err != nil {
		return err
	}
	return nil
}

func (command *Xlsx2Yaml) formatOutput(rawArray [][][]string) (interface{}, error) {
	switch command.format {
	case "raw":
		if command.isOutputAllSheet {
			return rawArray, nil
		} else {
			return rawArray[0], nil
		}
	case "grid":
		o := [][]map[string]string{}
		for _, s := range rawArray {
			o1, err := kmgExcel.TitleArrayToGrid(s)
			if err != nil {
				return nil, err
			}
			o = append(o, o1)
		}
		if command.isOutputAllSheet {
			return o, nil
		} else {
			return o[0], nil
		}
	default:
		return nil, fmt.Errorf("not support output format: %s", command.format)
	}
}
