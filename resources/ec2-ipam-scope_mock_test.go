package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2IPAMScope_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeIpamScopes", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeIpamScopesOutput{
				IpamScopes: []ec2types.IpamScope{
					{
						IpamScopeId: ptr.String("ipam-scope-11111111111111111"),
						IpamArn:     ptr.String("arn:aws:ec2::123456789012:ipam/ipam-aaa"),
						IsDefault:   ptr.Bool(false),
						State:       ec2types.IpamScopeStateCreateComplete,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-scope")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2IPAMScopeLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	scope := resources[0].(*EC2IPAMScope)
	assertions.Equal("ipam-scope-11111111111111111", *scope.IpamScopeID)
	assertions.Equal(false, *scope.IsDefault)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMScope_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeIpamScopes", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeIpamScopesOutput{
				IpamScopes: []ec2types.IpamScope{},
			}, nil,
		)

	lister := &EC2IPAMScopeLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMScope_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.IpamScope, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.IpamScope{
			IpamScopeId: ptr.String(fmt.Sprintf("ipam-scope-%d", i)),
			IpamArn:     ptr.String(fmt.Sprintf("arn:aws:ec2::123456789012:ipam/ipam-%d", i)),
			IsDefault:   ptr.Bool(false),
			State:       ec2types.IpamScopeStateCreateComplete,
		}
	}

	mockClient.
		On(
			"DescribeIpamScopes",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeIpamScopesInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeIpamScopesOutput{
				IpamScopes: firstPageItems,
				NextToken:  ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeIpamScopes",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeIpamScopesInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeIpamScopesOutput{
				IpamScopes: []ec2types.IpamScope{
					{
						IpamScopeId: ptr.String("ipam-scope-100"),
						IpamArn:     ptr.String("arn:aws:ec2::123456789012:ipam/ipam-100"),
						IsDefault:   ptr.Bool(false),
						State:       ec2types.IpamScopeStateCreateComplete,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2IPAMScopeLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMScope_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	scope := &EC2IPAMScope{
		svc:         mockClient,
		IpamScopeID: ptr.String("ipam-scope-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteIpamScope",
			mock.Anything,
			&ec2.DeleteIpamScopeInput{
				IpamScopeId: scope.IpamScopeID,
			},
		).
		Return(&ec2.DeleteIpamScopeOutput{}, nil)

	err := scope.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMScope_Properties(t *testing.T) {
	assertions := assert.New(t)

	scope := EC2IPAMScope{
		IpamScopeID: ptr.String("ipam-scope-11111111111111111"),
		IpamID:      ptr.String("arn:aws:ec2::123456789012:ipam/ipam-aaa"),
		IsDefault:   ptr.Bool(false),
		State:       "create-complete",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := scope.Properties()

	assertions.Equal("ipam-scope-11111111111111111", properties.Get("IpamScopeId"))
	assertions.Equal("false", properties.Get("IsDefault"))
	assertions.Equal("create-complete", properties.Get("State"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2IPAMScope_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	scope := EC2IPAMScope{
		IpamScopeID: ptr.String("ipam-scope-99999999999999999"),
		IpamID:      ptr.String("arn:aws:ec2::123456789012:ipam/ipam-zzz"),
		IsDefault:   ptr.Bool(true),
		State:       "create-complete",
		Tags:        []ec2types.Tag{},
	}

	properties := scope.Properties()

	assertions.Equal("ipam-scope-99999999999999999", properties.Get("IpamScopeId"))
	assertions.Equal("true", properties.Get("IsDefault"))
}

func Test_Mock_EC2IPAMScope_String(t *testing.T) {
	assertions := assert.New(t)

	scope := EC2IPAMScope{
		IpamScopeID: ptr.String("ipam-scope-11111111111111111"),
	}

	assertions.Equal("ipam-scope-11111111111111111", scope.String())
}

func Test_Mock_EC2IPAMScope_Filter_ExcludesDefaultScope(t *testing.T) {
	assertions := assert.New(t)

	scope := EC2IPAMScope{
		IpamScopeID: ptr.String("ipam-scope-default"),
		IsDefault:   ptr.Bool(true),
		State:       "create-complete",
	}

	err := scope.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "default")
}

func Test_Mock_EC2IPAMScope_Filter_ExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	scope := EC2IPAMScope{
		IpamScopeID: ptr.String("ipam-scope-deleted"),
		IsDefault:   ptr.Bool(false),
		State:       string(ec2types.IpamScopeStateDeleteComplete),
	}

	err := scope.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2IPAMScope_Filter_PassesCustomScope(t *testing.T) {
	assertions := assert.New(t)

	scope := EC2IPAMScope{
		IpamScopeID: ptr.String("ipam-scope-custom"),
		IsDefault:   ptr.Bool(false),
		State:       string(ec2types.IpamScopeStateCreateComplete),
	}

	err := scope.Filter()
	assertions.NoError(err)
}
