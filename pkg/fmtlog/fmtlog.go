package fmtlog

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Config struct {
	OutFmt string
	ErrFmt string
	Args   []string
	Indent int
}

type Fmt struct {
	Text string
}

// indent returns a string that is indented as many spaces as in c.
func (c Config) indent() string {
	return strings.Repeat(" ", c.Indent)
}

func (c Config) repeat(s string) string {
	return strings.Repeat(s, c.Indent)
}

func (c Config) Run(out, errout io.WriteCloser) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(context.Background(), c.Args[0], c.Args[1:]...)
	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdoutpipe: %v", err)
	}
	errpipe, err := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start(): %v", err)
	}
	var fmap template.FuncMap = template.FuncMap{
		"r":  red,
		"w":  white,
		"g":  green,
		"n":  nl,
		"i":  c.indent,
		"ir": c.repeat,
	}
	// Add textual template for writing this stuff out.  Needs to be instantiated
	// here because any errors need to be reported back.
	errTmpl, err := template.New("err").Funcs(fmap).Parse(c.ErrFmt)
	if err != nil {
		return fmt.Errorf("errtpl: %v", err)
	}
	outTmpl, err := template.New("out").Funcs(fmap).Parse(c.OutFmt)
	if err != nil {
		return fmt.Errorf("outtpl: %v", err)
	}
	go readPipe(ctx, errpipe, errTmpl, errout)
	go readPipe(ctx, outpipe, outTmpl, out)
	if err := cmd.Wait(); err != nil {
		fmt.Printf("process: %+v", cmd.ProcessState)
		return err
	}
	return nil
}

func readPipe(ctx context.Context, pipe io.ReadCloser, tpl *template.Template, fwd io.WriteCloser) {
	s := bufio.NewScanner(pipe)
	defer fwd.Close()
	for s.Scan() {
		f := Fmt{
			Text: s.Text(),
		}
		if err := tpl.Execute(fwd, f); err != nil {
			fmt.Printf("error while reading: %v", err)
			os.Exit(200)
		}
	}
}

func red() string {
	return "\x1B[0;31m"
}

func white() string {
	return "\x1B[0;37m"
}

func green() string {
	return "\x1B[0;32m"
}

func nl() string {
	return "\n"
}
