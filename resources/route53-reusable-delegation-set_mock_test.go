package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func Test_Mock_Route53ReusableDelegationSet_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53Client)

	mockClient.On("ListReusableDelegationSets", mock.Anything, mock.Anything).
		Return(&route53.ListReusableDelegationSetsOutput{
			DelegationSets: []route53types.DelegationSet{
				{
					Id:              ptr.String("/delegationset/N12345"),
					CallerReference: ptr.String("my-ref"),
				},
			},
		}, nil)

	lister := &Route53ReusableDelegationSetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	delegationSet := resources[0].(*Route53ReusableDelegationSet)
	assertions.Equal("/delegationset/N12345", *delegationSet.ID)
	assertions.Equal("my-ref", *delegationSet.CallerReference)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53ReusableDelegationSet_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53Client)

	mockClient.On("ListReusableDelegationSets", mock.Anything, mock.Anything).
		Return(&route53.ListReusableDelegationSetsOutput{
			DelegationSets: []route53types.DelegationSet{},
		}, nil)

	lister := &Route53ReusableDelegationSetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53ReusableDelegationSet_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53Client)

	delegationSet := &Route53ReusableDelegationSet{
		svc: mockClient,
		ID:  ptr.String("/delegationset/N12345"),
	}

	mockClient.On("DeleteReusableDelegationSet", mock.Anything, &route53.DeleteReusableDelegationSetInput{
		Id: delegationSet.ID,
	}).Return(&route53.DeleteReusableDelegationSetOutput{}, nil)

	err := delegationSet.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53ReusableDelegationSet_Properties(t *testing.T) {
	assertions := assert.New(t)

	delegationSet := Route53ReusableDelegationSet{
		ID:              ptr.String("/delegationset/N12345"),
		CallerReference: ptr.String("my-ref"),
	}

	properties := delegationSet.Properties()
	assertions.Equal("/delegationset/N12345", properties.Get("Id"))
	assertions.Equal("my-ref", properties.Get("CallerReference"))
}

func Test_Mock_Route53ReusableDelegationSet_String(t *testing.T) {
	assertions := assert.New(t)
	delegationSet := Route53ReusableDelegationSet{ID: ptr.String("/delegationset/N12345")}
	assertions.Equal("/delegationset/N12345", delegationSet.String())
}
