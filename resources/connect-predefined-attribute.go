package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectPredefinedAttributeResource = "ConnectPredefinedAttribute"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectPredefinedAttributeResource,
		Scope:    nuke.Account,
		Resource: &ConnectPredefinedAttribute{},
		Lister:   &ConnectPredefinedAttributeLister{},
	})
}

type ConnectPredefinedAttributeLister struct {
	svc ConnectClient
}

func (l *ConnectPredefinedAttributeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = connect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	instanceParams := &connect.ListInstancesInput{}
	for {
		instanceOutput, err := svc.ListInstances(ctx, instanceParams)
		if err != nil {
			return nil, err
		}

		for _, instance := range instanceOutput.InstanceSummaryList {
			attrParams := &connect.ListPredefinedAttributesInput{
				InstanceId: instance.Id,
			}
			for {
				attrOutput, err := svc.ListPredefinedAttributes(ctx, attrParams)
				if err != nil {
					return nil, err
				}

				for _, predefinedAttribute := range attrOutput.PredefinedAttributeSummaryList {
					resources = append(resources, &ConnectPredefinedAttribute{
						svc:        svc,
						Name:       predefinedAttribute.Name,
						InstanceID: instance.Id,
					})
				}

				if attrOutput.NextToken == nil {
					break
				}
				attrParams.NextToken = attrOutput.NextToken
			}
		}

		if instanceOutput.NextToken == nil {
			break
		}
		instanceParams.NextToken = instanceOutput.NextToken
	}

	return resources, nil
}

type ConnectPredefinedAttribute struct {
	svc        ConnectClient
	Name       *string
	InstanceID *string `property:"name=InstanceId"`
}

func (r *ConnectPredefinedAttribute) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePredefinedAttribute(ctx, &connect.DeletePredefinedAttributeInput{
		InstanceId: r.InstanceID,
		Name:       r.Name,
	})
	return err
}

func (r *ConnectPredefinedAttribute) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectPredefinedAttribute) String() string {
	return *r.Name
}
