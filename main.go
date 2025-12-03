package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/yottta/cobra_tofu/commands"
)

func main() {
	debugToFile()
	cmd := commands.CobraCommands()
	if commands.ExecuteLegacyAutocompletion(cmd) {
		return
	}
	prepareArgs()

	exitCode, rootCause := commands.ExtractExitCode(cmd.ExecuteContext(context.Background()))
	if rootCause != nil {
		fmt.Printf("Error executing CLI: %s\n", rootCause.Error())
		os.Exit(commands.DefaultErrorExitCode)
	}
	os.Exit(exitCode)
}

// prepareArgs is meant to transform any given flag from a single dash form (e.g.: "-flag=value") to
// a double dash form (e.g.: "--flag=value").
// This way we allow both, go standard flag package and cobra flag parsing, to work with these.
func prepareArgs() {
	alteredArgs := make([]string, len(os.Args))
	for i, arg := range os.Args {
		alteredArgs[i] = arg
		if strings.HasPrefix(arg, "--") {
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 2 {
			alteredArgs[i] = fmt.Sprintf("-%s", arg)
		}
	}
	os.Args = alteredArgs
}

func debugToFile() {
	if path := os.Getenv("DEBUG_TO_FILE"); path != "" {
		var b bytes.Buffer
		gatherEnviron(&b)
		gatherArgs(&b)
		if err := os.WriteFile(path, b.Bytes(), 0644); err != nil {
			fmt.Println("error writing env file to", path, err)
		}
	}
}

func gatherEnviron(b *bytes.Buffer) {
	b.WriteString("================ environment variables ================\n")
	for i, s := range os.Environ() {
		b.WriteString(fmt.Sprintf("%d: %s\n", i, s))
	}
}

func gatherArgs(b *bytes.Buffer) {
	b.WriteString("================ args ================\n")
	for i, s := range os.Args {
		b.WriteString(fmt.Sprintf("%d: %s\n", i, s))
	}
}
