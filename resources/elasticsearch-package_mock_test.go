package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice/types"
)

func Test_Mock_ElasticsearchPackage_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticsearchserviceClient)

	mockClient.On("DescribePackages", mock.Anything, mock.Anything).
		Return(&elasticsearchservice.DescribePackagesOutput{
			PackageDetailsList: []types.PackageDetails{
				{PackageID: ptr.String("test-value"), PackageName: ptr.String("test-value")},
			},
		}, nil)

	lister := &ElasticsearchPackageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testElasticsearchserviceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*ElasticsearchPackage)
	a.Equal("test-value", *r.PackageID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticsearchPackage_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticsearchserviceClient)

	mockClient.On("DescribePackages", mock.Anything, mock.Anything).
		Return(&elasticsearchservice.DescribePackagesOutput{
			PackageDetailsList: []types.PackageDetails{},
		}, nil)

	lister := &ElasticsearchPackageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testElasticsearchserviceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticsearchPackage_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticsearchserviceClient)

	r := &ElasticsearchPackage{
		svc:       mockClient,
		PackageID: ptr.String("test-pkg-01"),
	}

	mockClient.On("DeletePackage", mock.Anything,
		&elasticsearchservice.DeletePackageInput{
			PackageID: r.PackageID,
		}).Return(&elasticsearchservice.DeletePackageOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticsearchPackage_Properties(t *testing.T) {
	a := assert.New(t)
	r := &ElasticsearchPackage{
		PackageID:   ptr.String("test-pkg-01"),
		PackageName: ptr.String("test-pkg-name"),
	}
	props := r.Properties()
	a.Equal("test-pkg-01", props.Get("PackageID"))
	a.Equal("test-pkg-name", props.Get("PackageName"))
}

func Test_Mock_ElasticsearchPackage_String(t *testing.T) {
	a := assert.New(t)
	r := &ElasticsearchPackage{
		PackageID: ptr.String("test-pkg-01"),
	}
	a.Equal("test-pkg-01", r.String())
}
