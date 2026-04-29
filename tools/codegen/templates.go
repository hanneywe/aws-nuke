package main

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// templateFuncs returns the FuncMap used by all code generation templates.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"scope":                scopeFunc,
		"listerStruct":         listerStructFunc,
		"flatList":             flatListFunc,
		"nestedList":           nestedListFunc,
		"singletonList":        singletonListFunc,
		"resourceStruct":       resourceStructFunc,
		"filterMethod":         filterMethodFunc,
		"removeMethod":         removeMethodFunc,
		"stringMethod":         stringMethodFunc,
		"svcFieldType":         svcFieldTypeFunc,
		"ucfirst":              ucfirst,
		"lower":                strings.ToLower,
		"goIdiomatic":          goIdiomaticName,
		"needsSDKTypes":        needsSDKTypesFunc,
		"mockListSetup":        mockListSetupFunc,
		"mockListAssertField":  mockListAssertFieldFunc,
		"mockListEmptySetup":   mockListEmptySetupFunc,
		"exportedStringFields": exportedStringFieldsFunc,
		"stringFields":         stringFieldsFunc,
		"expectedString":       expectedStringFunc,
		"lcfirst":              lcfirst,
		"integInputVal":        integInputValFunc,
	}
}

func scopeFunc(scope string) string {
	switch scope {
	case "account":
		return "Account"
	case "region":
		return "Region"
	default:
		return ucfirst(scope)
	}
}

