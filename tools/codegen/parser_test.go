package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_MinimalDSL(t *testing.T) {
	a := assert.New(t)

	def, err := Parse("testdata/minimal.yaml")
	a.NoError(err)
	a.NotNil(def)

	a.Equal("mailmanager", def.Service)
	a.Equal("relay", def.Resource)
	a.Equal("MailManagerRelay", def.ResourceName)
	a.Equal("account", def.Scope)
	a.Equal("mailmanager", def.SDKPackage)
	a.Empty(def.SvcType)

	a.Equal("flat", def.List.Strategy)
	a.Equal("ListRelays", def.List.Operation)
	a.Equal("paginator", def.List.Pagination)
	a.Equal("Relays", def.List.ItemsField)
	a.Nil(def.List.Describe)
	a.Nil(def.List.Tags)
	a.Empty(def.List.Levels)

	a.Equal("DeleteRelay", def.Delete.Operation)
	a.Equal([]string{"RelayId"}, def.Delete.InputFields)
	a.Empty(def.Delete.Override)

	a.Len(def.Fields, 2)
	a.Equal("RelayId", def.Fields[0].Name)
	a.Equal("*string", def.Fields[0].Type)
	a.Equal("RelayId", def.Fields[0].FromList)
	a.Equal("RelayName", def.Fields[1].Name)

	a.Equal("RelayName", def.StringRep.Field)
	a.Empty(def.StringRep.Format)
	a.Nil(def.StringRep.Conditional)

	a.Empty(def.Filters)
	a.Empty(def.Settings)
	a.Empty(def.Dependencies)
	a.Empty(def.PreDeletion)
	a.Nil(def.IntegrationTest)
}

func TestParse_ComplexDSL(t *testing.T) {
	a := assert.New(t)

	def, err := Parse("testdata/complex.yaml")
	a.NoError(err)
	a.NotNil(def)

	a.Equal("mediapackagev2", def.Service)
	a.Equal("origin-endpoint", def.Resource)
	a.Equal("MediaPackageV2OriginEndpoint", def.ResourceName)
	a.Equal("concrete", def.SvcType)

	// Nested list with 3 levels
	a.Equal("nested", def.List.Strategy)
	a.Len(def.List.Levels, 3)
	a.Equal("ListChannelGroups", def.List.Levels[0].Operation)
	a.Equal("cg", def.List.Levels[0].IteratorVar)
	a.Empty(def.List.Levels[0].ParentInputMapping)
	a.Equal("ListChannels", def.List.Levels[1].Operation)
	a.Equal("ch", def.List.Levels[1].IteratorVar)
	a.Equal("cg.ChannelGroupName", def.List.Levels[1].ParentInputMapping["ChannelGroupName"])
	a.Equal("ListOriginEndpoints", def.List.Levels[2].Operation)
	a.Equal("ep", def.List.Levels[2].IteratorVar)
	a.Len(def.List.Levels[2].ParentInputMapping, 2)

	// Delete
	a.Equal("DeleteOriginEndpoint", def.Delete.Operation)
	a.Equal([]string{"ChannelGroupName", "ChannelName", "OriginEndpointName"}, def.Delete.InputFields)

	// Fields
	a.Len(def.Fields, 3)
	a.Equal("ChannelGroupName", def.Fields[0].Name)
	a.Equal("cg.ChannelGroupName", def.Fields[0].FromList)

	// Filters
	a.Len(def.Filters, 1)
	a.Equal("Status", def.Filters[0].Field)
	a.Equal("equals", def.Filters[0].Operator)
	a.Equal([]string{"DeleteInProgress", "Canceled"}, def.Filters[0].Values)
	a.Equal("already being deleted or canceled", def.Filters[0].Message)

	// Settings
	a.Len(def.Settings, 1)
	a.Equal("DisableDeletionProtection", def.Settings[0].Name)
	a.Equal("protection", def.Settings[0].ProtectionField)
	a.Equal("UpdateEndpointConfig", def.Settings[0].DisableOperation)

	// Dependencies
	a.Equal([]string{"MediaPackageV2Channel"}, def.Dependencies)

	// PreDeletion
	a.Len(def.PreDeletion, 1)
	a.Equal("listAndBatchDelete", def.PreDeletion[0].Type)
	a.Equal("ListTargets", def.PreDeletion[0].ListOperation)

	// Overrides
	a.NotEmpty(def.Delete.Override)
	a.Contains(def.Delete.Override, "DeleteOriginEndpoint")
	a.Len(def.ExtraImports, 1)

	// Integration test
	a.NotNil(def.IntegrationTest)
	a.Equal("CreateOriginEndpoint", def.IntegrationTest.Create.Operation)
}

