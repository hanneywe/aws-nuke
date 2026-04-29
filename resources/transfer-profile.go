package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/transfer"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const TransferProfileResource = "TransferProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     TransferProfileResource,
		Scope:    nuke.Account,
		Resource: &TransferProfile{},
		Lister:   &TransferProfileLister{},
	})
}

type TransferProfileLister struct {
	svc TransferClient
}

func (l *TransferProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = transfer.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &transfer.ListProfilesInput{}
	for {
		listOutput, err := svc.ListProfiles(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, profile := range listOutput.Profiles {
			resources = append(resources, &TransferProfile{
				svc:       svc,
				ProfileID: profile.ProfileId,
				As2ID:     profile.As2Id,
			})
		}

		if listOutput.NextToken == nil {
			break
		}
		params.NextToken = listOutput.NextToken
	}

	return resources, nil
}

type TransferProfile struct {
	svc       TransferClient
	ProfileID *string
	As2ID     *string
}

func (r *TransferProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteProfile(ctx, &transfer.DeleteProfileInput{
		ProfileId: r.ProfileID,
	})
	return err
}

func (r *TransferProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *TransferProfile) String() string {
	return *r.ProfileID
}
