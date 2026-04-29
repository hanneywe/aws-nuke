package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type cliFlags struct {
	dryRun                bool
	noForce               bool
	validate              bool
	outputManifest        bool
	noBuild               bool
	noCommit              bool
	forceMockTests        bool
	forceIntegrationTests bool
	outputDir             string
}

func main() {
	flags, yamlPath, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := run(flags, yamlPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() (*cliFlags, string, error) {
	flags := &cliFlags{}

	flag.BoolVar(&flags.dryRun, "dry-run", false, "Print output without writing files")
	flag.BoolVar(&flags.noForce, "no-force", false, "Do not overwrite existing resource files")
	flag.BoolVar(&flags.validate, "validate", false, "Check DSL only, no generation")
	flag.BoolVar(&flags.outputManifest, "output-manifest", false, "Write JSON manifest of generated files")
	flag.BoolVar(&flags.noBuild, "no-build", false, "Skip running goimports, go build, golangci-lint after generation")
	flag.BoolVar(&flags.noCommit, "no-commit", false, "Skip staging and committing generated files")
	flag.BoolVar(&flags.forceMockTests, "force-mock-tests", false, "Regenerate mock test files even if they already exist")
	flag.BoolVar(&flags.forceIntegrationTests, "force-integration-tests", false,
		"Regenerate integration test files even if they already exist")
	flag.StringVar(&flags.outputDir, "output-dir", ".", "Override output directory (default: current directory)")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: codegen [flags] <yaml-dsl-file>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return nil, "", fmt.Errorf("exactly one YAML DSL file path is required")
	}

	return flags, args[0], nil
}

func run(flags *cliFlags, yamlPath string) error {
	definition, err := Parse(yamlPath)
	if err != nil {
		return err
	}

	validationErrors := Validate(definition)

	if flags.validate {
		printValidationResult(validationErrors)
		if len(validationErrors) > 0 {
			return fmt.Errorf("validation failed")
		}
		return nil
	}

	if len(validationErrors) > 0 {
		printValidationErrors(validationErrors)
		return fmt.Errorf("validation failed")
	}

	generatedFiles, err := GenerateAll(definition, flags.outputDir, GenerateOpts{
		ForceMockTests:        flags.forceMockTests,
		ForceIntegrationTests: flags.forceIntegrationTests,
	})
	if err != nil {
		return fmt.Errorf("generating files: %w", err)
	}

	manifest, err := WriteFiles(generatedFiles, WriteOpts{
		DryRun:    flags.dryRun,
		Force:     !flags.noForce,
		OutputDir: flags.outputDir,
	})
	if err != nil {
		return fmt.Errorf("writing files: %w", err)
	}

	if flags.outputManifest {
		if err := printManifest(manifest); err != nil {
			return fmt.Errorf("writing manifest: %w", err)
		}
	}

	if err := runPostAndCommit(flags, &DefaultExecutor{}, manifest, definition.Service, definition.ResourceName); err != nil {
		return err
	}

	if !flags.dryRun {
		printChecklist()
	}

	return nil
}

func runPostAndCommit(flags *cliFlags, executor CommandExecutor, manifest *Manifest, serviceName, resourceName string) error {
	ctx := context.Background()

	if !flags.noBuild {
		fmt.Fprintln(os.Stderr, "Running build steps...")
		if err := RunPostSteps(ctx, executor); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Build steps passed.")
	}

	if !flags.noCommit {
		fmt.Fprintln(os.Stderr, "Committing generated files...")
		filePaths := collectFilePaths(manifest)
		if err := CommitFiles(ctx, executor, filePaths, serviceName, resourceName); err != nil {
			return fmt.Errorf("commit failed: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Committed successfully.")
	}

	return nil
}

func printValidationResult(validationErrors []ValidationError) {
	type jsonResult struct {
		Valid  bool              `json:"valid"`
		Errors []ValidationError `json:"errors"`
	}

	result := jsonResult{
		Valid:  len(validationErrors) == 0,
		Errors: validationErrors,
	}
	if result.Errors == nil {
		result.Errors = []ValidationError{}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(result)
}

func printValidationErrors(validationErrors []ValidationError) {
	fmt.Fprintln(os.Stderr, "Validation errors:")
	for _, validationErr := range validationErrors {
		fmt.Fprintf(os.Stderr, "  - %s: %s\n", validationErr.Field, validationErr.Message)
	}
}

func printManifest(manifest *Manifest) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(manifest)
}

func collectFilePaths(manifest *Manifest) []string {
	var paths []string
	for _, entry := range manifest.Files {
		if entry.Status != "skipped" {
			paths = append(paths, entry.Path)
		}
	}
	return paths
}

func printChecklist() {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Post-generation checklist:")
	fmt.Fprintln(os.Stdout, "  1. go build ./resources/...")
	fmt.Fprintln(os.Stdout, "  2. go fmt ./resources/...")
	fmt.Fprintln(os.Stdout, "  3. golangci-lint run ./resources/...")
}
