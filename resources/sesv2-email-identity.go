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

const SESv2EmailIdentityResource = "SESv2EmailIdentity"

func init() {
	registry.Register(&registry.Registration{
		Name:     SESv2EmailIdentityResource,
		Scope:    nuke.Account,
		Resource: &SESv2EmailIdentity{},
		Lister:   &SESv2EmailIdentityLister{},
	})
}

type SESv2EmailIdentityLister struct {
	svc SESv2Client
}

func (l *SESv2EmailIdentityLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = sesv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := sesv2.NewListEmailIdentitiesPaginator(svc, &sesv2.ListEmailIdentitiesInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, id := range resp.EmailIdentities {
			resources = append(resources, &SESv2EmailIdentity{
				svc:          svc,
				IdentityName: id.IdentityName,
				IdentityType: aws.String(string(id.IdentityType)),
			})
		}
	}
	return resources, nil
}

type SESv2EmailIdentity struct {
	svc          SESv2Client
	IdentityName *string
	IdentityType *string
}

func (r *SESv2EmailIdentity) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{
		EmailIdentity: r.IdentityName,
	})
	return err
}

func (r *SESv2EmailIdentity) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SESv2EmailIdentity) String() string {
	return *r.IdentityName
}
