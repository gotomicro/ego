package flag

import (
	"os"
)

var defaultFlags = []Flag{
	// HelpFlag prints usage of application.
	&BoolFlag{
		Name:  "help",
		Usage: "--help, show help information",
		Action: func(name string, fs *FlagSet) {
			fs.PrintDefaults()
			os.Exit(0)
		},
	},
}
