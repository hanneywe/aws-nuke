package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/panorama"
	panoramatypes "github.com/aws/aws-sdk-go-v2/service/panorama/types"
)

func Test_Mock_PanoramaPackage_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPanoramaClient)
	mockClient.On("ListPackages", mock.Anything, mock.Anything).
		Return(&panorama.ListPackagesOutput{
			Packages: []panoramatypes.PackageListItem{
				{PackageId: ptr.String("pkg-12345"), PackageName: ptr.String("my-package")},
			},
		}, nil)
	lister := &PanoramaPackageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPanoramaListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	pkg := resources[0].(*PanoramaPackage)
	a.Equal("pkg-12345", *pkg.PackageID)
	a.Equal("my-package", *pkg.PackageName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_PanoramaPackage_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPanoramaClient)
	mockClient.On("ListPackages", mock.Anything, mock.Anything).
		Return(&panorama.ListPackagesOutput{Packages: []panoramatypes.PackageListItem{}}, nil)
	lister := &PanoramaPackageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPanoramaListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_PanoramaPackage_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPanoramaClient)
	pkg := &PanoramaPackage{svc: mockClient, PackageID: ptr.String("pkg-12345")}
	mockClient.On("DeletePackage", mock.Anything, &panorama.DeletePackageInput{PackageId: pkg.PackageID, ForceDelete: true}).
		Return(&panorama.DeletePackageOutput{}, nil)
	a.NoError(pkg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_PanoramaPackage_Properties(t *testing.T) {
	a := assert.New(t)
	pkg := PanoramaPackage{PackageID: ptr.String("pkg-12345"), PackageName: ptr.String("my-package")}
	a.Equal("pkg-12345", pkg.Properties().Get("PackageId"))
	a.Equal("my-package", pkg.Properties().Get("PackageName"))
}

func Test_Mock_PanoramaPackage_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-package", (&PanoramaPackage{PackageName: ptr.String("my-package")}).String())
}
