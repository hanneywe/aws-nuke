package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2VpcEncryptionControlResource = "EC2VpcEncryptionControl"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2VpcEncryptionControlResource,
		Scope:    nuke.Account,
		Resource: &EC2VpcEncryptionControl{},
		Lister:   &EC2VpcEncryptionControlLister{},
	})
}

type EC2VpcEncryptionControlLister struct {
	svc EC2Client
}

func (l *EC2VpcEncryptionControlLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	input := &ec2.DescribeVpcEncryptionControlsInput{
		MaxResults: aws.Int32(100),
	}

	for {
		resp, err := svc.DescribeVpcEncryptionControls(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, control := range resp.VpcEncryptionControls {
			resources = append(resources, &EC2VpcEncryptionControl{
				svc:                    svc,
				VpcID:                  control.VpcId,
				VpcEncryptionControlID: control.VpcEncryptionControlId,
			})
		}

		if resp.NextToken == nil {
			break
		}
		input.NextToken = resp.NextToken
	}

	return resources, nil
}

type EC2VpcEncryptionControl struct {
	svc                    EC2Client
	VpcID                  *string `property:"name=VpcId"`
	VpcEncryptionControlID *string `property:"name=VpcEncryptionControlId"`
}

func (r *EC2VpcEncryptionControl) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteVpcEncryptionControl(ctx, &ec2.DeleteVpcEncryptionControlInput{
		VpcEncryptionControlId: r.VpcEncryptionControlID,
	})
	return err
}

func (r *EC2VpcEncryptionControl) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2VpcEncryptionControl) String() string {
	return *r.VpcEncryptionControlID
}
