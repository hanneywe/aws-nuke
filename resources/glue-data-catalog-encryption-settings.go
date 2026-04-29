package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GlueDataCatalogEncryptionSettingsResource = "GlueDataCatalogEncryptionSettings"

func init() {
	registry.Register(&registry.Registration{
		Name:     GlueDataCatalogEncryptionSettingsResource,
		Scope:    nuke.Account,
		Resource: &GlueDataCatalogEncryptionSettings{},
		Lister:   &GlueDataCatalogEncryptionSettingsLister{},
	})
}

type GlueDataCatalogEncryptionSettingsLister struct {
	svc GlueV2Client
}

func (l *GlueDataCatalogEncryptionSettingsLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = glue.NewFromConfig(*opts.Config)
	}

	resp, err := svc.GetDataCatalogEncryptionSettings(ctx, &glue.GetDataCatalogEncryptionSettingsInput{})
	if err != nil {
		return nil, err
	}

	if resp.DataCatalogEncryptionSettings == nil {
		return nil, nil
	}

	settings := resp.DataCatalogEncryptionSettings

	encryptionAtRest := settings.EncryptionAtRest != nil &&
		settings.EncryptionAtRest.CatalogEncryptionMode != gluetypes.CatalogEncryptionModeDisabled

	connectionPasswordEncryption := settings.ConnectionPasswordEncryption != nil &&
		settings.ConnectionPasswordEncryption.ReturnConnectionPasswordEncrypted

	if !encryptionAtRest && !connectionPasswordEncryption {
		return nil, nil
	}

	return []resource.Resource{
		&GlueDataCatalogEncryptionSettings{
			svc:                          svc,
			EncryptionAtRest:             &encryptionAtRest,
			ConnectionPasswordEncryption: &connectionPasswordEncryption,
		},
	}, nil
}

type GlueDataCatalogEncryptionSettings struct {
	svc                          GlueV2Client
	EncryptionAtRest             *bool
	ConnectionPasswordEncryption *bool
}

func (r *GlueDataCatalogEncryptionSettings) Remove(ctx context.Context) error {
	_, err := r.svc.PutDataCatalogEncryptionSettings(ctx, &glue.PutDataCatalogEncryptionSettingsInput{
		DataCatalogEncryptionSettings: &gluetypes.DataCatalogEncryptionSettings{
			ConnectionPasswordEncryption: &gluetypes.ConnectionPasswordEncryption{
				ReturnConnectionPasswordEncrypted: false,
			},
			EncryptionAtRest: &gluetypes.EncryptionAtRest{
				CatalogEncryptionMode: gluetypes.CatalogEncryptionModeDisabled,
			},
		},
	})
	return err
}

func (r *GlueDataCatalogEncryptionSettings) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GlueDataCatalogEncryptionSettings) String() string {
	return "GlueDataCatalogEncryptionSettings"
}
