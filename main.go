package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/mattn/go-encoding"
	"github.com/mattn/go-isatty"
)

var enc = flag.String("e", "cp932", "encoding")

func run() int {
	flag.Parse()
	ioenc := encoding.GetEncoding(*enc)
	if ioenc == nil {
		fmt.Fprintln(os.Stderr, "unknown encoding")
		return 1
	}
	var args []string
	if runtime.GOOS == "windows" {
		args = append([]string{"cmd", "/c"}, flag.Args()...)
	} else {
		args = append([]string{"sh", "-c"}, flag.Args()...)
	}
	cmd := exec.Command(args[0], args[1:]...)
	if isatty.IsTerminal(os.Stdout.Fd()) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	} else {
		cmd.Stdout = ioenc.NewEncoder().Writer(os.Stdout)
		cmd.Stderr = ioenc.NewEncoder().Writer(os.Stderr)
		cmd.Stdin = os.Stdin
	}
	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			if status, ok := err.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			} else {
				panic(errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus."))
			}
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}
