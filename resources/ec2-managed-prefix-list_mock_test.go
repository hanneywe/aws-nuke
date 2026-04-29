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

func Test_Mock_EC2ManagedPrefixList_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeManagedPrefixLists", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeManagedPrefixListsOutput{
				PrefixLists: []ec2types.ManagedPrefixList{
					{
						PrefixListId:   ptr.String("pl-11111111111111111"),
						PrefixListName: ptr.String("my-prefix-list"),
						OwnerId:        ptr.String("123456789012"),
						AddressFamily:  ptr.String("IPv4"),
						State:          ec2types.PrefixListStateCreateComplete,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-pl")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2ManagedPrefixListLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	prefixList := resources[0].(*EC2ManagedPrefixList)
	assertions.Equal("pl-11111111111111111", *prefixList.PrefixListID)
	assertions.Equal("my-prefix-list", *prefixList.PrefixListName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2ManagedPrefixList_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeManagedPrefixLists", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeManagedPrefixListsOutput{
				PrefixLists: []ec2types.ManagedPrefixList{},
			}, nil,
		)

	lister := &EC2ManagedPrefixListLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2ManagedPrefixList_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.ManagedPrefixList, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.ManagedPrefixList{
			PrefixListId:   ptr.String(fmt.Sprintf("pl-%d", i)),
			PrefixListName: ptr.String(fmt.Sprintf("prefix-list-%d", i)),
			OwnerId:        ptr.String("123456789012"),
			AddressFamily:  ptr.String("IPv4"),
			State:          ec2types.PrefixListStateCreateComplete,
		}
	}

	mockClient.
		On(
			"DescribeManagedPrefixLists",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeManagedPrefixListsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeManagedPrefixListsOutput{
				PrefixLists: firstPageItems,
				NextToken:   ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeManagedPrefixLists",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeManagedPrefixListsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeManagedPrefixListsOutput{
				PrefixLists: []ec2types.ManagedPrefixList{
					{
						PrefixListId:   ptr.String("pl-100"),
						PrefixListName: ptr.String("prefix-list-100"),
						OwnerId:        ptr.String("123456789012"),
						AddressFamily:  ptr.String("IPv6"),
						State:          ec2types.PrefixListStateCreateComplete,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2ManagedPrefixListLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2ManagedPrefixList_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	prefixList := &EC2ManagedPrefixList{
		svc:          mockClient,
		PrefixListID: ptr.String("pl-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteManagedPrefixList",
			mock.Anything,
			&ec2.DeleteManagedPrefixListInput{
				PrefixListId: prefixList.PrefixListID,
			},
		).
		Return(&ec2.DeleteManagedPrefixListOutput{}, nil)

	err := prefixList.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2ManagedPrefixList_Properties(t *testing.T) {
	assertions := assert.New(t)

	prefixList := EC2ManagedPrefixList{
		PrefixListID:   ptr.String("pl-11111111111111111"),
		PrefixListName: ptr.String("my-prefix-list"),
		OwnerID:        ptr.String("123456789012"),
		AddressFamily:  ptr.String("IPv4"),
		State:          "create-complete",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := prefixList.Properties()

	assertions.Equal("pl-11111111111111111", properties.Get("PrefixListId"))
	assertions.Equal("my-prefix-list", properties.Get("PrefixListName"))
	assertions.Equal("123456789012", properties.Get("OwnerId"))
	assertions.Equal("IPv4", properties.Get("AddressFamily"))
	assertions.Equal("create-complete", properties.Get("State"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2ManagedPrefixList_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	prefixList := EC2ManagedPrefixList{
		PrefixListID:   ptr.String("pl-99999999999999999"),
		PrefixListName: ptr.String("no-tags-pl"),
		OwnerID:        ptr.String("123456789012"),
		AddressFamily:  ptr.String("IPv6"),
		State:          "create-complete",
		Tags:           []ec2types.Tag{},
	}

	properties := prefixList.Properties()

	assertions.Equal("pl-99999999999999999", properties.Get("PrefixListId"))
	assertions.Equal("IPv6", properties.Get("AddressFamily"))
}

func Test_Mock_EC2ManagedPrefixList_String(t *testing.T) {
	assertions := assert.New(t)

	prefixList := EC2ManagedPrefixList{
		PrefixListID: ptr.String("pl-11111111111111111"),
	}

	assertions.Equal("pl-11111111111111111", prefixList.String())
}

func Test_Mock_EC2ManagedPrefixList_Filter_ExcludesAWSManaged(t *testing.T) {
	assertions := assert.New(t)

	prefixList := EC2ManagedPrefixList{
		PrefixListID: ptr.String("pl-aws-managed"),
		OwnerID:      ptr.String("amazon"),
		State:        "create-complete",
		accountID:    ptr.String("123456789012"),
	}

	err := prefixList.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "AWS-managed")
}

func Test_Mock_EC2ManagedPrefixList_Filter_ExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	prefixList := EC2ManagedPrefixList{
		PrefixListID: ptr.String("pl-deleted"),
		OwnerID:      ptr.String("123456789012"),
		State:        string(ec2types.PrefixListStateDeleteComplete),
		accountID:    ptr.String("123456789012"),
	}

	err := prefixList.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2ManagedPrefixList_Filter_PassesCustomerOwned(t *testing.T) {
	assertions := assert.New(t)

	prefixList := EC2ManagedPrefixList{
		PrefixListID: ptr.String("pl-customer"),
		OwnerID:      ptr.String("123456789012"),
		State:        "create-complete",
		accountID:    ptr.String("123456789012"),
	}

	err := prefixList.Filter()
	assertions.NoError(err)
}
