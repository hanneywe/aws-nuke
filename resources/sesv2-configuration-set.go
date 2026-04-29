package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SESv2ConfigurationSetResource = "SESv2ConfigurationSet"

func init() {
	registry.Register(&registry.Registration{
		Name:     SESv2ConfigurationSetResource,
		Scope:    nuke.Account,
		Resource: &SESv2ConfigurationSet{},
		Lister:   &SESv2ConfigurationSetLister{},
	})
}

type SESv2ConfigurationSetLister struct {
	svc SESv2Client
}

func (l *SESv2ConfigurationSetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = sesv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := sesv2.NewListConfigurationSetsPaginator(svc, &sesv2.ListConfigurationSetsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, name := range resp.ConfigurationSets {
			resources = append(resources, &SESv2ConfigurationSet{
				svc:                  svc,
				ConfigurationSetName: aws.String(name),
			})
		}
	}
	return resources, nil
}

type SESv2ConfigurationSet struct {
	svc                  SESv2Client
	ConfigurationSetName *string
}

func (r *SESv2ConfigurationSet) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConfigurationSet(ctx, &sesv2.DeleteConfigurationSetInput{
		ConfigurationSetName: r.ConfigurationSetName,
	})
	return err
}

func (r *SESv2ConfigurationSet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SESv2ConfigurationSet) String() string {
	return *r.ConfigurationSetName
}
