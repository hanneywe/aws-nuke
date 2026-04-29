package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// minimalFlatDef returns a minimal flat-list ResourceDef for testing.
func minimalFlatDef() *ResourceDef {
	return &ResourceDef{
		Service:      "mailmanager",
		Resource:     "relay",
		ResourceName: "MailManagerRelay",
		Scope:        "account",
		SDKPackage:   "mailmanager",
		List: ListDef{
			Strategy:   "flat",
			Operation:  "ListRelays",
			Pagination: "paginator",
			ItemsField: "Relays",
			ItemType:   "Relay",
		},
		Delete: DeleteDef{
			Operation:   "DeleteRelay",
			InputFields: []string{"RelayId"},
		},
		Fields: []FieldDef{
			{Name: "RelayId", Type: "*string", FromList: "RelayId"},
			{Name: "RelayName", Type: "*string", FromList: "RelayName"},
		},
		StringRep: StringRepDef{Field: "RelayName"},
	}
}

func TestGenerateAll_FlatListProducesExpectedFiles(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	// Should produce 4 files: resource, interface, mock, mock_test
	a.Len(files, 4)

	expectedPaths := []string{
		"resources/mailmanager-relay.go",
		"resources/mailmanager.go",
		"resources/mailmanager_mock_test.go",
		"resources/mailmanager-relay_mock_test.go",
	}
	for _, p := range expectedPaths {
		_, ok := files[p]
		a.True(ok, "expected file %s not found in generated files", p)
	}

	// Resource file should be Create mode
	a.Equal(Create, files["resources/mailmanager-relay.go"].Mode)
	// Interface file should be Create mode (no existing file)
	a.Equal(Create, files["resources/mailmanager.go"].Mode)
	// Mock file should be Create mode (no existing file)
	a.Equal(Create, files["resources/mailmanager_mock_test.go"].Mode)
	// Mock test file should be Create mode
	a.Equal(Create, files["resources/mailmanager-relay_mock_test.go"].Mode)
}

func TestGenerateAll_FlatListWithIntegrationTestProduces5Files(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	def.IntegrationTest = &IntegTestDef{
		Create: &IntegCreateDef{
			Operation: "CreateRelay",
			Input:     map[string]interface{}{"RelayName": "test-relay"},
		},
	}

	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)
	a.Len(files, 5)

	_, ok := files["resources/mailmanager-relay_test.go"]
	a.True(ok, "integration test file not found")
	a.Equal(Create, files["resources/mailmanager-relay_test.go"].Mode)
}

