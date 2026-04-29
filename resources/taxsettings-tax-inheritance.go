package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/taxsettings"
	taxsettingstypes "github.com/aws/aws-sdk-go-v2/service/taxsettings/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const TaxSettingsTaxInheritanceResource = "TaxSettingsTaxInheritance"

func init() {
	registry.Register(&registry.Registration{
		Name:     TaxSettingsTaxInheritanceResource,
		Scope:    nuke.Account,
		Resource: &TaxSettingsTaxInheritance{},
		Lister:   &TaxSettingsTaxInheritanceLister{},
	})
}

type TaxSettingsTaxInheritanceLister struct {
	svc TaxSettingsClient
}

func (l *TaxSettingsTaxInheritanceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = taxsettings.NewFromConfig(*opts.Config)
	}

	resp, err := svc.GetTaxInheritance(ctx, &taxsettings.GetTaxInheritanceInput{})
	if err != nil {
		return nil, err
	}

	var resources []resource.Resource
	if resp.HeritageStatus == taxsettingstypes.HeritageStatusOptIn {
		resources = append(resources, &TaxSettingsTaxInheritance{
			svc:    svc,
			Status: resp.HeritageStatus,
		})
	}

	return resources, nil
}

type TaxSettingsTaxInheritance struct {
	svc    TaxSettingsClient
	Status taxsettingstypes.HeritageStatus
}

func (r *TaxSettingsTaxInheritance) Remove(ctx context.Context) error {
	_, err := r.svc.PutTaxInheritance(ctx, &taxsettings.PutTaxInheritanceInput{
		HeritageStatus: taxsettingstypes.HeritageStatusOptOut,
	})
	return err
}

func (r *TaxSettingsTaxInheritance) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *TaxSettingsTaxInheritance) String() string {
	return string(r.Status)
}
