package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lstypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LightsailContactMethodResource = "LightsailContactMethod"

func init() {
	registry.Register(&registry.Registration{
		Name:     LightsailContactMethodResource,
		Scope:    nuke.Account,
		Resource: &LightsailContactMethod{},
		Lister:   &LightsailContactMethodLister{},
	})
}

type LightsailContactMethodLister struct {
	svc LightsailClient
}

func (l *LightsailContactMethodLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = lightsail.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	resp, err := svc.GetContactMethods(ctx, &lightsail.GetContactMethodsInput{})
	if err != nil {
		return nil, err
	}

	for i := range resp.ContactMethods {
		cm := &resp.ContactMethods[i]
		resources = append(resources, &LightsailContactMethod{
			svc:             svc,
			Protocol:        aws.String(string(cm.Protocol)),
			ContactEndpoint: cm.ContactEndpoint,
		})
	}

	return resources, nil
}

type LightsailContactMethod struct {
	svc             LightsailClient
	Protocol        *string
	ContactEndpoint *string
}

func (r *LightsailContactMethod) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteContactMethod(ctx, &lightsail.DeleteContactMethodInput{
		Protocol: lstypes.ContactProtocol(*r.Protocol),
	})
	return err
}

func (r *LightsailContactMethod) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LightsailContactMethod) String() string {
	return *r.ContactEndpoint
}
