package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	codeartifacttypes "github.com/aws/aws-sdk-go-v2/service/codeartifact/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CodeArtifactPackageVersionResource = "CodeArtifactPackageVersion"

func init() {
	registry.Register(&registry.Registration{
		Name:     CodeArtifactPackageVersionResource,
		Scope:    nuke.Account,
		Resource: &CodeArtifactPackageVersion{},
		Lister:   &CodeArtifactPackageVersionLister{},
	})
}

type CodeArtifactPackageVersionLister struct {
	svc CodeartifactClient
}

func (l *CodeArtifactPackageVersionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = codeartifact.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	domainPaginator := codeartifact.NewListDomainsPaginator(svc, &codeartifact.ListDomainsInput{})
	for domainPaginator.HasMorePages() {
		domainResp, err := domainPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, domain := range domainResp.Domains {
			repoPaginator := codeartifact.NewListRepositoriesInDomainPaginator(svc, &codeartifact.ListRepositoriesInDomainInput{
				Domain: domain.Name,
			})
			for repoPaginator.HasMorePages() {
				repoResp, err := repoPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, repo := range repoResp.Repositories {
					pkgPaginator := codeartifact.NewListPackagesPaginator(svc, &codeartifact.ListPackagesInput{
						Domain:     domain.Name,
						Repository: repo.Name,
					})
					for pkgPaginator.HasMorePages() {
						pkgResp, err := pkgPaginator.NextPage(ctx)
						if err != nil {
							return nil, err
						}
						for _, pkg := range pkgResp.Packages {
							versionPaginator := codeartifact.NewListPackageVersionsPaginator(svc, &codeartifact.ListPackageVersionsInput{
								Domain:     domain.Name,
								Repository: repo.Name,
								Package:    pkg.Package,
								Format:     pkg.Format,
								Namespace:  pkg.Namespace,
							})
							for versionPaginator.HasMorePages() {
								versionResp, err := versionPaginator.NextPage(ctx)
								if err != nil {
									return nil, err
								}
								for _, version := range versionResp.Versions {
									resources = append(resources, &CodeArtifactPackageVersion{
										svc:        svc,
										Domain:     domain.Name,
										Repository: repo.Name,
										Package:    pkg.Package,
										Format:     aws.String(string(pkg.Format)),
										Namespace:  pkg.Namespace,
										Version:    version.Version,
										Status:     aws.String(string(version.Status)),
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return resources, nil
}

type CodeArtifactPackageVersion struct {
	svc        CodeartifactClient
	Domain     *string
	Repository *string
	Package    *string
	Format     *string
	Namespace  *string
	Version    *string
	Status     *string
}

func (r *CodeArtifactPackageVersion) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePackageVersions(ctx, &codeartifact.DeletePackageVersionsInput{
		Domain:     r.Domain,
		Repository: r.Repository,
		Package:    r.Package,
		Format:     codeartifacttypes.PackageFormat(aws.ToString(r.Format)),
		Namespace:  r.Namespace,
		Versions:   []string{aws.ToString(r.Version)},
	})
	return err
}

func (r *CodeArtifactPackageVersion) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CodeArtifactPackageVersion) String() string {
	return fmt.Sprintf("%s:%s/%s@%s", *r.Domain, *r.Repository, *r.Package, *r.Version)
}
