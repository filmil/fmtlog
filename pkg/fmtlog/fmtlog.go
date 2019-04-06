package fmtlog

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"
)

type Config struct {
	OutFmt string
	ErrFmt string
	Args   []string
}

type Fmt struct {
	Text string
}

func readPipe(ctx context.Context, pipe io.ReadCloser, tpl *template.Template, fwd io.WriteCloser) {
	s := bufio.NewScanner(pipe)
	defer fwd.Close()
	for s.Scan() {
		t := fmt.Sprintf("%s\n", s.Text())
		f := Fmt{
			Text: t,
		}
		if err := tpl.Execute(fwd, f); err != nil {
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

func Run(c Config, out, errout io.WriteCloser) error {
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
	// Add textual template for writing this stuff out.  Needs to be instantiated
	// here because any errors need to be reported back.
	errTmpl, err := template.New("err").Funcs(
		template.FuncMap{
			"red":   red,
			"white": white,
		},
	).Parse(c.ErrFmt)
	if err != nil {
		return fmt.Errorf("errtpl: %v", err)
	}
	outTmpl, err := template.New("out").Funcs(
		template.FuncMap{
			"red":   red,
			"white": white,
		},
	).Parse(c.OutFmt)
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
