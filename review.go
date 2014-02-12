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

type ReviewCommand struct{}

func (b *ReviewCommand) Execute(args []string) error {
	fmt.Println("REVIEW")
	fmt.Println(args)
	return nil
}
