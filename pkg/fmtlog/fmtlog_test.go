package fmtlog

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type strbuild struct {
	strings.Builder
}

func (s *strbuild) Close() error {
	return nil
}

func TestFmtLog(t *testing.T) {
	tests := []struct {
		name                      string
		stdout, stderr, stdoutEnd string
		cfg                       Config
		expectedOut               string
		expectedErr               string
	}{
		{
			name:   "basic",
			stdout: "hello world\n",
			stderr: "hello error\n",
			cfg: Config{
				OutFmt: "[stdout]{{.Text}}{{n}}",
				ErrFmt: "[error]{{.Text}}{{n}}",
				Args:   []string{os.Args[0], "--test.run=TestHelperProcess", "--"},
			},
			expectedOut: "[stdout]hello world\n",
			expectedErr: "[error]hello error\n",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("FMTLOG_STDOUT", test.stdout)
			os.Setenv("FMTLOG_STDERR", test.stderr)
			os.Setenv("FMTLOG_STDOUT_END", test.stdoutEnd)
			os.Setenv("FMTLOG_WANT_HELPER_PROCESS", "1")
			defer os.Setenv("FMTLOG_WANT_HELPER_PROCESS", "")

			var (
				outBuf, errBuf strbuild
			)

			err := test.cfg.Run(&outBuf, &errBuf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if outBuf.String() != test.expectedOut {
				t.Errorf("stdout mismatch: want: %q, was: %q", test.expectedOut, outBuf.String())
			}
			if errBuf.String() != test.expectedErr {
				t.Errorf("stderr mismatch: want %q, was: %q", test.expectedErr, errBuf.String())
			}
		})
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("FMTLOG_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	stdout := os.Getenv("FMTLOG_STDOUT")
	if stdout != "" {
		fmt.Fprint(os.Stdout, stdout)
	}
	stderr := os.Getenv("FMTLOG_STDERR")
	if stderr != "" {
		fmt.Fprint(os.Stderr, stderr)
	}
	stdoutEnd := os.Getenv("FMTLOG_STDOUT_END")
	if stdoutEnd != "" {
		fmt.Fprint(os.Stdout, stdoutEnd)
	}
}
