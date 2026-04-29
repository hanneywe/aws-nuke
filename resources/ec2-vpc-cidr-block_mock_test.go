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

func Test_Mock_EC2VPCCIDRBlock_List_IPv4(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeVpcs", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						VpcId:     ptr.String("vpc-11111111111111111"),
						CidrBlock: ptr.String("10.0.0.0/16"),
						CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
							{
								CidrBlock:     ptr.String("10.0.0.0/16"),
								AssociationId: ptr.String("vpc-cidr-assoc-primary"),
								CidrBlockState: &ec2types.VpcCidrBlockState{
									State: ec2types.VpcCidrBlockStateCodeAssociated,
								},
							},
							{
								CidrBlock:     ptr.String("10.1.0.0/16"),
								AssociationId: ptr.String("vpc-cidr-assoc-secondary"),
								CidrBlockState: &ec2types.VpcCidrBlockState{
									State: ec2types.VpcCidrBlockStateCodeAssociated,
								},
							},
						},
					},
				},
			}, nil,
		)

	lister := &EC2VPCCIDRBlockLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	cidrBlock := resources[0].(*EC2VPCCIDRBlock)
	assertions.Equal("10.1.0.0/16", *cidrBlock.CidrBlock)
	assertions.Equal("vpc-cidr-assoc-secondary", *cidrBlock.AssociationID)
	assertions.Equal("vpc-11111111111111111", *cidrBlock.VpcID)
	assertions.Equal(false, *cidrBlock.IsIPv6)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VPCCIDRBlock_List_IPv6(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeVpcs", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						VpcId:     ptr.String("vpc-22222222222222222"),
						CidrBlock: ptr.String("10.0.0.0/16"),
						CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
							{
								CidrBlock:     ptr.String("10.0.0.0/16"),
								AssociationId: ptr.String("vpc-cidr-assoc-primary"),
								CidrBlockState: &ec2types.VpcCidrBlockState{
									State: ec2types.VpcCidrBlockStateCodeAssociated,
								},
							},
						},
						Ipv6CidrBlockAssociationSet: []ec2types.VpcIpv6CidrBlockAssociation{
							{
								Ipv6CidrBlock: ptr.String("2600:1f18::/56"),
								AssociationId: ptr.String("vpc-cidr-assoc-ipv6"),
								Ipv6CidrBlockState: &ec2types.VpcCidrBlockState{
									State: ec2types.VpcCidrBlockStateCodeAssociated,
								},
							},
						},
					},
				},
			}, nil,
		)

	lister := &EC2VPCCIDRBlockLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	cidrBlock := resources[0].(*EC2VPCCIDRBlock)
	assertions.Equal("2600:1f18::/56", *cidrBlock.CidrBlock)
	assertions.Equal("vpc-cidr-assoc-ipv6", *cidrBlock.AssociationID)
	assertions.Equal(true, *cidrBlock.IsIPv6)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VPCCIDRBlock_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeVpcs", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						VpcId:     ptr.String("vpc-33333333333333333"),
						CidrBlock: ptr.String("10.0.0.0/16"),
						CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
							{
								CidrBlock:     ptr.String("10.0.0.0/16"),
								AssociationId: ptr.String("vpc-cidr-assoc-primary"),
								CidrBlockState: &ec2types.VpcCidrBlockState{
									State: ec2types.VpcCidrBlockStateCodeAssociated,
								},
							},
						},
					},
				},
			}, nil,
		)

	lister := &EC2VPCCIDRBlockLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VPCCIDRBlock_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On(
			"DescribeVpcs",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeVpcsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						VpcId:     ptr.String("vpc-page1"),
						CidrBlock: ptr.String("10.0.0.0/16"),
						CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
							{
								CidrBlock:     ptr.String("10.0.0.0/16"),
								AssociationId: ptr.String("vpc-cidr-assoc-primary"),
								CidrBlockState: &ec2types.VpcCidrBlockState{
									State: ec2types.VpcCidrBlockStateCodeAssociated,
								},
							},
							{
								CidrBlock:     ptr.String("10.1.0.0/16"),
								AssociationId: ptr.String("vpc-cidr-assoc-secondary"),
								CidrBlockState: &ec2types.VpcCidrBlockState{
									State: ec2types.VpcCidrBlockStateCodeAssociated,
								},
							},
						},
					},
				},
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeVpcs",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeVpcsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						VpcId:     ptr.String("vpc-page2"),
						CidrBlock: ptr.String("172.16.0.0/16"),
						Ipv6CidrBlockAssociationSet: []ec2types.VpcIpv6CidrBlockAssociation{
							{
								Ipv6CidrBlock: ptr.String("2600:1f18::/56"),
								AssociationId: ptr.String("vpc-cidr-assoc-ipv6"),
								Ipv6CidrBlockState: &ec2types.VpcCidrBlockState{
									State: ec2types.VpcCidrBlockStateCodeAssociated,
								},
							},
						},
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2VPCCIDRBlockLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VPCCIDRBlock_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	cidrBlock := &EC2VPCCIDRBlock{
		svc:           mockClient,
		VpcID:         ptr.String("vpc-11111111111111111"),
		AssociationID: ptr.String("vpc-cidr-assoc-secondary"),
	}

	mockClient.
		On(
			"DisassociateVpcCidrBlock",
			mock.Anything,
			&ec2.DisassociateVpcCidrBlockInput{
				AssociationId: cidrBlock.AssociationID,
			},
		).
		Return(&ec2.DisassociateVpcCidrBlockOutput{}, nil)

	err := cidrBlock.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VPCCIDRBlock_Properties(t *testing.T) {
	assertions := assert.New(t)

	cidrBlock := EC2VPCCIDRBlock{
		VpcID:         ptr.String("vpc-11111111111111111"),
		CidrBlock:     ptr.String("10.1.0.0/16"),
		AssociationID: ptr.String("vpc-cidr-assoc-secondary"),
		State:         "associated",
		IsIPv6:        ptr.Bool(false),
	}

	properties := cidrBlock.Properties()

	assertions.Equal("vpc-11111111111111111", properties.Get("VpcId"))
	assertions.Equal("10.1.0.0/16", properties.Get("CidrBlock"))
	assertions.Equal("vpc-cidr-assoc-secondary", properties.Get("AssociationId"))
	assertions.Equal("associated", properties.Get("State"))
	assertions.Equal("false", properties.Get("IsIPv6"))
}

