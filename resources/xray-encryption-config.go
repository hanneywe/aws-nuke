package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/xray"
	xraytypes "github.com/aws/aws-sdk-go-v2/service/xray/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const XRayEncryptionConfigResource = "XRayEncryptionConfig"

func init() {
	registry.Register(&registry.Registration{
		Name:     XRayEncryptionConfigResource,
		Scope:    nuke.Account,
		Resource: &XRayEncryptionConfig{},
		Lister:   &XRayEncryptionConfigLister{},
	})
}

type XRayEncryptionConfigLister struct {
	svc XRayClient
}

func (l *XRayEncryptionConfigLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = xray.NewFromConfig(*opts.Config)
	}

	resp, err := svc.GetEncryptionConfig(ctx, &xray.GetEncryptionConfigInput{})
	if err != nil {
		return nil, err
	}

	var resources []resource.Resource
	if resp.EncryptionConfig != nil && resp.EncryptionConfig.Type == xraytypes.EncryptionTypeKms {
		resources = append(resources, &XRayEncryptionConfig{
			svc:   svc,
			Type:  resp.EncryptionConfig.Type,
			KeyID: resp.EncryptionConfig.KeyId,
		})
	}

	return resources, nil
}

type XRayEncryptionConfig struct {
	svc   XRayClient
	Type  xraytypes.EncryptionType
	KeyID *string
}

func (r *XRayEncryptionConfig) Remove(ctx context.Context) error {
	_, err := r.svc.PutEncryptionConfig(ctx, &xray.PutEncryptionConfigInput{
		Type: xraytypes.EncryptionTypeNone,
	})
	return err
}

func (r *XRayEncryptionConfig) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *XRayEncryptionConfig) String() string {
	return string(r.Type)
}
