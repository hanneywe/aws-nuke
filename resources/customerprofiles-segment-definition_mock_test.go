package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/customerprofiles"
	customerprofilestypes "github.com/aws/aws-sdk-go-v2/service/customerprofiles/types"
)

func Test_Mock_CustomerProfilesSegmentDefinition_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	mockClient.
		On("ListDomains", mock.Anything, mock.Anything).
		Return(&customerprofiles.ListDomainsOutput{
			Items: []customerprofilestypes.ListDomainItem{
				{
					DomainName: ptr.String("my-domain"),
				},
			},
		}, nil)

	mockClient.
		On("ListSegmentDefinitions", mock.Anything, mock.Anything).
		Return(&customerprofiles.ListSegmentDefinitionsOutput{
			Items: []customerprofilestypes.SegmentDefinitionItem{
				{
					SegmentDefinitionName: ptr.String("my-segment"),
					DisplayName:           ptr.String("My Segment"),
				},
			},
		}, nil)

	lister := &CustomerProfilesSegmentDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCustomerProfilesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	seg := resources[0].(*CustomerProfilesSegmentDefinition)
	a.Equal("my-domain", *seg.DomainName)
	a.Equal("my-segment", *seg.SegmentDefinitionName)
	a.Equal("My Segment", *seg.DisplayName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesSegmentDefinition_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	mockClient.
		On("ListDomains", mock.Anything, mock.Anything).
		Return(&customerprofiles.ListDomainsOutput{
			Items: []customerprofilestypes.ListDomainItem{},
		}, nil)

	lister := &CustomerProfilesSegmentDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCustomerProfilesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesSegmentDefinition_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	seg := &CustomerProfilesSegmentDefinition{
		svc:                   mockClient,
		DomainName:            ptr.String("my-domain"),
		SegmentDefinitionName: ptr.String("my-segment"),
		DisplayName:           ptr.String("My Segment"),
	}

	mockClient.
		On("DeleteSegmentDefinition", mock.Anything, &customerprofiles.DeleteSegmentDefinitionInput{
			DomainName:            seg.DomainName,
			SegmentDefinitionName: seg.SegmentDefinitionName,
		}).
		Return(&customerprofiles.DeleteSegmentDefinitionOutput{}, nil)

	err := seg.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesSegmentDefinition_Properties(t *testing.T) {
	a := assert.New(t)

	seg := CustomerProfilesSegmentDefinition{
		DomainName:            ptr.String("my-domain"),
		SegmentDefinitionName: ptr.String("my-segment"),
		DisplayName:           ptr.String("My Segment"),
	}

	props := seg.Properties()
	a.Equal("my-domain", props.Get("DomainName"))
	a.Equal("my-segment", props.Get("SegmentDefinitionName"))
	a.Equal("My Segment", props.Get("DisplayName"))
}

func Test_Mock_CustomerProfilesSegmentDefinition_String(t *testing.T) {
	a := assert.New(t)

	seg := CustomerProfilesSegmentDefinition{
		SegmentDefinitionName: ptr.String("my-segment"),
	}

	a.Equal("my-segment", seg.String())
}
