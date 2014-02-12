package main

import (
	"fmt"
	flags "github.com/zenoss/go-flags"
	"os"
	"os/exec"
)

type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}

var (
	opts       = &Options{}
	parser     = flags.NewParser(opts, flags.Default)
	gitbin     string
	gitflowbin string
	verbose    bool
)

func verifyDeps() {
	var err error
	gitbin, err = exec.LookPath("git")
	if err != nil {
		fmt.Println("Unable to find git in PATH")
		os.Exit(1)
	}
	gitflowbin, err = exec.LookPath("git-flow")
	if err != nil {
		fmt.Println("Unable to find git-flow in PATH")
		os.Exit(1)
	}
}

func execCommand(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

func git(args ...string) error {
	return execCommand(gitbin, args...)
}

func gitflow(args ...string) error {
	return execCommand(gitflowbin, args...)
}

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
