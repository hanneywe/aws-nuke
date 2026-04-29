package main

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func renderMockTestTemplate(t *testing.T, data *MockTestTemplateData) string {
	t.Helper()

	content, err := templateFS.ReadFile("templates/mock_test.go.tmpl")
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	tmpl, err := template.New("mock_test.go.tmpl").Funcs(templateFuncs()).Parse(string(content))
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	return buf.String()
}

func TestMockTestTemplate_FlatList(t *testing.T) {
	a := assert.New(t)

	def := &ResourceDef{
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

	data := &MockTestTemplateData{
		ResourceDef:       def,
		MockStructName:    "mockMailManagerClient",
		TestListerOptsVar: "testMailManagerListerOpts",
		SDKPackage:        "mailmanager",
	}

	output := renderMockTestTemplate(t, data)

	// Verify package declaration
	a.Contains(output, "package resources")

	// Verify imports
	a.Contains(output, `"github.com/gotidy/ptr"`)
	a.Contains(output, `"github.com/stretchr/testify/assert"`)
	a.Contains(output, `"github.com/stretchr/testify/mock"`)
	a.Contains(output, `"github.com/aws/aws-sdk-go-v2/service/mailmanager"`)

	// Verify Test_Mock_MailManagerRelay_List
	a.Contains(output, "func Test_Mock_MailManagerRelay_List(t *testing.T)")
	a.Contains(output, `mockClient.On("ListRelays", mock.Anything, mock.Anything)`)
	a.Contains(output, "MailManagerRelayLister{svc: mockClient}")
	a.Contains(output, "testMailManagerListerOpts")
	a.Contains(output, "a.Len(resources, 1)")
	a.Contains(output, `r := resources[0].(*MailManagerRelay)`)

	// Verify Test_Mock_MailManagerRelay_List_Empty
	a.Contains(output, "func Test_Mock_MailManagerRelay_List_Empty(t *testing.T)")
	a.Contains(output, "a.Len(resources, 0)")

	// Verify Test_Mock_MailManagerRelay_Remove
	a.Contains(output, "func Test_Mock_MailManagerRelay_Remove(t *testing.T)")
	a.Contains(output, `mockClient.On("DeleteRelay"`)
	a.Contains(output, "mailmanager.DeleteRelayInput")
	a.Contains(output, "a.NoError(r.Remove(context.TODO()))")

	// Verify Test_Mock_MailManagerRelay_Properties
	a.Contains(output, "func Test_Mock_MailManagerRelay_Properties(t *testing.T)")
	a.Contains(output, "r.Properties()")
	a.Contains(output, `props.Get("RelayId")`)
	a.Contains(output, `props.Get("RelayName")`)

	// Verify Test_Mock_MailManagerRelay_String
	a.Contains(output, "func Test_Mock_MailManagerRelay_String(t *testing.T)")
	a.Contains(output, `"test-relayname"`)
}

func TestMockTestTemplate_NestedList(t *testing.T) {
	a := assert.New(t)

	def := &ResourceDef{
		Service:      "mediapackagev2",
		Resource:     "origin-endpoint",
		ResourceName: "MediaPackageV2OriginEndpoint",
		Scope:        "account",
		SDKPackage:   "mediapackagev2",
		List: ListDef{
			Strategy: "nested",
			Levels: []NestedLevelDef{
				{
					Operation: "ListChannelGroups", Pagination: "paginator",
					ItemsField: "Items", ItemType: "ChannelGroupListConfiguration",
					IteratorVar: "cg",
				},
				{
					Operation: "ListChannels", Pagination: "paginator",
					ItemsField: "Items", ItemType: "ChannelListConfiguration",
					IteratorVar: "ch",
					ParentInputMapping: map[string]string{
						"ChannelGroupName": "cg.ChannelGroupName",
					},
				},
				{
					Operation: "ListOriginEndpoints", Pagination: "paginator",
					ItemsField: "Items", ItemType: "OriginEndpointListConfiguration",
					IteratorVar: "ep",
					ParentInputMapping: map[string]string{
						"ChannelGroupName": "cg.ChannelGroupName",
						"ChannelName":      "ch.ChannelName",
					},
				},
			},
		},
		Delete: DeleteDef{
			Operation:   "DeleteOriginEndpoint",
			InputFields: []string{"ChannelGroupName", "ChannelName", "OriginEndpointName"},
		},
		Fields: []FieldDef{
			{Name: "ChannelGroupName", Type: "*string", FromList: "cg.ChannelGroupName"},
			{Name: "ChannelName", Type: "*string", FromList: "ch.ChannelName"},
			{Name: "OriginEndpointName", Type: "*string", FromList: "ep.OriginEndpointName"},
		},
		StringRep: StringRepDef{Field: "OriginEndpointName"},
	}

	data := &MockTestTemplateData{
		ResourceDef:       def,
		MockStructName:    "mockMediaPackageV2Client",
		TestListerOptsVar: "testMediaPackageV2ListerOpts",
		SDKPackage:        "mediapackagev2",
	}

	output := renderMockTestTemplate(t, data)

	// Verify nested list mocks all levels
	a.Contains(output, `mockClient.On("ListChannelGroups"`)
	a.Contains(output, `mockClient.On("ListChannels"`)
	a.Contains(output, `mockClient.On("ListOriginEndpoints"`)

	// Verify empty list only mocks first level
	emptyIdx := strings.Index(output, "func Test_Mock_MediaPackageV2OriginEndpoint_List_Empty")
	a.Greater(emptyIdx, 0)
	emptySection := output[emptyIdx:]
	removeIdx := strings.Index(emptySection, "func Test_Mock_MediaPackageV2OriginEndpoint_Remove")
	a.Greater(removeIdx, 0)
	emptySection = emptySection[:removeIdx]
	a.Contains(emptySection, `mockClient.On("ListChannelGroups"`)
	a.NotContains(emptySection, `mockClient.On("ListChannels"`)

	// Verify Remove mocks all delete input fields
	a.Contains(output, "ChannelGroupName: r.ChannelGroupName")
	a.Contains(output, "ChannelName: r.ChannelName")
	a.Contains(output, "OriginEndpointName: r.OriginEndpointName")

	// Verify String test
	a.Contains(output, `"test-originendpointname"`)
}

func TestMockTestTemplate_FormatString(t *testing.T) {
	a := assert.New(t)

	def := &ResourceDef{
		Service:      "myservice",
		Resource:     "myresource",
		ResourceName: "MyResource",
		Scope:        "account",
		SDKPackage:   "myservice",
		List: ListDef{
			Strategy:   "flat",
			Operation:  "ListResources",
			Pagination: "none",
			ItemsField: "Resources",
			ItemType:   "Resource",
		},
		Delete: DeleteDef{
			Operation:   "DeleteResource",
			InputFields: []string{"ResourceId"},
		},
		Fields: []FieldDef{
			{Name: "ResourceId", Type: "*string", FromList: "ResourceId"},
			{Name: "Name", Type: "*string", FromList: "Name"},
		},
		StringRep: StringRepDef{
			Format: "%s (%s)",
			Fields: []string{"ResourceId", "Name"},
		},
	}

	data := &MockTestTemplateData{
		ResourceDef:       def,
		MockStructName:    "mockMyServiceClient",
		TestListerOptsVar: "testMyServiceListerOpts",
		SDKPackage:        "myservice",
	}

	output := renderMockTestTemplate(t, data)

	// Verify String test uses fmt.Sprintf for format strings
	a.Contains(output, `fmt.Sprintf("%s (%s)"`)
}

func TestMockTestTemplate_Singleton(t *testing.T) {
	a := assert.New(t)

	def := &ResourceDef{
		Service:      "emr",
		Resource:     "block-public-access",
		ResourceName: "EMRBlockPublicAccess",
		Scope:        "account",
		SDKPackage:   "emr",
		List: ListDef{
			Strategy:      "singleton",
			Operation:     "GetBlockPublicAccessConfiguration",
			ResponseField: "BlockPublicAccessConfiguration",
			NilCheck:      true,
		},
		Delete: DeleteDef{
			Operation:   "PutBlockPublicAccessConfiguration",
			InputFields: []string{},
		},
		Fields: []FieldDef{
			{Name: "Status", Type: "*string", FromList: "Status"},
		},
		StringRep: StringRepDef{Field: "Status"},
	}

	data := &MockTestTemplateData{
		ResourceDef:       def,
		MockStructName:    "mockEMRClient",
		TestListerOptsVar: "testEMRListerOpts",
		SDKPackage:        "emr",
	}

	output := renderMockTestTemplate(t, data)

	// Verify singleton list mock
	a.Contains(output, `mockClient.On("GetBlockPublicAccessConfiguration"`)
	a.Contains(output, "a.Len(resources, 1)")

	// Verify empty singleton returns empty output
	a.Contains(output, "func Test_Mock_EMRBlockPublicAccess_List_Empty")
	a.Contains(output, "a.Len(resources, 0)")

	// Singleton with ResponseField needs SDK types import
	a.Contains(output, "emrtypes")
}