func TestGenerateAll_ExistingServiceAppendsToInterfaceAndMock(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	// Create existing interface file with one method already declared
	resourcesDir := filepath.Join(outputDir, "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	a.NoError(err)

	existingInterface := `package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"
)

type MailmanagerClient interface {
	ListRelays(ctx context.Context, params *mailmanager.ListRelaysInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.ListRelaysOutput, error)
}
`
	err = os.WriteFile(filepath.Join(resourcesDir, "mailmanager.go"), []byte(existingInterface), 0600)
	a.NoError(err)

	// Create existing mock file with one method already declared
	existingMock := `package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockMailmanagerClient struct {
	mock.Mock
}

func (m *mockMailmanagerClient) ListRelays(
	ctx context.Context, params *mailmanager.ListRelaysInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.ListRelaysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.ListRelaysOutput), args.Error(1)
}

var testMailmanagerListerOpts = &nuke.ListerOpts{}
`
	err = os.WriteFile(filepath.Join(resourcesDir, "mailmanager_mock_test.go"), []byte(existingMock), 0600)
	a.NoError(err)

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	// Interface file should be Append mode since ListRelays already exists but DeleteRelay is new
	ifaceFile, ok := files["resources/mailmanager.go"]
	a.True(ok, "interface file should be present")
	a.Equal(Append, ifaceFile.Mode)
	a.Contains(ifaceFile.Content, "DeleteRelay")
	a.NotContains(ifaceFile.Content, "ListRelays")

	// Mock file should be AppendBeforeVar mode
	mockFile, ok := files["resources/mailmanager_mock_test.go"]
	a.True(ok, "mock file should be present")
	a.Equal(AppendBeforeVar, mockFile.Mode)
	a.Contains(mockFile.Content, "DeleteRelay")
	a.NotContains(mockFile.Content, "ListRelays")
}

func TestGenerateAll_OverridesUsesOverrideContent(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	def.Delete.Override = `_, err := r.svc.CustomDelete(ctx, &mailmanager.CustomDeleteInput{})
	return err`
	def.ExtraImports = []string{`mailmanagertypes "github.com/aws/aws-sdk-go-v2/service/mailmanager/types"`}

	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	resourceContent := files["resources/mailmanager-relay.go"].Content
	a.Contains(resourceContent, "CustomDelete")
	a.Contains(resourceContent, "mailmanagertypes")
}

func TestGenerateAll_NoFiltersOmitsFilterMethod(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	// Ensure no filters
	def.Filters = nil

	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	resourceContent := files["resources/mailmanager-relay.go"].Content
	a.NotContains(resourceContent, "func (r *MailManagerRelay) Filter()")
}

func TestGenerateAll_WithFiltersIncludesFilterMethod(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	def.Filters = []FilterDef{
		{
			Field:    "Status",
			Operator: "equals",
			Values:   []string{"DeleteInProgress"},
			Message:  "already being deleted",
		},
	}

	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	resourceContent := files["resources/mailmanager-relay.go"].Content
	a.Contains(resourceContent, "func (r *MailManagerRelay) Filter() error")
	a.Contains(resourceContent, "DeleteInProgress")
	a.Contains(resourceContent, "already being deleted")
}

func TestGenerateAll_WithSettingsIncludesSettingsInRegistration(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	def.Settings = []SettingDef{
		{
			Name:             "DisableDeletionProtection",
			ProtectionField:  "protection",
			DisableOperation: "UpdateRelayConfig",
			DisableInput:     map[string]interface{}{"Name": "Name", "Protection": false},
		},
	}

	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	resourceContent := files["resources/mailmanager-relay.go"].Content

	// Settings should appear in registration
	a.Contains(resourceContent, `Settings: []string{`)
	a.Contains(resourceContent, `"DisableDeletionProtection"`)

	// Settings method should be generated
	a.Contains(resourceContent, "func (r *MailManagerRelay) Settings(setting *libsettings.Setting)")

	// libsettings import should be present
	a.Contains(resourceContent, `libsettings "github.com/ekristen/libnuke/pkg/settings"`)

	// DisableOperation should be in required methods -> interface should include it
	ifaceContent := files["resources/mailmanager.go"].Content
	a.Contains(ifaceContent, "UpdateRelayConfig")
}

func TestGenerateAll_ResourceFileContainsExpectedStructure(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	content := files["resources/mailmanager-relay.go"].Content

	// Const declaration
	a.Contains(content, `const MailManagerRelayResource = "MailManagerRelay"`)

	// Init function with registry
	a.Contains(content, "func init()")
	a.Contains(content, "registry.Register")
	a.Contains(content, "MailManagerRelayResource")

	// Lister
	a.Contains(content, "MailManagerRelayLister")
	a.Contains(content, "func (l *MailManagerRelayLister) List(ctx context.Context, o interface{})")

	// Resource struct
	a.Contains(content, "type MailManagerRelay struct")

	// Remove method
	a.Contains(content, "func (r *MailManagerRelay) Remove(ctx context.Context) error")
	a.Contains(content, "DeleteRelay")

	// Properties method
	a.Contains(content, "func (r *MailManagerRelay) Properties() types.Properties")
	a.Contains(content, "types.NewPropertiesFromStruct(r)")

	// String method
	a.Contains(content, "func (r *MailManagerRelay) String() string")
}

func TestGenerateAll_ConcreteTypeSkipsInterfaceAndMock(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	def.SvcType = "concrete"

	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	// Concrete type should only produce the resource file
	a.Len(files, 1)
	_, ok := files["resources/mailmanager-relay.go"]
	a.True(ok)
}

func TestGenerateAll_AllExistingMethodsSkipsInterfaceAndMock(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	// Create existing interface file with ALL required methods already declared
	resourcesDir := filepath.Join(outputDir, "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	a.NoError(err)

	existingInterface := `package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"
)

type MailmanagerClient interface {
	ListRelays(ctx context.Context, params *mailmanager.ListRelaysInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.ListRelaysOutput, error)
	DeleteRelay(ctx context.Context, params *mailmanager.DeleteRelayInput,
		optFns ...func(*mailmanager.Options)) (*mailmanager.DeleteRelayOutput, error)
}
`
	err = os.WriteFile(filepath.Join(resourcesDir, "mailmanager.go"), []byte(existingInterface), 0600)
	a.NoError(err)

	existingMock := `package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockMailmanagerClient struct {
	mock.Mock
}

func (m *mockMailmanagerClient) ListRelays(
	ctx context.Context, params *mailmanager.ListRelaysInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.ListRelaysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.ListRelaysOutput), args.Error(1)
}

func (m *mockMailmanagerClient) DeleteRelay(
	ctx context.Context, params *mailmanager.DeleteRelayInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.DeleteRelayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.DeleteRelayOutput), args.Error(1)
}

var testMailmanagerListerOpts = &nuke.ListerOpts{}
`
	err = os.WriteFile(filepath.Join(resourcesDir, "mailmanager_mock_test.go"), []byte(existingMock), 0600)
	a.NoError(err)

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	// Interface and mock should NOT be in the output since all methods exist
	_, hasIface := files["resources/mailmanager.go"]
	a.False(hasIface, "interface file should not be generated when all methods exist")

	_, hasMock := files["resources/mailmanager_mock_test.go"]
	a.False(hasMock, "mock file should not be generated when all methods exist")

	// Resource and mock test files should still be generated
	_, hasResource := files["resources/mailmanager-relay.go"]
	a.True(hasResource)
	_, hasMockTest := files["resources/mailmanager-relay_mock_test.go"]
	a.True(hasMockTest)
}

func TestComputeRequiredMethods_FlatList(t *testing.T) {
	a := assert.New(t)

	def := minimalFlatDef()
	methods := computeRequiredMethods(def)

	a.Contains(methods, "ListRelays")
	a.Contains(methods, "DeleteRelay")
	a.Len(methods, 2)
}

func TestComputeRequiredMethods_WithDescribeAndTags(t *testing.T) {
	a := assert.New(t)

	def := minimalFlatDef()
	def.List.Describe = &DescribeDef{
		Operation:    "DescribeRelay",
		InputMapping: map[string]string{"RelayId": "RelayId"},
	}
	def.List.Tags = &TagFetchDef{
		Operation: "ListTagsForResource",
		ArnField:  "RelayArn",
	}

	methods := computeRequiredMethods(def)

	a.Contains(methods, "ListRelays")
	a.Contains(methods, "DescribeRelay")
	a.Contains(methods, "ListTagsForResource")
	a.Contains(methods, "DeleteRelay")
	a.Len(methods, 4)
}

func TestComputeRequiredMethods_NestedList(t *testing.T) {
	a := assert.New(t)

	def := &ResourceDef{
		List: ListDef{
			Strategy: "nested",
			Levels: []NestedLevelDef{
				{Operation: "ListParents"},
				{Operation: "ListChildren"},
			},
		},
		Delete: DeleteDef{Operation: "DeleteChild"},
	}

	methods := computeRequiredMethods(def)

	a.Contains(methods, "ListParents")
	a.Contains(methods, "ListChildren")
	a.Contains(methods, "DeleteChild")
	a.Len(methods, 3)
}

func TestComputeRequiredMethods_WithSettingsAndPreDeletion(t *testing.T) {
	a := assert.New(t)

	def := minimalFlatDef()
	def.Settings = []SettingDef{
		{Name: "DisableDeletionProtection", DisableOperation: "UpdateRelayConfig"},
	}
	def.PreDeletion = []PreDeletionDef{
		{
			Type:            "listAndBatchDelete",
			ListOperation:   "ListTargets",
			DeleteOperation: "DeregisterTargets",
		},
	}

	methods := computeRequiredMethods(def)

	a.Contains(methods, "ListRelays")
	a.Contains(methods, "DeleteRelay")
	a.Contains(methods, "UpdateRelayConfig")
	a.Contains(methods, "ListTargets")
	a.Contains(methods, "DeregisterTargets")
	a.Len(methods, 5)
}

func TestGenerateAll_MockTestFileContainsExpectedTests(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	content := files["resources/mailmanager-relay_mock_test.go"].Content

	a.Contains(content, "Test_Mock_MailManagerRelay_List")
	a.Contains(content, "Test_Mock_MailManagerRelay_List_Empty")
	a.Contains(content, "Test_Mock_MailManagerRelay_Remove")
	a.Contains(content, "Test_Mock_MailManagerRelay_Properties")
	a.Contains(content, "Test_Mock_MailManagerRelay_String")
}

func TestGenerateAll_InterfaceFileContainsAllMethods(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	content := files["resources/mailmanager.go"].Content

	a.Contains(content, "MailmanagerClient")
	a.Contains(content, "ListRelays")
	a.Contains(content, "DeleteRelay")
	a.True(strings.Contains(content, "interface"))
}

func TestGenerateAll_MockFileContainsAllMethods(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	content := files["resources/mailmanager_mock_test.go"].Content

	a.Contains(content, "mockMailmanagerClient")
	a.Contains(content, "ListRelays")
	a.Contains(content, "DeleteRelay")
	a.Contains(content, "m.Called")
}

func TestGenerateAll_ExistingMockTestIsNotOverwritten(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	// Create the mock test file before generation
	resourcesDir := filepath.Join(outputDir, "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	a.NoError(err)

	existingContent := "package resources\n// hand-edited mock test\n"
	mockTestPath := filepath.Join(resourcesDir, "mailmanager-relay_mock_test.go")
	err = os.WriteFile(mockTestPath, []byte(existingContent), 0600)
	a.NoError(err)

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	// Mock test file should NOT be in the generated output
	_, hasMockTest := files["resources/mailmanager-relay_mock_test.go"]
	a.False(hasMockTest, "mock test file should not be regenerated when it already exists")

	// Resource file should still be generated
	_, hasResource := files["resources/mailmanager-relay.go"]
	a.True(hasResource)
}

func TestGenerateAll_ExistingIntegrationTestIsNotOverwritten(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	// Create the integration test file before generation
	resourcesDir := filepath.Join(outputDir, "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	a.NoError(err)

	existingContent := "//go:build integration\npackage resources\n// hand-edited integration test\n"
	integTestPath := filepath.Join(resourcesDir, "mailmanager-relay_test.go")
	err = os.WriteFile(integTestPath, []byte(existingContent), 0600)
	a.NoError(err)

	def := minimalFlatDef()
	def.IntegrationTest = &IntegTestDef{
		Create: &IntegCreateDef{
			Operation: "CreateRelay",
			Input:     map[string]interface{}{"RelayName": "test-relay"},
		},
	}

	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	// Integration test file should NOT be in the generated output
	_, hasIntegTest := files["resources/mailmanager-relay_test.go"]
	a.False(hasIntegTest, "integration test file should not be regenerated when it already exists")

	// Resource file should still be generated
	_, hasResource := files["resources/mailmanager-relay.go"]
	a.True(hasResource)
}

func TestGenerateAll_MockTestIsGeneratedWhenMissing(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{})
	a.NoError(err)

	// Mock test file should be generated when it doesn't exist
	_, hasMockTest := files["resources/mailmanager-relay_mock_test.go"]
	a.True(hasMockTest, "mock test file should be generated when it doesn't exist")
}

func TestGenerateAll_ForceMockTestsOverwritesExisting(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	resourcesDir := filepath.Join(outputDir, "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	a.NoError(err)

	existingContent := "package resources\n// hand-edited mock test\n"
	mockTestPath := filepath.Join(resourcesDir, "mailmanager-relay_mock_test.go")
	err = os.WriteFile(mockTestPath, []byte(existingContent), 0600)
	a.NoError(err)

	def := minimalFlatDef()
	files, err := GenerateAll(def, outputDir, GenerateOpts{ForceMockTests: true})
	a.NoError(err)

	// Mock test file SHOULD be regenerated with --force-mock-tests
	generated, hasMockTest := files["resources/mailmanager-relay_mock_test.go"]
	a.True(hasMockTest, "mock test file should be regenerated with ForceMockTests")
	a.Contains(generated.Content, "Test_Mock_MailManagerRelay_List")
}

func TestGenerateAll_ForceIntegrationTestsOverwritesExisting(t *testing.T) {
	a := assert.New(t)
	outputDir := t.TempDir()

	resourcesDir := filepath.Join(outputDir, "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	a.NoError(err)

	existingContent := "//go:build integration\npackage resources\n// hand-edited\n"
	integTestPath := filepath.Join(resourcesDir, "mailmanager-relay_test.go")
	err = os.WriteFile(integTestPath, []byte(existingContent), 0600)
	a.NoError(err)

	def := minimalFlatDef()
	def.IntegrationTest = &IntegTestDef{
		Create: &IntegCreateDef{
			Operation: "CreateRelay",
			Input:     map[string]interface{}{"RelayName": "test-relay"},
		},
	}

	files, err := GenerateAll(def, outputDir, GenerateOpts{ForceIntegrationTests: true})
	a.NoError(err)

	// Integration test file SHOULD be regenerated with --force-integration-tests
	generated, hasIntegTest := files["resources/mailmanager-relay_test.go"]
	a.True(hasIntegTest, "integration test file should be regenerated with ForceIntegrationTests")
	a.Contains(generated.Content, "TestMailManagerRelay")
}