func ucfirst(input string) string {
	if input == "" {
		return input
	}
	runes := []rune(input)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// goIdiomaticName converts an AWS SDK field name to Go-idiomatic form.
// e.g. "UsagePlanId" -> "UsagePlanID", "UpstreamRegistryUrl" -> "UpstreamRegistryURL"
// Returns the original name if no conversion is needed.
func goIdiomaticName(name string) string {
	replacements := []struct {
		suffix      string
		replacement string
	}{
		{"Ids", "IDs"},
		{"Urls", "URLs"},
		{"Arns", "ARNs"},
		{"Id", "ID"},
		{"Url", "URL"},
		{"Arn", "ARN"},
		{"Api", "API"},
		{"Ip", "IP"},
		{"Ui", "UI"},
	}

	for _, rule := range replacements {
		if strings.HasSuffix(name, rule.suffix) {
			prefix := name[:len(name)-len(rule.suffix)]
			return goIdiomaticName(prefix) + rule.replacement
		}
	}
	return name
}

func lcfirst(input string) string {
	if input == "" {
		return input
	}
	runes := []rune(input)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// integInputValFunc converts an integration test input value to a Go expression.
// String values containing {{timestamp}} are converted to fmt.Sprintf with time.Now().UnixNano().
// Other string values become ptr.String("...").
// Boolean values become aws.Bool(...).
func integInputValFunc(value interface{}) string {
	switch typedValue := value.(type) {
	case string:
		if strings.Contains(typedValue, "{{timestamp}}") {
			replaced := strings.ReplaceAll(typedValue, "{{timestamp}}", "%d")
			return fmt.Sprintf("ptr.String(fmt.Sprintf(%q, time.Now().UnixNano()))", replaced)
		}
		return fmt.Sprintf("ptr.String(%q)", typedValue)
	case bool:
		return fmt.Sprintf("aws.Bool(%t)", typedValue)
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}

// svcFieldTypeFunc returns the type string for the svc field on the resource struct.
func svcFieldTypeFunc(definition *ResourceDef) string {
	if definition.SvcType == SvcTypeConcrete {
		return fmt.Sprintf("*%s.Client", definition.SDKPackage)
	}
	return resolveClientInterfaceName(definition)
}

// resolveClientInterfaceName returns the override from the DSL if set, otherwise
// derives the name from the SDK package.
func resolveClientInterfaceName(definition *ResourceDef) string {
	if definition.ClientInterfaceName != "" {
		return definition.ClientInterfaceName
	}
	return clientInterfaceName(definition.SDKPackage)
}

// clientInterfaceName returns the interface name for a given SDK package.
// e.g. "mediapackagev2" -> "MediaPackageV2Client", "eks" -> "EKSClient"
func clientInterfaceName(sdkPackage string) string {
	upper := strings.ToUpper(sdkPackage)
	switch upper {
	case "EKS", "EMR", "IAM", "EC2", "ECS", "RDS", "SNS", "SQS", "SES", "S3", "SSM", "KMS":
		return upper + "Client"
	case "IOT":
		return "IoTClient"
	}
	return ucfirst(sdkPackage) + "Client"
}

func listerStructFunc(definition *ResourceDef) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "type %sLister struct {\n", definition.ResourceName)
	if definition.SvcType != SvcTypeConcrete {
		fmt.Fprintf(&builder, "\tsvc %s\n", svcFieldTypeFunc(definition))
	}
	builder.WriteString("}")

	return builder.String()
}

func flatListFunc(definition *ResourceDef) string {
	var builder strings.Builder
	sdkPackage := definition.SDKPackage

	switch definition.List.Pagination {
	case "paginator":
		builder.WriteString(flatListPaginator(definition, sdkPackage))
	case "nextToken":
		builder.WriteString(flatListNextToken(definition, sdkPackage))
	case "none", "":
		builder.WriteString(flatListNoPagination(definition, sdkPackage))
	}

	return builder.String()
}

func flatListPaginator(definition *ResourceDef, sdkPackage string) string {
	var builder strings.Builder
	operation := definition.List.Operation

	fmt.Fprintf(&builder, "\n\tpaginator := %s.New%sPaginator(svc, &%s.%sInput{})\n", sdkPackage, operation, sdkPackage, operation)
	builder.WriteString("\n\tfor paginator.HasMorePages() {\n")
	builder.WriteString("\t\tresp, err := paginator.NextPage(ctx)\n")
	builder.WriteString("\t\tif err != nil {\n")
	builder.WriteString("\t\t\treturn nil, err\n")
	builder.WriteString("\t\t}\n")

	builder.WriteString(flatListItemLoop(definition, sdkPackage, "resp"))

	builder.WriteString("\t}\n")
	return builder.String()
}

func flatListNextToken(definition *ResourceDef, sdkPackage string) string {
	var builder strings.Builder
	operation := definition.List.Operation

	fmt.Fprintf(&builder, "\n\tparams := &%s.%sInput{}\n", sdkPackage, operation)
	builder.WriteString("\n\tfor {\n")
	fmt.Fprintf(&builder, "\t\tresp, err := svc.%s(ctx, params)\n", operation)
	builder.WriteString("\t\tif err != nil {\n")
	builder.WriteString("\t\t\treturn nil, err\n")
	builder.WriteString("\t\t}\n")

	builder.WriteString(flatListItemLoop(definition, sdkPackage, "resp"))

	builder.WriteString("\t\tif resp.NextToken == nil {\n")
	builder.WriteString("\t\t\tbreak\n")
	builder.WriteString("\t\t}\n")
	builder.WriteString("\t\tparams.NextToken = resp.NextToken\n")
	builder.WriteString("\t}\n")
	return builder.String()
}

func flatListNoPagination(definition *ResourceDef, sdkPackage string) string {
	var builder strings.Builder
	operation := definition.List.Operation

	fmt.Fprintf(&builder, "\n\tresp, err := svc.%s(ctx, &%s.%sInput{})\n", operation, sdkPackage, operation)
	builder.WriteString("\tif err != nil {\n")
	builder.WriteString("\t\treturn nil, err\n")
	builder.WriteString("\t}\n")

	builder.WriteString(flatListItemLoop(definition, sdkPackage, "resp"))

	return builder.String()
}

func flatListItemLoop(definition *ResourceDef, sdkPackage, responseVar string) string {
	var builder strings.Builder
	itemsField := definition.List.ItemsField
	iteratorVar := "item"

	fmt.Fprintf(&builder, "\n\t\tfor i := range %s.%s {\n", responseVar, itemsField)
	fmt.Fprintf(&builder, "\t\t\t%s := &%s.%s[i]\n", iteratorVar, responseVar, itemsField)

	if definition.List.Describe != nil {
		builder.WriteString(describeCall(definition, sdkPackage, iteratorVar))
	}

	if definition.List.Tags != nil {
		builder.WriteString(tagFetchCall(definition, sdkPackage, iteratorVar))
	}

	builder.WriteString(resourceAppend(definition, sdkPackage, iteratorVar))

	builder.WriteString("\t\t}\n")
	return builder.String()
}

func describeCall(definition *ResourceDef, sdkPackage, iteratorVar string) string {
	var builder strings.Builder
	describe := definition.List.Describe

	fmt.Fprintf(&builder, "\t\t\tdcResp, err := svc.%s(ctx, &%s.%sInput{\n", describe.Operation, sdkPackage, describe.Operation)
	for param, source := range describe.InputMapping {
		if strings.Contains(source, ".") {
			fmt.Fprintf(&builder, "\t\t\t\t%s: %s,\n", param, source)
		} else {
			fmt.Fprintf(&builder, "\t\t\t\t%s: %s.%s,\n", param, iteratorVar, source)
		}
	}
	builder.WriteString("\t\t\t})\n")
	builder.WriteString("\t\t\tif err != nil {\n")
	builder.WriteString("\t\t\t\treturn nil, err\n")
	builder.WriteString("\t\t\t}\n")

	return builder.String()
}

func tagFetchCall(definition *ResourceDef, sdkPackage, iteratorVar string) string {
	var builder strings.Builder
	tags := definition.List.Tags

	builder.WriteString("\t\t\tvar tags map[string]string\n")
	fmt.Fprintf(&builder, "\t\t\tif %s.%s != nil {\n", iteratorVar, tags.ArnField)
	fmt.Fprintf(&builder, "\t\t\t\ttagsResp, err := svc.%s(ctx, &%s.%sInput{\n", tags.Operation, sdkPackage, tags.Operation)
	fmt.Fprintf(&builder, "\t\t\t\t\tResourceArn: %s.%s,\n", iteratorVar, tags.ArnField)
	builder.WriteString("\t\t\t\t})\n")
	builder.WriteString("\t\t\t\tif err != nil {\n")
	builder.WriteString("\t\t\t\t\topts.Logger.Warnf(\"unable to fetch tags: %%s\", *" + iteratorVar + "." + tags.ArnField + ")\n")
	builder.WriteString("\t\t\t\t} else {\n")
	builder.WriteString("\t\t\t\t\ttags = tagsResp.Tags\n")
	builder.WriteString("\t\t\t\t}\n")
	builder.WriteString("\t\t\t}\n")

	return builder.String()
}

func resourceAppend(definition *ResourceDef, _, iteratorVar string) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "\t\t\tresources = append(resources, &%s{\n", definition.ResourceName)
	builder.WriteString("\t\t\t\tsvc: svc,\n")

	for _, field := range definition.Fields {
		goName := goIdiomaticName(field.Name)

		if field.FromTags {
			fmt.Fprintf(&builder, "\t\t\t\t%s: tags,\n", goName)
			continue
		}

		if field.FromDescribe != "" {
			describeResponse := "dcResp"
			if definition.List.Describe != nil && definition.List.Describe.ResponseField != "" {
				describeResponse = fmt.Sprintf("dcResp.%s", definition.List.Describe.ResponseField)
			}
			fmt.Fprintf(&builder, "\t\t\t\t%s: %s.%s,\n", goName, describeResponse, field.FromDescribe)
			continue
		}

		if field.FromList != "" {
			source := field.FromList
			if strings.Contains(source, ".") {
				fmt.Fprintf(&builder, "\t\t\t\t%s: %s,\n", goName, source)
			} else {
				fmt.Fprintf(&builder, "\t\t\t\t%s: %s.%s,\n", goName, iteratorVar, source)
			}
			continue
		}
	}

	builder.WriteString("\t\t\t})\n")
	return builder.String()
}

