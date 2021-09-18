# gosrcs

`gosrcs` is a tool to print all the source code a given go package depends on.

The original motivation of this tool is to integrate go builds into other build systems with
proper dependency tracking, but the only real limit is your imagination.

# Usage

```
Usage: gosrcs [options] <module dir>

Options:
  -all
        Print all source code.
  -allow-no-module
        Also print source files that don't have a module (like the stdlib).
  -deps
        Also print sources of dependencies.
  -only-go-files
        Don't print non go sources.
  -pattern string
        Go package pattern to gather sources for. (default ".")
  -relative
        Print the source path in relative form.
```

## Example

```
$ gosrcs ~/src/gosrcs
/home/ac/src/gosrcs/go.mod
/home/ac/src/gosrcs/main.go
$ cd ~/src/gosrcs
$ gosrcs -deps -relative
...
../../go/pkg/mod/golang.org/x/tools@v0.1.6/go/packages/visit.go
...
go.mod
main.go
```
