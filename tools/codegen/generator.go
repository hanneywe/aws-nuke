package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// FileMode indicates how a generated file should be written.
type FileMode int

const (
	// Create indicates a new file should be created (error if exists with --no-force).
	Create FileMode = iota
	// Append indicates content should be appended before the last closing brace (for interfaces).
	Append
	// AppendBeforeVar indicates content should be appended before the var declaration (for mock files).
	AppendBeforeVar
	// Skip indicates the file should not be written.
	Skip
)

// String constants for list strategies and svc types.
const (
	StrategyFlat      = "flat"
	StrategyNested    = "nested"
	StrategySingleton = "singleton"
	SvcTypeConcrete   = "concrete"
	TypeStringPtr     = "*string"
)

// GeneratedFile holds the path, content, and write mode for a generated file.
type GeneratedFile struct {
	Path    string
	Content string
	Mode    FileMode
}

// GenerateOpts controls optional behavior during code generation.
type GenerateOpts struct {
	ForceMockTests        bool
	ForceIntegrationTests bool
}

// GenerateAll produces all file contents from a ResourceDef.
// outputDir is the base directory where resources/ files live (e.g., the repo root).
// Returns a map of file path -> GeneratedFile.
func GenerateAll(definition *ResourceDef, outputDir string, opts GenerateOpts) (map[string]GeneratedFile, error) {
	templates, err := template.New("").Funcs(templateFuncs()).ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	resourceFile := filepath.Join("resources", fmt.Sprintf("%s-%s.go", definition.Service, definition.Resource))

	if definition.SvcType == SvcTypeConcrete {
		return generateConcreteResource(templates, definition, resourceFile)
	}

	return generateInterfaceResource(templates, definition, outputDir, resourceFile, opts)
}

// generateConcreteResource generates only the resource file for concrete svc types.
func generateConcreteResource(
	templates *template.Template, definition *ResourceDef, resourceFile string,
) (map[string]GeneratedFile, error) {
	files := make(map[string]GeneratedFile)

	content, err := executeTemplate(templates, "resource.go.tmpl", definition)
	if err != nil {
		return nil, fmt.Errorf("executing resource template: %w", err)
	}

	files[resourceFile] = GeneratedFile{Path: resourceFile, Content: content, Mode: Create}
	return files, nil
}

// generateInterfaceResource generates all files for resources using
// an interface-based svc type.
//
//nolint:cyclop
func generateInterfaceResource(
	templates *template.Template, definition *ResourceDef,
	outputDir, resourceFile string, opts GenerateOpts,
) (map[string]GeneratedFile, error) {
	files := make(map[string]GeneratedFile)
	requiredMethods := computeRequiredMethods(definition)

	interfaceFile := filepath.Join("resources", definition.Service+".go")
	mockFile := filepath.Join("resources", definition.Service+"_mock_test.go")
	mockTestFile := filepath.Join(
		"resources", definition.Service+"-"+definition.Resource+"_mock_test.go",
	)

	interfaceName := resolveClientInterfaceName(definition)
	mockStructName := "mock" + interfaceName

	var testOptsVar string
	if definition.ClientInterfaceName != "" {
		root := strings.TrimSuffix(definition.ClientInterfaceName, "Client")
		testOptsVar = fmt.Sprintf("test%sListerOpts", root)
	} else {
		testOptsVar = fmt.Sprintf("test%sListerOpts", ucfirst(definition.SDKPackage))
	}

	// Client interface file
	interfaceFullPath := filepath.Join(outputDir, interfaceFile)
	err := generateInterfaceFile(
		templates, definition, files, interfaceFile, interfaceFullPath,
		interfaceName, requiredMethods,
	)
	if err != nil {
		return nil, err
	}

	// Mock file
	mockFullPath := filepath.Join(outputDir, mockFile)
	err = generateMockFile(
		templates, definition, files, mockFile, mockFullPath,
		mockStructName, testOptsVar, requiredMethods,
	)
	if err != nil {
		return nil, err
	}

	// Resource file (always Create)
	content, err := executeTemplate(templates, "resource.go.tmpl", definition)
	if err != nil {
		return nil, fmt.Errorf("executing resource template: %w", err)
	}
	files[resourceFile] = GeneratedFile{Path: resourceFile, Content: content, Mode: Create}

	// Mock test file (only generated if it does not already exist, unless --force-mock-tests)
	mockTestFullPath := filepath.Join(outputDir, mockTestFile)
	if opts.ForceMockTests || os.IsNotExist(statError(mockTestFullPath)) {
		mockTestData := &MockTestTemplateData{
			ResourceDef:       definition,
			MockStructName:    mockStructName,
			TestListerOptsVar: testOptsVar,
			SDKPackage:        definition.SDKPackage,
			Methods:           requiredMethods,
		}

		content, err = executeTemplate(templates, "mock_test.go.tmpl", mockTestData)
		if err != nil {
			return nil, fmt.Errorf("executing mock_test template: %w", err)
		}
		files[mockTestFile] = GeneratedFile{Path: mockTestFile, Content: content, Mode: Create}
	}

	// Integration test file (only generated if defined in DSL and does not already exist, unless --force-integration-tests)
	if definition.IntegrationTest != nil {
		integTestFile := filepath.Join("resources", fmt.Sprintf("%s-%s_test.go", definition.Service, definition.Resource))
		integTestFullPath := filepath.Join(outputDir, integTestFile)

		if opts.ForceIntegrationTests || os.IsNotExist(statError(integTestFullPath)) {
			content, err = executeTemplate(templates, "integration_test.go.tmpl", definition)
			if err != nil {
				return nil, fmt.Errorf("executing integration_test template: %w", err)
			}
			files[integTestFile] = GeneratedFile{Path: integTestFile, Content: content, Mode: Create}
		}
	}

	return files, nil
}

