package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/b2bi"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const B2BIProfileResource = "B2BIProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     B2BIProfileResource,
		Scope:    nuke.Account,
		Resource: &B2BIProfile{},
		Lister:   &B2BIProfileLister{},
	})
}

type B2BIProfileLister struct {
	svc B2BIClient
}

func (l *B2BIProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = b2bi.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := b2bi.NewListProfilesPaginator(svc, &b2bi.ListProfilesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Profiles {
			resources = append(resources, &B2BIProfile{
				svc:       svc,
				ProfileID: item.ProfileId,
				Name:      item.Name,
			})
		}
	}

	return resources, nil
}

type B2BIProfile struct {
	svc       B2BIClient
	ProfileID *string `property:"name=ProfileId"`
	Name      *string
}

func (r *B2BIProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteProfile(ctx, &b2bi.DeleteProfileInput{
		ProfileId: r.ProfileID,
	})
	return err
}

func (r *B2BIProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *B2BIProfile) String() string {
	return *r.ProfileID
}
