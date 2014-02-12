package main

import (
	"errors"
)

func init() {
	begin := &BeginCommand{}
	parser.AddCommand("begin", "Begin working on a feature", "", begin)
}

func initialize(r *Repository) {
	if exit := r.gitflow("init", "-d"); exit == 0 {
		r.git("push", "origin", "--all")
	}
}

func begin(r *Repository, feature string) {
	r.git("stash")
	defer r.git("stash", "pop")
	result := r.gitflow("feature", "start", feature)
	if result == 1 {
		// Branch already exists
		r.git("checkout", "feature/"+feature)
	}
	r.gitflow("feature", "publish", feature)
}

type BeginCommand struct{}

func (b *BeginCommand) Execute(args []string) error {
	verifyDeps()
	if len(args) > 0 {
		repo := NewRepository("")
		initialize(repo)
		feature := args[0]
		begin(repo, feature)
	} else {
		return errors.New("Please specify feature")
	}
	return nil
}