func TestParse_InvalidYAML(t *testing.T) {
	a := assert.New(t)

	def, err := Parse("testdata/invalid.yaml")
	a.Error(err)
	a.Nil(def)
	a.Contains(err.Error(), "parsing YAML")
}

func TestParse_NonexistentFile(t *testing.T) {
	a := assert.New(t)

	def, err := Parse("testdata/does-not-exist.yaml")
	a.Error(err)
	a.Nil(def)
	a.Contains(err.Error(), "reading DSL file")
}

func TestParse_EmptyFile(t *testing.T) {
	a := assert.New(t)

	tmpDir := t.TempDir()
	emptyFile := filepath.Join(tmpDir, "empty.yaml")
	err := os.WriteFile(emptyFile, []byte(""), 0600)
	a.NoError(err)

	def, err := Parse(emptyFile)
	a.NoError(err)
	a.NotNil(def)
	a.Empty(def.Service)
}

func TestRoundTrip_Minimal(t *testing.T) {
	a := assert.New(t)

	original, err := Parse("testdata/minimal.yaml")
	a.NoError(err)

	yamlBytes, err := PrettyPrint(original)
	a.NoError(err)
	a.NotEmpty(yamlBytes)

	// Write to temp file and parse again
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "roundtrip.yaml")
	err = os.WriteFile(tmpFile, yamlBytes, 0600)
	a.NoError(err)

	roundTripped, err := Parse(tmpFile)
	a.NoError(err)

	a.Equal(original.Service, roundTripped.Service)
	a.Equal(original.Resource, roundTripped.Resource)
	a.Equal(original.ResourceName, roundTripped.ResourceName)
	a.Equal(original.Scope, roundTripped.Scope)
	a.Equal(original.SDKPackage, roundTripped.SDKPackage)
	a.Equal(original.List, roundTripped.List)
	a.Equal(original.Delete, roundTripped.Delete)
	a.Equal(original.Fields, roundTripped.Fields)
	a.Equal(original.StringRep, roundTripped.StringRep)
}

func TestRoundTrip_Complex(t *testing.T) {
	a := assert.New(t)

	original, err := Parse("testdata/complex.yaml")
	a.NoError(err)

	yamlBytes, err := PrettyPrint(original)
	a.NoError(err)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "roundtrip.yaml")
	err = os.WriteFile(tmpFile, yamlBytes, 0600)
	a.NoError(err)

	roundTripped, err := Parse(tmpFile)
	a.NoError(err)

	a.Equal(original.Service, roundTripped.Service)
	a.Equal(original.ResourceName, roundTripped.ResourceName)
	a.Equal(original.SvcType, roundTripped.SvcType)
	a.Equal(original.List.Strategy, roundTripped.List.Strategy)
	a.Equal(original.List.Levels, roundTripped.List.Levels)
	a.Equal(original.Delete, roundTripped.Delete)
	a.Equal(original.Fields, roundTripped.Fields)
	a.Equal(original.Filters, roundTripped.Filters)
	a.Equal(original.Settings, roundTripped.Settings)
	a.Equal(original.Dependencies, roundTripped.Dependencies)
	a.Equal(original.PreDeletion, roundTripped.PreDeletion)
	a.Equal(original.Delete.Override, roundTripped.Delete.Override)
	a.Equal(original.ExtraImports, roundTripped.ExtraImports)
	a.Equal(original.StringRep, roundTripped.StringRep)
	a.Equal(original.IntegrationTest.Create.Operation, roundTripped.IntegrationTest.Create.Operation)
}

func TestRoundTrip_FullResourceDef(t *testing.T) {
	a := assert.New(t)

	original := fullResourceDef()

	yamlBytes, err := PrettyPrint(&original)
	a.NoError(err)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "roundtrip.yaml")
	err = os.WriteFile(tmpFile, yamlBytes, 0600)
	a.NoError(err)

	roundTripped, err := Parse(tmpFile)
	a.NoError(err)

	a.Equal(original, *roundTripped)
}
