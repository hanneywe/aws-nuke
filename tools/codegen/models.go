package main

// ResourceDef is the top-level structure parsed from a YAML DSL file.
type ResourceDef struct {
	Service      string `yaml:"service"`
	Resource     string `yaml:"resource"`
	ResourceName string `yaml:"resourceName"`
	Scope        string `yaml:"scope"`
	SDKPackage   string `yaml:"sdkPackage"`
	SvcType      string `yaml:"svcType,omitempty"`
	// ClientInterfaceName optionally overrides the derived client interface name.
	// When set, the mock struct is "mock"+ClientInterfaceName and the test lister
	// opts var is "test"+strings.TrimSuffix(ClientInterfaceName,"Client")+"ListerOpts".
	// Use this when the existing interface file (resources/<service>.go) already
	// declares an interface whose name doesn't match the name derived from sdkPackage.
	ClientInterfaceName string `yaml:"clientInterfaceName,omitempty"`

	List            ListDef          `yaml:"list"`
	Delete          DeleteDef        `yaml:"delete"`
	Fields          []FieldDef       `yaml:"fields"`
	Filters         []FilterDef      `yaml:"filters,omitempty"`
	Settings        []SettingDef     `yaml:"settings,omitempty"`
	Dependencies    []string         `yaml:"dependencies,omitempty"`
	PreDeletion     []PreDeletionDef `yaml:"preDeletion,omitempty"`
	StringRep       StringRepDef     `yaml:"stringRepresentation"`
	ExtraImports    []string         `yaml:"extraImports,omitempty"`
	FilterOverride  string           `yaml:"filterOverride,omitempty"`
	IntegrationTest *IntegTestDef    `yaml:"integrationTest,omitempty"`
}

// ListDef describes how to list resources.
type ListDef struct {
	Strategy      string           `yaml:"strategy"`
	Operation     string           `yaml:"operation,omitempty"`
	Pagination    string           `yaml:"pagination,omitempty"`
	ItemsField    string           `yaml:"itemsField,omitempty"`
	ItemType      string           `yaml:"itemType,omitempty"`
	Describe      *DescribeDef     `yaml:"describe,omitempty"`
	Tags          *TagFetchDef     `yaml:"tags,omitempty"`
	Levels        []NestedLevelDef `yaml:"levels,omitempty"`
	ResponseField string           `yaml:"responseField,omitempty"`
	NilCheck      bool             `yaml:"nilCheck,omitempty"`
	Override      string           `yaml:"override,omitempty"`
}

// NestedLevelDef describes one level in a nested list strategy.
type NestedLevelDef struct {
	Operation          string            `yaml:"operation"`
	Pagination         string            `yaml:"pagination"`
	ItemsField         string            `yaml:"itemsField"`
	ItemType           string            `yaml:"itemType,omitempty"`
	IteratorVar        string            `yaml:"iteratorVar"`
	ParentInputMapping map[string]string `yaml:"parentInputMapping,omitempty"`
}

// DescribeDef describes an optional describe call to enrich listed items.
type DescribeDef struct {
	Operation     string            `yaml:"operation"`
	InputMapping  map[string]string `yaml:"inputMapping"`
	ResponseField string            `yaml:"responseField,omitempty"`
}

// TagFetchDef describes an optional tag-fetching call.
type TagFetchDef struct {
	Operation string `yaml:"operation"`
	ArnField  string `yaml:"arnField"`
}

// DeleteDef describes how to delete a resource.
type DeleteDef struct {
	Operation   string   `yaml:"operation"`
	InputFields []string `yaml:"inputFields"`
	Override    string   `yaml:"override,omitempty"`
}

// FieldDef describes a single field on the resource struct.
type FieldDef struct {
	Name         string `yaml:"name"`
	Type         string `yaml:"type"`
	FromList     string `yaml:"fromList,omitempty"`
	FromDescribe string `yaml:"fromDescribe,omitempty"`
	FromTags     bool   `yaml:"fromTags,omitempty"`
	Exported     *bool  `yaml:"exported,omitempty"`
	PropertyTag  string `yaml:"propertyTag,omitempty"`
}

// FilterDef describes a filter condition for excluding resources from deletion.
type FilterDef struct {
	Field    string   `yaml:"field"`
	Operator string   `yaml:"operator"`
	Values   []string `yaml:"values,omitempty"`
	Message  string   `yaml:"message"`
}

// SettingDef describes a user-configurable setting for the resource.
type SettingDef struct {
	Name             string                 `yaml:"name"`
	ProtectionField  string                 `yaml:"protectionField,omitempty"`
	DisableOperation string                 `yaml:"disableOperation,omitempty"`
	DisableInput     map[string]interface{} `yaml:"disableInput,omitempty"`
}

// PreDeletionDef describes a pre-deletion step that runs before the main delete call.
type PreDeletionDef struct {
	Type             string            `yaml:"type"`
	ListOperation    string            `yaml:"listOperation,omitempty"`
	ListInput        map[string]string `yaml:"listInput,omitempty"`
	ListItemsField   string            `yaml:"listItemsField,omitempty"`
	DeleteOperation  string            `yaml:"deleteOperation,omitempty"`
	DeleteInput      map[string]string `yaml:"deleteInput,omitempty"`
	DeleteItemsField string            `yaml:"deleteItemsField,omitempty"`
	ItemMapping      map[string]string `yaml:"itemMapping,omitempty"`
	Condition        string            `yaml:"condition,omitempty"`
	Operation        string            `yaml:"operation,omitempty"`
	Input            map[string]string `yaml:"input,omitempty"`
}

// StringRepDef describes how the resource's String() method should be generated.
type StringRepDef struct {
	Field       string          `yaml:"field,omitempty"`
	Format      string          `yaml:"format,omitempty"`
	Fields      []string        `yaml:"fields,omitempty"`
	Conditional *ConditionalStr `yaml:"conditional,omitempty"`
	Override    string          `yaml:"override,omitempty"`
}

// ConditionalStr describes a conditional string representation with nil-check logic.
type ConditionalStr struct {
	Field    string        `yaml:"field"`
	IfNil    string        `yaml:"ifNil"`
	IfNotNil *StringRepDef `yaml:"ifNotNil,omitempty"`
}

// IntegTestDef describes integration test configuration.
type IntegTestDef struct {
	Create *IntegCreateDef `yaml:"create,omitempty"`
}

// IntegCreateDef describes the create operation for integration test setup.
type IntegCreateDef struct {
	Operation string                 `yaml:"operation"`
	Input     map[string]interface{} `yaml:"input"`
}

// MockTestTemplateData holds the data passed to the mock_test.go.tmpl template.
// It embeds ResourceDef and adds computed metadata fields.
type MockTestTemplateData struct {
	*ResourceDef
	MockStructName    string
	TestListerOptsVar string
	SDKPackage        string
	Methods           []string
}

// MockTemplateData holds the data passed to the mock.go.tmpl template.
type MockTemplateData struct {
	MockStructName    string
	TestListerOptsVar string
	SDKPackage        string
	Methods           []string
}

// InterfaceTemplateData holds the data passed to the client_interface.go.tmpl template.
type InterfaceTemplateData struct {
	InterfaceName string
	SDKPackage    string
	Methods       []string
}
