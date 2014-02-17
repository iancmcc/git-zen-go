package main

import (
	"bytes"
	"fmt"
	flags "github.com/zenoss/go-flags"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

func execute(pwd string, bin string, args ...string) (int, string) {
	b := &bytes.Buffer{}
	if pwd == "" {
		pwd, _ = os.Getwd()
	}
	cmd := exec.Command(bin, args...)
	cmd.Dir = pwd
	cmd.Stdout = b
	cmd.Stderr = b
	if opts.Verbose {
		cmd.Stdout = io.MultiWriter(cmd.Stdout, os.Stdout)
		cmd.Stderr = io.MultiWriter(cmd.Stderr, os.Stderr)
	}
	err := cmd.Run()
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// There is no plattform independent way to retrieve
		// the exit code, but the following will work on Unix
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), b.String()
		}
	}
	return 0, b.String()
}

type Repository struct {
	path string
}

func NewRepository(path string) *Repository {
	return &Repository{path}
}

func (r *Repository) Git(args ...string) (int, string) {
	return execute(r.path, gitbin, args...)
}

func (r *Repository) Gitflow(args ...string) (int, string) {
	return execute(r.path, gitflowbin, args...)
}

func (r *Repository) Branch() string {
	_, out := r.Git("status", "-s", "-b")
	return strings.Split(string(out[3:]), "\n")[0]
}

func (r *Repository) HasBranch(branch string) bool {
	_, out := r.Git("branch", "--list", branch)
	return len(out) > 0
}

func (r *Repository) HasChanges() bool {
	_, out := r.Git("diff-index", "--quiet", "HEAD", "--")
	fmt.Println(out)
	return len(out) > 0
}

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
