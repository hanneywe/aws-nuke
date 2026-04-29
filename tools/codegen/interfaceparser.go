package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// parseExistingMethods parses a Go source file and extracts method names.
// If isInterface is true, it extracts interface method names; otherwise it extracts receiver method names.
func parseExistingMethods(path string, isInterface bool) ([]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	fileSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing Go file %s: %w", path, err)
	}

	if isInterface {
		return extractInterfaceMethods(parsedFile), nil
	}
	return extractReceiverMethods(parsedFile), nil
}

// extractInterfaceMethods walks the AST to find all method names declared in interfaces.
func extractInterfaceMethods(parsedFile *ast.File) []string {
	seen := make(map[string]bool)
	var methods []string

	for _, declaration := range parsedFile.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			ifaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			for _, method := range ifaceType.Methods.List {
				for _, methodName := range method.Names {
					if !seen[methodName.Name] {
						seen[methodName.Name] = true
						methods = append(methods, methodName.Name)
					}
				}
			}
		}
	}

	return methods
}

// extractReceiverMethods walks the AST to find all method names on receiver types (mock methods).
func extractReceiverMethods(parsedFile *ast.File) []string {
	seen := make(map[string]bool)
	var methods []string

	for _, declaration := range parsedFile.Decls {
		funcDecl, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil {
			continue
		}

		methodName := funcDecl.Name.Name
		if !seen[methodName] {
			seen[methodName] = true
			methods = append(methods, methodName)
		}
	}

	return methods
}

// filterNewMethods returns methods from required that are not in existing.
func filterNewMethods(required, existing []string) []string {
	existingSet := make(map[string]bool, len(existing))
	for _, method := range existing {
		existingSet[method] = true
	}

	var newMethods []string
	for _, method := range required {
		if !existingSet[method] {
			newMethods = append(newMethods, method)
		}
	}
	return newMethods
}
