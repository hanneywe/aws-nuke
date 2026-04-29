package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CognitoManagedLoginBrandingResource = "CognitoManagedLoginBranding"

func init() {
	registry.Register(&registry.Registration{
		Name:     CognitoManagedLoginBrandingResource,
		Scope:    nuke.Account,
		Resource: &CognitoManagedLoginBranding{},
		Lister:   &CognitoManagedLoginBrandingLister{},
	})
}

type CognitoManagedLoginBrandingLister struct {
	svc CognitoClient
}

func (l *CognitoManagedLoginBrandingLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cognitoidentityprovider.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// Step 1: List all user pools
	poolParams := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int32(60),
	}
	for {
		poolOutput, err := svc.ListUserPools(ctx, poolParams)
		if err != nil {
			return nil, err
		}

		// Step 2: For each pool, list clients
		for _, pool := range poolOutput.UserPools {
			clientParams := &cognitoidentityprovider.ListUserPoolClientsInput{
				UserPoolId: pool.Id,
				MaxResults: aws.Int32(60),
			}
			for {
				clientOutput, err := svc.ListUserPoolClients(ctx, clientParams)
				if err != nil {
					return nil, err
				}

				// Step 3: For each client, describe branding
				for _, client := range clientOutput.UserPoolClients {
					brandingOutput, err := svc.DescribeManagedLoginBrandingByClient(ctx,
						&cognitoidentityprovider.DescribeManagedLoginBrandingByClientInput{
							ClientId:   client.ClientId,
							UserPoolId: pool.Id,
						})
					if err != nil {
						// If no branding exists for this client, skip
						continue
					}
					if brandingOutput.ManagedLoginBranding != nil && brandingOutput.ManagedLoginBranding.ManagedLoginBrandingId != nil {
						resources = append(resources, &CognitoManagedLoginBranding{
							svc:                    svc,
							ManagedLoginBrandingID: brandingOutput.ManagedLoginBranding.ManagedLoginBrandingId,
							UserPoolID:             pool.Id,
						})
					}
				}

				if clientOutput.NextToken == nil {
					break
				}
				clientParams.NextToken = clientOutput.NextToken
			}
		}

		if poolOutput.NextToken == nil {
			break
		}
		poolParams.NextToken = poolOutput.NextToken
	}

	return resources, nil
}

type CognitoManagedLoginBranding struct {
	svc                    CognitoClient
	ManagedLoginBrandingID *string `property:"name=ManagedLoginBrandingId"`
	UserPoolID             *string `property:"name=UserPoolId"`
}

func (r *CognitoManagedLoginBranding) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteManagedLoginBranding(ctx, &cognitoidentityprovider.DeleteManagedLoginBrandingInput{
		ManagedLoginBrandingId: r.ManagedLoginBrandingID,
		UserPoolId:             r.UserPoolID,
	})
	return err
}

func (r *CognitoManagedLoginBranding) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CognitoManagedLoginBranding) String() string {
	return *r.ManagedLoginBrandingID
}
