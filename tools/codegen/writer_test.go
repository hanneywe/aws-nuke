package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFiles_CreateMode(t *testing.T) {
	a := assert.New(t)
	dir := t.TempDir()

	files := map[string]GeneratedFile{
		"test.go": {
			Path:    "test.go",
			Content: "package resources\n",
			Mode:    Create,
		},
	}

	manifest, err := WriteFiles(files, WriteOpts{OutputDir: dir})
	a.NoError(err)
	a.Len(manifest.Files, 1)
	a.Equal("created", manifest.Files[0].Status)

	content, err := os.ReadFile(filepath.Join(dir, "test.go"))
	a.NoError(err)
	a.Equal("package resources\n", string(content))
}

func TestWriteFiles_AppendMode(t *testing.T) {
	a := assert.New(t)
	dir := t.TempDir()

	existing := "type MyInterface interface {\n\tFoo(ctx context.Context) error\n}\n"
	err := os.WriteFile(filepath.Join(dir, "iface.go"), []byte(existing), 0o600)
	a.NoError(err)

	files := map[string]GeneratedFile{
		"iface.go": {
			Path:    "iface.go",
			Content: "\tBar(ctx context.Context) error",
			Mode:    Append,
		},
	}

	manifest, err := WriteFiles(files, WriteOpts{OutputDir: dir})
	a.NoError(err)
	a.Len(manifest.Files, 1)
	a.Equal("updated", manifest.Files[0].Status)

	content, err := os.ReadFile(filepath.Join(dir, "iface.go"))
	a.NoError(err)
	a.Contains(string(content), "Bar(ctx context.Context) error")
	a.Contains(string(content), "Foo(ctx context.Context) error")
}

func TestWriteFiles_SkipMode(t *testing.T) {
	a := assert.New(t)
	dir := t.TempDir()

	files := map[string]GeneratedFile{
		"skip.go": {
			Path:    "skip.go",
			Content: "should not be written",
			Mode:    Skip,
		},
	}

	manifest, err := WriteFiles(files, WriteOpts{OutputDir: dir})
	a.NoError(err)
	a.Len(manifest.Files, 1)
	a.Equal("skipped", manifest.Files[0].Status)

	_, err = os.Stat(filepath.Join(dir, "skip.go"))
	a.True(os.IsNotExist(err))
}

func TestWriteFiles_CreateSkipsExistingWithoutForce(t *testing.T) {
	a := assert.New(t)
	dir := t.TempDir()

	originalContent := "original content"
	err := os.WriteFile(filepath.Join(dir, "existing.go"), []byte(originalContent), 0o600)
	a.NoError(err)

	files := map[string]GeneratedFile{
		"existing.go": {
			Path:    "existing.go",
			Content: "new content",
			Mode:    Create,
		},
	}

	manifest, err := WriteFiles(files, WriteOpts{OutputDir: dir})
	a.NoError(err)
	a.Len(manifest.Files, 1)
	a.Equal("skipped", manifest.Files[0].Status)

	content, err := os.ReadFile(filepath.Join(dir, "existing.go"))
	a.NoError(err)
	a.Equal(originalContent, string(content))
}

func TestWriteFiles_ForceOverwritesExisting(t *testing.T) {
	a := assert.New(t)
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "existing.go"), []byte("old"), 0o600)
	a.NoError(err)

	files := map[string]GeneratedFile{
		"existing.go": {
			Path:    "existing.go",
			Content: "new content",
			Mode:    Create,
		},
	}

	manifest, err := WriteFiles(files, WriteOpts{OutputDir: dir, Force: true})
	a.NoError(err)
	a.Len(manifest.Files, 1)
	a.Equal("created", manifest.Files[0].Status)

	content, err := os.ReadFile(filepath.Join(dir, "existing.go"))
	a.NoError(err)
	a.Equal("new content", string(content))
}

func TestWriteFiles_DryRunProducesNoFiles(t *testing.T) {
	a := assert.New(t)
	dir := t.TempDir()

	files := map[string]GeneratedFile{
		"new.go": {
			Path:    "new.go",
			Content: "package resources\n",
			Mode:    Create,
		},
	}

	manifest, err := WriteFiles(files, WriteOpts{OutputDir: dir, DryRun: true})
	a.NoError(err)
	a.Len(manifest.Files, 1)
	a.Equal("created", manifest.Files[0].Status)

	_, err = os.Stat(filepath.Join(dir, "new.go"))
	a.True(os.IsNotExist(err))
}

func TestWriteFiles_MissingOutputDirReturnsError(t *testing.T) {
	a := assert.New(t)

	files := map[string]GeneratedFile{
		"test.go": {
			Path:    "test.go",
			Content: "package resources\n",
			Mode:    Create,
		},
	}

	_, err := WriteFiles(files, WriteOpts{OutputDir: "/nonexistent/path"})
	a.Error(err)
	a.Contains(err.Error(), "output directory does not exist")
}

func TestWriteFiles_ManifestJSON(t *testing.T) {
	a := assert.New(t)

	manifest := &Manifest{
		Files: []ManifestEntry{
			{Path: "resources/svc-res.go", Status: "created"},
			{Path: "resources/svc.go", Status: "updated"},
			{Path: "resources/svc_mock_test.go", Status: "skipped"},
		},
	}

	data, err := json.Marshal(manifest)
	a.NoError(err)

	var decoded Manifest
	err = json.Unmarshal(data, &decoded)
	a.NoError(err)
	a.Len(decoded.Files, 3)
	a.Equal("resources/svc-res.go", decoded.Files[0].Path)
	a.Equal("created", decoded.Files[0].Status)
	a.Equal("updated", decoded.Files[1].Status)
	a.Equal("skipped", decoded.Files[2].Status)
}

func TestInsertBeforeLastClosingBrace(t *testing.T) {
	a := assert.New(t)

	existing := "type MyInterface interface {\n\tFoo() error\n}\n"
	newContent := "\tBar() error"

	result := insertBeforeLastClosingBrace(existing, newContent)
	a.Contains(result, "Foo() error")
	a.Contains(result, "Bar() error")
	a.Contains(result, "}")
}

func TestInsertBeforeLastClosingBrace_NoBrace(t *testing.T) {
	a := assert.New(t)

	existing := "package resources\n"
	newContent := "func Foo() {}"

	result := insertBeforeLastClosingBrace(existing, newContent)
	a.Contains(result, "package resources")
	a.Contains(result, "func Foo() {}")
}

func TestInsertBeforeVarDecl(t *testing.T) {
	a := assert.New(t)

	existing := `package resources

func (m *mockFooClient) ListFoos() {}

var testFooListerOpts = &nuke.ListerOpts{}
`
	newContent := `
func (m *mockFooClient) DeleteFoo() {}
`

	result := insertBeforeVarDecl(existing, newContent)
	a.Contains(result, "ListFoos")
	a.Contains(result, "DeleteFoo")
	a.Contains(result, "var testFooListerOpts")

	// DeleteFoo should appear before the var declaration
	deleteIdx := strings.Index(result, "DeleteFoo")
	varIdx := strings.Index(result, "var testFooListerOpts")
	a.Greater(varIdx, deleteIdx, "new method should be inserted before var declaration")
}

func TestInsertBeforeVarDecl_NoVar(t *testing.T) {
	a := assert.New(t)

	existing := "package resources\n\nfunc (m *mockFooClient) ListFoos() {}\n"
	newContent := "\nfunc (m *mockFooClient) DeleteFoo() {}\n"

	result := insertBeforeVarDecl(existing, newContent)
	a.Contains(result, "ListFoos")
	a.Contains(result, "DeleteFoo")
}