func nestedListFunc(definition *ResourceDef) string {
	var builder strings.Builder
	sdkPackage := definition.SDKPackage
	levels := definition.List.Levels

	builder.WriteString("\n")
	builder.WriteString(nestedLevel(definition, sdkPackage, levels, 0, 1))

	return builder.String()
}

type nestedLevelCtx struct {
	definition  *ResourceDef
	sdkPackage  string
	levels      []NestedLevelDef
	levelIndex  int
	indent      string
	indentBase  int
	iteratorVar string
	inputStr    string
	isLastLevel bool
}

func nestedLevel(definition *ResourceDef, sdkPackage string, levels []NestedLevelDef, levelIndex, indentBase int) string {
	level := levels[levelIndex]
	iteratorVar := level.IteratorVar

	inputStr := ""
	if len(level.ParentInputMapping) > 0 {
		var params []string
		for param, source := range level.ParentInputMapping {
			params = append(params, fmt.Sprintf("\t\t%s: %s,\n", param, source))
		}
		inputStr = strings.Join(params, "")
	}

	context := &nestedLevelCtx{
		definition:  definition,
		sdkPackage:  sdkPackage,
		levels:      levels,
		levelIndex:  levelIndex,
		indent:      strings.Repeat("\t", indentBase),
		indentBase:  indentBase,
		iteratorVar: iteratorVar,
		inputStr:    inputStr,
		isLastLevel: levelIndex == len(levels)-1,
	}

	switch level.Pagination {
	case "paginator":
		return nestedLevelPaginator(context)
	case "nextToken":
		return nestedLevelNextToken(context)
	default:
		return nestedLevelNoPagination(context)
	}
}

func nestedLevelChildContent(levelCtx *nestedLevelCtx, extraIndent int) string {
	if levelCtx.isLastLevel {
		return nestedResourceAppend(levelCtx.definition, levelCtx.sdkPackage, levelCtx.indent+strings.Repeat("\t", extraIndent))
	}
	return nestedLevel(levelCtx.definition, levelCtx.sdkPackage, levelCtx.levels, levelCtx.levelIndex+1, levelCtx.indentBase+extraIndent)
}

func nestedLevelPaginator(levelCtx *nestedLevelCtx) string {
	var builder strings.Builder
	level := levelCtx.levels[levelCtx.levelIndex]
	paginatorVar := fmt.Sprintf("%sPaginator", levelCtx.iteratorVar)

	builder.WriteString(fmt.Sprintf("%s%s := %s.New%sPaginator(svc, &%s.%sInput{\n", //nolint:staticcheck
		levelCtx.indent, paginatorVar, levelCtx.sdkPackage, level.Operation, levelCtx.sdkPackage, level.Operation))
	if levelCtx.inputStr != "" {
		builder.WriteString(levelCtx.inputStr)
	}
	fmt.Fprintf(&builder, "%s})\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%sfor %s.HasMorePages() {\n", levelCtx.indent, paginatorVar)

	responseVar := fmt.Sprintf("%sResp", levelCtx.iteratorVar)
	fmt.Fprintf(&builder, "%s\t%s, err := %s.NextPage(ctx)\n", levelCtx.indent, responseVar, paginatorVar)
	fmt.Fprintf(&builder, "%s\tif err != nil {\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s\t\treturn nil, err\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s\t}\n", levelCtx.indent)
	indexVar := "i" + goIdiomaticName(ucfirst(levelCtx.iteratorVar))
	fmt.Fprintf(&builder, "%s\tfor %s := range %s.%s {\n", levelCtx.indent, indexVar, responseVar, level.ItemsField)
	fmt.Fprintf(&builder, "%s\t\t%s := &%s.%s[%s]\n", levelCtx.indent, levelCtx.iteratorVar, responseVar, level.ItemsField, indexVar)

	builder.WriteString(nestedLevelChildContent(levelCtx, 2))

	fmt.Fprintf(&builder, "%s\t}\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s}\n", levelCtx.indent)
	return builder.String()
}

func nestedLevelNextToken(levelCtx *nestedLevelCtx) string {
	var builder strings.Builder
	level := levelCtx.levels[levelCtx.levelIndex]
	paramsVar := fmt.Sprintf("%sParams", levelCtx.iteratorVar)

	fmt.Fprintf(&builder, "%s%s := &%s.%sInput{\n", levelCtx.indent, paramsVar, levelCtx.sdkPackage, level.Operation)
	if levelCtx.inputStr != "" {
		builder.WriteString(levelCtx.inputStr)
	}
	fmt.Fprintf(&builder, "%s}\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%sfor {\n", levelCtx.indent)

	responseVar := fmt.Sprintf("%sResp", levelCtx.iteratorVar)
	fmt.Fprintf(&builder, "%s\t%s, err := svc.%s(ctx, %s)\n", levelCtx.indent, responseVar, level.Operation, paramsVar)
	fmt.Fprintf(&builder, "%s\tif err != nil {\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s\t\treturn nil, err\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s\t}\n", levelCtx.indent)
	indexVar := "i" + goIdiomaticName(ucfirst(levelCtx.iteratorVar))
	fmt.Fprintf(&builder, "%s\tfor %s := range %s.%s {\n", levelCtx.indent, indexVar, responseVar, level.ItemsField)
	fmt.Fprintf(&builder, "%s\t\t%s := &%s.%s[%s]\n", levelCtx.indent, levelCtx.iteratorVar, responseVar, level.ItemsField, indexVar)

	builder.WriteString(nestedLevelChildContent(levelCtx, 2))

	fmt.Fprintf(&builder, "%s\t}\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s\tif %s.NextToken == nil {\n", levelCtx.indent, responseVar)
	fmt.Fprintf(&builder, "%s\t\tbreak\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s\t}\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s\t%s.NextToken = %s.NextToken\n", levelCtx.indent, paramsVar, responseVar)
	fmt.Fprintf(&builder, "%s}\n", levelCtx.indent)
	return builder.String()
}

