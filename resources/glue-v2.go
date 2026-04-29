package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/glue"
)

// GlueV2Client is the interface for the Glue SDK v2 client methods.
type GlueV2Client interface {
	GetCatalogs(ctx context.Context, params *glue.GetCatalogsInput,
		optFns ...func(*glue.Options)) (*glue.GetCatalogsOutput, error)
	DeleteCatalog(ctx context.Context, params *glue.DeleteCatalogInput,
		optFns ...func(*glue.Options)) (*glue.DeleteCatalogOutput, error)
	ListCustomEntityTypes(ctx context.Context, params *glue.ListCustomEntityTypesInput,
		optFns ...func(*glue.Options)) (*glue.ListCustomEntityTypesOutput, error)
	DeleteCustomEntityType(ctx context.Context, params *glue.DeleteCustomEntityTypeInput,
		optFns ...func(*glue.Options)) (*glue.DeleteCustomEntityTypeOutput, error)
	ListRegistries(ctx context.Context, params *glue.ListRegistriesInput,
		optFns ...func(*glue.Options)) (*glue.ListRegistriesOutput, error)
	DeleteRegistry(ctx context.Context, params *glue.DeleteRegistryInput,
		optFns ...func(*glue.Options)) (*glue.DeleteRegistryOutput, error)
	ListDataQualityRulesets(ctx context.Context, params *glue.ListDataQualityRulesetsInput,
		optFns ...func(*glue.Options)) (*glue.ListDataQualityRulesetsOutput, error)
	DeleteDataQualityRuleset(ctx context.Context, params *glue.DeleteDataQualityRulesetInput,
		optFns ...func(*glue.Options)) (*glue.DeleteDataQualityRulesetOutput, error)
	ListUsageProfiles(ctx context.Context, params *glue.ListUsageProfilesInput,
		optFns ...func(*glue.Options)) (*glue.ListUsageProfilesOutput, error)
	DeleteUsageProfile(ctx context.Context, params *glue.DeleteUsageProfileInput,
		optFns ...func(*glue.Options)) (*glue.DeleteUsageProfileOutput, error)
	GetDatabases(ctx context.Context, params *glue.GetDatabasesInput,
		optFns ...func(*glue.Options)) (*glue.GetDatabasesOutput, error)
	GetUserDefinedFunctions(ctx context.Context, params *glue.GetUserDefinedFunctionsInput,
		optFns ...func(*glue.Options)) (*glue.GetUserDefinedFunctionsOutput, error)
	DeleteUserDefinedFunction(ctx context.Context, params *glue.DeleteUserDefinedFunctionInput,
		optFns ...func(*glue.Options)) (*glue.DeleteUserDefinedFunctionOutput, error)
	GetDataCatalogEncryptionSettings(ctx context.Context, params *glue.GetDataCatalogEncryptionSettingsInput,
		optFns ...func(*glue.Options)) (*glue.GetDataCatalogEncryptionSettingsOutput, error)
	PutDataCatalogEncryptionSettings(ctx context.Context, params *glue.PutDataCatalogEncryptionSettingsInput,
		optFns ...func(*glue.Options)) (*glue.PutDataCatalogEncryptionSettingsOutput, error)

	GetTables(ctx context.Context, params *glue.GetTablesInput,
		optFns ...func(*glue.Options)) (*glue.GetTablesOutput, error)
	DeleteTable(ctx context.Context, params *glue.DeleteTableInput,
		optFns ...func(*glue.Options)) (*glue.DeleteTableOutput, error)
}
