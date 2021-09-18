package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

const mode packages.LoadMode = packages.NeedName |
	packages.NeedDeps |
	packages.NeedImports |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedModule

func die(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

func main() {

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintln(out, "Usage: gosrcs [options] <module dir>\n")
		fmt.Fprintln(out, "Options:")
		flag.PrintDefaults()
	}

	deps := flag.Bool("deps", false, "Print source paths of dependencies.")
	moduleRequired := flag.Bool("module-required", true, "Only print source paths with corresponding go modules (ignores the stdlib by default).")
	relative := flag.Bool("relative", false, "Print source paths relative to the working directory.")
	pattern := flag.String("pattern", ".", "Go package pattern to gather source paths for.")

	flag.Parse()

	modDir := ""

	switch flag.NArg() {
	case 0:
		// nothing
	case 1:
		modDir = flag.Args()[0]
	default:
		die("Expecting a single argument: directory of module")
	}

	cfg := &packages.Config{Mode: mode, Dir: modDir}
	pkgs, err := packages.Load(cfg, *pattern)
	if err != nil {
		die(err.Error())
	}

	printedModFiles := make(map[string]struct{})

	cwd, err := os.Getwd()
	if err != nil {
		die(err.Error())
	}

	printPath := func(p string) {
		if !*relative {
			_, _ = fmt.Println(p)
		} else {
			p, err := filepath.Rel(cwd, p)
			if err != nil {
				die(err.Error())
			}
			_, _ = fmt.Println(p)
		}
	}

	printSrcs := func(pkg *packages.Package) {
		if !*deps && pkg.Module != nil && !pkg.Module.Main {
			return
		}
		if pkg.Module != nil {
			if _, printed := printedModFiles[pkg.Module.GoMod]; !printed {
				printPath(pkg.Module.GoMod)
				printedModFiles[pkg.Module.GoMod] = struct{}{}
			}
		} else if *moduleRequired {
			return
		}
		for _, p := range pkg.GoFiles {
			printPath(p)
		}
		for _, p := range pkg.OtherFiles {
			printPath(p)
		}
	}

	packages.Visit(pkgs, nil, printSrcs)
}