func nestedLevelNoPagination(levelCtx *nestedLevelCtx) string {
	var builder strings.Builder
	level := levelCtx.levels[levelCtx.levelIndex]
	responseVar := fmt.Sprintf("%sResp", levelCtx.iteratorVar)

	fmt.Fprintf(&builder, "%s%s, err := svc.%s(ctx, &%s.%sInput{\n",
		levelCtx.indent, responseVar, level.Operation,
		levelCtx.sdkPackage, level.Operation)
	if levelCtx.inputStr != "" {
		builder.WriteString(levelCtx.inputStr)
	}
	fmt.Fprintf(&builder, "%s})\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%sif err != nil {\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s\treturn nil, err\n", levelCtx.indent)
	fmt.Fprintf(&builder, "%s}\n", levelCtx.indent)
	indexVar := "i" + goIdiomaticName(ucfirst(levelCtx.iteratorVar))
	fmt.Fprintf(&builder, "%sfor %s := range %s.%s {\n", levelCtx.indent, indexVar, responseVar, level.ItemsField)
	fmt.Fprintf(&builder, "%s\t%s := &%s.%s[%s]\n", levelCtx.indent, levelCtx.iteratorVar, responseVar, level.ItemsField, indexVar)

	builder.WriteString(nestedLevelChildContent(levelCtx, 1))

	fmt.Fprintf(&builder, "%s}\n", levelCtx.indent)
	return builder.String()
}

func nestedResourceAppend(definition *ResourceDef, _, indent string) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "%sresources = append(resources, &%s{\n", indent, definition.ResourceName)
	fmt.Fprintf(&builder, "%s\tsvc: svc,\n", indent)

	for _, field := range definition.Fields {
		goName := goIdiomaticName(field.Name)

		if field.FromTags {
			fmt.Fprintf(&builder, "%s\t%s: tags,\n", indent, goName)
			continue
		}
		if field.FromList != "" {
			fmt.Fprintf(&builder, "%s\t%s: %s,\n", indent, goName, field.FromList)
			continue
		}
		if field.FromDescribe != "" {
			fmt.Fprintf(&builder, "%s\t%s: dcResp.%s,\n", indent, goName, field.FromDescribe)
			continue
		}
	}

	fmt.Fprintf(&builder, "%s})\n", indent)
	return builder.String()
}

func singletonListFunc(definition *ResourceDef) string {
	var builder strings.Builder
	sdkPackage := definition.SDKPackage
	operation := definition.List.Operation

	fmt.Fprintf(&builder, "\n\tresp, err := svc.%s(ctx, &%s.%sInput{})\n", operation, sdkPackage, operation)
	builder.WriteString("\tif err != nil {\n")
	builder.WriteString("\t\treturn nil, err\n")
	builder.WriteString("\t}\n")

	if definition.List.NilCheck && definition.List.ResponseField != "" {
		fmt.Fprintf(&builder, "\n\tif resp.%s == nil {\n", definition.List.ResponseField)
		builder.WriteString("\t\treturn nil, nil\n")
		builder.WriteString("\t}\n")
	}

	fmt.Fprintf(&builder, "\n\tresources = append(resources, &%s{\n", definition.ResourceName)
	builder.WriteString("\t\tsvc: svc,\n")

	for _, field := range definition.Fields {
		if field.FromList != "" {
			source := field.FromList
			if definition.List.ResponseField != "" && !strings.Contains(source, ".") {
				source = fmt.Sprintf("resp.%s.%s", definition.List.ResponseField, source)
			} else {
				source = fmt.Sprintf("resp.%s", source)
			}
			fmt.Fprintf(&builder, "\t\t%s: %s,\n", field.Name, source)
			continue
		}
	}

	builder.WriteString("\t})\n")
	return builder.String()
}

func resourceStructFunc(definition *ResourceDef) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "type %s struct {\n", definition.ResourceName)

	// svc field
	fmt.Fprintf(&builder, "\tsvc %s\n", svcFieldTypeFunc(definition))

	// Exported fields first
	for _, field := range definition.Fields {
		if field.Exported != nil && !*field.Exported {
			continue
		}

		goName := goIdiomaticName(field.Name)
		tag := ""
		if field.PropertyTag != "" {
			tag = fmt.Sprintf(" `property:%q`", field.PropertyTag)
		} else if goName != field.Name {
			tag = fmt.Sprintf(" `property:\"name=%s\"`", field.Name)
		}
		fmt.Fprintf(&builder, "\t%s %s%s\n", goName, field.Type, tag)
	}

	// Settings and protection fields
	if len(definition.Settings) > 0 {
		builder.WriteString("\tsettings *libsettings.Setting\n")
	}

	// Unexported fields
	for _, field := range definition.Fields {
		if field.Exported != nil && !*field.Exported {
			tag := ""
			if field.PropertyTag != "" {
				tag = fmt.Sprintf(" `property:\"%s\"`", field.PropertyTag) //nolint:gocritic
			}
			fmt.Fprintf(&builder, "\t%s %s%s\n", field.Name, field.Type, tag)
		}
	}

	builder.WriteString("}")
	return builder.String()
}

