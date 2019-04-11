package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/filmil/fmtlog/pkg/fmtlog"
)

var c fmtlog.Config

const (
	outfmtEnv = "FMTLOG_OUTFMT"
	errfmtEnv = "FMTLOG_ERRFMT"
	indentEnv = "FMTLOG_INDENT"
)

// getenvDef returns the environment variable value by name, or a given default
// if none is provided.
func getenvDef(name, def string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return v
}

func getenvIntDef(name string, def int) int {
	v, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func main() {
	flag.StringVar(&c.OutFmt, "outfmt",
		getenvDef(outfmtEnv, "[stdout] {{.Text}}{{n}}"),
		"format string to use to display stdout")
	flag.StringVar(&c.ErrFmt, "errfmt",
		getenvDef(errfmtEnv, "[stderr] {{.Text}}{{n}}"),
		"format string to use to display stderr")
	flag.IntVar(&c.Indent, "indent",
		getenvIntDef(indentEnv, 4),
		"Set indent value")
	flag.Parse()
	c.Args = flag.Args()
	if len(c.Args) == 0 {
		// Nothing to do.
		os.Exit(0)
	}
	if err := c.Run(os.Stdout, os.Stderr); err != nil {
		fmt.Printf("errr: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
