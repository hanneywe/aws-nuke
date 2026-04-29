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

func Test_Mock_Route53CidrCollection_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53Client)

	mockClient.
		On("ListCidrCollections", mock.Anything, mock.Anything).
		Return(&route53.ListCidrCollectionsOutput{
			CidrCollections: []route53types.CollectionSummary{
				{
					Id:   ptr.String("collection-1"),
					Name: ptr.String("my-cidr-collection"),
				},
			},
		}, nil)

	lister := &Route53CidrCollectionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testRoute53ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	cidrCollection := resources[0].(*Route53CidrCollection)
	assertions.Equal("collection-1", *cidrCollection.ID)
	assertions.Equal("my-cidr-collection", *cidrCollection.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53CidrCollection_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53Client)

	mockClient.
		On("ListCidrCollections", mock.Anything, mock.Anything).
		Return(&route53.ListCidrCollectionsOutput{
			CidrCollections: []route53types.CollectionSummary{},
		}, nil)

	lister := &Route53CidrCollectionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testRoute53ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53CidrCollection_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53Client)

	cidrCollection := &Route53CidrCollection{
		svc:  mockClient,
		ID:   ptr.String("collection-1"),
		Name: ptr.String("my-cidr-collection"),
	}

	mockClient.
		On("DeleteCidrCollection", mock.Anything, &route53.DeleteCidrCollectionInput{
			Id: cidrCollection.ID,
		}).
		Return(&route53.DeleteCidrCollectionOutput{}, nil)

	err := cidrCollection.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53CidrCollection_Properties(t *testing.T) {
	assertions := assert.New(t)

	cidrCollection := Route53CidrCollection{
		ID:   ptr.String("collection-1"),
		Name: ptr.String("my-cidr-collection"),
	}

	properties := cidrCollection.Properties()

	assertions.Equal("collection-1", properties.Get("ID"))
	assertions.Equal("my-cidr-collection", properties.Get("Name"))
}

func Test_Mock_Route53CidrCollection_String(t *testing.T) {
	assertions := assert.New(t)

	cidrCollection := Route53CidrCollection{
		Name: ptr.String("my-cidr-collection"),
	}

	assertions.Equal("my-cidr-collection", cidrCollection.String())
}
