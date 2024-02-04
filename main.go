package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var form = flag.String("f", "%.3f", "format, or 'str' for string")
var fd = flag.Uint64("fd", 2, "output file descriptor")
var onlyReal = flag.Bool("r", false, "only print real")

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	end := time.Now()

	real := end.Sub(start)
	user := cmd.ProcessState.UserTime()
	sys := cmd.ProcessState.SystemTime()

	out := os.NewFile(uintptr(*fd), "out")
	defer out.Close()

	if *onlyReal {
		fmt.Fprintln(out, p(real))
		return
	}
	fmt.Fprintf(out, "real\t%v\n", p(real))
	fmt.Fprintf(out, "user\t%v\n", p(user))
	fmt.Fprintf(out, "sys \t%v\n", p(sys))
}

func p(d time.Duration) string {
	if *form == "str" {
		return d.String()
	}
	return fmt.Sprintf(*form, d.Seconds())
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v [flags...] command [args...]\n\n", os.Args[0])
	flag.PrintDefaults()
}
