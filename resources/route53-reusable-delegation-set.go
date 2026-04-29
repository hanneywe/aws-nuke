package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Route53ReusableDelegationSetResource = "Route53ReusableDelegationSet"

func init() {
	registry.Register(&registry.Registration{
		Name:     Route53ReusableDelegationSetResource,
		Scope:    nuke.Account,
		Resource: &Route53ReusableDelegationSet{},
		Lister:   &Route53ReusableDelegationSetLister{},
	})
}

type Route53ReusableDelegationSetLister struct {
	svc Route53Client
}

func (l *Route53ReusableDelegationSetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = route53.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	output, err := svc.ListReusableDelegationSets(ctx, &route53.ListReusableDelegationSetsInput{})
	if err != nil {
		return nil, err
	}

	for _, delegationSet := range output.DelegationSets {
		resources = append(resources, &Route53ReusableDelegationSet{
			svc:             svc,
			ID:              delegationSet.Id,
			CallerReference: delegationSet.CallerReference,
		})
	}

	return resources, nil
}

type Route53ReusableDelegationSet struct {
	svc             Route53Client
	ID              *string `property:"name=Id"`
	CallerReference *string
}

func (r *Route53ReusableDelegationSet) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteReusableDelegationSet(ctx, &route53.DeleteReusableDelegationSetInput{
		Id: r.ID,
	})
	return err
}

func (r *Route53ReusableDelegationSet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Route53ReusableDelegationSet) String() string {
	return *r.ID
}
