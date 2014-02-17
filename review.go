package main

import (
	"errors"
)

var (
	review *ReviewCommand
)

func init() {
	review = &ReviewCommand{}
	parser.AddCommand("review", "Request a review of a feature", "", review)
}

func checkChanges(r *Repository) error {
	if r.HasChanges() {
		return errors.New("You have uncommitted changes. Commit them.")
	}
	return nil
}

type ReviewCommand struct{}

func (b *ReviewCommand) Execute(args []string) error {
	repo := NewRepository("")
	err := checkChanges(repo)
	if err != nil {
		return err
	}
	return nil
}
