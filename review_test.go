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
