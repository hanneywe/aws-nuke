package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/signer"
	signertypes "github.com/aws/aws-sdk-go-v2/service/signer/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SignerSigningProfileResource = "SignerSigningProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     SignerSigningProfileResource,
		Scope:    nuke.Account,
		Resource: &SignerSigningProfile{},
		Lister:   &SignerSigningProfileLister{},
	})
}

type SignerSigningProfileLister struct {
	svc SignerClient
}

func (l *SignerSigningProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = signer.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &signer.ListSigningProfilesInput{}
	for {
		output, err := svc.ListSigningProfiles(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, profile := range output.Profiles {
			profileStatus := string(profile.Status)
			resources = append(resources, &SignerSigningProfile{
				svc:            svc,
				ProfileName:    profile.ProfileName,
				ProfileVersion: profile.ProfileVersion,
				Status:         &profileStatus,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type SignerSigningProfile struct {
	svc            SignerClient
	ProfileName    *string
	ProfileVersion *string
	Status         *string
}

func (r *SignerSigningProfile) Remove(ctx context.Context) error {
	_, err := r.svc.CancelSigningProfile(ctx, &signer.CancelSigningProfileInput{
		ProfileName: r.ProfileName,
	})
	return err
}

func (r *SignerSigningProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SignerSigningProfile) String() string {
	return *r.ProfileName
}

func (r *SignerSigningProfile) Filter() error {
	if r.Status != nil {
		status := signertypes.SigningProfileStatus(*r.Status)
		if status == signertypes.SigningProfileStatusCanceled || status == signertypes.SigningProfileStatusRevoked {
			return fmt.Errorf("already %s", *r.Status)
		}
	}
	return nil
}