func filterMethodFunc(definition *ResourceDef) string {
	if definition.FilterOverride != "" {
		var builder strings.Builder
		fmt.Fprintf(&builder, "func (r *%s) Filter() error {\n", definition.ResourceName)
		builder.WriteString(definition.FilterOverride)
		builder.WriteString("\n}")
		return builder.String()
	}

	fieldTypes := make(map[string]string)
	for _, field := range definition.Fields {
		fieldTypes[field.Name] = field.Type
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "func (r *%s) Filter() error {\n", definition.ResourceName)

	for _, filter := range definition.Filters {
		goName := goIdiomaticName(filter.Field)
		isPointer := strings.HasPrefix(fieldTypes[filter.Field], "*")

		switch filter.Operator {
		case "equals":
			for _, value := range filter.Values {
				if isPointer {
					fmt.Fprintf(&builder, "\tif r.%s != nil && *r.%s == \"%s\" {\n", goName, goName, value)
				} else {
					fmt.Fprintf(&builder, "\tif string(r.%s) == \"%s\" {\n", goName, value)
				}
				fmt.Fprintf(&builder, "\t\treturn fmt.Errorf(\"%s\")\n", filter.Message)
				builder.WriteString("\t}\n")
			}
		case "boolFalse":
			fmt.Fprintf(&builder, "\tif r.%s != nil && !*r.%s {\n", goName, goName)
			fmt.Fprintf(&builder, "\t\treturn fmt.Errorf(\"%s\")\n", filter.Message)
			builder.WriteString("\t}\n")
		case "boolTrue":
			fmt.Fprintf(&builder, "\tif r.%s != nil && *r.%s {\n", goName, goName)
			fmt.Fprintf(&builder, "\t\treturn fmt.Errorf(\"%s\")\n", filter.Message)
			builder.WriteString("\t}\n")
		}
	}

	builder.WriteString("\treturn nil\n")
	builder.WriteString("}")
	return builder.String()
}

func removeMethodFunc(definition *ResourceDef) string {
	var builder strings.Builder
	sdkPackage := definition.SDKPackage

	// Settings-based deletion protection
	for _, setting := range definition.Settings {
		if setting.ProtectionField == "" || setting.DisableOperation == "" {
			continue
		}

		fmt.Fprintf(&builder, "\tif ptr.ToBool(r.%s) && r.settings.GetBool(\"%s\") {\n", goIdiomaticName(setting.ProtectionField), setting.Name)
		fmt.Fprintf(&builder, "\t\t_, err := r.svc.%s(ctx, &%s.%sInput{\n", setting.DisableOperation, sdkPackage, setting.DisableOperation)

		for param, value := range setting.DisableInput {
			switch typedValue := value.(type) {
			case string:
				fmt.Fprintf(&builder, "\t\t\t%s: r.%s,\n", param, goIdiomaticName(typedValue))
			case bool:
				fmt.Fprintf(&builder, "\t\t\t%s: aws.Bool(%t),\n", param, typedValue)
			}
		}

		builder.WriteString("\t\t})\n")
		builder.WriteString("\t\tif err != nil {\n")
		builder.WriteString("\t\t\treturn err\n")
		builder.WriteString("\t\t}\n")
		builder.WriteString("\t}\n\n")
	}

	// Pre-deletion steps
	for i := range definition.PreDeletion {
		preDeletion := &definition.PreDeletion[i]

		switch preDeletion.Type {
		case "listAndBatchDelete":
			builder.WriteString(preDeletionListAndBatchDelete(definition, sdkPackage, preDeletion))
		case "conditional":
			builder.WriteString(preDeletionConditional(definition, sdkPackage, preDeletion))
		}
	}

	// Main delete call
	fmt.Fprintf(&builder, "\t_, err := r.svc.%s(ctx, &%s.%sInput{\n", definition.Delete.Operation, sdkPackage, definition.Delete.Operation)
	for _, inputField := range definition.Delete.InputFields {
		fmt.Fprintf(&builder, "\t\t%s: r.%s,\n", inputField, goIdiomaticName(inputField))
	}
	builder.WriteString("\t})\n")
	builder.WriteString("\treturn err")

	return builder.String()
}

func preDeletionListAndBatchDelete(definition *ResourceDef, sdkPackage string, preDeletion *PreDeletionDef) string {
	var builder strings.Builder
	typesAlias := definition.SDKPackage + "types"

	fmt.Fprintf(&builder, "\tlistResp, err := r.svc.%s(ctx, &%s.%sInput{\n", preDeletion.ListOperation, sdkPackage, preDeletion.ListOperation)
	for param, fieldName := range preDeletion.ListInput {
		fmt.Fprintf(&builder, "\t\t%s: r.%s,\n", param, fieldName)
	}
	builder.WriteString("\t})\n")
	fmt.Fprintf(&builder, "\tif err == nil && len(listResp.%s) > 0 {\n", preDeletion.ListItemsField)

	itemType := singularize(preDeletion.DeleteItemsField)
	fmt.Fprintf(&builder, "\t\tvar items []%s.%s\n", typesAlias, itemType)
	fmt.Fprintf(&builder, "\t\tfor _, target := range listResp.%s {\n", preDeletion.ListItemsField)
	fmt.Fprintf(&builder, "\t\t\titems = append(items, %s.%s{\n", typesAlias, itemType)
	for targetField, sourceField := range preDeletion.ItemMapping {
		fmt.Fprintf(&builder, "\t\t\t\t%s: target.%s,\n", targetField, sourceField)
	}
	builder.WriteString("\t\t\t})\n")
	builder.WriteString("\t\t}\n")

	fmt.Fprintf(&builder, "\t\t_, err = r.svc.%s(ctx, &%s.%sInput{\n", preDeletion.DeleteOperation, sdkPackage, preDeletion.DeleteOperation)
	for param, fieldName := range preDeletion.DeleteInput {
		fmt.Fprintf(&builder, "\t\t\t%s: r.%s,\n", param, fieldName)
	}
	fmt.Fprintf(&builder, "\t\t\t%s: items,\n", preDeletion.DeleteItemsField)
	builder.WriteString("\t\t})\n")
	builder.WriteString("\t\tif err != nil {\n")
	builder.WriteString("\t\t\treturn err\n")
	builder.WriteString("\t\t}\n")
	builder.WriteString("\t}\n\n")

	return builder.String()
}

