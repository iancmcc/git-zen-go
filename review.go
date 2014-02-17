package main

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	review      *ReviewCommand
	githubre    *regexp.Regexp
	extensionre *regexp.Regexp
)

func init() {
	review = &ReviewCommand{}
	parser.AddCommand("review", "Request a review of the current feature branch",
		"", review)
	githubre = regexp.MustCompile(`^(?:https?:\/\/|git:\/\/)?(?:[^@]+@)?(github.com)[:\/]([^\/]+\/[^\/]+?|[0-9]+)$`)
	extensionre = regexp.MustCompile(`\.git$`)
}

func checkChanges(r *Repository) error {
	if r.HasChanges() {
		return errors.New("You have uncommitted changes. Commit them.")
	}
	return nil
}

func pullRequestUrl(repourl string) string {
	repourl = extensionre.ReplaceAllString(repourl, "")
	split := githubre.FindStringSubmatch(repourl)
	path := split[2]
	return fmt.Sprintf("/repos/%s/pulls", path)
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
