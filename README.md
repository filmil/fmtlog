# fmtlog: a log formatter program (like `awk` but lamer)

`fmtlog` is a program that reformats the standard output and standard error for
including in an enclosing log file.

# Example

```console
$ fmtlog --outfmt="[stdout] {{.Text}}{{n}}" \
  --errfmt="[stderr] {{.Text}}{{n}}" \
  bash -c "echo hello 1>&2"
[stderr] hello
$ fmtlog --outfmt="[stdout] {{.Text}}{{n}}" \
  --errfmt="[stderr] {{.Text}}{{n}}" \
  bash -c "echo hello 2>&1"
[stdout] hello
```

`fmtlog` allows you to reformat each line of the standard error and standard
output of the program, and allows you to specify different formats for the
standard error and standard output.  You can get a similar effect using say
`awk`, but `logfmt` is way simpler and more focused, which in some settings is
an advantage.

# Options

```console
$ fmtlog --help
Usage of fmtlog:
  -errfmt string
        format string to use to display stderr (default "[stderr] {{.Text}}{{n}}")
  -indent int
        Set indent value (default 4)
  -outfmt string
        format string to use to display stdout (default "[stdout] {{.Text}}{{n}}")
```

# Default settings

Default settings can be controlled via environment variables too, so that
multiple invocations of `fmtlog` can share parameters:

* `FMTLOG_OUTFMT`: the default value for `--outfmt=...`.
* `FMTLOG_ERRFMT`: the defualt value for `--errfmt=...`.
* `FMTLOG_INDENT`: the default value for `--indent=...`.

# Formatting

`fmtlog` uses the `text/template` engine from the go programming language.  The
following special format strings are defined at the moment:

* `{{.Text}}`: this is the content of the log line
* `{{g}}`: turn on green text
* `{{ir "xyz"}}:`: add indent by repeating `xyz` as many times as specified in `--indent=...`
* `{{i}}`: add indent at this spot, as many spaces as specified in `--indent=...`
* `{{n}}`: print a "newline" (\n)
* `{{r}}`: turn on red text
* `{{w}}`: turn on white text

Thus for example, the following will prefix all standard output lines with green
`----|` prefix, and all standard error lines with red `----|` prefix. Note that
a newline (`{{n}}`) must be added explicitly.

```
$ fmtlog --outfmt="{{g}}{{ir \"-\"}}|{{w}} {{.Text}}{{n}}" \
  --errfmt="{{r}}{{ir \"-\"}}|{{w}} {{.Text}}{{n}}" \
  ls -laR
```

A similar effect can be obtained by offloading the flag parameters into env
variables:

```
$ export FMTLOG_OUTFMT="{{g}}{{ir \"-\"}}|{{w}} {{.Text}}{{n}}"
$ export FMTLOG_ERRFMT="{{r}}{{ir \"-\"}}|{{w}} {{.Text}}{{n}}"
$ fmtlog ls -laR
```
