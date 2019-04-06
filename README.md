# fmtlog: a log formatter program

`fmtlog` is a program that reformats the standard output and standard error for
including in an enclosing log file.

# Example

```console
fmtlog --stdout-prefix="[output]" -stderr-prefix="[error ]" -- someprg
[output]Hello world
[error ]This is an error.
```
