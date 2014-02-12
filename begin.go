package main

import (
	"errors"
)

var (
	begin *BeginCommand
)

func init() {
	begin = &BeginCommand{}
	parser.AddCommand("begin", "Begin working on a feature", "", begin)
}

type BeginCommand struct{}

func (b *BeginCommand) Execute(args []string) error {
	verifyDeps()
	if len(args) > 0 {
		feature := args[0]
		gitflow("init", "-d")
		gitflow("feature", "start", feature)
		git("stash")
		gitflow("feature", "publish", feature)
		git("stash", "apply")
	} else {
		return errors.New("Please specify feature")
	}
	return nil
}
