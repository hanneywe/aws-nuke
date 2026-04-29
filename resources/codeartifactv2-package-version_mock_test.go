package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	codeartifacttypes "github.com/aws/aws-sdk-go-v2/service/codeartifact/types"
)

func Test_Mock_CodeArtifactPackageVersion_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeartifactClient)

	mockClient.On("ListDomains", mock.Anything, mock.Anything).
		Return(&codeartifact.ListDomainsOutput{
			Domains: []codeartifacttypes.DomainSummary{
				{Name: ptr.String("test-domain")},
			},
		}, nil)

	mockClient.On("ListRepositoriesInDomain", mock.Anything, mock.Anything).
		Return(&codeartifact.ListRepositoriesInDomainOutput{
			Repositories: []codeartifacttypes.RepositorySummary{
				{Name: ptr.String("test-repo")},
			},
		}, nil)

	mockClient.On("ListPackages", mock.Anything, mock.Anything).
		Return(&codeartifact.ListPackagesOutput{
			Packages: []codeartifacttypes.PackageSummary{
				{
					Package:   ptr.String("test-package"),
					Format:    codeartifacttypes.PackageFormatNpm,
					Namespace: ptr.String("test-namespace"),
				},
			},
		}, nil)

	mockClient.On("ListPackageVersions", mock.Anything, mock.Anything).
		Return(&codeartifact.ListPackageVersionsOutput{
			Versions: []codeartifacttypes.PackageVersionSummary{
				{
					Version: ptr.String("1.0.0"),
					Status:  codeartifacttypes.PackageVersionStatusPublished,
				},
			},
		}, nil)

	lister := &CodeArtifactPackageVersionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeartifactListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CodeArtifactPackageVersion)
	a.Equal("test-domain", *r.Domain)
	a.Equal("test-repo", *r.Repository)
	a.Equal("test-package", *r.Package)
	a.Equal("npm", *r.Format)
	a.Equal("test-namespace", *r.Namespace)
	a.Equal("1.0.0", *r.Version)
	a.Equal("Published", *r.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeArtifactPackageVersion_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeartifactClient)

	mockClient.On("ListDomains", mock.Anything, mock.Anything).
		Return(&codeartifact.ListDomainsOutput{
			Domains: []codeartifacttypes.DomainSummary{},
		}, nil)

	lister := &CodeArtifactPackageVersionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeartifactListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeArtifactPackageVersion_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeartifactClient)

	r := &CodeArtifactPackageVersion{
		svc:        mockClient,
		Domain:     ptr.String("test-domain"),
		Repository: ptr.String("test-repository"),
		Package:    ptr.String("test-package"),
		Format:     ptr.String("npm"),
		Namespace:  ptr.String("test-namespace"),
		Version:    ptr.String("1.0.0"),
	}

	mockClient.On("DeletePackageVersions", mock.Anything,
		&codeartifact.DeletePackageVersionsInput{
			Domain:     r.Domain,
			Repository: r.Repository,
			Package:    r.Package,
			Format:     codeartifacttypes.PackageFormatNpm,
			Namespace:  r.Namespace,
			Versions:   []string{"1.0.0"},
		}).Return(&codeartifact.DeletePackageVersionsOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeArtifactPackageVersion_Properties(t *testing.T) {
	a := assert.New(t)
	r := &CodeArtifactPackageVersion{
		Domain:     ptr.String("test-domain"),
		Repository: ptr.String("test-repository"),
		Package:    ptr.String("test-package"),
		Format:     ptr.String("test-format"),
		Namespace:  ptr.String("test-namespace"),
		Version:    ptr.String("test-version"),
		Status:     ptr.String("test-status"),
	}
	props := r.Properties()
	a.Equal("test-domain", props.Get("Domain"))
	a.Equal("test-repository", props.Get("Repository"))
	a.Equal("test-package", props.Get("Package"))
	a.Equal("test-format", props.Get("Format"))
	a.Equal("test-namespace", props.Get("Namespace"))
	a.Equal("test-version", props.Get("Version"))
	a.Equal("test-status", props.Get("Status"))
}

func Test_Mock_CodeArtifactPackageVersion_String(t *testing.T) {
	a := assert.New(t)
	r := &CodeArtifactPackageVersion{
		Domain:     ptr.String("test-domain"),
		Repository: ptr.String("test-repository"),
		Package:    ptr.String("test-package"),
		Version:    ptr.String("test-version"),
	}
	a.Equal(fmt.Sprintf("%s:%s/%s@%s", "test-domain", "test-repository", "test-package", "test-version"), r.String())
}
