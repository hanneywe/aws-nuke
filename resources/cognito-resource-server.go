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

const CognitoResourceServerResource = "CognitoResourceServer"

func init() {
	registry.Register(&registry.Registration{
		Name:     CognitoResourceServerResource,
		Scope:    nuke.Account,
		Resource: &CognitoResourceServer{},
		Lister:   &CognitoResourceServerLister{},
	})
}

type CognitoResourceServerLister struct {
	svc CognitoClient
}

func (l *CognitoResourceServerLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = cognitoidentityprovider.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// List all user pools
	poolParams := &cognitoidentityprovider.ListUserPoolsInput{MaxResults: aws.Int32(60)}
	for {
		poolResp, err := svc.ListUserPools(ctx, poolParams)
		if err != nil {
			return nil, err
		}
		for _, pool := range poolResp.UserPools {
			// List resource servers per pool
			rsParams := &cognitoidentityprovider.ListResourceServersInput{
				UserPoolId: pool.Id,
			}
			for {
				rsResp, err := svc.ListResourceServers(ctx, rsParams)
				if err != nil {
					return nil, err
				}
				for _, rs := range rsResp.ResourceServers {
					resources = append(resources, &CognitoResourceServer{
						svc:        svc,
						Identifier: rs.Identifier,
						Name:       rs.Name,
						UserPoolID: pool.Id,
					})
				}
				if rsResp.NextToken == nil {
					break
				}
				rsParams.NextToken = rsResp.NextToken
			}
		}
		if poolResp.NextToken == nil {
			break
		}
		poolParams.NextToken = poolResp.NextToken
	}
	return resources, nil
}

type CognitoResourceServer struct {
	svc        CognitoClient
	Identifier *string
	Name       *string
	UserPoolID *string `property:"name=UserPoolId"`
}

func (r *CognitoResourceServer) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteResourceServer(ctx, &cognitoidentityprovider.DeleteResourceServerInput{
		Identifier: r.Identifier,
		UserPoolId: r.UserPoolID,
	})
	return err
}

func (r *CognitoResourceServer) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CognitoResourceServer) String() string {
	return *r.Name
}
