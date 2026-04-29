package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/taxsettings"
)

type TaxSettingsClient interface {
	GetTaxInheritance(ctx context.Context, params *taxsettings.GetTaxInheritanceInput,
		optFns ...func(*taxsettings.Options)) (*taxsettings.GetTaxInheritanceOutput, error)
	PutTaxInheritance(ctx context.Context, params *taxsettings.PutTaxInheritanceInput,
		optFns ...func(*taxsettings.Options)) (*taxsettings.PutTaxInheritanceOutput, error)
}
