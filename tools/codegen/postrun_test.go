package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockExecutor records calls and returns configurable results.
type mockExecutor struct {
	calls   []mockCall
	results map[string]mockResult
}

type mockCall struct {
	name string
	args []string
}

type mockResult struct {
	output []byte
	err    error
}

func newMockExecutor() *mockExecutor {
	return &mockExecutor{
		results: make(map[string]mockResult),
	}
}

func (m *mockExecutor) setResult(name string, output []byte, err error) {
	m.results[name] = mockResult{output: output, err: err}
}

func (m *mockExecutor) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	m.calls = append(m.calls, mockCall{name: name, args: args})
	key := name
	if len(args) > 0 {
		key = name + " " + args[0]
	}
	if r, ok := m.results[key]; ok {
		return r.output, r.err
	}
	return nil, nil
}

func Test_RunPostSteps_Success(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()

	err := RunPostSteps(context.TODO(), mock)
	a.NoError(err)
	a.Len(mock.calls, 3)
	a.Equal("goimports", mock.calls[0].name)
	a.Equal([]string{"-w", "./resources/"}, mock.calls[0].args)
	a.Equal("go", mock.calls[1].name)
	a.Equal([]string{"build", "./resources/..."}, mock.calls[1].args)
	a.Equal("golangci-lint", mock.calls[2].name)
	a.Equal([]string{"run", "./resources/..."}, mock.calls[2].args)
}

func Test_RunPostSteps_BuildFailStopsExecution(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()
	mock.setResult("go build", []byte("compile error: undefined reference"), fmt.Errorf("exit status 1"))

	err := RunPostSteps(context.TODO(), mock)
	a.Error(err)
	a.Contains(err.Error(), "go build ./resources/... failed")
	a.Contains(err.Error(), "compile error: undefined reference")

	// Should have run goimports and build, but NOT golangci-lint
	a.Len(mock.calls, 2)
	a.Equal("goimports", mock.calls[0].name)
	a.Equal([]string{"-w", "./resources/"}, mock.calls[0].args)
	a.Equal("go", mock.calls[1].name)
	a.Equal([]string{"build", "./resources/..."}, mock.calls[1].args)
}

func Test_CommitFiles_Success(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()

	files := []string{"resources/eks-clusters.go", "resources/eks.go"}
	err := CommitFiles(context.TODO(), mock, files, "eks", "EKSCluster")
	a.NoError(err)
	a.Len(mock.calls, 2)

	// git add with all files
	a.Equal("git", mock.calls[0].name)
	a.Equal([]string{"add", "resources/eks-clusters.go", "resources/eks.go"}, mock.calls[0].args)

	// git commit with conventional message
	a.Equal("git", mock.calls[1].name)
	a.Equal([]string{"commit", "-m", "feat(eks): add EKSCluster resource"}, mock.calls[1].args)
}

func Test_CommitFiles_GitAddFails(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()
	mock.setResult("git add", []byte("fatal: not a git repository"), fmt.Errorf("exit status 128"))

	err := CommitFiles(context.TODO(), mock, []string{"resources/foo.go"}, "foo", "Foo")
	a.Error(err)
	a.Contains(err.Error(), "git add failed")
	a.Contains(err.Error(), "fatal: not a git repository")

	// Should not have attempted git commit
	a.Len(mock.calls, 1)
}

func Test_CommitFiles_CommitMessageFormat(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()

	err := CommitFiles(context.TODO(), mock, []string{"resources/s3-bucket.go"}, "s3", "S3Bucket")
	a.NoError(err)

	a.Equal("commit", mock.calls[1].args[0])
	a.Equal("-m", mock.calls[1].args[1])
	a.Equal("feat(s3): add S3Bucket resource", mock.calls[1].args[2])
}

