package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	origin string
	clone  string
	repo   *Repository
	dirs   = []string{}
)

func testGitRepo() string {
	td, _ := ioutil.TempDir("", "gztest")
	repo := NewRepository(td)
	repo.Git("init")
	cmd := exec.Command("touch", "testfile")
	cmd.Dir = td
	cmd.Run()
	repo.Git("add", "testfile")
	repo.Git("commit", "-m", "Commit")
	dirs = append(dirs, td)
	return td
}

func cloneGitRepo(origin string) string {
	td, _ := ioutil.TempDir("", "gztest")
	repo := NewRepository(td)
	repo.Git("clone", origin, td)
	dirs = append(dirs, td)
	return td
}

func cleanup() {
	for i := 0; i < len(dirs); i++ {
		os.RemoveAll(dirs[i])
	}
}

func setup() {
	//opts.Verbose = true
	verifyDeps()
	origin = testGitRepo()
	clone = cloneGitRepo(origin)
	repo = NewRepository(clone)
	initialize(repo)
}

func assertHasBranch(t *testing.T, path, branch string) {
	if !NewRepository(path).HasBranch(branch) {
		t.Fatalf("Repository %s doesn't have branch %s", path, branch)
	}
}

func assertCurrentBranch(t *testing.T, path, branch string) {
	if NewRepository(path).Branch() != branch {
		t.Fatalf("Repository %s is not on branch %s", path, branch)
	}
}

func assertUncommittedChanges(t *testing.T, path string) {
	if !NewRepository(path).HasChanges() {
		t.Fatalf("Repository %s has no uncommitted changes", path)
	}
}

func TestInitialization(t *testing.T) {
	setup()
	defer cleanup()

	assertHasBranch(t, origin, "master")
	assertHasBranch(t, origin, "develop")
	assertHasBranch(t, clone, "master")
	assertHasBranch(t, clone, "develop")

	// Initialize again
	initialize(NewRepository(clone))

	assertHasBranch(t, origin, "master")
	assertHasBranch(t, origin, "develop")
	assertHasBranch(t, clone, "master")
	assertHasBranch(t, clone, "develop")

}

func TestBegin(t *testing.T) {
	setup()
	defer cleanup()

	begin(repo, "myfeature")

	assertHasBranch(t, origin, "feature/myfeature")
	assertHasBranch(t, clone, "feature/myfeature")
	assertCurrentBranch(t, clone, "feature/myfeature")
}

func makeLocalChange(path string) string {
	testfile := filepath.Join(path, "testfile")
	exec.Command("bash", "-c", fmt.Sprintf("echo 1 >> %s", testfile)).Run()
	NewRepository(path).Git("add", "testfile")
	return testfile
}

func commitChanges(path string) {
	NewRepository(path).Git("commit", "-m", "x")
}

func TestBeginWithLocalChanges(t *testing.T) {
	setup()
	defer cleanup()

	testfile := makeLocalChange(clone)
	commitChanges(clone)

	exec.Command("bash", "-c", fmt.Sprintf("echo 1 >> %s", testfile)).Run()
	data1, _ := ioutil.ReadFile(testfile)

	begin(repo, "myfeature")

	data2, _ := ioutil.ReadFile(testfile)

	if string(data1) != string(data2) {
		t.Fatalf("Data isn't equal before and after")
	}

	assertUncommittedChanges(t, clone)
}

func TestBranchAlreadyExists(t *testing.T) {
	setup()
	defer cleanup()

	begin(repo, "myfeature")
	begin(repo, "otherfeature")
	begin(repo, "myfeature")

	assertHasBranch(t, clone, "feature/myfeature")
	assertCurrentBranch(t, clone, "feature/myfeature")
}