// generateInterfaceFile generates or appends to the client interface file.
func generateInterfaceFile(
	templates *template.Template, definition *ResourceDef,
	files map[string]GeneratedFile,
	relPath, fullPath, interfaceName string,
	requiredMethods []string,
) error {
	existingMethods, err := parseExistingMethods(fullPath, true)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("parsing existing interface file %s: %w", relPath, err)
	}

	data := &InterfaceTemplateData{
		InterfaceName: interfaceName,
		SDKPackage:    definition.SDKPackage,
	}

	if len(existingMethods) == 0 {
		data.Methods = requiredMethods

		content, err := executeTemplate(templates, "client_interface.go.tmpl", data)
		if err != nil {
			return fmt.Errorf("executing client_interface template: %w", err)
		}
		files[relPath] = GeneratedFile{Path: relPath, Content: content, Mode: Create}
	} else if newMethods := filterNewMethods(requiredMethods, existingMethods); len(newMethods) > 0 {
		data.Methods = newMethods

		content, err := executeTemplate(templates, "client_interface_methods.go.tmpl", data)
		if err != nil {
			return fmt.Errorf("executing client_interface_methods template: %w", err)
		}
		files[relPath] = GeneratedFile{Path: relPath, Content: content, Mode: Append}
	}

	return nil
}

// generateMockFile generates or appends to the mock file.
func generateMockFile(
	templates *template.Template, definition *ResourceDef,
	files map[string]GeneratedFile,
	relPath, fullPath, mockStructName, testOptsVar string,
	requiredMethods []string,
) error {
	existingMethods, err := parseExistingMethods(fullPath, false)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("parsing existing mock file %s: %w", relPath, err)
	}

	data := &MockTemplateData{
		MockStructName:    mockStructName,
		TestListerOptsVar: testOptsVar,
		SDKPackage:        definition.SDKPackage,
	}

	if len(existingMethods) == 0 {
		data.Methods = requiredMethods

		content, err := executeTemplate(templates, "mock.go.tmpl", data)
		if err != nil {
			return fmt.Errorf("executing mock template: %w", err)
		}
		files[relPath] = GeneratedFile{Path: relPath, Content: content, Mode: Create}
	} else if newMethods := filterNewMethods(requiredMethods, existingMethods); len(newMethods) > 0 {
		data.Methods = newMethods

		content, err := executeTemplate(templates, "mock_methods.go.tmpl", data)
		if err != nil {
			return fmt.Errorf("executing mock_methods template: %w", err)
		}
		files[relPath] = GeneratedFile{Path: relPath, Content: content, Mode: AppendBeforeVar}
	}

	return nil
}

// computeRequiredMethods determines all SDK methods needed based on the DSL definition.
func computeRequiredMethods(definition *ResourceDef) []string {
	seen := make(map[string]bool)
	var methods []string

	addMethod := func(name string) {
		if name != "" && !seen[name] {
			seen[name] = true
			methods = append(methods, name)
		}
	}

	switch definition.List.Strategy {
	case StrategyFlat:
		addMethod(definition.List.Operation)
	case StrategyNested:
		for _, level := range definition.List.Levels {
			addMethod(level.Operation)
		}
	case StrategySingleton:
		addMethod(definition.List.Operation)
	}

	if definition.List.Describe != nil {
		addMethod(definition.List.Describe.Operation)
	}
	if definition.List.Tags != nil {
		addMethod(definition.List.Tags.Operation)
	}

	addMethod(definition.Delete.Operation)

	for i := range definition.PreDeletion {
		addMethod(definition.PreDeletion[i].ListOperation)
		addMethod(definition.PreDeletion[i].DeleteOperation)
		addMethod(definition.PreDeletion[i].Operation)
	}

	for i := range definition.Settings {
		addMethod(definition.Settings[i].DisableOperation)
	}

	return methods
}

// executeTemplate executes a named template and returns the rendered content.
func executeTemplate(templates *template.Template, name string, data interface{}) (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, name, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// statError returns the error from os.Stat, discarding the FileInfo.
func statError(path string) error {
	_, err := os.Stat(path)
	return err
}
