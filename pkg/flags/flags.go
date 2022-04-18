package flags

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/whytehack/goftp/pkg/logs"
)

var parser *flags.Parser

// Parse parses the provided parameters and fills the out parameter while returning the
// positional arguments as an array
func Parse(out interface{}, params []string) ([]string, error) {
	parser = flags.NewParser(out, flags.PrintErrors|flags.PassDoubleDash|flags.IgnoreUnknown)
	positionalArgs, err := parser.ParseArgs(params) // No need to handle errors, as they will be written by the library
	if err != nil {
		logs.ERROR.Printf("Error parsing parameters %v, %v \n", params, err)
		return nil, err
	}
	return positionalArgs[1:], nil
}

// ShowHelp displays help
func ShowHelp() {
	parser.WriteHelp(os.Stdout)
}
