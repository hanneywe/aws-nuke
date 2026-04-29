package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/customerprofiles"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CustomerProfilesProfileResource = "CustomerProfilesProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     CustomerProfilesProfileResource,
		Scope:    nuke.Account,
		Resource: &CustomerProfilesProfile{},
		Lister:   &CustomerProfilesProfileLister{},
	})
}

type CustomerProfilesProfileLister struct {
	svc CustomerProfilesClient
}

func (l *CustomerProfilesProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = customerprofiles.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	domainParams := &customerprofiles.ListDomainsInput{}
	for {
		domainResp, err := svc.ListDomains(ctx, domainParams)
		if err != nil {
			return nil, err
		}

		for _, domain := range domainResp.Items {
			searchParams := &customerprofiles.SearchProfilesInput{
				DomainName: domain.DomainName,
				KeyName:    aws.String("_profileId"),
				Values:     []string{"*"},
				MaxResults: aws.Int32(100),
			}
			for {
				searchResp, err := svc.SearchProfiles(ctx, searchParams)
				if err != nil {
					return nil, err
				}

				for i := range searchResp.Items {
					profile := &searchResp.Items[i]
					resources = append(resources, &CustomerProfilesProfile{
						svc:        svc,
						DomainName: domain.DomainName,
						ProfileID:  profile.ProfileId,
						FirstName:  profile.FirstName,
						LastName:   profile.LastName,
					})
				}

				if searchResp.NextToken == nil {
					break
				}
				searchParams.NextToken = searchResp.NextToken
			}
		}

		if domainResp.NextToken == nil {
			break
		}
		domainParams.NextToken = domainResp.NextToken
	}

	return resources, nil
}

type CustomerProfilesProfile struct {
	svc        CustomerProfilesClient
	DomainName *string
	ProfileID  *string
	FirstName  *string
	LastName   *string
}

func (r *CustomerProfilesProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteProfile(ctx, &customerprofiles.DeleteProfileInput{
		DomainName: r.DomainName,
		ProfileId:  r.ProfileID,
	})
	return err
}

func (r *CustomerProfilesProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CustomerProfilesProfile) String() string {
	return *r.ProfileID
}
