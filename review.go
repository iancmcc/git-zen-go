package main

import (
	"fmt"
)

var (
	review *ReviewCommand
)

func init() {
	review = &ReviewCommand{}
	parser.AddCommand("review", "Request a review of a feature", "", review)
}

func checkChanges(r *Repository) {
	if r.HasChanges() {
		fmt.Println("You have uncommitted changes. Commit them.")
	}
}

type ReviewCommand struct{}

func (b *ReviewCommand) Execute(args []string) error {
	repo := NewRepository("")
	checkChanges(repo)
	return nil
}
