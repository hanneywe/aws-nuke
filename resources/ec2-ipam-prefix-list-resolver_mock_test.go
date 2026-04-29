package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2IPAMPrefixListResolver_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeIpamPrefixListResolvers", mock.Anything, mock.Anything).
		Return(&ec2.DescribeIpamPrefixListResolversOutput{
			IpamPrefixListResolvers: []ec2types.IpamPrefixListResolver{
				{
					IpamPrefixListResolverId: ptr.String("ipam-plr-11111111111111111"),
					IpamArn:                  ptr.String("arn:aws:ec2::123456789012:ipam/ipam-1"),
					OwnerId:                  ptr.String("123456789012"),
					State:                    ec2types.IpamPrefixListResolverStateCreateComplete,
					Tags: []ec2types.Tag{
						{Key: ptr.String("Name"), Value: ptr.String("test-resolver")},
					},
				},
			},
		}, nil)

	lister := &EC2IPAMPrefixListResolverLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*EC2IPAMPrefixListResolver)
	a.Equal("ipam-plr-11111111111111111", *r.IpamPrefixListResolverID)
	a.Equal("123456789012", *r.OwnerID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMPrefixListResolver_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeIpamPrefixListResolvers", mock.Anything, mock.Anything).
		Return(&ec2.DescribeIpamPrefixListResolversOutput{
			IpamPrefixListResolvers: []ec2types.IpamPrefixListResolver{},
		}, nil)

	lister := &EC2IPAMPrefixListResolverLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMPrefixListResolver_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)

	r := &EC2IPAMPrefixListResolver{
		svc:                      mockClient,
		IpamPrefixListResolverID: ptr.String("ipam-plr-11111111111111111"),
	}

	mockClient.
		On("DeleteIpamPrefixListResolver", mock.Anything, &ec2.DeleteIpamPrefixListResolverInput{
			IpamPrefixListResolverId: r.IpamPrefixListResolverID,
		}).
		Return(&ec2.DeleteIpamPrefixListResolverOutput{}, nil)

	err := r.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMPrefixListResolver_Properties(t *testing.T) {
	a := assert.New(t)

	r := EC2IPAMPrefixListResolver{
		IpamPrefixListResolverID: ptr.String("ipam-plr-11111111111111111"),
		IpamArn:                  ptr.String("arn:aws:ec2::123456789012:ipam/ipam-1"),
		OwnerID:                  ptr.String("123456789012"),
		State:                    "create-complete",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Env"), Value: ptr.String("prod")},
		},
	}

	props := r.Properties()
	a.Equal("ipam-plr-11111111111111111", props.Get("IpamPrefixListResolverId"))
	a.Equal("arn:aws:ec2::123456789012:ipam/ipam-1", props.Get("IpamArn"))
	a.Equal("123456789012", props.Get("OwnerId"))
	a.Equal("create-complete", props.Get("State"))
	a.Equal("prod", props.Get("tag:Env"))
}

func Test_Mock_EC2IPAMPrefixListResolver_String(t *testing.T) {
	a := assert.New(t)

	r := EC2IPAMPrefixListResolver{
		IpamPrefixListResolverID: ptr.String("ipam-plr-11111111111111111"),
	}

	a.Equal("ipam-plr-11111111111111111", r.String())
}

func Test_Mock_EC2IPAMPrefixListResolver_Filter_ExcludesDeleted(t *testing.T) {
	a := assert.New(t)

	r := EC2IPAMPrefixListResolver{
		IpamPrefixListResolverID: ptr.String("ipam-plr-deleted"),
		State:                    string(ec2types.IpamPrefixListResolverStateDeleteComplete),
	}

	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2IPAMPrefixListResolver_Filter_PassesActive(t *testing.T) {
	a := assert.New(t)

	r := EC2IPAMPrefixListResolver{
		IpamPrefixListResolverID: ptr.String("ipam-plr-active"),
		State:                    string(ec2types.IpamPrefixListResolverStateCreateComplete),
	}

	err := r.Filter()
	a.NoError(err)
}
