package main

import "fmt"

// ValidationError represents a single validation issue with a field path and message.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate checks a ResourceDef for completeness and consistency.
// Returns a list of ValidationError with field paths and messages.
func Validate(definition *ResourceDef) []ValidationError {
	var errors []ValidationError

	// Required top-level fields
	if definition.Service == "" {
		errors = append(errors, ValidationError{Field: "service", Message: "service is required"})
	}
	if definition.Resource == "" {
		errors = append(errors, ValidationError{Field: "resource", Message: "resource is required"})
	}
	if definition.ResourceName == "" {
		errors = append(errors, ValidationError{Field: "resourceName", Message: "resourceName is required"})
	}
	if definition.Scope == "" {
		errors = append(errors, ValidationError{Field: "scope", Message: "scope is required"})
	}
	if definition.SDKPackage == "" {
		errors = append(errors, ValidationError{Field: "sdkPackage", Message: "sdkPackage is required"})
	}
	if len(definition.Fields) == 0 {
		errors = append(errors, ValidationError{Field: "fields", Message: "fields must not be empty"})
	}

	// stringRepresentation must have at least one of field, format, or conditional
	if definition.StringRep.Field == "" && definition.StringRep.Format == "" && definition.StringRep.Conditional == nil {
		errors = append(errors, ValidationError{Field: "stringRepresentation", Message: "stringRepresentation is required"})
	}

	// delete validation
	if definition.Delete.Operation == "" {
		errors = append(errors, ValidationError{Field: "delete.operation", Message: "delete.operation is required"})
	}
	if len(definition.Delete.InputFields) == 0 {
		errors = append(errors, ValidationError{Field: "delete.inputFields", Message: "delete.inputFields is required"})
	}

	// list.strategy validation
	errors = append(errors, validateListStrategy(definition)...)

	// override blocks validation
	errors = append(errors, validateOverrides(definition)...)

	return errors
}

func validateListStrategy(definition *ResourceDef) []ValidationError {
	var errors []ValidationError

	switch definition.List.Strategy {
	case StrategyFlat:
		errors = append(errors, validateFlatList(definition)...)
	case StrategyNested:
		errors = append(errors, validateNestedList(definition)...)
	case StrategySingleton:
		if definition.List.Operation == "" {
			errors = append(errors, ValidationError{Field: "list.operation", Message: "list.operation is required for singleton strategy"})
		}
	case "":
		errors = append(errors, ValidationError{Field: "list.strategy", Message: "list.strategy is required"})
	default:
		errors = append(errors, ValidationError{Field: "list.strategy", Message: "list.strategy must be flat, nested, or singleton"})
	}

	return errors
}

func validateFlatList(definition *ResourceDef) []ValidationError {
	var errors []ValidationError
	if definition.List.Operation == "" {
		errors = append(errors, ValidationError{Field: "list.operation", Message: "list.operation is required for flat strategy"})
	}
	if definition.List.Pagination == "" {
		errors = append(errors, ValidationError{Field: "list.pagination", Message: "list.pagination is required for flat strategy"})
	}
	hasItemsNoType := definition.List.ItemsField != "" && definition.List.ItemType == ""
	notOverridden := definition.List.Override == "" && definition.SvcType != SvcTypeConcrete
	if hasItemsNoType && notOverridden {
		errors = append(errors, ValidationError{Field: "list.itemType", Message: "list.itemType is required when itemsField is set"})
	}
	return errors
}

func validateNestedList(definition *ResourceDef) []ValidationError {
	var errors []ValidationError
	if len(definition.List.Levels) < 2 {
		errors = append(errors, ValidationError{Field: "list.levels", Message: "nested strategy requires at least 2 levels"})
	}
	if definition.List.Override == "" && definition.SvcType != SvcTypeConcrete {
		for i, level := range definition.List.Levels {
			if level.ItemsField != "" && level.ItemType == "" {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("list.levels[%d].itemType", i),
					Message: fmt.Sprintf("itemType is required for level %d when itemsField is set", i),
				})
			}
		}
	}
	return errors
}

func validateOverrides(definition *ResourceDef) []ValidationError {
	var errors []ValidationError

	if definition.Delete.Override != "" && !hasBalancedBraces(definition.Delete.Override) {
		errors = append(errors, ValidationError{Field: "delete.override", Message: "override block has unbalanced braces"})
	}
	if definition.List.Override != "" && !hasBalancedBraces(definition.List.Override) {
		errors = append(errors, ValidationError{Field: "list.override", Message: "override block has unbalanced braces"})
	}
	if definition.FilterOverride != "" && !hasBalancedBraces(definition.FilterOverride) {
		errors = append(errors, ValidationError{Field: "filterOverride", Message: "override block has unbalanced braces"})
	}
	if definition.StringRep.Override != "" && !hasBalancedBraces(definition.StringRep.Override) {
		errors = append(errors, ValidationError{Field: "stringRepresentation.override", Message: "override block has unbalanced braces"})
	}

	return errors
}

func hasBalancedBraces(code string) bool {
	depth := 0
	for _, char := range code {
		switch char {
		case '{':
			depth++
		case '}':
			depth--
			if depth < 0 {
				return false
			}
		}
	}
	return depth == 0
}