func preDeletionConditional(_ *ResourceDef, sdkPackage string, preDeletion *PreDeletionDef) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "\tif %s {\n", preDeletion.Condition)
	fmt.Fprintf(&builder, "\t\t_, err := r.svc.%s(ctx, &%s.%sInput{\n", preDeletion.Operation, sdkPackage, preDeletion.Operation)
	for param, fieldName := range preDeletion.Input {
		fmt.Fprintf(&builder, "\t\t\t%s: r.%s,\n", param, fieldName)
	}
	builder.WriteString("\t\t})\n")
	builder.WriteString("\t\tif err != nil {\n")
	builder.WriteString("\t\t\treturn err\n")
	builder.WriteString("\t\t}\n")
	builder.WriteString("\t}\n\n")

	return builder.String()
}

// singularize converts a plural word to singular form.
func singularize(word string) string {
	if strings.HasSuffix(word, "ies") {
		return word[:len(word)-3] + "y"
	}
	if strings.HasSuffix(word, "s") {
		return word[:len(word)-1]
	}
	return word
}

// resolveItemType returns the explicit itemType. Falls back to singularize only for
// pre-deletion fields where the type is derived from the delete items field name.
func resolveItemType(itemType, itemsField string) string {
	if itemType != "" {
		return itemType
	}
	return itemsField
}

func stringMethodFunc(definition *ResourceDef) string {
	var builder strings.Builder
	stringRep := definition.StringRep

	fmt.Fprintf(&builder, "func (r *%s) String() string {\n", definition.ResourceName)

	if stringRep.Override != "" {
		builder.WriteString(stringRep.Override)
		builder.WriteString("\n}")
		return builder.String()
	}

	if stringRep.Conditional != nil {
		conditional := stringRep.Conditional

		fmt.Fprintf(&builder, "\tif r.%s != nil {\n", goIdiomaticName(conditional.Field))

		if conditional.IfNotNil != nil && conditional.IfNotNil.Format != "" {
			formatArgs := make([]string, len(conditional.IfNotNil.Fields))
			for i, fieldName := range conditional.IfNotNil.Fields {
				formatArgs[i] = fmt.Sprintf("*r.%s", goIdiomaticName(fieldName))
			}
			fmt.Fprintf(&builder, "\t\treturn fmt.Sprintf(\"%s\", %s)\n", conditional.IfNotNil.Format, strings.Join(formatArgs, ", "))
		} else if conditional.IfNotNil != nil && conditional.IfNotNil.Field != "" {
			fmt.Fprintf(&builder, "\t\treturn *r.%s\n", goIdiomaticName(conditional.IfNotNil.Field))
		} else {
			fmt.Fprintf(&builder, "\t\treturn *r.%s\n", goIdiomaticName(conditional.Field))
		}

		builder.WriteString("\t}\n")
		fmt.Fprintf(&builder, "\treturn *r.%s\n", goIdiomaticName(conditional.IfNil))
	} else if stringRep.Format != "" && len(stringRep.Fields) > 0 {
		formatArgs := make([]string, len(stringRep.Fields))
		for i, fieldName := range stringRep.Fields {
			formatArgs[i] = fmt.Sprintf("*r.%s", goIdiomaticName(fieldName))
		}
		fmt.Fprintf(&builder, "\treturn fmt.Sprintf(\"%s\", %s)\n", stringRep.Format, strings.Join(formatArgs, ", "))
	} else if stringRep.Field != "" {
		fmt.Fprintf(&builder, "\treturn *r.%s\n", goIdiomaticName(stringRep.Field))
	}

	builder.WriteString("}")
	return builder.String()
}

// needsSDKTypesFunc returns true if the mock test template needs to import SDK types.
func needsSDKTypesFunc(data *MockTestTemplateData) bool {
	definition := data.ResourceDef

	if definition.List.Strategy == StrategyNested && len(definition.List.Levels) > 0 {
		return true
	}
	if definition.List.Strategy == StrategyFlat && definition.List.ItemsField != "" {
		return true
	}
	if definition.List.Strategy == StrategySingleton && definition.List.ResponseField != "" {
		return true
	}
	if definition.List.Describe != nil && definition.List.Describe.ResponseField != "" {
		return true
	}
	return false
}

// mockListSetupFunc generates the mock expectations for the List test.
func mockListSetupFunc(data *MockTestTemplateData) string {
	definition := data.ResourceDef
	var builder strings.Builder
	sdkPackage := definition.SDKPackage

	switch definition.List.Strategy {
	case StrategyFlat:
		builder.WriteString(mockFlatListSetup(definition, sdkPackage))
	case StrategyNested:
		builder.WriteString(mockNestedListSetup(definition, sdkPackage))
	case StrategySingleton:
		builder.WriteString(mockSingletonListSetup(definition, sdkPackage))
	}

	return builder.String()
}

