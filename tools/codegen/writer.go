package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteOpts controls how files are written to disk.
type WriteOpts struct {
	DryRun    bool
	Force     bool
	OutputDir string
}

// Manifest records the outcome of a file-writing operation.
type Manifest struct {
	Files []ManifestEntry `json:"files"`
}

// ManifestEntry records the path and status of a single generated file.
type ManifestEntry struct {
	Path   string `json:"path"`
	Status string `json:"status"` // "created", "updated", "skipped"
}

// WriteFiles writes generated files to disk, respecting the write mode,
// --dry-run, and --no-force flags. It returns a Manifest describing what happened.
func WriteFiles(files map[string]GeneratedFile, opts WriteOpts) (*Manifest, error) {
	if info, err := os.Stat(opts.OutputDir); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("output directory does not exist: %s", opts.OutputDir)
	}

	manifest := &Manifest{}

	for _, generatedFile := range files {
		entry, err := writeFile(generatedFile, opts)
		if err != nil {
			return nil, err
		}
		manifest.Files = append(manifest.Files, entry)
	}

	return manifest, nil
}

func writeFile(generatedFile GeneratedFile, opts WriteOpts) (ManifestEntry, error) {
	fullPath := filepath.Join(opts.OutputDir, generatedFile.Path)

	switch generatedFile.Mode {
	case Create:
		return writeCreate(generatedFile, fullPath, opts)
	case Append:
		return writeAppend(generatedFile, fullPath, opts)
	case AppendBeforeVar:
		return writeAppendBeforeVar(generatedFile, fullPath, opts)
	case Skip:
		fmt.Fprintf(os.Stderr, "warning: skipping %s\n", generatedFile.Path)
		return ManifestEntry{Path: generatedFile.Path, Status: "skipped"}, nil
	default:
		return ManifestEntry{}, fmt.Errorf("unknown file mode for %s", generatedFile.Path)
	}
}

func writeCreate(generatedFile GeneratedFile, fullPath string, opts WriteOpts) (ManifestEntry, error) {
	if opts.DryRun {
		fmt.Fprintf(os.Stdout, "--- %s (create) ---\n%s\n", generatedFile.Path, generatedFile.Content)
		return ManifestEntry{Path: generatedFile.Path, Status: "created"}, nil
	}

	if _, err := os.Stat(fullPath); err == nil {
		if !opts.Force {
			fmt.Fprintf(os.Stderr, "warning: %s already exists, skipping (remove --no-force to overwrite)\n", generatedFile.Path)
			return ManifestEntry{Path: generatedFile.Path, Status: "skipped"}, nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return ManifestEntry{}, fmt.Errorf("creating directory for %s: %w", generatedFile.Path, err)
	}

	if err := os.WriteFile(fullPath, []byte(generatedFile.Content), 0o600); err != nil {
		return ManifestEntry{}, fmt.Errorf("writing %s: %w", generatedFile.Path, err)
	}

	return ManifestEntry{Path: generatedFile.Path, Status: "created"}, nil
}

func writeAppend(generatedFile GeneratedFile, fullPath string, opts WriteOpts) (ManifestEntry, error) {
	existingContent, err := os.ReadFile(fullPath)
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("reading existing file %s for append: %w", generatedFile.Path, err)
	}

	merged := insertBeforeLastClosingBrace(string(existingContent), generatedFile.Content)

	if opts.DryRun {
		fmt.Fprintf(os.Stdout, "--- %s (append) ---\n%s\n", generatedFile.Path, merged)
		return ManifestEntry{Path: generatedFile.Path, Status: "updated"}, nil
	}

	if err := os.WriteFile(fullPath, []byte(merged), 0o600); err != nil { //nolint:gosec
		return ManifestEntry{}, fmt.Errorf("writing %s: %w", generatedFile.Path, err)
	}

	return ManifestEntry{Path: generatedFile.Path, Status: "updated"}, nil
}

// insertBeforeLastClosingBrace finds the last `}` in the existing content
// (which closes the interface or mock struct) and inserts the new content before it.
func insertBeforeLastClosingBrace(existingContent, newContent string) string {
	braceIndex := strings.LastIndex(existingContent, "}")
	if braceIndex < 0 {
		return existingContent + "\n" + newContent
	}

	before := existingContent[:braceIndex]
	after := existingContent[braceIndex:]
	trimmed := strings.TrimRight(before, " \t\n")

	return trimmed + "\n" + newContent + "\n" + after
}

func writeAppendBeforeVar(generatedFile GeneratedFile, fullPath string, opts WriteOpts) (ManifestEntry, error) {
	existingContent, err := os.ReadFile(fullPath)
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("reading existing file %s for append: %w", generatedFile.Path, err)
	}

	merged := insertBeforeVarDecl(string(existingContent), generatedFile.Content)

	if opts.DryRun {
		fmt.Fprintf(os.Stdout, "--- %s (append) ---\n%s\n", generatedFile.Path, merged)
		return ManifestEntry{Path: generatedFile.Path, Status: "updated"}, nil
	}

	if err := os.WriteFile(fullPath, []byte(merged), 0o600); err != nil { //nolint:gosec
		return ManifestEntry{}, fmt.Errorf("writing %s: %w", generatedFile.Path, err)
	}

	return ManifestEntry{Path: generatedFile.Path, Status: "updated"}, nil
}

// insertBeforeVarDecl finds the last `var ` declaration in the existing content
// and inserts the new content before it. This is used for mock files where
// new receiver methods should be added before the var testXxxListerOpts line.
func insertBeforeVarDecl(existingContent, newContent string) string {
	varIndex := strings.LastIndex(existingContent, "\nvar ")
	if varIndex < 0 {
		return strings.TrimRight(existingContent, " \t\n") + "\n" + newContent
	}

	before := existingContent[:varIndex]
	after := existingContent[varIndex:]
	trimmed := strings.TrimRight(before, " \t\n")

	return trimmed + "\n" + newContent + after
}
