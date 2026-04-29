package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	emrtypes "github.com/aws/aws-sdk-go-v2/service/emr/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EMRBlockPublicAccessConfigurationResource = "EMRBlockPublicAccessConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     EMRBlockPublicAccessConfigurationResource,
		Scope:    nuke.Account,
		Resource: &EMRBlockPublicAccessConfiguration{},
		Lister:   &EMRBlockPublicAccessConfigurationLister{},
	})
}

type EMRBlockPublicAccessConfigurationLister struct {
	svc EMRV2Client
}

func (l *EMRBlockPublicAccessConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = emr.NewFromConfig(*opts.Config)
	}

	resp, err := svc.GetBlockPublicAccessConfiguration(ctx, &emr.GetBlockPublicAccessConfigurationInput{})
	if err != nil {
		return nil, err
	}

	if resp.BlockPublicAccessConfiguration == nil {
		return nil, nil
	}

	return []resource.Resource{
		&EMRBlockPublicAccessConfiguration{
			svc:                           svc,
			BlockPublicSecurityGroupRules: resp.BlockPublicAccessConfiguration.BlockPublicSecurityGroupRules,
		},
	}, nil
}

type EMRBlockPublicAccessConfiguration struct {
	svc                           EMRV2Client
	BlockPublicSecurityGroupRules *bool
}

func (r *EMRBlockPublicAccessConfiguration) Filter() error {
	if r.BlockPublicSecurityGroupRules != nil && !*r.BlockPublicSecurityGroupRules {
		return fmt.Errorf("already at default configuration")
	}
	return nil
}

func (r *EMRBlockPublicAccessConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.PutBlockPublicAccessConfiguration(ctx, &emr.PutBlockPublicAccessConfigurationInput{
		BlockPublicAccessConfiguration: &emrtypes.BlockPublicAccessConfiguration{
			BlockPublicSecurityGroupRules: aws.Bool(false),
		},
	})
	return err
}

func (r *EMRBlockPublicAccessConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EMRBlockPublicAccessConfiguration) String() string {
	return "EMRBlockPublicAccessConfiguration"
}