func mockFlatListSetup(definition *ResourceDef, sdkPackage string) string {
	var builder strings.Builder
	operation := definition.List.Operation
	itemsField := definition.List.ItemsField
	typesAlias := sdkPackage + "types"
	itemType := resolveItemType(definition.List.ItemType, itemsField)

	fmt.Fprintf(&builder, "\n\tmockClient.On(\"%s\", mock.Anything, mock.Anything).\n", operation)
	fmt.Fprintf(&builder, "\t\tReturn(&%s.%sOutput{\n", sdkPackage, operation)
	fmt.Fprintf(&builder, "\t\t\t%s: []%s.%s{\n", itemsField, typesAlias, itemType)
	fmt.Fprintf(&builder, "\t\t\t\t{")

	first := true
	for _, field := range definition.Fields {
		if field.FromList == "" || strings.Contains(field.FromList, ".") {
			continue
		}
		if !first {
			builder.WriteString(", ")
		}
		first = false
		fmt.Fprintf(&builder, "%s: ptr.String(\"test-value\")", field.FromList)
	}

	builder.WriteString("},\n")
	fmt.Fprintf(&builder, "\t\t\t},\n")
	fmt.Fprintf(&builder, "\t\t}, nil)\n")

	if definition.List.Describe != nil {
		describe := definition.List.Describe

		fmt.Fprintf(&builder, "\n\tmockClient.On(\"%s\", mock.Anything, mock.Anything).\n", describe.Operation)
		fmt.Fprintf(&builder, "\t\tReturn(&%s.%sOutput{", sdkPackage, describe.Operation)

		if describe.ResponseField != "" {
			fmt.Fprintf(&builder, "\n\t\t\t%s: &%s.%s{\n", describe.ResponseField, typesAlias, describe.ResponseField)
			for _, field := range definition.Fields {
				if field.FromDescribe != "" {
					switch field.Type {
					case TypeStringPtr:
						fmt.Fprintf(&builder, "\t\t\t\t%s: ptr.String(\"test-%s\"),\n", field.FromDescribe, strings.ToLower(field.Name))
					case "map[string]string":
						fmt.Fprintf(&builder, "\t\t\t\t%s: map[string]string{\"key\": \"value\"},\n", field.FromDescribe)
					}
				}
			}
			builder.WriteString("\t\t\t},\n\t\t")
		}
		builder.WriteString("}, nil)\n")
	}

	if definition.List.Tags != nil {
		tags := definition.List.Tags

		fmt.Fprintf(&builder, "\n\tmockClient.On(\"%s\", mock.Anything, mock.Anything).\n", tags.Operation)
		fmt.Fprintf(&builder, "\t\tReturn(&%s.%sOutput{\n", sdkPackage, tags.Operation)
		fmt.Fprintf(&builder, "\t\t\tTags: map[string]string{\"key\": \"value\"},\n")
		fmt.Fprintf(&builder, "\t\t}, nil)\n")
	}

	return builder.String()
}

func mockNestedListSetup(definition *ResourceDef, sdkPackage string) string {
	var builder strings.Builder
	typesAlias := sdkPackage + "types"

	for levelIndex, level := range definition.List.Levels {
		itemType := resolveItemType(level.ItemType, level.ItemsField)

		fmt.Fprintf(&builder, "\n\tmockClient.On(\"%s\", mock.Anything, mock.Anything).\n", level.Operation)
		fmt.Fprintf(&builder, "\t\tReturn(&%s.%sOutput{\n", sdkPackage, level.Operation)
		fmt.Fprintf(&builder, "\t\t\t%s: []%s.%s{\n", level.ItemsField, typesAlias, itemType)
		fmt.Fprintf(&builder, "\t\t\t\t{%s},\n", nestedLevelFields(definition, levelIndex))
		fmt.Fprintf(&builder, "\t\t\t},\n")
		fmt.Fprintf(&builder, "\t\t}, nil)\n")
	}

	return builder.String()
}

// nestedLevelFields builds the struct field initializers for a nested level's mock item.
func nestedLevelFields(definition *ResourceDef, levelIndex int) string {
	var builder strings.Builder
	level := definition.List.Levels[levelIndex]
	first := true

	for _, field := range definition.Fields {
		if field.FromList == "" {
			continue
		}

		parts := strings.SplitN(field.FromList, ".", 2)
		if len(parts) == 2 && parts[0] == level.IteratorVar {
			if !first {
				builder.WriteString(", ")
			}
			first = false
			fmt.Fprintf(&builder, "%s: ptr.String(\"test-%s\")", parts[1], strings.ToLower(field.Name))
		}
	}

	if levelIndex < len(definition.List.Levels)-1 {
		addParentMappingFields(&builder, definition, levelIndex, first)
	}

	return builder.String()
}

// addParentMappingFields adds fields from child levels' parentInputMapping that aren't already present.
func addParentMappingFields(builder *strings.Builder, definition *ResourceDef, levelIndex int, first bool) bool {
	level := definition.List.Levels[levelIndex]

	for childIndex := levelIndex + 1; childIndex < len(definition.List.Levels); childIndex++ {
		for _, source := range definition.List.Levels[childIndex].ParentInputMapping {
			parts := strings.SplitN(source, ".", 2)
			if len(parts) != 2 || parts[0] != level.IteratorVar {
				continue
			}

			fieldName := parts[1]
			if isFieldAlreadyMapped(definition, source) {
				continue
			}

			if !first {
				builder.WriteString(", ")
			}
			first = false
			fmt.Fprintf(builder, "%s: ptr.String(\"test-%s\")", fieldName, strings.ToLower(fieldName))
		}
	}

	return first
}

// isFieldAlreadyMapped checks if a fromList source is already mapped by a field definition.
func isFieldAlreadyMapped(definition *ResourceDef, source string) bool {
	for _, field := range definition.Fields {
		if field.FromList == source {
			return true
		}
	}
	return false
}

func mockSingletonListSetup(definition *ResourceDef, sdkPackage string) string {
	var builder strings.Builder
	operation := definition.List.Operation

	fmt.Fprintf(&builder, "\n\tmockClient.On(\"%s\", mock.Anything, mock.Anything).\n", operation)
	fmt.Fprintf(&builder, "\t\tReturn(&%s.%sOutput{", sdkPackage, operation)

	if definition.List.ResponseField != "" {
		typesAlias := sdkPackage + "types"

		fmt.Fprintf(&builder, "\n\t\t\t%s: &%s.%s{\n", definition.List.ResponseField, typesAlias, definition.List.ResponseField)
		for _, field := range definition.Fields {
			if field.FromList != "" {
				if field.Type == TypeStringPtr {
					fmt.Fprintf(&builder, "\t\t\t\t%s: ptr.String(\"test-value\"),\n", field.FromList)
				}
			}
		}
		builder.WriteString("\t\t\t},\n\t\t")
	} else {
		builder.WriteString("\n")
		for _, field := range definition.Fields {
			if field.FromList != "" {
				if field.Type == TypeStringPtr {
					fmt.Fprintf(&builder, "\t\t\t%s: ptr.String(\"test-value\"),\n", field.FromList)
				}
			}
		}
		builder.WriteString("\t\t")
	}

	builder.WriteString("}, nil)\n")

	return builder.String()
}

