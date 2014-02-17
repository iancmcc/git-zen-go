package main

import (
	"testing"
)

func TestCheckChanges(t *testing.T) {
	setup()
	defer cleanup()

	begin(repo, "myfeature")

	makeLocalChange(clone)
	commitChanges(clone)
	makeLocalChange(clone)

	err := checkChanges(repo)

	if err == nil {
		t.Fatalf("Failed to detect that a repo has changes")
	}
}

func TestPullRequestUrl(t *testing.T) {
	urls := []string{
		"git@github.com:iancmcc/git-zen.git",
		"git@github.com:iancmcc/git-zen",
		"http://github.com/iancmcc/git-zen",
		"http://github.com/iancmcc/git-zen.git",
	}
	for _, s := range urls {
		u := pullRequestUrl(s)
		if u != "/repos/iancmcc/git-zen/pulls" {
			t.Fatalf("Pull request url for %s generated incorrectly as %s", s, u)
		}
	}
}