func Test_Mock_EC2VPCCIDRBlock_String(t *testing.T) {
	assertions := assert.New(t)

	cidrBlock := EC2VPCCIDRBlock{
		VpcID:     ptr.String("vpc-11111111111111111"),
		CidrBlock: ptr.String("10.1.0.0/16"),
	}

	assertions.Equal("10.1.0.0/16 (vpc-11111111111111111)", cidrBlock.String())
}

func Test_Mock_EC2VPCCIDRBlock_Filter_ExcludesDisassociatedState(t *testing.T) {
	assertions := assert.New(t)

	cidrBlock := EC2VPCCIDRBlock{
		VpcID:         ptr.String("vpc-disassociated"),
		CidrBlock:     ptr.String("10.1.0.0/16"),
		AssociationID: ptr.String("vpc-cidr-assoc-old"),
		State:         string(ec2types.VpcCidrBlockStateCodeDisassociated),
	}

	err := cidrBlock.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already disassociated")
}

func Test_Mock_EC2VPCCIDRBlock_Filter_PassesAssociatedState(t *testing.T) {
	assertions := assert.New(t)

	cidrBlock := EC2VPCCIDRBlock{
		VpcID:         ptr.String("vpc-active"),
		CidrBlock:     ptr.String("10.1.0.0/16"),
		AssociationID: ptr.String("vpc-cidr-assoc-active"),
		State:         string(ec2types.VpcCidrBlockStateCodeAssociated),
	}

	err := cidrBlock.Filter()
	assertions.NoError(err)
}
