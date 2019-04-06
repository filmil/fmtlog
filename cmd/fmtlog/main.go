package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/filmil/fmtlog/pkg/fmtlog"
)

var c fmtlog.Config

func main() {
	flag.StringVar(&c.OutFmt, "outfmt", "[stdout] {{.Text}}", "out format")
	flag.StringVar(&c.ErrFmt, "errfmt", "[stderr] {{.Text}}", "err format")
	flag.Parse()
	c.Args = flag.Args()
	if len(c.Args) == 0 {
		// Nothing to do.
		os.Exit(0)
	}
	if err := fmtlog.Run(c, os.Stdout, os.Stderr); err != nil {
		fmt.Printf("errr: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
