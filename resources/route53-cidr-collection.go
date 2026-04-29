package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Route53CidrCollectionResource = "Route53CidrCollection"

func init() {
	registry.Register(&registry.Registration{
		Name:     Route53CidrCollectionResource,
		Scope:    nuke.Account,
		Resource: &Route53CidrCollection{},
		Lister:   &Route53CidrCollectionLister{},
	})
}

type Route53CidrCollectionLister struct {
	svc Route53Client
}

func (l *Route53CidrCollectionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = route53.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &route53.ListCidrCollectionsInput{}
	for {
		listOutput, err := svc.ListCidrCollections(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, cidrCollection := range listOutput.CidrCollections {
			resources = append(resources, &Route53CidrCollection{
				svc:  svc,
				ID:   cidrCollection.Id,
				Name: cidrCollection.Name,
			})
		}

		if listOutput.NextToken == nil {
			break
		}
		params.NextToken = listOutput.NextToken
	}

	return resources, nil
}

type Route53CidrCollection struct {
	svc  Route53Client
	ID   *string
	Name *string
}

func (r *Route53CidrCollection) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCidrCollection(ctx, &route53.DeleteCidrCollectionInput{
		Id: r.ID,
	})
	return err
}

func (r *Route53CidrCollection) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Route53CidrCollection) String() string {
	return *r.Name
}
