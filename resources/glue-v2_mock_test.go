package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockGlueV2Client struct {
	mock.Mock
}

func (m *mockGlueV2Client) GetCatalogs(
	ctx context.Context, params *glue.GetCatalogsInput,
	_ ...func(*glue.Options),
) (*glue.GetCatalogsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.GetCatalogsOutput), args.Error(1)
}

func (m *mockGlueV2Client) DeleteCatalog(
	ctx context.Context, params *glue.DeleteCatalogInput,
	_ ...func(*glue.Options),
) (*glue.DeleteCatalogOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.DeleteCatalogOutput), args.Error(1)
}

func (m *mockGlueV2Client) ListCustomEntityTypes(
	ctx context.Context, params *glue.ListCustomEntityTypesInput,
	_ ...func(*glue.Options),
) (*glue.ListCustomEntityTypesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.ListCustomEntityTypesOutput), args.Error(1)
}

func (m *mockGlueV2Client) DeleteCustomEntityType(
	ctx context.Context, params *glue.DeleteCustomEntityTypeInput,
	_ ...func(*glue.Options),
) (*glue.DeleteCustomEntityTypeOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.DeleteCustomEntityTypeOutput), args.Error(1)
}

func (m *mockGlueV2Client) ListRegistries(
	ctx context.Context, params *glue.ListRegistriesInput,
	_ ...func(*glue.Options),
) (*glue.ListRegistriesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.ListRegistriesOutput), args.Error(1)
}

func (m *mockGlueV2Client) DeleteRegistry(
	ctx context.Context, params *glue.DeleteRegistryInput,
	_ ...func(*glue.Options),
) (*glue.DeleteRegistryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.DeleteRegistryOutput), args.Error(1)
}

func (m *mockGlueV2Client) ListDataQualityRulesets(
	ctx context.Context, params *glue.ListDataQualityRulesetsInput,
	_ ...func(*glue.Options),
) (*glue.ListDataQualityRulesetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.ListDataQualityRulesetsOutput), args.Error(1)
}

func (m *mockGlueV2Client) DeleteDataQualityRuleset(
	ctx context.Context, params *glue.DeleteDataQualityRulesetInput,
	_ ...func(*glue.Options),
) (*glue.DeleteDataQualityRulesetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.DeleteDataQualityRulesetOutput), args.Error(1)
}

func (m *mockGlueV2Client) ListUsageProfiles(
	ctx context.Context, params *glue.ListUsageProfilesInput,
	_ ...func(*glue.Options),
) (*glue.ListUsageProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.ListUsageProfilesOutput), args.Error(1)
}

func (m *mockGlueV2Client) DeleteUsageProfile(
	ctx context.Context, params *glue.DeleteUsageProfileInput,
	_ ...func(*glue.Options),
) (*glue.DeleteUsageProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.DeleteUsageProfileOutput), args.Error(1)
}

func (m *mockGlueV2Client) GetDatabases(
	ctx context.Context, params *glue.GetDatabasesInput,
	_ ...func(*glue.Options),
) (*glue.GetDatabasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.GetDatabasesOutput), args.Error(1)
}

func (m *mockGlueV2Client) GetUserDefinedFunctions(
	ctx context.Context, params *glue.GetUserDefinedFunctionsInput,
	_ ...func(*glue.Options),
) (*glue.GetUserDefinedFunctionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.GetUserDefinedFunctionsOutput), args.Error(1)
}

func (m *mockGlueV2Client) DeleteUserDefinedFunction(
	ctx context.Context, params *glue.DeleteUserDefinedFunctionInput,
	_ ...func(*glue.Options),
) (*glue.DeleteUserDefinedFunctionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.DeleteUserDefinedFunctionOutput), args.Error(1)
}

func (m *mockGlueV2Client) GetDataCatalogEncryptionSettings(
	ctx context.Context, params *glue.GetDataCatalogEncryptionSettingsInput,
	_ ...func(*glue.Options),
) (*glue.GetDataCatalogEncryptionSettingsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.GetDataCatalogEncryptionSettingsOutput), args.Error(1)
}

func (m *mockGlueV2Client) PutDataCatalogEncryptionSettings(
	ctx context.Context, params *glue.PutDataCatalogEncryptionSettingsInput,
	_ ...func(*glue.Options),
) (*glue.PutDataCatalogEncryptionSettingsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.PutDataCatalogEncryptionSettingsOutput), args.Error(1)
}

func (m *mockGlueV2Client) GetTables(
	ctx context.Context, params *glue.GetTablesInput,
	_ ...func(*glue.Options),
) (*glue.GetTablesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.GetTablesOutput), args.Error(1)
}

func (m *mockGlueV2Client) DeleteTable(
	ctx context.Context, params *glue.DeleteTableInput,
	_ ...func(*glue.Options),
) (*glue.DeleteTableOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.DeleteTableOutput), args.Error(1)
}

var testGlueV2ListerOpts = &nuke.ListerOpts{}