func Test_RunPostAndCommit_DefaultRunsBuildAndCommit(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()

	f := &cliFlags{}
	manifest := &Manifest{
		Files: []ManifestEntry{
			{Path: "resources/svc-res.go", Status: "created"},
		},
	}

	err := runPostAndCommit(f, mock, manifest, "svc", "SvcRes")
	a.NoError(err)

	// Should run all 3 build steps + git add + git commit = 5 calls
	a.Len(mock.calls, 5)
	a.Equal("goimports", mock.calls[0].name)
	a.Equal("go", mock.calls[1].name)
	a.Equal("golangci-lint", mock.calls[2].name)
	a.Equal("git", mock.calls[3].name)
	a.Equal([]string{"add", "resources/svc-res.go"}, mock.calls[3].args)
	a.Equal("git", mock.calls[4].name)
	a.Equal("commit", mock.calls[4].args[0])
}

func Test_RunPostAndCommit_NoBuildSkipsBuildOnly(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()

	f := &cliFlags{noBuild: true}
	manifest := &Manifest{
		Files: []ManifestEntry{
			{Path: "resources/svc-res.go", Status: "created"},
		},
	}

	err := runPostAndCommit(f, mock, manifest, "svc", "SvcRes")
	a.NoError(err)

	// Should only run git add + git commit = 2 calls
	a.Len(mock.calls, 2)
	a.Equal("git", mock.calls[0].name)
	a.Equal([]string{"add", "resources/svc-res.go"}, mock.calls[0].args)
	a.Equal("git", mock.calls[1].name)
	a.Equal("commit", mock.calls[1].args[0])
}

func Test_RunPostAndCommit_NoCommitSkipsCommitOnly(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()

	f := &cliFlags{noCommit: true}
	manifest := &Manifest{
		Files: []ManifestEntry{
			{Path: "resources/svc-res.go", Status: "created"},
		},
	}

	err := runPostAndCommit(f, mock, manifest, "svc", "SvcRes")
	a.NoError(err)

	// Should run all 3 build steps only
	a.Len(mock.calls, 3)
	a.Equal("goimports", mock.calls[0].name)
	a.Equal("go", mock.calls[1].name)
	a.Equal("golangci-lint", mock.calls[2].name)
}

func Test_RunPostAndCommit_NoBuildAndNoCommitSkipsBoth(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()

	f := &cliFlags{noBuild: true, noCommit: true}
	manifest := &Manifest{
		Files: []ManifestEntry{
			{Path: "resources/svc-res.go", Status: "created"},
		},
	}

	err := runPostAndCommit(f, mock, manifest, "svc", "SvcRes")
	a.NoError(err)

	// Should run nothing
	a.Len(mock.calls, 0)
}

func Test_RunPostAndCommit_BuildFailureStopsBeforeCommit(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()
	mock.setResult("go build", []byte("compile error"), fmt.Errorf("exit status 1"))

	f := &cliFlags{}
	manifest := &Manifest{
		Files: []ManifestEntry{
			{Path: "resources/svc-res.go", Status: "created"},
		},
	}

	err := runPostAndCommit(f, mock, manifest, "svc", "SvcRes")
	a.Error(err)
	a.Contains(err.Error(), "build failed")

	// goimports + go build = 2 calls, no git
	a.Len(mock.calls, 2)
	a.Equal("goimports", mock.calls[0].name)
	a.Equal("go", mock.calls[1].name)
}

func Test_RunPostAndCommit_SkippedFilesNotCommitted(t *testing.T) {
	a := assert.New(t)
	mock := newMockExecutor()

	f := &cliFlags{noBuild: true}
	manifest := &Manifest{
		Files: []ManifestEntry{
			{Path: "resources/svc-res.go", Status: "created"},
			{Path: "resources/svc.go", Status: "skipped"},
			{Path: "resources/svc-res_mock_test.go", Status: "updated"},
		},
	}

	err := runPostAndCommit(f, mock, manifest, "svc", "SvcRes")
	a.NoError(err)

	// git add should only include non-skipped files
	a.Equal("git", mock.calls[0].name)
	a.Equal([]string{"add", "resources/svc-res.go", "resources/svc-res_mock_test.go"}, mock.calls[0].args)
}
