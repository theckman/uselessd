package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
)

type commandLine struct {
	Host    string `short:"H" long:"host" default:"127.0.0.1" description:"address to bind service to"`
	Version bool   `short:"V" long:"version" description:"print uselessd version"`
}

// Parse is a function used to parse the command line args. If the only parameter
// is nil os.Args is used by default.
func (a *commandLine) Parse(args []string) (string, error) {
	if args == nil {
		args = os.Args
	}

	parser := flags.NewParser(a, flags.HelpFlag|flags.PassDoubleDash)

	// args[0] is mir's binary path
	// we don't need to parse that
	_, err := parser.ParseArgs(args[1:])

	// determine if there was a parsing error
	if err != nil {
		// determine whether this was a help message by doing a type
		// assertion of err to *flags.Error and check the error type
		// if it was a help message, do not return an error
		if errType, ok := err.(*flags.Error); ok {
			if errType.Type == flags.ErrHelp {
				return err.Error(), nil
			}
		}

		return "", err
	}

	if a.Version {
		out := fmt.Sprintf(
			"uselessd v%s built with %s\nCopyright 2015 Tim Heckman\n",
			appVersion, runtime.Version(),
		)
		return out, nil
	}

	return "", nil
}
