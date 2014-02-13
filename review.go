package main

import ()

var (
	review *ReviewCommand
)

func init() {
	review = &ReviewCommand{}
	parser.AddCommand("review", "Request a review of a feature", "", review)
}

func checkChanges(r *Repository) {
	r.git("diff")
}

type ReviewCommand struct{}

func (b *ReviewCommand) Execute(args []string) error {
	repo := NewRepository("")
	checkChanges(repo)
	return nil
}