// mockListAssertFieldFunc returns the assertion expression for the first resource field.
func mockListAssertFieldFunc(data *MockTestTemplateData) string {
	definition := data.ResourceDef

	for _, field := range definition.Fields {
		if field.Type == TypeStringPtr && (field.Exported == nil || *field.Exported) {
			return fmt.Sprintf("*r.%s", goIdiomaticName(field.Name))
		}
	}
	return "r.String()"
}

// mockListEmptySetupFunc generates mock expectations for the empty list test.
func mockListEmptySetupFunc(data *MockTestTemplateData) string {
	definition := data.ResourceDef
	var builder strings.Builder
	sdkPackage := definition.SDKPackage

	switch definition.List.Strategy {
	case StrategyFlat:
		operation := definition.List.Operation
		itemsField := definition.List.ItemsField
		typesAlias := sdkPackage + "types"
		itemType := resolveItemType(definition.List.ItemType, itemsField)

		fmt.Fprintf(&builder, "\n\tmockClient.On(\"%s\", mock.Anything, mock.Anything).\n", operation)
		fmt.Fprintf(&builder, "\t\tReturn(&%s.%sOutput{\n", sdkPackage, operation)
		fmt.Fprintf(&builder, "\t\t\t%s: []%s.%s{},\n", itemsField, typesAlias, itemType)
		fmt.Fprintf(&builder, "\t\t}, nil)\n")

	case StrategyNested:
		level := definition.List.Levels[0]
		typesAlias := sdkPackage + "types"
		itemType := resolveItemType(level.ItemType, level.ItemsField)

		fmt.Fprintf(&builder, "\n\tmockClient.On(\"%s\", mock.Anything, mock.Anything).\n", level.Operation)
		fmt.Fprintf(&builder, "\t\tReturn(&%s.%sOutput{\n", sdkPackage, level.Operation)
		fmt.Fprintf(&builder, "\t\t\t%s: []%s.%s{},\n", level.ItemsField, typesAlias, itemType)
		fmt.Fprintf(&builder, "\t\t}, nil)\n")

	case StrategySingleton:
		operation := definition.List.Operation
		fmt.Fprintf(&builder, "\n\tmockClient.On(\"%s\", mock.Anything, mock.Anything).\n", operation)
		fmt.Fprintf(&builder, "\t\tReturn(&%s.%sOutput{}, nil)\n", sdkPackage, operation)
	}

	return builder.String()
}

// exportedStringFieldsFunc returns exported *string fields for the Properties test.
func exportedStringFieldsFunc(data *MockTestTemplateData) []FieldDef {
	definition := data.ResourceDef
	var fields []FieldDef

	for _, field := range definition.Fields {
		if field.Exported != nil && !*field.Exported {
			continue
		}
		if field.Type == TypeStringPtr {
			fields = append(fields, field)
		}
	}
	return fields
}

// stringFieldsFunc returns the fields needed for the String test.
func stringFieldsFunc(data *MockTestTemplateData) []FieldDef {
	stringRep := data.StringRep
	neededNames := stringFieldNames(&stringRep)
	return filterFieldsByName(data.Fields, neededNames)
}

// stringFieldNames returns the set of field names referenced by the string representation.
func stringFieldNames(stringRep *StringRepDef) map[string]bool {
	seen := make(map[string]bool)

	if stringRep.Conditional != nil {
		seen[stringRep.Conditional.Field] = true
		seen[stringRep.Conditional.IfNil] = true
		if stringRep.Conditional.IfNotNil != nil {
			for _, fieldName := range stringRep.Conditional.IfNotNil.Fields {
				seen[fieldName] = true
			}
		}
	} else if stringRep.Format != "" && len(stringRep.Fields) > 0 {
		for _, fieldName := range stringRep.Fields {
			seen[fieldName] = true
		}
	} else if stringRep.Field != "" {
		seen[stringRep.Field] = true
	}

	return seen
}

// filterFieldsByName returns *string fields whose names are in the needed set.
func filterFieldsByName(fields []FieldDef, needed map[string]bool) []FieldDef {
	var result []FieldDef
	for _, field := range fields {
		if needed[field.Name] && field.Type == TypeStringPtr {
			result = append(result, field)
		}
	}
	return result
}

// expectedStringFunc returns the expected string value for the String test.
func expectedStringFunc(data *MockTestTemplateData) string {
	definition := data.ResourceDef
	stringRep := definition.StringRep

	if stringRep.Conditional != nil {
		if stringRep.Conditional.IfNotNil != nil && stringRep.Conditional.IfNotNil.Format != "" {
			formatArgs := make([]string, len(stringRep.Conditional.IfNotNil.Fields))
			for i, fieldName := range stringRep.Conditional.IfNotNil.Fields {
				formatArgs[i] = "test-" + strings.ToLower(fieldName)
			}
			return `fmt.Sprintf("` + stringRep.Conditional.IfNotNil.Format + `", ` + quoteJoin(formatArgs) + `)`
		}
		return `"test-` + strings.ToLower(stringRep.Conditional.Field) + `"`
	}

	if stringRep.Format != "" && len(stringRep.Fields) > 0 {
		formatArgs := make([]string, len(stringRep.Fields))
		for i, fieldName := range stringRep.Fields {
			formatArgs[i] = "test-" + strings.ToLower(fieldName)
		}
		return `fmt.Sprintf("` + stringRep.Format + `", ` + quoteJoin(formatArgs) + `)`
	}

	if stringRep.Field != "" {
		return `"test-` + strings.ToLower(stringRep.Field) + `"`
	}

	return `""`
}

func quoteJoin(values []string) string {
	quoted := make([]string, len(values))
	for i, value := range values {
		quoted[i] = fmt.Sprintf("%q", value)
	}
	return strings.Join(quoted, ", ")
}
